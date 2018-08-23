package app

import (
	"context"
	"github.com/archaron/secimage/app/api"
	"github.com/chapsuk/mserv"
	"github.com/chapsuk/worker"
	"github.com/im-kulikov/helium"
	"github.com/im-kulikov/helium/grace"
	"github.com/im-kulikov/helium/logger"
	"github.com/im-kulikov/helium/module"
	"github.com/im-kulikov/helium/settings"
	"github.com/im-kulikov/helium/web"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"net/http"
)

type (
	JSON = map[string]string

	// Workers interface
	Workers interface {
		Add(worker ...*worker.Worker)
		Run()
		Stop()
	}

	// App instance
	App struct {
		log     *zap.Logger
		servers mserv.Server
		//		workers  Workers
		receiver chan JSON
	}

	Health struct {
		v   *viper.Viper
		log *zap.SugaredLogger
	}

	Params struct {
		dig.In

		//		Workers *worker.Group
		Servers mserv.Server
		Logger  *zap.Logger
	}

	HealthParams struct {
		dig.In

		Logger *zap.Logger
		Viper  *viper.Viper
	}

	JobsParams struct {
		dig.In
	}
)

var (
	HealthModule = module.Module{
		// App specific modules
		{Constructor: health},
	}.
		Append(grace.Module).    // Graceful
		Append(settings.Module). // Settings
		Append(logger.Module)    // Logger

	ServeModule = module.Module{
		// App specific modules
		{Constructor: New},
	}.
		Append(grace.Module).      // graceful context
		Append(settings.Module).   // settings
		Append(logger.Module).     // logger
		Append(web.ServersModule). // web-servers
		Append(api.Module)         // API

)

// New creates instance
func New(params Params) helium.App {
	return &App{
		log:     params.Logger,
		servers: params.Servers,
	}
}

func (a *App) Run(ctx context.Context) error {

	a.log.Info("running servers")
	a.servers.Start()

	a.log.Info("app successfully runned")
	<-ctx.Done()

	a.log.Info("stopping http servers...")
	a.servers.Stop()

	a.log.Info("gracefully stopped")
	return nil
}

func health(params HealthParams) helium.App {
	return &Health{
		v:   params.Viper,
		log: params.Logger.Sugar(),
	}
}

func (h Health) Run(ctx context.Context) error {
	baseURL := h.v.GetString("api.address")

	req, err := http.NewRequest(http.MethodGet, "http://"+baseURL+"/version", nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		h.log.Panicw("bad status code", "status", resp.StatusCode)
	}

	return nil
}
