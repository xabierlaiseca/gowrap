package commands

import (
	"github.com/akamensky/argparse"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
	"github.com/xabierlaiseca/gowrap/pkg/versionsfile"
)

func NewVersionsFileCommand(parent *argparse.Command) (*argparse.Command, func() error) {
	cmd := parent.NewCommand("versions-file", argparse.DisableDescription)
	generateCmd, generateCmdHandler := newVersionsFileGenerateCommand(cmd)

	return cmd, func() error {
		if generateCmd.Happened() {
			return generateCmdHandler()
		}

		return customerrors.New("unexpected error: subcommand for 'versions-file' not found")
	}
}

func newVersionsFileGenerateCommand(parent *argparse.Command) (*argparse.Command, func() error) {
	cmd := parent.NewCommand("generate", argparse.DisableDescription)
	file := cmd.String("f", "file", &argparse.Options{
		Required: false,
		Default:  "versions.json",
	})

	return cmd, func() error {
		return versionsfile.Generate(*file)
	}
}
