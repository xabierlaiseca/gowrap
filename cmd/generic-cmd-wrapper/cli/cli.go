package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xabierlaiseca/gowrap/pkg/config"

	"github.com/xabierlaiseca/gowrap/pkg/project"
	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
	"github.com/xabierlaiseca/gowrap/pkg/versions"
)

func GenerateSubCommand(gowrapHome, wd, wrappedCmd string, args []string) (*SubCommand, error) {
	version, err := findVersionToUse(gowrapHome, wd)
	if customerrors.IsNotFound(err) {
		return nil, customerrors.Errorf("No suitable version found")
	} else if err != nil {
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
		return upgradeVersionIfConfigured(gowrapHome, detectedVersion)
	}

	candidate, err := versions.FindLatestAvailable(detectedVersion.Defined)
	if err != nil {
		return "", err
	}

	if accepted, err := installIfAccepted(gowrapHome, candidate, "No suitable version installed found"); err != nil {
		return "", err
	} else if !accepted {
		return candidate, nil
	}

	return "", customerrors.Errorf("no versions available for go %s installed", detectedVersion.Defined)
}

func upgradeVersionIfConfigured(gowrapHome string, version *project.Version) (string, error) {
	if semver.IsFullVersion(version.Defined) {
		return version.Installed, nil
	}

	c, err := config.Load(gowrapHome)
	if err != nil {
		return "", err
	}

	if c.Upgrades == config.UpgradesDisabled {
		return version.Installed, nil
	}

	candidate, err := versions.FindLatestAvailable(version.Defined)
	if err != nil {
		return "", err
	}

	if !semver.IsLessThan(version.Installed, candidate) {
		return version.Installed, nil
	}

	if c.Upgrades == config.UpgradesAuto {
		_, err := versions.InstallIfNotInstalled(gowrapHome, candidate)
		return candidate, err
	}

	accepted, err := installIfAccepted(gowrapHome, candidate, fmt.Sprintf("Upgrade found for version %s", version.Defined))
	if err != nil {
		return "", err
	} else if accepted {
		return candidate, nil
	}

	return version.Installed, nil
}

func installIfAccepted(gowrapHome, candidate, messagePrefix string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	accepted, err := askForInstallingVersion(reader, candidate, messagePrefix)
	if err != nil {
		return false, err
	}

	if accepted {
		_, err := versions.InstallIfNotInstalled(gowrapHome, candidate)
		return true, err
	}

	return false, nil
}

func askForInstallingVersion(reader *bufio.Reader, candidate, messagePrefix string) (bool, error) {
	fmt.Printf("%s, would you like to install %s? (Y/n): ", messagePrefix, candidate)
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
		return askForInstallingVersion(reader, candidate, messagePrefix)
	}
}
