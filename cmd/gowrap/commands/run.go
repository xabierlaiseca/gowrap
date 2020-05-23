package commands

import (
	"fmt"

	"github.com/alecthomas/kingpin"
)

func RunCli(gowrapVersion, gowrapHome, wd string, args []string) error {
	app := kingpin.New("gowrap", "Utility to manage installed go versions")

	newConfigureCommand(app, gowrapHome)
	newInstallCommand(app, gowrapHome)
	newListCommand(app, gowrapHome)
	newProjectCommand(app, gowrapHome, wd)
	newUninstallCommand(app, gowrapHome)
	newVersionsFileCommand(app)

	app.Command("version", "Prints the gowrap version").
		Action(func(context *kingpin.ParseContext) error {
			fmt.Println(gowrapVersion)
			return nil
		})

	_, err := app.Parse(args)
	return err
}
