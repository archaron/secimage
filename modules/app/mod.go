package app

import (
	"github.com/archaron/secimage/modules/api"
	"github.com/go-helium/echo"
	"github.com/im-kulikov/helium/grace"
	"github.com/im-kulikov/helium/logger"
	"github.com/im-kulikov/helium/module"
	"github.com/im-kulikov/helium/settings"
	"github.com/im-kulikov/helium/web"
)

// Module application
var Module = module.New(New).Append(
	grace.Module,      // graceful context
	settings.Module,   // settings
	logger.Module,     // logger
	web.ServersModule, // web-servers
	echo.Module,       // echo module
	api.Module,        // API
)
