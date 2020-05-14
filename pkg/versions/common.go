package versions

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/xabierlaiseca/gowrap/pkg/common"
	"github.com/xabierlaiseca/gowrap/pkg/semver"
)

const versionsDir = "versions"

func GetVersionsDir() (string, error) {
	gowrapDir, err := common.GetGowrapDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(gowrapDir, versionsDir)
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
