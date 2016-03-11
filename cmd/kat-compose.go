package main

import (
	"github.com/conductant/gohm/pkg/command"
	"github.com/conductant/gohm/pkg/runtime"
	"github.com/conductant/kat-compose/pkg/compose"
)

func main() {

	command.Register("load",
		func() (command.Module, command.ErrorHandling) {
			return &compose.Load{
				Dump: true,
			}, command.PanicOnError
		})

	runtime.Main()
}
