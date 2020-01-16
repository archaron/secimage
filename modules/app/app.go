package app

import (
	"context"
	"net/http"

	"github.com/im-kulikov/helium"
	"github.com/im-kulikov/helium/web"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type (
	// App instance
	App struct {
		log *zap.Logger
		srv web.Service
	}

	// Params to create new app instance
	Params struct {
		dig.In

		Handler http.Handler
		Logger  *zap.Logger
		Service web.Service
	}
)

// New creates instance
func New(params Params) helium.App {
	return &App{
		log: params.Logger,
		srv: params.Service,
	}
}

func (a *App) Run(ctx context.Context) error {
	a.log.Info("running servers")
	if err := a.srv.Start(); err != nil {
		return err
	}

	a.log.Info("app successfully started")
	<-ctx.Done()

	a.log.Info("stopping http servers",
		zap.Error(a.srv.Stop()))

	a.log.Info("gracefully stopped")
	return nil
}
