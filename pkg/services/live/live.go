package live

import (
	"fmt"
	"sync"

	"github.com/centrifugal/centrifuge"
	"github.com/openinsight-project/grafinsight/pkg/api/routing"
	"github.com/openinsight-project/grafinsight/pkg/infra/log"
	"github.com/openinsight-project/grafinsight/pkg/models"
	"github.com/openinsight-project/grafinsight/pkg/plugins"
	"github.com/openinsight-project/grafinsight/pkg/registry"
	"github.com/openinsight-project/grafinsight/pkg/services/live/features"
	"github.com/openinsight-project/grafinsight/pkg/setting"
	"github.com/openinsight-project/grafinsight/pkg/tsdb/cloudwatch"
)

var (
	logger   = log.New("live")
	loggerCF = log.New("live.centrifuge")
)

func init() {
	registry.RegisterService(&GrafinsightLive{
		channels:   make(map[string]models.ChannelHandler),
		channelsMu: sync.RWMutex{},
		GrafinsightScope: CoreGrafinsightScope{
			Features: make(map[string]models.ChannelHandlerFactory),
		},
	})
}

// CoreGrafinsightScope list of core features
type CoreGrafinsightScope struct {
	Features map[string]models.ChannelHandlerFactory

	// The generic service to advertise dashboard changes
	Dashboards models.DashboardActivityChannel
}

// GrafinsightLive pretends to be the server
type GrafinsightLive struct {
	Cfg           *setting.Cfg            `inject:""`
	RouteRegister routing.RouteRegister   `inject:""`
	LogsService   *cloudwatch.LogsService `inject:""`
	node          *centrifuge.Node

	// The websocket handler
	WebsocketHandler interface{}

	// Full channel handler
	channels   map[string]models.ChannelHandler
	channelsMu sync.RWMutex

	// The core internal features
	GrafinsightScope CoreGrafinsightScope
}

// Init initializes the instance.
// Required to implement the registry.Service interface.
func (g *GrafinsightLive) Init() error {
	logger.Debug("GrafinsightLive initialization")

	if !g.IsEnabled() {
		logger.Debug("GrafinsightLive feature not enabled, skipping initialization")
		return nil
	}

	// We use default config here as starting point. Default config contains
	// reasonable values for available options.
	cfg := centrifuge.DefaultConfig

	// cfg.LogLevel = centrifuge.LogLevelDebug
	cfg.LogHandler = handleLog

	// Node is the core object in Centrifuge library responsible for many useful
	// things. For example Node allows to publish messages to channels from server
	// side with its Publish method, but in this example we will publish messages
	// only from client side.
	node, err := centrifuge.New(cfg)
	if err != nil {
		return err
	}
	g.node = node

	// Initialize the main features
	dash := &features.DashboardHandler{
		Publisher: g.Publish,
	}

	g.GrafinsightScope.Dashboards = dash
	g.GrafinsightScope.Features["dashboard"] = dash
	g.GrafinsightScope.Features["testdata"] = &features.TestDataSupplier{
		Publisher: g.Publish,
	}
	g.GrafinsightScope.Features["broadcast"] = &features.BroadcastRunner{}
	g.GrafinsightScope.Features["measurements"] = &features.MeasurementsRunner{}

	// Set ConnectHandler called when client successfully connected to Node. Your code
	// inside handler must be synchronized since it will be called concurrently from
	// different goroutines (belonging to different client connections). This is also
	// true for other event handlers.
	node.OnConnect(func(client *centrifuge.Client) {
		logger.Debug("Client connected", "user", client.UserID())

		client.OnSubscribe(func(e centrifuge.SubscribeEvent, cb centrifuge.SubscribeCallback) {
			handler, err := g.GetChannelHandler(e.Channel)
			if err != nil {
				cb(centrifuge.SubscribeReply{}, err)
			} else {
				cb(handler.OnSubscribe(client, e))
			}
		})

		// Called when a client publishes to the websocket channel.
		// In general, we should prefer writing to the HTTP API, but this
		// allows some simple prototypes to work quickly.
		client.OnPublish(func(e centrifuge.PublishEvent, cb centrifuge.PublishCallback) {
			handler, err := g.GetChannelHandler(e.Channel)
			if err != nil {
				cb(centrifuge.PublishReply{}, err)
			} else {
				cb(handler.OnPublish(client, e))
			}
		})
	})

	// Run node. This method does not block.
	if err := node.Run(); err != nil {
		return err
	}

	// Use a pure websocket transport.
	wsHandler := centrifuge.NewWebsocketHandler(node, centrifuge.WebsocketConfig{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})

	g.WebsocketHandler = func(ctx *models.ReqContext) {
		user := ctx.SignedInUser
		if user == nil {
			ctx.Resp.WriteHeader(401)
			return
		}

		// Centrifuge expects Credentials in context with a current user ID.
		cred := &centrifuge.Credentials{
			UserID: fmt.Sprintf("%d", user.UserId),
		}
		newCtx := centrifuge.SetCredentials(ctx.Req.Context(), cred)

		r := ctx.Req.Request
		r = r.WithContext(newCtx) // Set a user ID.

		wsHandler.ServeHTTP(ctx.Resp, r)
	}

	g.RouteRegister.Get("/live/ws", g.WebsocketHandler)

	return nil
}

// GetChannelHandler gives threadsafe access to the channel
func (g *GrafinsightLive) GetChannelHandler(channel string) (models.ChannelHandler, error) {
	g.channelsMu.RLock()
	c, ok := g.channels[channel]
	g.channelsMu.RUnlock() // defer? but then you can't lock further down
	if ok {
		return c, nil
	}

	// Parse the identifier ${scope}/${namespace}/${path}
	addr := ParseChannelAddress(channel)
	if !addr.IsValid() {
		return nil, fmt.Errorf("invalid channel: %q", channel)
	}
	logger.Info("initChannel", "channel", channel, "address", addr)

	g.channelsMu.Lock()
	defer g.channelsMu.Unlock()
	c, ok = g.channels[channel] // may have filled in while locked
	if ok {
		return c, nil
	}

	getter, err := g.GetChannelHandlerFactory(addr.Scope, addr.Namespace)
	if err != nil {
		return nil, err
	}

	// First access will initialize
	c, err = getter.GetHandlerForPath(addr.Path)
	if err != nil {
		return nil, err
	}

	g.channels[channel] = c
	return c, nil
}

// GetChannelHandlerFactory gets a ChannelHandlerFactory for a namespace.
// It gives threadsafe access to the channel.
func (g *GrafinsightLive) GetChannelHandlerFactory(scope string, name string) (models.ChannelHandlerFactory, error) {
	if scope == "grafinsight" {
		p, ok := g.GrafinsightScope.Features[name]
		if ok {
			return p, nil
		}
		return nil, fmt.Errorf("unknown feature: %q", name)
	}

	if scope == "ds" {
		return nil, fmt.Errorf("todo... look up datasource: %q", name)
	}

	if scope == "plugin" {
		// Temporary hack until we have a more generic solution later on
		if name == "cloudwatch" {
			return &cloudwatch.LogQueryRunnerSupplier{
				Publisher: g.Publish,
				Service:   g.LogsService,
			}, nil
		}

		p, ok := plugins.Plugins[name]
		if ok {
			h := &PluginHandler{
				Plugin: p,
			}
			return h, nil
		}
		return nil, fmt.Errorf("unknown plugin: %q", name)
	}

	return nil, fmt.Errorf("invalid scope: %q", scope)
}

// Publish sends the data to the channel without checking permissions etc
func (g *GrafinsightLive) Publish(channel string, data []byte) error {
	_, err := g.node.Publish(channel, data)
	return err
}

// IsEnabled returns true if the Grafinsight Live feature is enabled.
func (g *GrafinsightLive) IsEnabled() bool {
	return g.Cfg.IsLiveEnabled()
}

// Write to the standard log15 logger
func handleLog(msg centrifuge.LogEntry) {
	arr := make([]interface{}, 0)
	for k, v := range msg.Fields {
		if v == nil {
			v = "<nil>"
		} else if v == "" {
			v = "<empty>"
		}
		arr = append(arr, k, v)
	}

	switch msg.Level {
	case centrifuge.LogLevelDebug:
		loggerCF.Debug(msg.Message, arr...)
	case centrifuge.LogLevelError:
		loggerCF.Error(msg.Message, arr...)
	case centrifuge.LogLevelInfo:
		loggerCF.Info(msg.Message, arr...)
	default:
		loggerCF.Debug(msg.Message, arr...)
	}
}
