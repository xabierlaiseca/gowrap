package commands

import (
	"github.com/akamensky/argparse"
	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
	"github.com/xabierlaiseca/gowrap/pkg/versions"
)

func NewInstallCommand(parent *argparse.Command) (*argparse.Command, func() error) {
	cmd := parent.NewCommand("install", "install go version")
	version := cmd.String("v", "version", &argparse.Options{
		Required: true,
		Validate: validateVersion,
	})

	return cmd, func() error {
		return versions.Install(*version)
	}
}

func validateVersion(version []string) error {
	if semver.IsValid(version[0]) {
		return nil
	}
	return customerrors.Errorf("Invalid version provided: %s", version[0])
}
