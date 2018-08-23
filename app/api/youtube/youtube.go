package youtube

import (
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
)

// , r *chan map[string]string
func Attach(e *echo.Echo, log *zap.SugaredLogger, v *viper.Viper) error {
	cli := &http.Client{}

	// Image request:
	e.GET("/:size/:id/:file", Vi(log, cli, v.GetString("youtube.save_path"), v.GetString("youtube.cache_path"), v.GetStringSlice("youtube.allowed_sizes"), v.GetInt("youtube.quality")))

	return nil
}
