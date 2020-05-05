package versions

import (
	"github.com/xabierlaiseca/gowrap/pkg/versionsfile"
)

func ListAvailable() error {
	versionGoArchives, err := versionsfile.Load()
	if err != nil {
		return err
	}

	var versions []string
	for version := range versionGoArchives {
		versions = append(versions, version)
	}

	return printSortedVersions(versions)
}
