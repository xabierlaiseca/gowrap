// +build integration

package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	genericcli "github.com/xabierlaiseca/gowrap/cmd/generic-cmd-wrapper/cli"
	gowrapcmds "github.com/xabierlaiseca/gowrap/cmd/gowrap/commands"
)

func Test_CLIs(t *testing.T) {
	testCases := map[string]struct {
		init               func(testDir string) (string, error)
		gowrapExecutions   [][]string
		wrapperCommand     []string
		expectedSubCommand genericcli.SubCommand
	}{
		"InstallOneVersion": {
			init: func(testDir string) (string, error) {
				if err := createGoMod(testDir, "1.14"); err != nil {
					return "", err
				}

				return testDir, nil
			},
			gowrapExecutions: [][]string{
				{"install", "1.14.1"},
			},
			wrapperCommand: []string{"go", "version"},
			expectedSubCommand: genericcli.SubCommand{
				Binary: filepath.Join("versions", "1.14.1", "bin", "go"),
				Args:   []string{"go", "version"},
			},
		},
		"InstallMultipleVersionsSameMayorAndMinor": {
			init: func(testDir string) (string, error) {
				if err := createGoMod(testDir, "1.14"); err != nil {
					return "", err
				}

				return testDir, nil
			},
			gowrapExecutions: [][]string{
				{"install", "1.14.2"},
				{"install", "1.14.1"},
			},
			wrapperCommand: []string{"go", "version"},
			expectedSubCommand: genericcli.SubCommand{
				Binary: filepath.Join("versions", "1.14.2", "bin", "go"),
				Args:   []string{"go", "version"},
			},
		},
		"InstallMultipleVersionsDifferentMinor": {
			init: func(testDir string) (string, error) {
				if err := createGoMod(testDir, "1.12"); err != nil {
					return "", err
				}

				return testDir, nil
			},
			gowrapExecutions: [][]string{
				{"install", "1.12.8"},
				{"install", "1.13.10"},
			},
			wrapperCommand: []string{"go", "version"},
			expectedSubCommand: genericcli.SubCommand{
				Binary: filepath.Join("versions", "1.12.8", "bin", "go"),
				Args:   []string{"go", "version"},
			},
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
			wd, err := testCase.init(testDir)
			require.NoError(t, err)

			for _, execution := range testCase.gowrapExecutions {
				require.NoError(t, gowrapcmds.RunCli(gowrapHome, wd, execution))
			}

			actualSubCommand, err := genericcli.GenerateSubCommand(gowrapHome, wd, testCase.wrapperCommand[0], testCase.wrapperCommand[1:])
			assert.NoError(t, err)

			expectedSubCommand := testCase.expectedSubCommand
			expectedSubCommand.Binary = filepath.Join(gowrapHome, expectedSubCommand.Binary)
			assert.Equal(t, &expectedSubCommand, actualSubCommand)
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
