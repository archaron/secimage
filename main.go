package main

import (
	"github.com/archaron/secimage/misc"
	"github.com/archaron/secimage/mod/app"
	"github.com/davecgh/go-spew/spew"
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

	spew.Dump(err)
	check(err)
	check(h.Run())
}
