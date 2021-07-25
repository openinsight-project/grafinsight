package testinfra

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/openinsight-project/grafinsight/pkg/infra/fs"
	"github.com/openinsight-project/grafinsight/pkg/registry"
	"github.com/openinsight-project/grafinsight/pkg/server"
	"github.com/openinsight-project/grafinsight/pkg/services/sqlstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/ini.v1"
)

// StartGrafana starts a Grafana server.
// The server address is returned.
func StartGrafana(t *testing.T, grafDir, cfgPath string, sqlStore *sqlstore.SQLStore) string {
	t.Helper()

	origSQLStore := registry.GetService(sqlstore.ServiceName)
	t.Cleanup(func() {
		registry.Register(origSQLStore)
	})
	registry.Register(&registry.Descriptor{
		Name:         sqlstore.ServiceName,
		Instance:     sqlStore,
		InitPriority: sqlstore.InitPriority,
	})

	t.Logf("Registered SQL store %p", sqlStore)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	server, err := server.New(server.Config{
		ConfigFile: cfgPath,
		HomePath:   grafDir,
		Listener:   listener,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		// Have to reset the route register between tests, since it doesn't get re-created
		server.HTTPServer.RouteRegister.Reset()
	})

	go func() {
		// When the server runs, it will also build and initialize the service graph
		if err := server.Run(); err != nil {
			t.Log("Server exited uncleanly", "error", err)
		}
	}()
	t.Cleanup(func() {
		server.Shutdown("")
	})

	// Wait for Grafana to be ready
	addr := listener.Addr().String()
	resp, err := http.Get(fmt.Sprintf("http://%s/api/health", addr))
	require.NoError(t, err)
	require.NotNil(t, resp)
	t.Cleanup(func() {
		err := resp.Body.Close()
		assert.NoError(t, err)
	})
	require.Equal(t, 200, resp.StatusCode)

	t.Logf("Grafana is listening on %s", addr)

	return addr
}

// SetUpDatabase sets up the Grafana database.
func SetUpDatabase(t *testing.T, grafDir string) *sqlstore.SQLStore {
	t.Helper()

	sqlStore := sqlstore.InitTestDB(t, sqlstore.InitTestDBOpt{
		EnsureDefaultOrgAndUser: true,
	})
	// We need the main org, since it's used for anonymous access
	org, err := sqlStore.GetOrgByName(sqlstore.MainOrgName)
	require.NoError(t, err)
	require.NotNil(t, org)

	// Make sure changes are synced with other goroutines
	err = sqlStore.Sync()
	require.NoError(t, err)

	return sqlStore
}

// CreateGrafDir creates the Grafana directory.
func CreateGrafDir(t *testing.T, opts ...GrafanaOpts) (string, string) {
	t.Helper()

	tmpDir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	t.Cleanup(func() {
		err := os.RemoveAll(tmpDir)
		assert.NoError(t, err)
	})

	// Search upwards in directory tree for project root
	var rootDir string
	found := false
	for i := 0; i < 20; i++ {
		rootDir = filepath.Join(rootDir, "..")
		exists, err := fs.Exists(filepath.Join(rootDir, "public", "views"))
		require.NoError(t, err)
		if exists {
			found = true
			break
		}
	}
	require.True(t, found, "Couldn't detect project root directory")

	cfgDir := filepath.Join(tmpDir, "conf")
	err = os.MkdirAll(cfgDir, 0750)
	require.NoError(t, err)
	dataDir := filepath.Join(tmpDir, "data")
	// nolint:gosec
	err = os.MkdirAll(dataDir, 0750)
	require.NoError(t, err)
	logsDir := filepath.Join(tmpDir, "logs")
	pluginsDir := filepath.Join(tmpDir, "plugins")
	publicDir := filepath.Join(tmpDir, "public")
	err = os.MkdirAll(publicDir, 0750)
	require.NoError(t, err)
	viewsDir := filepath.Join(publicDir, "views")
	err = fs.CopyRecursive(filepath.Join(rootDir, "public", "views"), viewsDir)
	require.NoError(t, err)
	// Copy index template to index.html, since Grafana will try to use the latter
	err = fs.CopyFile(filepath.Join(rootDir, "public", "views", "index-template.html"),
		filepath.Join(viewsDir, "index.html"))
	require.NoError(t, err)
	// Copy error template to error.html, since Grafana will try to use the latter
	err = fs.CopyFile(filepath.Join(rootDir, "public", "views", "error-template.html"),
		filepath.Join(viewsDir, "error.html"))
	require.NoError(t, err)
	emailsDir := filepath.Join(publicDir, "emails")
	err = fs.CopyRecursive(filepath.Join(rootDir, "public", "emails"), emailsDir)
	require.NoError(t, err)
	provDir := filepath.Join(cfgDir, "provisioning")
	provDSDir := filepath.Join(provDir, "datasources")
	err = os.MkdirAll(provDSDir, 0750)
	require.NoError(t, err)
	provNotifiersDir := filepath.Join(provDir, "notifiers")
	err = os.MkdirAll(provNotifiersDir, 0750)
	require.NoError(t, err)
	provPluginsDir := filepath.Join(provDir, "plugins")
	err = os.MkdirAll(provPluginsDir, 0750)
	require.NoError(t, err)
	provDashboardsDir := filepath.Join(provDir, "dashboards")
	err = os.MkdirAll(provDashboardsDir, 0750)
	require.NoError(t, err)

	cfg := ini.Empty()
	dfltSect := cfg.Section("")
	_, err = dfltSect.NewKey("app_mode", "development")
	require.NoError(t, err)

	pathsSect, err := cfg.NewSection("paths")
	require.NoError(t, err)
	_, err = pathsSect.NewKey("data", dataDir)
	require.NoError(t, err)
	_, err = pathsSect.NewKey("logs", logsDir)
	require.NoError(t, err)
	_, err = pathsSect.NewKey("plugins", pluginsDir)
	require.NoError(t, err)

	logSect, err := cfg.NewSection("log")
	require.NoError(t, err)
	_, err = logSect.NewKey("level", "debug")
	require.NoError(t, err)

	serverSect, err := cfg.NewSection("server")
	require.NoError(t, err)
	_, err = serverSect.NewKey("port", "0")
	require.NoError(t, err)

	anonSect, err := cfg.NewSection("auth.anonymous")
	require.NoError(t, err)
	_, err = anonSect.NewKey("enabled", "true")
	require.NoError(t, err)

	for _, o := range opts {
		if o.EnableCSP {
			securitySect, err := cfg.NewSection("security")
			require.NoError(t, err)
			_, err = securitySect.NewKey("content_security_policy", "true")
			require.NoError(t, err)
		}
	}

	cfgPath := filepath.Join(cfgDir, "test.ini")
	err = cfg.SaveTo(cfgPath)
	require.NoError(t, err)

	err = fs.CopyFile(filepath.Join(rootDir, "conf", "defaults.ini"), filepath.Join(cfgDir, "defaults.ini"))
	require.NoError(t, err)

	return tmpDir, cfgPath
}

type GrafanaOpts struct {
	EnableCSP bool
}
