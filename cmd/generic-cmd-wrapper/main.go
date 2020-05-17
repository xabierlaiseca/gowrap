package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/xabierlaiseca/gowrap/pkg/project"
	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
	"github.com/xabierlaiseca/gowrap/pkg/versions"
	"github.com/xabierlaiseca/gowrap/pkg/versionsfile"
)

var wrappedCmd = "<not-set>"

func main() {
	wd, err := os.Getwd()
	exitOnError(err)

	version, err := findVersionToUse(wd)
	exitOnError(err)

	versionsDir, err := versions.GetVersionsDir()
	exitOnError(err)

	goBin := filepath.Join(versionsDir, version, "bin", wrappedCmd)
	wrapperArgs := os.Args
	execArgs := make([]string, len(wrapperArgs))
	copy(execArgs, wrapperArgs)
	execArgs[0] = wrappedCmd
	exitOnError(syscall.Exec(goBin, execArgs, os.Environ()))
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func findVersionToUse(wd string) (string, error) {
	projectVersion, err := project.DetectVersion(wd)
	if err != nil {
		return "", err
	}

	installedVersion, err := versions.FindLatestInstalledForPrefix(projectVersion)
	if customerrors.IsNotFound(err) {
		return installIfAccepted(projectVersion)
	} else if err != nil {
		return "", err
	}

	return installedVersion, nil
}

func installIfAccepted(version string) (string, error) {
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
		_, err := versions.InstallIfNotInstalled(candidate)
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
