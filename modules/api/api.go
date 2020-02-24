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

	p.Viper.SetDefault("youtube.image_type", "hqdefault")
	p.Viper.SetDefault("youtube.input_format", "jpg")
	p.Viper.SetDefault("youtube.jpeg_quality", 50)
	p.Viper.SetDefault("youtube.webp_quality", 0.5)
	p.Viper.SetDefault("youtube.webp_lossless", false)

	// Image request:
	handler, err := vi(viParams{
		log:          p.Logger,
		cli:          cli,
		savePath:     p.Viper.GetString("youtube.save_path"),
		cachePath:    p.Viper.GetString("youtube.cache_path"),
		allowedSizes: p.Viper.GetStringSlice("youtube.allowed_sizes"),
		jpegQuality:  p.Viper.GetInt("youtube.jpeg_quality"),
		imageType:    p.Viper.GetString("youtube.image_type"),
		webpLossless: p.Viper.GetBool("youtube.webp_lossless"),
		inputFormat:  p.Viper.GetString("youtube.input_format"),
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
