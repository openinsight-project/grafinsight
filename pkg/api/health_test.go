package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/openinsight-project/grafinsight/pkg/bus"
	"github.com/openinsight-project/grafinsight/pkg/infra/localcache"
	"github.com/openinsight-project/grafinsight/pkg/models"
	"github.com/openinsight-project/grafinsight/pkg/setting"
	"github.com/stretchr/testify/require"
	macaron "gopkg.in/macaron.v1"
)

func TestHealthAPI_Version(t *testing.T) {
	m, _ := setupHealthAPITestEnvironment(t)
	setting.BuildVersion = "7.4.0"
	setting.BuildCommit = "59906ab1bf"

	bus.AddHandler("test", func(query *models.GetDBHealthQuery) error {
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	m.ServeHTTP(rec, req)

	require.Equal(t, 200, rec.Code)
	expectedBody := `
		{
			"database": "ok",
			"version": "7.4.0",
			"commit": "59906ab1bf"
		}
	`
	require.JSONEq(t, expectedBody, rec.Body.String())
}

func TestHealthAPI_AnonymousHideVersion(t *testing.T) {
	m, hs := setupHealthAPITestEnvironment(t)
	hs.Cfg.AnonymousHideVersion = true

	bus.AddHandler("test", func(query *models.GetDBHealthQuery) error {
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	m.ServeHTTP(rec, req)

	require.Equal(t, 200, rec.Code)
	expectedBody := `
		{
			"database": "ok"
		}
	`
	require.JSONEq(t, expectedBody, rec.Body.String())
}

func TestHealthAPI_DatabaseHealthy(t *testing.T) {
	const cacheKey = "db-healthy"

	m, hs := setupHealthAPITestEnvironment(t)
	hs.Cfg.AnonymousHideVersion = true

	bus.AddHandler("test", func(query *models.GetDBHealthQuery) error {
		return nil
	})

	healthy, found := hs.CacheService.Get(cacheKey)
	require.False(t, found)
	require.Nil(t, healthy)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	m.ServeHTTP(rec, req)

	require.Equal(t, 200, rec.Code)
	expectedBody := `
		{
			"database": "ok"
		}
	`
	require.JSONEq(t, expectedBody, rec.Body.String())

	healthy, found = hs.CacheService.Get(cacheKey)
	require.True(t, found)
	require.True(t, healthy.(bool))
}

func TestHealthAPI_DatabaseUnhealthy(t *testing.T) {
	const cacheKey = "db-healthy"

	m, hs := setupHealthAPITestEnvironment(t)
	hs.Cfg.AnonymousHideVersion = true

	bus.AddHandler("test", func(query *models.GetDBHealthQuery) error {
		return errors.New("bad")
	})

	healthy, found := hs.CacheService.Get(cacheKey)
	require.False(t, found)
	require.Nil(t, healthy)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	m.ServeHTTP(rec, req)

	require.Equal(t, 503, rec.Code)
	expectedBody := `
		{
			"database": "failing"
		}
	`
	require.JSONEq(t, expectedBody, rec.Body.String())

	healthy, found = hs.CacheService.Get(cacheKey)
	require.True(t, found)
	require.False(t, healthy.(bool))
}

func TestHealthAPI_DatabaseHealthCached(t *testing.T) {
	const cacheKey = "db-healthy"

	m, hs := setupHealthAPITestEnvironment(t)
	hs.Cfg.AnonymousHideVersion = true

	// Database is healthy.
	bus.AddHandler("test", func(query *models.GetDBHealthQuery) error {
		return nil
	})

	// Mock unhealthy database in cache.
	hs.CacheService.Set(cacheKey, false, 5*time.Minute)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	m.ServeHTTP(rec, req)

	require.Equal(t, 503, rec.Code)
	expectedBody := `
		{
			"database": "failing"
		}
	`
	require.JSONEq(t, expectedBody, rec.Body.String())

	// Purge cache and redo request.
	hs.CacheService.Delete(cacheKey)
	rec = httptest.NewRecorder()
	m.ServeHTTP(rec, req)

	require.Equal(t, 200, rec.Code)
	expectedBody = `
		{
			"database": "ok"
		}
	`
	require.JSONEq(t, expectedBody, rec.Body.String())

	healthy, found := hs.CacheService.Get(cacheKey)
	require.True(t, found)
	require.True(t, healthy.(bool))
}

func setupHealthAPITestEnvironment(t *testing.T) (*macaron.Macaron, *HTTPServer) {
	t.Helper()

	oldVersion := setting.BuildVersion
	oldCommit := setting.BuildCommit
	t.Cleanup(func() {
		setting.BuildVersion = oldVersion
		setting.BuildCommit = oldCommit
	})

	bus.ClearBusHandlers()
	t.Cleanup(bus.ClearBusHandlers)

	m := macaron.New()
	hs := &HTTPServer{
		CacheService: localcache.New(5*time.Minute, 10*time.Minute),
		Cfg:          setting.NewCfg(),
	}

	m.Get("/api/health", hs.apiHealthHandler)
	return m, hs
}
