package expr

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/openinsight-project/grafinsight/pkg/bus"
	"github.com/openinsight-project/grafinsight/pkg/components/simplejson"
	"github.com/openinsight-project/grafinsight/pkg/models"
	"github.com/openinsight-project/grafinsight/pkg/tsdb"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	expressionsQuerySummary *prometheus.SummaryVec
)

func init() {
	expressionsQuerySummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "expressions_queries_duration_milliseconds",
			Help:       "Expressions query summary",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"status"},
	)

	prometheus.MustRegister(expressionsQuerySummary)
}

// WrapTransformData creates and executes transform requests
func (s *Service) WrapTransformData(ctx context.Context, query *tsdb.TsdbQuery) (*tsdb.Response, error) {
	sdkReq := &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{
			OrgID: query.User.OrgId,
		},
		Queries: []backend.DataQuery{},
	}

	for _, q := range query.Queries {
		modelJSON, err := q.Model.MarshalJSON()
		if err != nil {
			return nil, err
		}
		sdkReq.Queries = append(sdkReq.Queries, backend.DataQuery{
			JSON:          modelJSON,
			Interval:      time.Duration(q.IntervalMs) * time.Millisecond,
			RefID:         q.RefId,
			MaxDataPoints: q.MaxDataPoints,
			QueryType:     q.QueryType,
			TimeRange: backend.TimeRange{
				From: query.TimeRange.GetFromAsTimeUTC(),
				To:   query.TimeRange.GetToAsTimeUTC(),
			},
		})
	}
	pbRes, err := s.TransformData(ctx, sdkReq)
	if err != nil {
		return nil, err
	}

	tR := &tsdb.Response{
		Results: make(map[string]*tsdb.QueryResult, len(pbRes.Responses)),
	}
	for refID, res := range pbRes.Responses {
		tRes := &tsdb.QueryResult{
			RefId:      refID,
			Dataframes: tsdb.NewDecodedDataFrames(res.Frames),
		}
		// if len(res.JsonMeta) != 0 {
		// 	tRes.Meta = simplejson.NewFromAny(res.JsonMeta)
		// }
		if res.Error != nil {
			tRes.Error = res.Error
			tRes.ErrorString = res.Error.Error()
		}
		tR.Results[refID] = tRes
	}

	return tR, nil
}

// TransformData takes Queries which are either expressions nodes
// or are datasource requests.
func (s *Service) TransformData(ctx context.Context, req *backend.QueryDataRequest) (r *backend.QueryDataResponse, err error) {
	if s.isDisabled() {
		return nil, status.Error(codes.PermissionDenied, "Expressions are disabled")
	}

	start := time.Now()
	defer func() {
		var respStatus string
		switch {
		case err == nil:
			respStatus = "success"
		default:
			respStatus = "failure"
		}
		duration := float64(time.Since(start).Nanoseconds()) / float64(time.Millisecond)
		expressionsQuerySummary.WithLabelValues(respStatus).Observe(duration)
	}()

	// Build the pipeline from the request, checking for ordering issues (e.g. loops)
	// and parsing graph nodes from the queries.
	pipeline, err := s.BuildPipeline(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Execute the pipeline
	responses, err := s.ExecutePipeline(ctx, pipeline)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	// Get which queries have the Hide property so they those queries' results
	// can be excluded from the response.
	hidden, err := hiddenRefIDs(req.Queries)
	if err != nil {
		return nil, status.Error((codes.Internal), err.Error())
	}

	if len(hidden) != 0 {
		filteredRes := backend.NewQueryDataResponse()
		for refID, res := range responses.Responses {
			if _, ok := hidden[refID]; !ok {
				filteredRes.Responses[refID] = res
			}
		}
		responses = filteredRes
	}

	return responses, nil
}

func hiddenRefIDs(queries []backend.DataQuery) (map[string]struct{}, error) {
	hidden := make(map[string]struct{})

	for _, query := range queries {
		hide := struct {
			Hide bool `json:"hide"`
		}{}

		if err := json.Unmarshal(query.JSON, &hide); err != nil {
			return nil, err
		}

		if hide.Hide {
			hidden[query.RefID] = struct{}{}
		}
	}
	return hidden, nil
}

// QueryData is called used to query datasources that are not expression commands, but are used
// alongside expressions and/or are the input of an expression command.
func QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	if len(req.Queries) == 0 {
		return nil, fmt.Errorf("zero queries found in datasource request")
	}

	datasourceID := int64(0)
	var datasourceUID string

	if req.PluginContext.DataSourceInstanceSettings != nil {
		datasourceID = req.PluginContext.DataSourceInstanceSettings.ID
		datasourceUID = req.PluginContext.DataSourceInstanceSettings.UID
	}

	getDsInfo := &models.GetDataSourceQuery{
		OrgId: req.PluginContext.OrgID,
		Id:    datasourceID,
		Uid:   datasourceUID,
	}

	if err := bus.Dispatch(getDsInfo); err != nil {
		return nil, fmt.Errorf("could not find datasource: %w", err)
	}

	// Convert plugin-model (datasource) queries to tsdb queries
	queries := make([]*tsdb.Query, len(req.Queries))
	for i, query := range req.Queries {
		sj, err := simplejson.NewJson(query.JSON)
		if err != nil {
			return nil, err
		}
		queries[i] = &tsdb.Query{
			RefId:         query.RefID,
			IntervalMs:    query.Interval.Milliseconds(),
			MaxDataPoints: query.MaxDataPoints,
			QueryType:     query.QueryType,
			DataSource:    getDsInfo.Result,
			Model:         sj,
		}
	}

	// For now take Time Range from first query.
	timeRange := tsdb.NewTimeRange(strconv.FormatInt(req.Queries[0].TimeRange.From.Unix()*1000, 10), strconv.FormatInt(req.Queries[0].TimeRange.To.Unix()*1000, 10))

	tQ := &tsdb.TsdbQuery{
		TimeRange: timeRange,
		Queries:   queries,
	}

	// Execute the converted queries
	tsdbRes, err := tsdb.HandleRequest(ctx, getDsInfo.Result, tQ)
	if err != nil {
		return nil, err
	}
	// Convert tsdb results (map) to plugin-model/datasource (slice) results.
	// Only error, tsdb.Series, and encoded Dataframes responses are mapped.
	responses := make(map[string]backend.DataResponse, len(tsdbRes.Results))
	for refID, res := range tsdbRes.Results {
		pRes := backend.DataResponse{}
		if res.Error != nil {
			pRes.Error = res.Error
		}

		if res.Dataframes != nil {
			decoded, err := res.Dataframes.Decoded()
			if err != nil {
				return nil, err
			}
			pRes.Frames = decoded
			responses[refID] = pRes
			continue
		}

		for _, series := range res.Series {
			frame, err := tsdb.SeriesToFrame(series)
			frame.RefID = refID
			if err != nil {
				return nil, err
			}
			pRes.Frames = append(pRes.Frames, frame)
		}

		responses[refID] = pRes
	}
	return &backend.QueryDataResponse{
		Responses: responses,
	}, nil
}
