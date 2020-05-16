package commands

import (
	"github.com/xabierlaiseca/gowrap/pkg/versions"
	"github.com/xabierlaiseca/gowrap/pkg/versionsfile"
)

func availableVersionCompletion() []string {
	return versionCompletionHelper(func(string) bool {
		return true
	})
}

func installedVersionCompletion() []string {
	installed, err := versions.ListInstalled()
	if err != nil {
		return []string{}
	}

	return installed
}

func notInstalledVersionCompletion() []string {
	installed, err := versions.ListInstalled()
	if err != nil {
		return []string{}
	}

	alreadyInstalled := make(map[string]bool)
	for _, v := range installed {
		alreadyInstalled[v] = true
	}

	return versionCompletionHelper(func(curr string) bool {
		_, found := alreadyInstalled[curr]
		return !found
	})
}

func versionCompletionHelper(filter func(string) bool) []string {
	versionGoArchives, err := versionsfile.Load()
	if err != nil {
		return []string{}
	}

	options := make([]string, 0, len(versionGoArchives))
	for version := range versionGoArchives {
		if filter(version) {
			options = append(options, version)
		}
	}
	return options
}
