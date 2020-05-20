package project

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/xabierlaiseca/gowrap/pkg/semver"

	"github.com/xabierlaiseca/gowrap/pkg/config"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
	"github.com/xabierlaiseca/gowrap/pkg/versions"
)

// Version contains details about the Go version to be used.
type Version struct {
	Defined   string
	Installed string
}

func (pv *Version) IsAvailable() bool {
	return len(pv.Installed) > 0
}

func (pv *Version) InProject() bool {
	return len(pv.Defined) > 0
}

func DetectVersion(gowrapHome, path string) (*Version, error) {
	p := filepath.Clean(path)
	info, err := os.Stat(p)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, customerrors.Errorf("provided path is not directory: %s", path)
	}

	projectRoot, err := findProjectRoot(p)
	if customerrors.IsNotFound(err) {
		return detectVersionOutsideProject(gowrapHome)
	} else if err != nil {
		return nil, err
	}

	definedVersion, err := findGoVersion(projectRoot)
	if err != nil {
		return nil, err
	}

	installedVersion, err := versions.FindLatestInstalledForPrefix(gowrapHome, definedVersion)
	if err != nil {
		return nil, err
	}

	return &Version{
		Defined:   definedVersion,
		Installed: installedVersion,
	}, nil
}

func detectVersionOutsideProject(gowrapHome string) (*Version, error) {
	configuration, err := config.Load(gowrapHome)
	if err != nil {
		return nil, err
	}

	installedVersionToUse := strings.TrimSpace(configuration.DefaultVersion)

	if !semver.IsValid(installedVersionToUse) {
		installedVersionToUse, err = versions.FindLatestInstalled(gowrapHome)
		if err != nil {
			return nil, err
		}
	}

	return &Version{
		Defined:   "",
		Installed: installedVersionToUse,
	}, nil
}
