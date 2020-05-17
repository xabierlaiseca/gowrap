package main

import (
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/xabierlaiseca/gowrap/cmd/gowrap/commands"
)

func main() {
	app := kingpin.New("gowrap", "Utility to manage installed go versions")
	commands.NewConfigureCommand(app)
	commands.NewInstallCommand(app)
	commands.NewListCommand(app)
	commands.NewProjectCommand(app)
	commands.NewUninstallCommand(app)
	commands.NewVersionsFileCommand(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
