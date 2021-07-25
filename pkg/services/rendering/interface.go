package rendering

import (
	"context"
	"errors"
	"time"

	"github.com/openinsight-project/grafinsight/pkg/models"
)

var ErrTimeout = errors.New("timeout error - you can set timeout in seconds with &timeout url parameter")
var ErrPhantomJSNotInstalled = errors.New("PhantomJS executable not found")

type Opts struct {
	Width             int
	Height            int
	Timeout           time.Duration
	OrgId             int64
	UserId            int64
	OrgRole           models.RoleType
	Path              string
	Encoding          string
	Timezone          string
	ConcurrentLimit   int
	DeviceScaleFactor float64
	Headers           map[string][]string
}

type RenderResult struct {
	FilePath string
}

type renderFunc func(ctx context.Context, renderKey string, options Opts) (*RenderResult, error)

type Service interface {
	IsAvailable() bool
	Render(ctx context.Context, opts Opts) (*RenderResult, error)
	RenderErrorImage(error error) (*RenderResult, error)
	GetRenderUser(key string) (*RenderUser, bool)
}
