package main

import (
	"github.com/conductant/gohm/pkg/command"
	"github.com/conductant/gohm/pkg/runtime"
	"github.com/conductant/kat-compose/pkg/aurora"
	"github.com/conductant/kat-compose/pkg/compose"
	"os"
)

func main() {

	command.Register("load",
		func() (command.Module, command.ErrorHandling) {
			return &compose.Load{
				Dump: true,
			}, command.PanicOnError
		})

	command.Register("convert",
		func() (command.Module, command.ErrorHandling) {
			// Defaults
			return &compose.Convert{
				Load: compose.Load{
					Dump: false,
				},
				Project: aurora.Project{
					Cluster:     "devcluster",
					Contact:     os.Getenv("USER") + "@localhost",
					Role:        os.Getenv("USER"),
					Environment: "production",
				},
			}, command.PanicOnError
		})

	runtime.Main()
}
