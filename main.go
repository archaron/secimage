package main

import (
	"github.com/archaron/secimage/app"
	"github.com/im-kulikov/helium"
	"github.com/im-kulikov/helium/module"
	"github.com/im-kulikov/helium/settings"
	"github.com/urfave/cli"
	"go.uber.org/dig"
	"os"
)

const (
	name        = "secimage"
	description = "Image processing service"
	config      = "config.yml"
	version     = "1.0.0"
	buildTime   = "now"
)

func run(mod module.Module) cli.ActionFunc {
	return func(*cli.Context) error {
		h, err := helium.New(&settings.App{
			File:         config,
			Name:         name,
			BuildTime:    version,
			BuildVersion: buildTime,
		}, mod)

		if err != nil {
			return err
		}

		return h.Run()
	}
}

func main() {
	c := cli.NewApp()
	c.Name = name
	c.Version = version
	c.Description = description
	c.Commands = cli.Commands{
		{
			Name:      "serve",
			ShortName: "s",
			Action:    run(app.ServeModule),
		},
	}

	if err := c.Run(os.Args); err != nil {
		panic(dig.RootCause(err))
	}
}
