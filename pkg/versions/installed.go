package versions

import (
	"io/ioutil"
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

func PrintInstalled() error {
	versions, err := ListInstalled()
	if err != nil {
		return nil
	}

	return printSortedVersions(versions)
}
