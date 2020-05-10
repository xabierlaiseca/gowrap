package commands

import (
	"github.com/alecthomas/kingpin"
	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
	"github.com/xabierlaiseca/gowrap/pkg/versions"
)

func NewInstallCommand(app *kingpin.Application) {
	cmd := app.Command("install", "install go version")
	version := cmd.Arg("version", "Version to install").
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
			return versions.InstallLatestForPrefix(*version)
		})
}
