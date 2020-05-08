package versions

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/xabierlaiseca/gowrap/pkg/semver"
)

const (
	gowrapDir   = ".gowrap"
	versionsDir = "versions"
)

func GetVersionsDir() (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(userHome, gowrapDir, versionsDir)
	return dir, os.MkdirAll(dir, 0755)
}

func printSortedVersions(versions []string) error {
	comparator, err := semver.SliceStableComparatorFor(versions)
	if err != nil {
		return err
	}

	sort.SliceStable(versions, comparator)
	for _, version := range versions {
		fmt.Println(version)
	}

	return nil
}
