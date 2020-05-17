package project

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
	"golang.org/x/mod/modfile"
)

const (
	goModFile     = "go.mod"
	goVersionFile = ".go-version"
)

func findProjectRoot(directory string) (string, error) {
	candidateGoModPath := filepath.Join(directory, goModFile)
	if goModExists, err := fileExists(candidateGoModPath); err != nil {
		return "", err
	} else if goModExists {
		return directory, nil
	}

	candidateGoVersionPath := filepath.Join(directory, goVersionFile)
	if goVersionExists, err := fileExists(candidateGoVersionPath); err != nil {
		return "", err
	} else if goVersionExists {
		return directory, nil
	}

	parent := filepath.Dir(directory)
	if parent == directory {
		return "", customerrors.NotFound()
	}

	return findProjectRoot(parent)
}

func fileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}

	return err == nil && !info.Mode().IsDir(), nil
}

func findGoVersion(projectRoot string) (string, error) {
	version, err := findVersionInGoVersionFile(projectRoot)
	if err == nil {
		return version, nil
	} else if customerrors.IsNotFound(err) {
		return findVersionInGoModFile(projectRoot)
	}

	return "", err
}

func findVersionInGoVersionFile(projectRoot string) (string, error) {
	goVersionPath := filepath.Join(projectRoot, goVersionFile)
	content, err := readFile(goVersionPath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

func findVersionInGoModFile(projectRoot string) (string, error) {
	goModPath := filepath.Join(projectRoot, goModFile)
	content, err := readFile(goModPath)
	if err != nil {
		return "", err
	}

	goMod, err := modfile.ParseLax(goModPath, content, nil)
	if err != nil {
		return "", err
	}

	return goMod.Go.Version, nil
}

func readFile(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, customerrors.NotFound()
	} else if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(path)
}
