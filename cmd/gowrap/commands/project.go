package commands

import (
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/xabierlaiseca/gowrap/pkg/project"
	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
)

func newProjectCommand(app *kingpin.Application, gowrapHome, wd string) {
	cmd := app.Command("project", "Project operations")
	newProjectPinCommand(cmd, wd)
	newProjectUnpinCommand(cmd, wd)
	newProjectVersionCommand(cmd, gowrapHome, wd)
}

func newProjectPinCommand(parent *kingpin.CmdClause, wd string) {
	cmd := parent.Command("pin", "Pin specific version for current project")
	version := cmd.Arg("version", "version to pin").
		Required().
		HintAction(availableVersionCompletion).
		String()

	cmd.
		Validate(func(*kingpin.CmdClause) error {
			if len(*version) == 0 || (semver.IsValid(*version) && semver.IsFullVersion(*version)) {
				return nil
			}
			return customerrors.Errorf("invalid version provided: %s, 'a.b.c' like version required", *version)
		}).
		Action(func(*kingpin.ParseContext) error {
			return project.PinVersion(wd, *version)
		})
}

func newProjectUnpinCommand(parent *kingpin.CmdClause, wd string) {
	parent.Command("unpin", "Unpin specific version for current project").
		Action(func(*kingpin.ParseContext) error {
			return project.UnpinVersion(wd)
		})
}

func newProjectVersionCommand(parent *kingpin.CmdClause, gowrapHome, wd string) {
	parent.Command("version", "Show the Go version used by the project").
		Action(func(*kingpin.ParseContext) error {
			detectedVersion, err := project.DetectVersion(gowrapHome, wd)
			if err != nil {
				return err
			}

			var message string
			switch {
			case detectedVersion.IsDefined() && detectedVersion.IsAvailable():
				message = fmt.Sprintf("%s (specific version to use: %s)", detectedVersion.Defined, detectedVersion.Installed)
			case detectedVersion.IsDefined():
				message = fmt.Sprintf("%s (no compatible installed version found)", detectedVersion.Defined)
			case detectedVersion.IsAvailable():
				message = detectedVersion.Installed
			default:
				message = "no versions installed found"
			}

			fmt.Println(message)
			return nil
		})
}
