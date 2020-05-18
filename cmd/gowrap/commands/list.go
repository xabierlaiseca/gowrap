package commands

import (
	"github.com/alecthomas/kingpin"
	"github.com/xabierlaiseca/gowrap/pkg/versions"
)

func newListCommand(app *kingpin.Application, gowrapHome string) {
	cmd := app.Command("list", "List operations")
	newListAvailableCommand(cmd)
	newListInstalledCommand(cmd, gowrapHome)
}

func newListAvailableCommand(parent *kingpin.CmdClause) {
	parent.Command("available", "Lists the available go versions to install").
		Action(func(*kingpin.ParseContext) error {
			return versions.PrintAvailable()
		})
}

func newListInstalledCommand(parent *kingpin.CmdClause, gowrapHome string) {
	parent.Command("installed", "Lists installed go versions").
		Action(func(*kingpin.ParseContext) error {
			return versions.PrintInstalled(gowrapHome)
		})
}
