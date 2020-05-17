package commands

import (
	"os"
	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/xabierlaiseca/gowrap/pkg/project"
	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
)

func NewProjectCommand(app *kingpin.Application) {
	cmd := app.Command("project", "Project operations")
	newProjectPinCommand(cmd)
}

func newProjectPinCommand(parent *kingpin.CmdClause) {
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
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			return project.PinVersion(wd, *version)
		})
}
