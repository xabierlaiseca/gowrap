package versions

import (
	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
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

func FindLatestAvailable(prefix string) (string, error) {
	availableVersions, err := versionsfile.Load()
	if err != nil {
		return "", err
	}

	var compatibleVersions []string
	for availableVersion := range availableVersions {
		if semver.HasPrefix(availableVersion, prefix) {
			compatibleVersions = append(compatibleVersions, availableVersion)
		}
	}

	if len(compatibleVersions) == 0 {
		return "", customerrors.NotFound()
	}

	return semver.Latest(compatibleVersions)
}
