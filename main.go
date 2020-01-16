package main

import (
	"github.com/archaron/secimage/misc"
	"github.com/archaron/secimage/modules/app"
	"github.com/im-kulikov/helium"
)

// Replace with helium.Catch for production:
var check = helium.CatchTrace

func main() {
	h, err := helium.New(&helium.Settings{
		File:         misc.Config,
		Prefix:       misc.Prefix,
		Name:         misc.Name,
		BuildTime:    misc.Build,
		BuildVersion: misc.Version,
	}, app.Module)

	check(err)
	check(h.Run())
}
