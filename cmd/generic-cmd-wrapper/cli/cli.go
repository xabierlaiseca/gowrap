package cli

import (
	"path/filepath"

	"github.com/xabierlaiseca/gowrap/pkg/config"
	"github.com/xabierlaiseca/gowrap/pkg/project"
	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
	"github.com/xabierlaiseca/gowrap/pkg/versions"
)

func GenerateSubCommand(gowrapHome, wd, wrappedCmd string, args []string) (*SubCommand, error) {
	version, err := findVersionToUse(gowrapHome, wd)
	if customerrors.IsNotFound(err) {
		return nil, customerrors.Errorf("No suitable version found")
	} else if err != nil {
		return nil, err
	}

	versionsDir, err := versions.GetVersionsDir(gowrapHome)
	if err != nil {
		return nil, err
	}

	binary := filepath.Join(versionsDir, version, "bin", wrappedCmd)
	scArgs := []string{wrappedCmd}
	scArgs = append(scArgs, args...)
	return &SubCommand{
		Binary: binary,
		Args:   scArgs,
	}, nil
}

type SubCommand struct {
	Binary string
	Args   []string
}

func findVersionToUse(gowrapHome, wd string) (string, error) {
	detectedVersion, err := project.DetectVersion(gowrapHome, wd)
	if err != nil && !customerrors.IsNotFound(err) {
		return "", err
	}

	installedVersion, err := autoInstallVersionIfConfigured(gowrapHome, detectedVersion)
	if err != nil {
		return "", err
	}

	if len(installedVersion) > 0 {
		return installedVersion, nil
	} else if detectedVersion.IsAvailable() {
		return detectedVersion.Installed, nil
	}
	return "", customerrors.Errorf("no versions available for go %s installed", detectedVersion.Defined)
}

func autoInstallVersionIfConfigured(gowrapHome string, version *project.Version) (string, error) {
	c, err := config.Load(gowrapHome)
	if err != nil {
		return "", err
	}

	if c.AutoInstall == config.AutoInstallDisabled || (c.AutoInstall == config.AutoInstallMissing && version.IsAvailable()) {
		return "", nil
	}

	candidate, err := versions.FindLatestAvailable(version.Defined)
	if err != nil {
		return "", err
	}

	if !semver.IsLessThan(version.Installed, candidate) {
		return "", nil
	}

	_, err = versions.InstallIfNotInstalled(gowrapHome, candidate)
	return candidate, err
}
