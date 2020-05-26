// +build integration

package integration

import (
	"fmt"
	"github.com/xabierlaiseca/gowrap/pkg/config"
	"github.com/xabierlaiseca/gowrap/pkg/semver"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	genericcli "github.com/xabierlaiseca/gowrap/cmd/generic-cmd-wrapper/cli"
	gowrapcmds "github.com/xabierlaiseca/gowrap/cmd/gowrap/commands"
)

var goVersionCommand = []string{"go", "version"}

func Test_CLIs(t *testing.T) {
	downloadsDir, err := ioutil.TempDir(os.TempDir(), "gowrap-integration-downloads-")
	require.NoError(t, err)
	defer os.RemoveAll(downloadsDir)

	require.NoError(t, os.Setenv("GOWRAP_DOWNLOADS_DIR", downloadsDir))

	testCases := map[string]struct {
		init             func(testDir string) (string, error)
		gowrapExecutions [][]string
		wrapperCommand   []string
		assertSubCommand func(*testing.T, *genericcli.SubCommand)
	}{
		"InstallOneVersion": {
			init: initGoProject("1.14"),
			gowrapExecutions: [][]string{
				{"install", "1.14.1"},
			},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertExactVersion("1.14.1", goVersionCommand),
		},
		"InstallMultipleVersionsSameMayorAndMinor": {
			init: initGoProject("1.14"),
			gowrapExecutions: [][]string{
				{"install", "1.14.2"},
				{"install", "1.14.1"},
			},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertExactVersion("1.14.2", goVersionCommand),
		},
		"InstallMultipleVersionsDifferentMinor": {
			init: initGoProject("1.12"),
			gowrapExecutions: [][]string{
				{"install", "1.12.8"},
				{"install", "1.13.10"},
			},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertExactVersion("1.12.8", goVersionCommand),
		},
		"InstallVersionAndConfigureDefaultLessThanInstalled_OutsideProject": {
			gowrapExecutions: [][]string{
				{"install", "1.13.10"},
				{"configure", "default", "1.12.8"},
			},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertExactVersion("1.12.8", goVersionCommand),
		},
		"ConfigureDefaultAndOverride_OutsideProject": {
			gowrapExecutions: [][]string{
				{"configure", "default", "1.13.10"},
				{"configure", "default", "1.12.8"},
			},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertExactVersion("1.12.8", goVersionCommand),
		},
		"InstallNotLatestVersionAndSetupAutoInstallEnabled_InProject": {
			init: initGoProject("1.14"),
			gowrapExecutions: [][]string{
				{"install", "1.14.2"},
				{"configure", "autoinstall", config.AutoInstallEnabled},
			},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertVersionGreaterThan("1.14.2", goVersionCommand),
		},
		"InstallNotLatestVersionAndSetupAutoInstallMissing_InProject": {
			init: initGoProject("1.14"),
			gowrapExecutions: [][]string{
				{"install", "1.14.2"},
				{"configure", "autoinstall", config.AutoInstallMissing},
			},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertExactVersion("1.14.2", goVersionCommand),
		},
		"InstallNotLatestVersionAndSetupAutoInstallDisabled_InProject": {
			init: initGoProject("1.14"),
			gowrapExecutions: [][]string{
				{"install", "1.14.2"},
				{"configure", "autoinstall", config.AutoInstallDisabled},
			},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertExactVersion("1.14.2", goVersionCommand),
		},
		"NoVersionInstalledAndSetupAutoInstallMissing_InProject": {
			init: initGoProject("1.14"),
			gowrapExecutions: [][]string{
				{"configure", "autoinstall", config.AutoInstallMissing},
			},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertVersionGreaterThan("1.14", goVersionCommand),
		},
		"InstallNotLatestVersionAndSetupAutoInstallEnabled_OutsideProject": {
			gowrapExecutions: [][]string{
				{"install", "1.14.2"},
				{"configure", "autoinstall", config.AutoInstallEnabled},
			},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertVersionGreaterThan("1.14.2", goVersionCommand),
		},
		"InstallNotLatestVersionAndSetupAutoInstallMissing_OutsideProject": {
			gowrapExecutions: [][]string{
				{"install", "1.14.2"},
				{"configure", "autoinstall", config.AutoInstallMissing},
			},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertExactVersion("1.14.2", goVersionCommand),
		},
		"InstallNotLatestVersionAndSetupAutoInstallDisabled_OutsideProject": {
			gowrapExecutions: [][]string{
				{"install", "1.14.2"},
				{"configure", "autoinstall", config.AutoInstallDisabled},
			},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertExactVersion("1.14.2", goVersionCommand),
		},
		"NoVersionInstalledAndSetupAutoInstallMissing_OutsideProject": {
			gowrapExecutions: [][]string{
				{"configure", "autoinstall", config.AutoInstallMissing},
			},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertVersionGreaterThan("1.14", goVersionCommand),
		},
		"NoVersionInstalledWhenGoVersionFileAvailable": {
			init: initGoProjectWithGoVersionFile("1.14", "1.14.1"),
			gowrapExecutions: [][]string{},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertExactVersion("1.14.1", goVersionCommand),
		},
		"InstallVersionLessThanVersionInGoVersionFile": {
			init: initGoProjectWithGoVersionFile("1.14", "1.14.2"),
			gowrapExecutions: [][]string{
				{"install", "1.14.1"},
			},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertExactVersion("1.14.2", goVersionCommand),
		},
		"InstallVersionGreaterThanVersionInGoVersionFile": {
			init: initGoProjectWithGoVersionFile("1.14", "1.14.1"),
			gowrapExecutions: [][]string{
				{"install", "1.14.2"},
			},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertExactVersion("1.14.1", goVersionCommand),
		},
		"InstallNotLatestVersionAndUpgradesNotConfigured": {
			init: initGoProject("1.14"),
			gowrapExecutions: [][]string{
				{"install", "1.14.2"},
			},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertExactVersion("1.14.2", goVersionCommand),
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			tmpDir, err := ioutil.TempDir("", fmt.Sprintf("test-%s-", testName))
			require.NoError(t, err)
			defer os.RemoveAll(tmpDir)

			gowrapHome := filepath.Join(tmpDir, ".gowrap")
			require.NoError(t, os.MkdirAll(gowrapHome, 0700))

			testDir := filepath.Join(tmpDir, "test-dir")
			require.NoError(t, os.MkdirAll(testDir, 0700))

			wd := testDir
			if testCase.init != nil {
				wd, err = testCase.init(testDir)
				require.NoError(t, err)
			}

			for _, execution := range testCase.gowrapExecutions {
				require.NoError(t, gowrapcmds.RunCli("0.0.1", gowrapHome, wd, execution))
			}

			actualSubCommand, err := genericcli.GenerateSubCommand(gowrapHome, wd, testCase.wrapperCommand[0], testCase.wrapperCommand[1:])
			assert.NoError(t, err)

			relBinary, err := filepath.Rel(gowrapHome, actualSubCommand.Binary)
			assert.NoError(t, err)
			actualSubCommand.Binary = relBinary
			testCase.assertSubCommand(t, actualSubCommand)
		})
	}
}

func createGoMod(dir, version string) error {
	return createFile(dir, "go.mod", fmt.Sprintf("go %s", version))
}

func createGoVersionFile(dir, version string) error {
	return createFile(dir, ".go-version", version)
}

func createFile(dir, name, content string) error {
	p := filepath.Join(dir, name)
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

func initGoProject(version string) func (testDir string) (string, error) {
	return func(testDir string) (s string, e error) {
		if err := createGoMod(testDir, version); err != nil {
			return "", err
		}

		return testDir, nil
	}
}

func initGoProjectWithGoVersionFile(modVersion, goVersionFileVersion string) func (testDir string) (string, error) {
	return func(testDir string) (s string, e error) {
		if err := createGoMod(testDir, modVersion); err != nil {
			return "", err
		}

		if err := createGoVersionFile(testDir, goVersionFileVersion); err != nil {
			return "", err
		}

		return testDir, nil
	}
}

func assertExactVersion(version string, args []string) func(*testing.T, *genericcli.SubCommand) {
	expected := &genericcli.SubCommand{
		Binary: filepath.Join("versions", version, "bin", "go"),
		Args:   args,
	}

	return func(t *testing.T, actual *genericcli.SubCommand) {
		assert.Equal(t, expected, actual)
	}
}

func assertVersionGreaterThan(version string, args []string) func(*testing.T, *genericcli.SubCommand) {
	return func(t *testing.T, actual *genericcli.SubCommand) {
		binaryVersion := extractBinaryVersion(t, actual.Binary)
		assert.True(t, semver.IsLessThan(version, binaryVersion), "%s is not less than %s", version, binaryVersion)

		assertExactVersion(binaryVersion, args)(t, actual)
	}
}

func assertVersionBetween(lowerBound, upperBound string, args []string) func(*testing.T, *genericcli.SubCommand) {
	return func(t *testing.T, actual *genericcli.SubCommand) {
		binaryVersion := extractBinaryVersion(t, actual.Binary)
		assert.True(t, semver.IsLessThan(lowerBound, binaryVersion), "lower bound %s is not less than %s", lowerBound, binaryVersion)
		assert.True(t, semver.IsLessThan(binaryVersion, upperBound), "%s is not less than upper bound %s", binaryVersion, upperBound)

		assertExactVersion(binaryVersion, args)(t, actual)
	}
}

func extractBinaryVersion(t *testing.T, binaryPath string) string {
	segments := strings.Split(binaryPath, string(os.PathSeparator))
	require.Greater(t, len(segments), 2, "binary (%s) does not have more than 2 segments", binaryPath)
	return segments[1]
}