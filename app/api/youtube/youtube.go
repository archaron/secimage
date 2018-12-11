package youtube

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Params struct {
	Engine *echo.Echo
	Logger *zap.SugaredLogger
	Viper  *viper.Viper
}

// , r *chan map[string]string
func Attach(p Params) error {
	// creates clientm and pass proxy from envs:
	// - HTTP_PROXY  - for http requests
	// - HTTPS_PROXY - for https requests
	cli := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}

	// Image request:
	handler, err := Vi(ViParams{
		log:          p.Logger,
		cli:          cli,
		savePath:     p.Viper.GetString("youtube.save_path"),
		cachePath:    p.Viper.GetString("youtube.cache_path"),
		allowedSizes: p.Viper.GetStringSlice("youtube.allowed_sizes"),
		quality:      p.Viper.GetInt("youtube.quality"),
		imageType:    p.Viper.GetString("youtube.type"),
	})

	if err != nil {
		return err
	}

	p.Engine.GET("/:size/:id/:file", handler)

	return nil
}
