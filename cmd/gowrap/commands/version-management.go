package commands

import (
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
	"github.com/xabierlaiseca/gowrap/pkg/versions"
)

func newInstallCommand(app *kingpin.Application, gowrapHome string) {
	versionManagementCommand(app, gowrapHome, "install", notInstalledVersionCompletion, installVersion)
}

func newUninstallCommand(app *kingpin.Application, gowrapHome string) {
	versionManagementCommand(app, gowrapHome, "uninstall", installedVersionCompletion, versions.Uninstall)
}

func versionManagementCommand(app *kingpin.Application, gowrapHome string, name string, hintFn func() []string, actionFn func(string, string) error) {
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
			return actionFn(gowrapHome, *version)
		})
}

func installVersion(gowrapHome string, prefix string) error {
	if installed, err := versions.InstallLatestIfNotInstalled(gowrapHome, prefix); err != nil {
		return err
	} else if !installed {
		fmt.Printf("version '%s' was already installed\n", prefix)
	}

	return nil
}
