package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/openinsight-project/grafinsight/pkg/api/response"
	"github.com/openinsight-project/grafinsight/pkg/api/routing"
	"github.com/openinsight-project/grafinsight/pkg/bus"
	"github.com/openinsight-project/grafinsight/pkg/infra/fs"
	"github.com/openinsight-project/grafinsight/pkg/infra/remotecache"
	"github.com/openinsight-project/grafinsight/pkg/models"
	"github.com/openinsight-project/grafinsight/pkg/registry"
	"github.com/openinsight-project/grafinsight/pkg/services/auth"
	"github.com/openinsight-project/grafinsight/pkg/services/contexthandler"
	"github.com/openinsight-project/grafinsight/pkg/services/rendering"
	"github.com/openinsight-project/grafinsight/pkg/services/sqlstore"
	"github.com/openinsight-project/grafinsight/pkg/setting"
	"github.com/stretchr/testify/require"
	"gopkg.in/macaron.v1"
)

func loggedInUserScenario(t *testing.T, desc string, url string, fn scenarioFunc) {
	loggedInUserScenarioWithRole(t, desc, "GET", url, url, models.ROLE_EDITOR, fn)
}

func loggedInUserScenarioWithRole(t *testing.T, desc string, method string, url string, routePattern string, role models.RoleType, fn scenarioFunc) {
	t.Run(fmt.Sprintf("%s %s", desc, url), func(t *testing.T) {
		t.Cleanup(bus.ClearBusHandlers)

		sc := setupScenarioContext(t, url)
		sc.defaultHandler = routing.Wrap(func(c *models.ReqContext) response.Response {
			sc.context = c
			sc.context.UserId = testUserID
			sc.context.OrgId = testOrgID
			sc.context.Login = testUserLogin
			sc.context.OrgRole = role
			if sc.handlerFunc != nil {
				return sc.handlerFunc(sc.context)
			}

			return nil
		})

		switch method {
		case "GET":
			sc.m.Get(routePattern, sc.defaultHandler)
		case "DELETE":
			sc.m.Delete(routePattern, sc.defaultHandler)
		}

		fn(sc)
	})
}

func anonymousUserScenario(t *testing.T, desc string, method string, url string, routePattern string, fn scenarioFunc) {
	t.Run(fmt.Sprintf("%s %s", desc, url), func(t *testing.T) {
		defer bus.ClearBusHandlers()

		sc := setupScenarioContext(t, url)
		sc.defaultHandler = routing.Wrap(func(c *models.ReqContext) response.Response {
			sc.context = c
			if sc.handlerFunc != nil {
				return sc.handlerFunc(sc.context)
			}

			return nil
		})

		switch method {
		case "GET":
			sc.m.Get(routePattern, sc.defaultHandler)
		case "DELETE":
			sc.m.Delete(routePattern, sc.defaultHandler)
		}

		fn(sc)
	})
}

func (sc *scenarioContext) fakeReq(method, url string) *scenarioContext {
	sc.resp = httptest.NewRecorder()
	req, err := http.NewRequest(method, url, nil)
	require.NoError(sc.t, err)
	sc.req = req

	return sc
}

func (sc *scenarioContext) fakeReqWithParams(method, url string, queryParams map[string]string) *scenarioContext {
	sc.resp = httptest.NewRecorder()
	req, err := http.NewRequest(method, url, nil)
	// TODO: Depend on sc.t
	if sc.t != nil {
		require.NoError(sc.t, err)
	} else if err != nil {
		panic(fmt.Sprintf("Making request failed: %s", err))
	}

	q := req.URL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	sc.req = req

	return sc
}

func (sc *scenarioContext) fakeReqNoAssertions(method, url string) *scenarioContext {
	sc.resp = httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, nil)
	sc.req = req

	return sc
}

func (sc *scenarioContext) fakeReqNoAssertionsWithCookie(method, url string, cookie http.Cookie) *scenarioContext {
	sc.resp = httptest.NewRecorder()
	http.SetCookie(sc.resp, &cookie)

	req, _ := http.NewRequest(method, url, nil)
	req.Header = http.Header{"Cookie": sc.resp.Header()["Set-Cookie"]}

	sc.req = req

	return sc
}

type scenarioContext struct {
	t                    *testing.T
	cfg                  *setting.Cfg
	m                    *macaron.Macaron
	context              *models.ReqContext
	resp                 *httptest.ResponseRecorder
	handlerFunc          handlerFunc
	defaultHandler       macaron.Handler
	req                  *http.Request
	url                  string
	userAuthTokenService *auth.FakeUserAuthTokenService
}

func (sc *scenarioContext) exec() {
	sc.m.ServeHTTP(sc.resp, sc.req)
}

type scenarioFunc func(c *scenarioContext)
type handlerFunc func(c *models.ReqContext) response.Response

func getContextHandler(t *testing.T, cfg *setting.Cfg) *contexthandler.ContextHandler {
	t.Helper()

	if cfg == nil {
		cfg = setting.NewCfg()
	}

	sqlStore := sqlstore.InitTestDB(t)
	remoteCacheSvc := &remotecache.RemoteCache{}
	cfg.RemoteCacheOptions = &setting.RemoteCacheOptions{
		Name: "database",
	}
	userAuthTokenSvc := auth.NewFakeUserAuthTokenService()
	renderSvc := &fakeRenderService{}
	ctxHdlr := &contexthandler.ContextHandler{}

	err := registry.BuildServiceGraph([]interface{}{cfg}, []*registry.Descriptor{
		{
			Name:     sqlstore.ServiceName,
			Instance: sqlStore,
		},
		{
			Name:     remotecache.ServiceName,
			Instance: remoteCacheSvc,
		},
		{
			Name:     auth.ServiceName,
			Instance: userAuthTokenSvc,
		},
		{
			Name:     rendering.ServiceName,
			Instance: renderSvc,
		},
		{
			Name:     contexthandler.ServiceName,
			Instance: ctxHdlr,
		},
	})
	require.NoError(t, err)

	return ctxHdlr
}

func setupScenarioContext(t *testing.T, url string) *scenarioContext {
	cfg := setting.NewCfg()
	sc := &scenarioContext{
		url: url,
		t:   t,
		cfg: cfg,
	}
	viewsPath, err := filepath.Abs("../../public/views")
	require.NoError(t, err)
	exists, err := fs.Exists(viewsPath)
	require.NoError(t, err)
	require.Truef(t, exists, "Views should be in %q", viewsPath)

	sc.m = macaron.New()
	sc.m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory: viewsPath,
		Delims:    macaron.Delims{Left: "[[", Right: "]]"},
	}))
	sc.m.Use(getContextHandler(t, cfg).Middleware)

	return sc
}

type fakeRenderService struct {
	rendering.Service
}

func (s *fakeRenderService) Init() error {
	return nil
}
