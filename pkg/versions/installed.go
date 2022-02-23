package versions

import (
	"io/ioutil"

	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
)

func ListInstalled(gowrapHome string) ([]string, error) {
	versionsDir, err := GetVersionsDir(gowrapHome)
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

func FindLatestInstalled(gowrapHome string) (string, error) {
	return FindLatestInstalledForPrefix(gowrapHome, "")
}

func FindLatestInstalledForPrefix(gowrapHome, prefix string) (string, error) {
	installedVersions, err := ListInstalled(gowrapHome)
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

func PrintInstalled(gowrapHome string) error {
	versions, err := ListInstalled(gowrapHome)
	if err != nil {
		return err
	}

	return printSortedVersions(versions)
}
