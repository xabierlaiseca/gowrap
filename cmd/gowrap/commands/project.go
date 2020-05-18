package commands

import (
	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/xabierlaiseca/gowrap/pkg/project"
	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
)

func newProjectCommand(app *kingpin.Application, wd string) {
	cmd := app.Command("project", "Project operations")
	newProjectPinCommand(cmd, wd)
	newProjectUnpinCommand(cmd, wd)
}

func newProjectPinCommand(parent *kingpin.CmdClause, wd string) {
	cmd := parent.Command("pin", "Pin specific version for current project")
	version := cmd.Arg("version", "version to pin").
		Required().
		HintAction(availableVersionCompletion).
		String()

	cmd.
		Validate(func(*kingpin.CmdClause) error {
			if len(*version) == 0 || (semver.IsValid(*version) && strings.Count(*version, ".") == 2) {
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
