package commands

import (
	"errors"
	"fmt"
	"sort"

	"github.com/akamensky/argparse"
	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/versionsfile"
)

func NewListCommand(parent *argparse.Command) (*argparse.Command, func() error) {
	cmd := parent.NewCommand("list", "List operations")
	availableCmd, availableCmdHandler := newListAvailableCommand(cmd)

	return cmd, func() error {
		if availableCmd.Happened() {
			return availableCmdHandler()
		}

		return errors.New("unexpected error: subcommand for 'list' not found")
	}
}

func newListAvailableCommand(parent *argparse.Command) (*argparse.Command, func() error) {
	cmd := parent.NewCommand("available", "Lists the available go versions to install")

	return cmd, func() error {
		versionGoArchives, err := versionsfile.Load()
		if err != nil {
			return err
		}

		var versions []string
		for version := range versionGoArchives {
			versions = append(versions, version)
		}

		comparator, err := semver.SliceStableComparatorFor(versions)
		if err != nil {
			return err
		}

		sort.SliceStable(versions, comparator)
		for _, version := range versions {
			fmt.Println(version)
		}

		return nil
	}
}
