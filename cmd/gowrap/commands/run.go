package commands

import "github.com/alecthomas/kingpin"

func RunCli(gowrapHome, wd string, args []string) error {
	app := kingpin.New("gowrap", "Utility to manage installed go versions")

	newConfigureCommand(app, gowrapHome)
	newInstallCommand(app, gowrapHome)
	newListCommand(app, gowrapHome)
	newProjectCommand(app, wd)
	newUninstallCommand(app, gowrapHome)
	newVersionsFileCommand(app)

	_, err := app.Parse(args)
	return err
}
