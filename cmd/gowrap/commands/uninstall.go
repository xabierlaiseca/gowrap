package commands

import (
	"github.com/alecthomas/kingpin"
	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
	"github.com/xabierlaiseca/gowrap/pkg/versions"
)

func NewUninstallCommand(app *kingpin.Application) {
	cmd := app.Command("uninstall", "uninstall go version")
	version := cmd.Arg("version", "version to uninstall").
		Required().
		String()

	cmd.
		Validate(func(*kingpin.CmdClause) error {
			if len(*version) == 0 || semver.IsValid(*version) {
				return nil
			}
			return customerrors.Errorf("Invalid version provided: %s", *version)
		}).
		Action(func(*kingpin.ParseContext) error {
			return versions.Uninstall(*version)
		})
}
