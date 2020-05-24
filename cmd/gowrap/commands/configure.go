package commands

import (
	"github.com/alecthomas/kingpin"
	"github.com/xabierlaiseca/gowrap/pkg/config"
	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
)

func newConfigureCommand(app *kingpin.Application, gowrapHome string) {
	cmd := app.Command("configure", "configuration related operations")
	newConfigureDefaultCommand(cmd, gowrapHome)
	newConfigurationUpgradesCommand(cmd, gowrapHome)
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

func newConfigurationUpgradesCommand(parent *kingpin.CmdClause, gowrapHome string) {
	cmd := parent.Command("upgrades", "Configure the behaviour on how to upgrade go versions")
	upgradesType := cmd.Arg("type", "how to upgrade go versions").
		Required().
		HintOptions(config.UpgradesAuto, config.UpgradesAsk, config.UpgradesDisabled).
		String()

	cmd.Action(func(*kingpin.ParseContext) error {
		c, err := config.Load(gowrapHome)
		if err != nil {
			return err
		}

		c.Upgrades = *upgradesType
		return c.Save()
	})
}

func newConfigurationSelfUpgradesCommand(parent *kingpin.CmdClause, gowrapHome string) {
	cmd := parent.Command("selfupgrades", "Configure the behaviour on how to upgrade gowrap versions")
	selfUpgradesType := cmd.Arg("type", "how to upgrade gowrap versions").
		Required().
		HintOptions(config.SelfUpgradesEnabled, config.SelfUpgradesDisabled).
		String()

	cmd.Action(func(*kingpin.ParseContext) error {
		c, err := config.Load(gowrapHome)
		if err != nil {
			return err
		}

		c.SelfUpgrades = *selfUpgradesType
		return c.Save()
	})
}
