package plugins

import (
	"net/http"
	"time"

	"github.com/openinsight-project/grafinsight/pkg/setting"
)

var (
	httpClient = http.Client{Timeout: 10 * time.Second}
)

type GrafinsightNetPlugin struct {
	Slug    string `json:"slug"`
	Version string `json:"version"`
}

type GithubLatest struct {
	Stable  string `json:"stable"`
	Testing string `json:"testing"`
}

func (pm *PluginManager) checkForUpdates() {
	if !setting.CheckForUpdates {
		return
	}
}
