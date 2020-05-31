package commands

import (
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/xabierlaiseca/gowrap/cmd/common"
)

func RunCli(gowrapVersion, gowrapHome, wd string, args []string) error {
	app := kingpin.New("gowrap", "Utility to manage installed go versions").
		Action(selfUpgradeAction(gowrapVersion, gowrapHome))

	app.UsageTemplate("{{CustomUsage .App .Context.SelectedCommand}}")
	app.UsageFuncs(map[string]interface{}{
		"CustomUsage": usage,
	})
	app.HelpFlag.Help("Show context-sensitive help")

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

func selfUpgradeAction(currentVersion, gowrapHome string) func(context *kingpin.ParseContext) error {
	return func(context *kingpin.ParseContext) error {
		if context.String() == "configure selfupgrade" {
			return nil
		}

		common.SelfUpgrade(gowrapHome, currentVersion)
		return nil
	}
}
