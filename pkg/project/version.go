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
	if customerrors.IsNotFound(err) {
		logrus.Warning("Cannot pin version, currently not in a Go project")
		return nil
	} else if err != nil {
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

func UnpinVersion(path string) error {
	projectRoot, err := findProjectRoot(path)
	if customerrors.IsNotFound(err) {
		logrus.Warning("Cannot unpin version, currently not in a Go project")
		return nil
	} else if err != nil {
		return err
	}

	goVersionPath := filepath.Join(projectRoot, goVersionFile)
	_, err = os.Stat(goVersionPath)
	if os.IsNotExist(err) {
		logrus.Warning("Cannot unpin version, no version pinned for current project")
		return nil
	} else if err != nil {
		return err
	}

	return os.RemoveAll(goVersionPath)
}
