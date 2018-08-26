package main

import (
	"github.com/archaron/secimage/app"
	"github.com/im-kulikov/helium"
	"github.com/im-kulikov/helium/settings"
	"go.uber.org/dig"
)

var (
	name      = "secimage"
	config    = "config.yml"
	Version   = "1.0.0"
	BuildTime = "now"
)

func main() {
	h, err := helium.New(&settings.App{
		File:         config,
		Name:         name,
		BuildTime:    BuildTime,
		BuildVersion: Version,
	}, app.Module)
	errCheck(err)

	err = h.Run()
	errCheck(err)
}

func errCheck(err error) {
	if err != nil {
		panic(dig.RootCause(err))
	}
}
