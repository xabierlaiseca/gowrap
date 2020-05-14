package main

import (
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/xabierlaiseca/gowrap/cmd/gowrap/commands"
)

func main() {
	app := kingpin.New("gowrap", "Utility to manage installed go versions")
	commands.NewVersionsFileCommand(app)
	commands.NewListCommand(app)
	commands.NewInstallCommand(app)
	commands.NewUninstallCommand(app)
	commands.NewConfigureCommand(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
