package web

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/openinsight-project/grafinsight/pkg/tests/testinfra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIndexView tests the Grafinsight index view.
func TestIndexView(t *testing.T) {
	t.Run("CSP enabled", func(t *testing.T) {
		grafDir, cfgPath := testinfra.CreateGrafDir(t, testinfra.GrafinsightOpts{
			EnableCSP: true,
		})
		sqlStore := testinfra.SetUpDatabase(t, grafDir)
		addr := testinfra.StartGrafinsight(t, grafDir, cfgPath, sqlStore)

		// nolint:bodyclose
		resp, html := makeRequest(t, addr)

		assert.Regexp(t, "script-src 'unsafe-eval' 'strict-dynamic' 'nonce-[^']+';object-src 'none';font-src 'self';style-src 'self' 'unsafe-inline';img-src 'self' data:;base-uri 'self';connect-src 'self' grafinsight.com;manifest-src 'self';media-src 'none';form-action 'self';", resp.Header.Get("Content-Security-Policy"))
		assert.Regexp(t, `<script nonce="[^"]+"`, html)
	})

	t.Run("CSP disabled", func(t *testing.T) {
		grafDir, cfgPath := testinfra.CreateGrafDir(t)
		sqlStore := testinfra.SetUpDatabase(t, grafDir)
		addr := testinfra.StartGrafinsight(t, grafDir, cfgPath, sqlStore)

		// nolint:bodyclose
		resp, html := makeRequest(t, addr)

		assert.Empty(t, resp.Header.Get("Content-Security-Policy"))
		assert.Regexp(t, `<script nonce=""`, html)
	})
}

func makeRequest(t *testing.T, addr string) (*http.Response, string) {
	t.Helper()

	u := fmt.Sprintf("http://%s", addr)
	t.Logf("Making GET request to %s", u)
	// nolint:gosec
	resp, err := http.Get(u)
	require.NoError(t, err)
	require.NotNil(t, resp)
	t.Cleanup(func() {
		err := resp.Body.Close()
		assert.NoError(t, err)
	})

	var b strings.Builder
	_, err = io.Copy(&b, resp.Body)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	return resp, b.String()
}
