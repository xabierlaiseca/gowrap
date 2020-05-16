package commands

import (
	"github.com/alecthomas/kingpin"
	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
	"github.com/xabierlaiseca/gowrap/pkg/versions"
)

func NewConfigureCommand(app *kingpin.Application) {
	cmd := app.Command("configure", "configuration related operations")
	newConfigureDefaultCommand(cmd)
}

func newConfigureDefaultCommand(parent *kingpin.CmdClause) {
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
		return versions.SetDefaultVersion(*version)
	})
}
