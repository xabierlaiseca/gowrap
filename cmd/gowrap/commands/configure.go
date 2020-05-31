package commands

import (
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/xabierlaiseca/gowrap/pkg/config"
	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
)

func newConfigureCommand(app *kingpin.Application, gowrapHome string) {
	cmd := app.Command("configure", "configuration related operations")
	newConfigureDefaultCommand(cmd, gowrapHome)
	newConfigurationAutoInstallCommand(cmd, gowrapHome)
	newConfigurationSelfUpgradesCommand(cmd, gowrapHome)
}

func newConfigureDefaultCommand(parent *kingpin.CmdClause, gowrapHome string) {
	cmd := parent.Command("default", "Configure the default go version to use")
	version := cmd.Arg("version", "version to use as default").
		Required().
		HintAction(availableVersionCompletion).
		String()

	cmd.
		Validate(func(*kingpin.CmdClause) error {
			if len(*version) == 0 || semver.IsValid(*version) {
				return nil
			}
			return customerrors.Errorf("invalid version provided: %s", *version)
		})

	cmd.Action(func(*kingpin.ParseContext) error {
		return config.SetDefaultVersion(gowrapHome, *version)
	})
}

func newConfigurationAutoInstallCommand(parent *kingpin.CmdClause, gowrapHome string) {
	cmd := parent.Command("autoinstall", "Configure when gowrap tooling can automatically install go versions").
		HelpLong(fmt.Sprintf(
			"If autoinstall is set to '%s' the tool will automatically install latest available go version for the required version if not already installed, "+
				"if set to '%s' it will automatically install a go version only when a valid version is missing, "+
				"if set to '%s' it will never automatically install a go version", config.AutoInstallEnabled, config.AutoInstallMissing, config.AutoInstallDisabled))

	typeArg := cmd.Arg("type", "how to upgrade go versions").Required()
	autoInstallType := asEnum(typeArg, config.AutoInstallEnabled, config.AutoInstallMissing, config.AutoInstallDisabled)

	cmd.Action(func(*kingpin.ParseContext) error {
		c, err := config.Load(gowrapHome)
		if err != nil {
			return err
		}

		c.AutoInstall = *autoInstallType
		return c.Save()
	})
}

func newConfigurationSelfUpgradesCommand(parent *kingpin.CmdClause, gowrapHome string) {
	cmd := parent.Command("selfupgrade", "Configure the behaviour on how to upgrade gowrap versions")

	typeArg := cmd.Arg("type", "how to upgrade gowrap versions").Required()
	selfUpgradesType := asEnum(typeArg, config.SelfUpgradesEnabled, config.SelfUpgradesDisabled)

	cmd.Action(func(*kingpin.ParseContext) error {
		c, err := config.Load(gowrapHome)
		if err != nil {
			return err
		}

		c.SelfUpgrade = *selfUpgradesType
		return c.Save()
	})
}
