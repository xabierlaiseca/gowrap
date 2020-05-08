package versions

import (
	"github.com/xabierlaiseca/gowrap/pkg/versionsfile"
)

func PrintAvailable() error {
	versionGoArchives, err := versionsfile.Load()
	if err != nil {
		return err
	}

	versions := make([]string, 0, len(versionGoArchives))
	for version := range versionGoArchives {
		versions = append(versions, version)
	}

	return printSortedVersions(versions)
}
