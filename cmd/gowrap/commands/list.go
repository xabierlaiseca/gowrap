package commands

import (
	"github.com/akamensky/argparse"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
	"github.com/xabierlaiseca/gowrap/pkg/versions"
)

func NewListCommand(parent *argparse.Command) (*argparse.Command, func() error) {
	cmd := parent.NewCommand("list", "List operations")
	availableCmd, availableCmdHandler := newListAvailableCommand(cmd)
	installedCmd, installedCmdHandler := newListInstalledCommand(cmd)

	return cmd, func() error {
		if availableCmd.Happened() {
			return availableCmdHandler()
		} else if installedCmd.Happened() {
			return installedCmdHandler()
		}

		return customerrors.New("unexpected error: subcommand for 'list' not found")
	}
}

func newListAvailableCommand(parent *argparse.Command) (*argparse.Command, func() error) {
	cmd := parent.NewCommand("available", "Lists the available go versions to install")

	return cmd, versions.PrintAvailable
}

func newListInstalledCommand(parent *argparse.Command) (*argparse.Command, func() error) {
	cmd := parent.NewCommand("installed", "Lists installed go versions")

	return cmd, versions.PrintInstalled
}
