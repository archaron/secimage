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
		Logger *zap.SugaredLogger
		Engine *echo.Echo
		Viper  *viper.Viper
		//Nats   *nats.Client
	}
)

var Module = module.Module{
	{Constructor: Router},
	{Constructor: web.NewValidator},
	{Constructor: web.NewBinder},
	{Constructor: web.NewEngine},
}

// New - Create new REST-API
func Router(params Params) (http.Handler, error) {
	// Version:
	params.Engine.GET("/version", version(JSON{
		"time":    params.App.BuildTime,
		"version": params.App.BuildVersion,
	}))

	// Youtube actions:
	if err := youtube.Attach(youtube.Params{
		Engine: params.Engine,
		Logger: params.Logger,
		Viper:  params.Viper,
	}); err != nil {
		return nil, err
	}

	return params.Engine, nil
}

// Version of the application
func version(data JSON) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, data)
	}
}
