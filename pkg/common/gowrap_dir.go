package common

import (
	"os"
	"path/filepath"
)

const gowrapDir = ".gowrap"

func GetGowrapDir() (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(userHome, gowrapDir)
	return dir, os.MkdirAll(dir, 0755)
}
