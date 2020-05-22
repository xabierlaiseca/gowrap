// +build integration

package integration

import (
	"fmt"
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
		"InstallNotLatestVersionAndSetupAutoUpgrades_InProject": {
			init: initGoProject("1.14"),
			gowrapExecutions: [][]string{
				{"install", "1.14.2"},
				{"configure", "upgrades", "auto"},
			},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertVersionGreaterThan("1.14.2", goVersionCommand),
		},
		"InstallNotLatestVersionAndSetupDisabledUpgrades_InProject": {
			init: initGoProject("1.14"),
			gowrapExecutions: [][]string{
				{"install", "1.14.2"},
				{"configure", "upgrades", "disabled"},
			},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertExactVersion("1.14.2", goVersionCommand),
		},
		"InstallNotLatestVersionAndSetupAutoUpgrades_OutsideProject": {
			gowrapExecutions: [][]string{
				{"install", "1.13.10"},
				{"configure", "default", "1.13"},
				{"configure", "upgrades", "auto"},
			},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertVersionBetween("1.13.10", "1.14.0", goVersionCommand),
		},
		"InstallNotLatestVersionAndSetupDisabledUpgrades_OutsideProject": {
			gowrapExecutions: [][]string{
				{"install", "1.14.2"},
				{"configure", "default", "1.14"},
				{"configure", "upgrades", "disabled"},
			},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertExactVersion("1.14.2", goVersionCommand),
		},
		"InstallNotLatestVersionAndSetupAutoUpgradesButNoDefaultVersion_OutsideProject": {
			gowrapExecutions: [][]string{
				{"install", "1.13.2"},
				{"configure", "upgrades", "auto"},
			},
			wrapperCommand:   goVersionCommand,
			assertSubCommand: assertVersionGreaterThan("1.13.2", goVersionCommand),
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
				require.NoError(t, gowrapcmds.RunCli(gowrapHome, wd, execution))
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
	goModPath := filepath.Join(dir, "go.mod")
	goModFile, err := os.OpenFile(goModPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer goModFile.Close()

	_, err = goModFile.WriteString(fmt.Sprintf("go %s", version))
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