package commands

import (
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
	"github.com/xabierlaiseca/gowrap/pkg/versions"
)

func NewInstallCommand(app *kingpin.Application) {
	versionManagementCommand(app, "install", notInstalledVersionCompletion, installVersion)
}

func NewUninstallCommand(app *kingpin.Application) {
	versionManagementCommand(app, "uninstall", installedVersionCompletion, versions.Uninstall)
}

func versionManagementCommand(app *kingpin.Application, name string, hintFn func() []string, actionFn func(string) error) {
	cmd := app.Command(name, fmt.Sprintf("%s go version", name))
	version := cmd.Arg("version", fmt.Sprintf("version to %s", name)).
		Required().
		HintAction(hintFn).
		String()

	cmd.
		Validate(func(*kingpin.CmdClause) error {
			if len(*version) == 0 || semver.IsValid(*version) {
				return nil
			}
			return customerrors.Errorf("invalid version provided: %s", *version)
		}).
		Action(func(*kingpin.ParseContext) error {
			return actionFn(*version)
		})
}

func installVersion(prefix string) error {
	if installed, err := versions.InstallLatestIfNotInstalled(prefix); err != nil {
		return err
	} else if !installed {
		fmt.Printf("version '%s' was already installed\n", prefix)
	}

	return nil
}
