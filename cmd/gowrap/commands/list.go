package commands

import (
	"errors"

	"github.com/akamensky/argparse"
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

		return errors.New("unexpected error: subcommand for 'list' not found")
	}
}

func newListAvailableCommand(parent *argparse.Command) (*argparse.Command, func() error) {
	cmd := parent.NewCommand("available", "Lists the available go versions to install")

	return cmd, func() error {
		return versions.ListAvailable()
	}
}

func newListInstalledCommand(parent *argparse.Command) (*argparse.Command, func() error) {
	cmd := parent.NewCommand("installed", "Lists installed go versions")

	return cmd, func() error {
		return versions.ListInstalled()
	}
}
