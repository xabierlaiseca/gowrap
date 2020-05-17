package project

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/xabierlaiseca/gowrap/pkg/config"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
	"github.com/xabierlaiseca/gowrap/pkg/versions"
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
