package project

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/xabierlaiseca/gowrap/pkg/config"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
	"github.com/xabierlaiseca/gowrap/pkg/versions"
	"golang.org/x/mod/modfile"
)

const (
	goModFile = "go.mod"
)

func DetectVersion(path string) (string, error) {
	p := filepath.Clean(path)
	info, err := os.Stat(p)
	if err != nil {
		return "", err
	}

	if !info.IsDir() {
		return "", customerrors.Errorf("provided path is not directory: %s", path)
	}

	projectRoot, err := findProjectRoot(p)
	if customerrors.IsNotFound(err) {
		configuration, err := config.Load()
		if err != nil {
			return "", err
		}

		if trimmedDefault := strings.TrimSpace(configuration.DefaultVersion); trimmedDefault != "" {
			return trimmedDefault, nil
		}

		return versions.FindLatestInstalled()
	} else if err != nil {
		return "", err
	}

	return findGoVersion(projectRoot)
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
		return "", customerrors.NotFound()
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
