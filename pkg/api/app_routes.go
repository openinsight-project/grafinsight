package api

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/openinsight-project/grafinsight/pkg/api/pluginproxy"
	"github.com/openinsight-project/grafinsight/pkg/infra/log"
	"github.com/openinsight-project/grafinsight/pkg/middleware"
	"github.com/openinsight-project/grafinsight/pkg/models"
	"github.com/openinsight-project/grafinsight/pkg/plugins"
	"github.com/openinsight-project/grafinsight/pkg/util"
	macaron "gopkg.in/macaron.v1"
)

var pluginProxyTransport *http.Transport

func (hs *HTTPServer) initAppPluginRoutes(r *macaron.Macaron) {
	pluginProxyTransport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: hs.Cfg.PluginsAppsSkipVerifyTLS,
			Renegotiation:      tls.RenegotiateFreelyAsClient,
		},
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	for _, plugin := range plugins.Apps {
		for _, route := range plugin.Routes {
			url := util.JoinURLFragments("/api/plugin-proxy/"+plugin.Id, route.Path)
			handlers := make([]macaron.Handler, 0)
			handlers = append(handlers, middleware.Auth(&middleware.AuthOptions{
				ReqSignedIn: true,
			}))

			if route.ReqRole != "" {
				if route.ReqRole == models.ROLE_ADMIN {
					handlers = append(handlers, middleware.RoleAuth(models.ROLE_ADMIN))
				} else if route.ReqRole == models.ROLE_EDITOR {
					handlers = append(handlers, middleware.RoleAuth(models.ROLE_EDITOR, models.ROLE_ADMIN))
				}
			}
			handlers = append(handlers, AppPluginRoute(route, plugin.Id, hs))
			r.Route(url, route.Method, handlers...)
			log.Debugf("Plugins: Adding proxy route %s", url)
		}
	}
}

func AppPluginRoute(route *plugins.AppPluginRoute, appID string, hs *HTTPServer) macaron.Handler {
	return func(c *models.ReqContext) {
		path := c.Params("*")

		proxy := pluginproxy.NewApiPluginProxy(c, path, route, appID, hs.Cfg)
		proxy.Transport = pluginProxyTransport
		proxy.ServeHTTP(c.Resp, c.Req.Request)
	}
}
