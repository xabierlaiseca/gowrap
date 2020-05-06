package project

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/versions"
	"golang.org/x/mod/modfile"
)

const (
	goModFile = "go.mod"
)

var (
	errProjectRootNotFound = errors.New("project root not found")
)

func Detect(path string) (string, error) {
	p := filepath.Clean(path)
	info, err := os.Stat(p)
	if err != nil {
		return "", err
	}

	if !info.IsDir() {
		return "", fmt.Errorf("provided path is not directory: %s", path)
	}

	installedVersions, err := versions.ListInstalled()
	if err != nil {
		return "", err
	} else if len(installedVersions) == 0 {
		return "", errors.New("no go versions installed")
	}

	projectRoot, err := findProjectRoot(p)
	if err == errProjectRootNotFound {
		return semver.Latest(installedVersions)
	} else if err != nil {
		return "", err
	}

	minorGoVersion, err := findGoVersion(projectRoot)
	if err != nil {
		return "", err
	}

	var compatibleVersions []string
	for _, installedVersion := range installedVersions {
		if strings.HasPrefix(installedVersion, minorGoVersion) {
			compatibleVersions = append(compatibleVersions, installedVersion)
		}
	}

	if len(compatibleVersions) == 0 {
		return "", fmt.Errorf("no suitable installed version for go %s", minorGoVersion)
	}

	return semver.Latest(compatibleVersions)
}

func findProjectRoot(directory string) (string, error) {
	candidateGoModPath := filepath.Join(directory, goModFile)
	info, err := os.Stat(candidateGoModPath)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	if err == nil && !info.Mode().IsDir() {
		return directory, nil
	}

	parent := filepath.Dir(directory)
	if parent == directory {
		return "", errProjectRootNotFound
	}

	return findProjectRoot(parent)
}

func findGoVersion(projectRoot string) (string, error) {
	goModPath := filepath.Join(projectRoot, "go.mod")
	content, err := ioutil.ReadFile(goModPath)
	if err != nil {
		return "", err
	}

	goMod, err := modfile.ParseLax(goModPath, content, nil)
	if err != nil {
		return "", err
	}

	return goMod.Go.Version, nil
}
