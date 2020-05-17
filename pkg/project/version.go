package project

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
)

func PinVersion(path string, version string) error {
	projectRoot, err := findProjectRoot(path)
	if err != nil {
		return err
	}

	if goModVersion, err := findVersionInGoModFile(projectRoot); err != nil && !customerrors.IsNotFound(err) {
		return err
	} else if err == nil && !semver.HasPrefix(version, goModVersion) {
		logrus.Warningf("Pinned version (%s) is not compatible with version in go.mod (%s)", version, goModVersion)
	}

	goVersionPath := filepath.Join(projectRoot, goVersionFile)
	file, err := os.OpenFile(goVersionPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(version))
	return err
}
