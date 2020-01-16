package api

import (
	"net/http"

	"github.com/archaron/secimage/misc"
	"github.com/im-kulikov/helium/module"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type (
	JSON   = map[string]string
	Params struct {
		dig.In

		Logger *zap.Logger
		Engine *echo.Echo
		Viper  *viper.Viper
		// Nats   *nats.Client
	}
)

var Module = module.New(Router)

// New - Create new REST-API
func Router(p Params) (http.Handler, error) {
	// Version:
	p.Engine.GET("/version", version(JSON{
		"time":    misc.Build,
		"version": misc.Version,
	}))

	cli := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}

	// Image request:
	handler, err := vi(viParams{
		log:          p.Logger,
		cli:          cli,
		savePath:     p.Viper.GetString("youtube.save_path"),
		cachePath:    p.Viper.GetString("youtube.cache_path"),
		allowedSizes: p.Viper.GetStringSlice("youtube.allowed_sizes"),
		quality:      p.Viper.GetInt("youtube.quality"),
		imageType:    p.Viper.GetString("youtube.type"),
	})
	if err != nil {
		return nil, err
	}

	p.Engine.GET("/:size/:id/:file", handler)

	return p.Engine, nil
}

// Version of the application
func version(data JSON) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, data)
	}
}
