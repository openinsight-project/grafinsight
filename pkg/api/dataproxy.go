package api

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/openinsight-project/grafinsight/pkg/api/datasource"
	"github.com/openinsight-project/grafinsight/pkg/api/pluginproxy"
	"github.com/openinsight-project/grafinsight/pkg/infra/metrics"
	"github.com/openinsight-project/grafinsight/pkg/models"
	"github.com/openinsight-project/grafinsight/pkg/plugins"
)

// ProxyDataSourceRequest proxies datasource requests
func (hs *HTTPServer) ProxyDataSourceRequest(c *models.ReqContext) {
	c.TimeRequest(metrics.MDataSourceProxyReqTimer)

	dsID := c.ParamsInt64(":id")
	ds, err := hs.DatasourceCache.GetDatasource(dsID, c.SignedInUser, c.SkipCache)
	if err != nil {
		if errors.Is(err, models.ErrDataSourceAccessDenied) {
			c.JsonApiErr(http.StatusForbidden, "Access denied to datasource", err)
			return
		}
		c.JsonApiErr(http.StatusInternalServerError, "Unable to load datasource meta data", err)
		return
	}

	err = hs.PluginRequestValidator.Validate(ds.Url, c.Req.Request)
	if err != nil {
		c.JsonApiErr(http.StatusForbidden, "Access denied", err)
		return
	}

	// find plugin
	plugin, ok := plugins.DataSources[ds.Type]
	if !ok {
		c.JsonApiErr(http.StatusInternalServerError, "Unable to find datasource plugin", err)
		return
	}

	proxyPath := getProxyPath(c)
	proxy, err := pluginproxy.NewDataSourceProxy(ds, plugin, c, proxyPath, hs.Cfg)
	if err != nil {
		if errors.Is(err, datasource.URLValidationError{}) {
			c.JsonApiErr(http.StatusBadRequest, fmt.Sprintf("Invalid data source URL: %q", ds.Url), err)
		} else {
			c.JsonApiErr(http.StatusInternalServerError, "Failed creating data source proxy", err)
		}
		return
	}
	proxy.HandleRequest()
}

var proxyPathRegexp = regexp.MustCompile(`^\/api\/datasources\/proxy\/[\d]+\/?`)

func extractProxyPath(originalRawPath string) string {
	return proxyPathRegexp.ReplaceAllString(originalRawPath, "")
}

func getProxyPath(c *models.ReqContext) string {
	return extractProxyPath(c.Req.URL.EscapedPath())
}
