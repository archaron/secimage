package app

import (
	"context"

	"github.com/archaron/secimage/app/api"
	"github.com/chapsuk/mserv"
	"github.com/im-kulikov/helium"
	"github.com/im-kulikov/helium/grace"
	"github.com/im-kulikov/helium/logger"
	"github.com/im-kulikov/helium/module"
	"github.com/im-kulikov/helium/settings"
	"github.com/im-kulikov/helium/web"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type (
	// App instance
	App struct {
		log     *zap.Logger
		servers mserv.Server
	}

	// Params to create new app instance
	Params struct {
		dig.In

		Servers mserv.Server
		Logger  *zap.Logger
	}
)

// Module application
var Module = module.New(New).Append(
	grace.Module,      // graceful context
	settings.Module,   // settings
	logger.Module,     // logger
	web.ServersModule, // web-servers
	api.Module,        // API
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
