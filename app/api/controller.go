package api

import (
	"net/http"

	"github.com/archaron/secimage/app/api/youtube"
	"github.com/im-kulikov/helium/module"
	"github.com/im-kulikov/helium/settings"
	"github.com/im-kulikov/helium/web"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type (
	JSON   = map[string]string
	Params struct {
		dig.In

		App    *settings.App
		Logger *zap.Logger
		Viper  *viper.Viper
		//Nats   *nats.Client
	}
)

var Module = module.Module{
	{Constructor: Router},
}

// New - Create new REST-API
func Router(params Params) (http.Handler, error) {
	v := web.NewValidator()
	b := web.NewBinder(v)
	e := web.NewEngine(web.EngineParams{
		Config:    params.Viper,
		Binder:    b,
		Logger:    params.Logger,
		Validator: v,
	})

	// Version:
	e.GET("/version", version(JSON{
		"time":    params.App.BuildTime,
		"version": params.App.BuildVersion,
	}))

	// Ajax:
	if err := youtube.Attach(e, params.Logger.Sugar(), params.Viper); err != nil {
		return nil, err
	}

	return e, nil
}

// Version of the application
func version(data JSON) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, data)
	}
}
