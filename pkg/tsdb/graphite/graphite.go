package graphite

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/context/ctxhttp"

	"github.com/openinsight-project/grafinsight/pkg/infra/log"
	"github.com/openinsight-project/grafinsight/pkg/models"
	"github.com/openinsight-project/grafinsight/pkg/setting"
	"github.com/openinsight-project/grafinsight/pkg/tsdb"
	"github.com/opentracing/opentracing-go"
)

type GraphiteExecutor struct {
	HttpClient *http.Client
}

func NewGraphiteExecutor(datasource *models.DataSource) (tsdb.TsdbQueryEndpoint, error) {
	return &GraphiteExecutor{}, nil
}

var glog = log.New("tsdb.graphite")

func init() {
	tsdb.RegisterTsdbQueryEndpoint("graphite", NewGraphiteExecutor)
}

func (e *GraphiteExecutor) Query(ctx context.Context, dsInfo *models.DataSource, tsdbQuery *tsdb.TsdbQuery) (*tsdb.Response, error) {
	result := &tsdb.Response{}

	// This logic is used when called from Dashboard Alerting.
	from := "-" + formatTimeRange(tsdbQuery.TimeRange.From)
	until := formatTimeRange(tsdbQuery.TimeRange.To)

	// This logic is used when called through server side expressions.
	if isTimeRangeNumeric(tsdbQuery.TimeRange) {
		var err error
		from, until, err = epochMStoGraphiteTime(tsdbQuery.TimeRange)
		if err != nil {
			return nil, err
		}
	}

	var target string

	formData := url.Values{
		"from":          []string{from},
		"until":         []string{until},
		"format":        []string{"json"},
		"maxDataPoints": []string{"500"},
	}

	emptyQueries := make([]string, 0)
	for _, query := range tsdbQuery.Queries {
		glog.Debug("graphite", "query", query.Model)
		currTarget := ""
		if fullTarget, err := query.Model.Get("targetFull").String(); err == nil {
			currTarget = fullTarget
		} else {
			currTarget = query.Model.Get("target").MustString()
		}
		if currTarget == "" {
			glog.Debug("graphite", "empty query target", query.Model)
			emptyQueries = append(emptyQueries, fmt.Sprintf("Query: %v has no target", query.Model))
			continue
		}
		target = fixIntervalFormat(currTarget)
	}

	if target == "" {
		glog.Error("No targets in query model", "models without targets", strings.Join(emptyQueries, "\n"))
		return nil, errors.New("no query target found for the alert rule")
	}

	formData["target"] = []string{target}

	if setting.Env == setting.Dev {
		glog.Debug("Graphite request", "params", formData)
	}

	req, err := e.createRequest(dsInfo, formData)
	if err != nil {
		return nil, err
	}

	httpClient, err := dsInfo.GetHttpClient()
	if err != nil {
		return nil, err
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "graphite query")
	span.SetTag("target", target)
	span.SetTag("from", from)
	span.SetTag("until", until)
	span.SetTag("datasource_id", dsInfo.Id)
	span.SetTag("org_id", dsInfo.OrgId)

	defer span.Finish()

	if err := opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header)); err != nil {
		return nil, err
	}

	res, err := ctxhttp.Do(ctx, httpClient, req)
	if err != nil {
		return nil, err
	}

	data, err := e.parseResponse(res)
	if err != nil {
		return nil, err
	}

	result.Results = make(map[string]*tsdb.QueryResult)
	queryRes := tsdb.NewQueryResult()

	for _, series := range data {
		queryRes.Series = append(queryRes.Series, &tsdb.TimeSeries{
			Name:   series.Target,
			Points: series.DataPoints,
		})

		if setting.Env == setting.Dev {
			glog.Debug("Graphite response", "target", series.Target, "datapoints", len(series.DataPoints))
		}
	}

	result.Results["A"] = queryRes
	return result, nil
}

func (e *GraphiteExecutor) parseResponse(res *http.Response) ([]TargetResponseDTO, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			glog.Warn("Failed to close response body", "err", err)
		}
	}()

	if res.StatusCode/100 != 2 {
		glog.Info("Request failed", "status", res.Status, "body", string(body))
		return nil, fmt.Errorf("request failed, status: %s", res.Status)
	}

	var data []TargetResponseDTO
	err = json.Unmarshal(body, &data)
	if err != nil {
		glog.Info("Failed to unmarshal graphite response", "error", err, "status", res.Status, "body", string(body))
		return nil, err
	}

	for si := range data {
		// Convert Response to timestamps MS
		for pi, point := range data[si].DataPoints {
			data[si].DataPoints[pi][1].Float64 = point[1].Float64 * 1000
		}
	}
	return data, nil
}

func (e *GraphiteExecutor) createRequest(dsInfo *models.DataSource, data url.Values) (*http.Request, error) {
	u, err := url.Parse(dsInfo.Url)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, "render")

	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		glog.Info("Failed to create request", "error", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if dsInfo.BasicAuth {
		req.SetBasicAuth(dsInfo.BasicAuthUser, dsInfo.DecryptedBasicAuthPassword())
	}

	return req, err
}

func formatTimeRange(input string) string {
	if input == "now" {
		return input
	}
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(input, "now", ""), "m", "min"), "M", "mon")
}

func fixIntervalFormat(target string) string {
	rMinute := regexp.MustCompile(`'(\d+)m'`)
	target = rMinute.ReplaceAllStringFunc(target, func(m string) string {
		return strings.ReplaceAll(m, "m", "min")
	})
	rMonth := regexp.MustCompile(`'(\d+)M'`)
	target = rMonth.ReplaceAllStringFunc(target, func(M string) string {
		return strings.ReplaceAll(M, "M", "mon")
	})
	return target
}

func isTimeRangeNumeric(tr *tsdb.TimeRange) bool {
	if _, err := strconv.ParseInt(tr.From, 10, 64); err != nil {
		return false
	}
	if _, err := strconv.ParseInt(tr.To, 10, 64); err != nil {
		return false
	}
	return true
}

func epochMStoGraphiteTime(tr *tsdb.TimeRange) (string, string, error) {
	from, err := strconv.ParseInt(tr.From, 10, 64)
	if err != nil {
		return "", "", err
	}

	to, err := strconv.ParseInt(tr.To, 10, 64)
	if err != nil {
		return "", "", err
	}

	return fmt.Sprintf("%d", from/1000), fmt.Sprintf("%d", to/1000), nil
}
