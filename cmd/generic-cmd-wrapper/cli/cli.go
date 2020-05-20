package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xabierlaiseca/gowrap/pkg/project"
	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
	"github.com/xabierlaiseca/gowrap/pkg/versions"
	"github.com/xabierlaiseca/gowrap/pkg/versionsfile"
)

func GenerateSubCommand(gowrapHome, wd, wrappedCmd string, args []string) (*SubCommand, error) {
	version, err := findVersionToUse(gowrapHome, wd)
	if err != nil {
		return nil, err
	}

	versionsDir, err := versions.GetVersionsDir(gowrapHome)
	if err != nil {
		return nil, err
	}

	binary := filepath.Join(versionsDir, version, "bin", wrappedCmd)
	scArgs := []string{wrappedCmd}
	scArgs = append(scArgs, args...)
	return &SubCommand{
		Binary: binary,
		Args:   scArgs,
	}, nil
}

type SubCommand struct {
	Binary string
	Args   []string
}

func findVersionToUse(gowrapHome, wd string) (string, error) {
	detectedVersion, err := project.DetectVersion(gowrapHome, wd)
	if err != nil {
		return "", err
	} else if detectedVersion.IsAvailable() {
		return detectedVersion.Installed, nil
	}

	return installIfAccepted(gowrapHome, detectedVersion.Defined)
}

func installIfAccepted(gowrapHome, version string) (string, error) {
	matchingVersions := make([]string, 0)
	availableVersions, err := versionsfile.Load()
	if err != nil {
		return "", err
	}

	for availableVersion := range availableVersions {
		if semver.HasPrefix(availableVersion, version) {
			matchingVersions = append(matchingVersions, availableVersion)
		}
	}

	if len(matchingVersions) == 0 {
		return "", customerrors.Errorf("no versions available for go %s installed or available", version)
	}

	candidate, err := semver.Latest(matchingVersions)
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(os.Stdin)
	accepted, err := askForInstallingVersion(reader, candidate)
	if err != nil {
		return "", err
	}

	if accepted {
		_, err := versions.InstallIfNotInstalled(gowrapHome, candidate)
		return candidate, err
	}

	return "", customerrors.Errorf("no versions available for go %s installed", version)
}

func askForInstallingVersion(reader *bufio.Reader, candidate string) (bool, error) {
	fmt.Printf("No suitable version installed found, would you like to install %s? (Y/n): ", candidate)
	text, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	text = strings.ToLower(strings.TrimSpace(text))

	if len(text) == 0 {
		text = "y"
	}

	switch {
	case strings.HasPrefix(text, "y"):
		return true, nil
	case strings.HasPrefix(text, "n"):
		return false, nil
	default:
		fmt.Printf("Unexpected option provided: %s\n", text)
		return askForInstallingVersion(reader, candidate)
	}
}
