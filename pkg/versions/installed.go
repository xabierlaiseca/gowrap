package versions

import (
	"io/ioutil"
)

func ListInstalled() error {
	versionsDir, err := getVersionsDir()
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(versionsDir)
	if err != nil {
		return nil
	}

	var versions []string
	for _, f := range files {
		if f.IsDir() {
			versions = append(versions, f.Name())
		}
	}

	return printSortedVersions(versions)
}
