package versions

import (
	"io/ioutil"

	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
)

func ListInstalled() ([]string, error) {
	versionsDir, err := GetVersionsDir()
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(versionsDir)
	if err != nil {
		return nil, err
	}

	var versions []string
	for _, f := range files {
		if f.IsDir() {
			versions = append(versions, f.Name())
		}
	}

	return versions, nil
}

func FindLatestInstalled() (string, error) {
	return FindLatestInstalledForPrefix("")
}

func FindLatestInstalledForPrefix(prefix string) (string, error) {
	installedVersions, err := ListInstalled()
	if err != nil {
		return "", err
	}

	var compatibleVersions []string
	for _, installedVersion := range installedVersions {
		if semver.HasPrefix(installedVersion, prefix) {
			compatibleVersions = append(compatibleVersions, installedVersion)
		}
	}

	if len(compatibleVersions) == 0 {
		return "", customerrors.NotFound()
	}

	return semver.Latest(compatibleVersions)
}

func PrintInstalled() error {
	versions, err := ListInstalled()
	if err != nil {
		return nil
	}

	return printSortedVersions(versions)
}
