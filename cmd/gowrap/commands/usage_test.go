package commands

import (
	"github.com/alecthomas/kingpin"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_appendFormattedRows(t *testing.T) {
	testCases := map[string]struct {
		lines            []string
		rows             [][]string
		spacesBeforeRows []int

		expected []string
	}{
		"OneRow": {
			rows:             [][]string{{"cell11", "cell12"}},
			spacesBeforeRows: []int{0, 0},
			expected:         []string{"cell11 cell12"},
		},
		"ManyRows": {
			rows: [][]string{
				{"cell11", "cell12"},
				{"cell21", "cell22"}},
			spacesBeforeRows: []int{0, 0},
			expected: []string{
				"cell11 cell12",
				"cell21 cell22"},
		},
		"DifferentLengthRows": {
			rows: [][]string{
				{"cell11", "longer cell 12", "cell13"},
				{"longer cell 21", "cell22", "cell23"},
				{"even longer cell 21", "cell22", "cell23"}},
			spacesBeforeRows: []int{0, 0, 0},
			expected: []string{
				"cell11              longer cell 12 cell13",
				"longer cell 21      cell22         cell23",
				"even longer cell 21 cell22         cell23"},
		},
		"SetSpacesBeforeRows": {
			rows:             [][]string{{"cell11", "cell12"}},
			spacesBeforeRows: []int{1, 5},
			expected:         []string{" cell11      cell12"},
		},
		"WithInputLines": {
			lines:            []string{"lines 1", "lines 2"},
			rows:             [][]string{{"cell11", "cell12"}},
			spacesBeforeRows: []int{0, 0},
			expected: []string{
				"lines 1",
				"lines 2",
				"cell11 cell12"},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			actual := appendFormattedRows(testCase.lines, testCase.rows, testCase.spacesBeforeRows)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func Test_usageFlags(t *testing.T) {
	testCases := map[string]struct {
		app *kingpin.ApplicationModel
		cmd *kingpin.CmdModel

		expected []string
	}{
		"NoFlags": {
			app: &kingpin.ApplicationModel{
				FlagGroupModel: &kingpin.FlagGroupModel{
					Flags: []*kingpin.FlagModel{},
				},
			},
			expected: nil,
		},
		"OneFlag": {
			app: &kingpin.ApplicationModel{
				FlagGroupModel: &kingpin.FlagGroupModel{
					Flags: []*kingpin.FlagModel{
						{Name: "foo", Help: "foo help"},
					},
				},
			},
			expected: []string{
				"Flags:",
				"    --foo FOO   foo help",
			},
		},
		"OneFlagWithShort": {
			app: &kingpin.ApplicationModel{
				FlagGroupModel: &kingpin.FlagGroupModel{
					Flags: []*kingpin.FlagModel{
						{Name: "foo", Short: 'f', Help: "foo help"},
					},
				},
			},
			expected: []string{
				"Flags:",
				"    -f, --foo FOO   foo help",
			},
		},
		"MultipleFlags": {
			app: &kingpin.ApplicationModel{
				FlagGroupModel: &kingpin.FlagGroupModel{
					Flags: []*kingpin.FlagModel{
						{Name: "foo", Short: 'f', Help: "foo help"},
						{Name: "bar", Help: "bar help"},
						{Name: "longer", Short: 'l', Help: "longer help"},
					},
				},
			},
			expected: []string{
				"Flags:",
				"    -f, --foo FOO         foo help",
				"        --bar BAR         bar help",
				"    -l, --longer LONGER   longer help",
			},
		},
		"SelectedCommandHasPreferenceOverApp": {
			app: &kingpin.ApplicationModel{
				FlagGroupModel: &kingpin.FlagGroupModel{
					Flags: []*kingpin.FlagModel{
						{Name: "foo", Short: 'f', Help: "foo help"},
					},
				},
			},
			cmd: &kingpin.CmdModel{
				FlagGroupModel: &kingpin.FlagGroupModel{
					Flags: []*kingpin.FlagModel{
						{Name: "bar", Help: "bar help"},
					},
				},
			},
			expected: []string{
				"Flags:",
				"    --bar BAR   bar help",
			},
		},
		"IgnoresHiddenFlags": {
			app: &kingpin.ApplicationModel{
				FlagGroupModel: &kingpin.FlagGroupModel{
					Flags: []*kingpin.FlagModel{
						{Name: "foo", Short: 'f', Help: "foo help"},
						{Name: "hidden", Hidden: true, Help: "longer help"},
					},
				},
			},
			expected: []string{
				"Flags:",
				"    -f, --foo FOO   foo help",
			},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			lines := usageFlags(testCase.app, testCase.cmd)

			assert.Equal(t, testCase.expected, lines)
		})
	}
}

func Test_usageSubCommands(t *testing.T) {
	testCases := map[string]struct {
		app *kingpin.ApplicationModel
		cmd *kingpin.CmdModel

		expected []string
	}{
		"NoSubCommands": {
			app: &kingpin.ApplicationModel{
				CmdGroupModel: &kingpin.CmdGroupModel{
					Commands: []*kingpin.CmdModel{},
				},
			},
			expected: nil,
		},
		"OneSubCommand": {
			app: &kingpin.ApplicationModel{
				CmdGroupModel: &kingpin.CmdGroupModel{
					Commands: []*kingpin.CmdModel{
						{Name: "sc1", Help: "sc1 help"},
					},
				},
			},
			expected: []string{
				"Commands:",
				"    sc1   sc1 help",
			},
		},
		"MultipleSubCommands": {
			app: &kingpin.ApplicationModel{
				CmdGroupModel: &kingpin.CmdGroupModel{
					Commands: []*kingpin.CmdModel{
						{Name: "sc1", Help: "sc1 help"},
						{Name: "sc2", Help: "sc2 help"},
						{Name: "longer", Help: "longer help"},
					},
				},
			},
			expected: []string{
				"Commands:",
				"    sc1      sc1 help",
				"    sc2      sc2 help",
				"    longer   longer help",
			},
		},
		"SelectedCommandHasPreferenceOverApp": {
			app: &kingpin.ApplicationModel{
				CmdGroupModel: &kingpin.CmdGroupModel{
					Commands: []*kingpin.CmdModel{
						{Name: "sc1", Help: "sc1 help"},
					},
				},
			},
			cmd: &kingpin.CmdModel{
				CmdGroupModel: &kingpin.CmdGroupModel{
					Commands: []*kingpin.CmdModel{
						{Name: "sc2", Help: "sc2 help"},
					},
				},
			},
			expected: []string{
				"Commands:",
				"    sc2   sc2 help",
			},
		},
		"IgnoresHiddenSubCommands": {
			app: &kingpin.ApplicationModel{
				CmdGroupModel: &kingpin.CmdGroupModel{
					Commands: []*kingpin.CmdModel{
						{Name: "sc1", Help: "sc1 help"},
						{Name: "hidden", Hidden: true, Help: "hidden help"},
					},
				},
			},
			expected: []string{
				"Commands:",
				"    sc1   sc1 help",
			},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			lines := usageSubCommands(testCase.app, testCase.cmd)

			assert.Equal(t, testCase.expected, lines)
		})
	}
}

func Test_usageCommand(t *testing.T) {
	testCases := map[string]struct {
		app *kingpin.ApplicationModel
		cmd *kingpin.CmdModel

		expected []string
	}{
		"Empty": {
			app: &kingpin.ApplicationModel{
				Name:           "gowrap",
				ArgGroupModel:  &kingpin.ArgGroupModel{},
				CmdGroupModel:  &kingpin.CmdGroupModel{},
				FlagGroupModel: &kingpin.FlagGroupModel{},
			},
			expected: []string{"Usage: gowrap"},
		},
		"WithSubCommand": {
			app: &kingpin.ApplicationModel{
				Name:          "gowrap",
				ArgGroupModel: &kingpin.ArgGroupModel{},
				CmdGroupModel: &kingpin.CmdGroupModel{
					Commands: []*kingpin.CmdModel{
						{FullCommand: "sc1"},
					},
				},
				FlagGroupModel: &kingpin.FlagGroupModel{},
			},
			expected: []string{"Usage: gowrap <command> [<flags> ...]"},
		},
		"WithFlags": {
			app: &kingpin.ApplicationModel{
				Name:          "gowrap",
				ArgGroupModel: &kingpin.ArgGroupModel{},
				CmdGroupModel: &kingpin.CmdGroupModel{},
				FlagGroupModel: &kingpin.FlagGroupModel{
					Flags: []*kingpin.FlagModel{
						{Name: "flag"},
					},
				},
			},
			expected: []string{"Usage: gowrap [<flags> ...]"},
		},
		"WithArgs": {
			app: &kingpin.ApplicationModel{
				Name: "gowrap",
				ArgGroupModel: &kingpin.ArgGroupModel{
					Args: []*kingpin.ArgModel{
						{Name: "name"},
					},
				},
				CmdGroupModel:  &kingpin.CmdGroupModel{},
				FlagGroupModel: &kingpin.FlagGroupModel{},
			},
			expected: []string{"Usage: gowrap <name>"},
		},
		"WithArgsWithPlaceholder": {
			app: &kingpin.ApplicationModel{
				Name: "gowrap",
				ArgGroupModel: &kingpin.ArgGroupModel{
					Args: []*kingpin.ArgModel{
						{Name: "name", PlaceHolder: "placeholder"},
					},
				},
				CmdGroupModel:  &kingpin.CmdGroupModel{},
				FlagGroupModel: &kingpin.FlagGroupModel{},
			},
			expected: []string{"Usage: gowrap <placeholder>"},
		},
		"SelectedCommandHasPreferenceOverApp": {
			app: &kingpin.ApplicationModel{
				Name:          "gowrap",
				ArgGroupModel: &kingpin.ArgGroupModel{},
				CmdGroupModel: &kingpin.CmdGroupModel{
					Commands: []*kingpin.CmdModel{
						{FullCommand: "sc1"},
					},
				},
				FlagGroupModel: &kingpin.FlagGroupModel{},
			},
			cmd: &kingpin.CmdModel{
				FullCommand: "sc1",
				ArgGroupModel: &kingpin.ArgGroupModel{
					Args: []*kingpin.ArgModel{
						{Name: "name"},
					},
				},
				CmdGroupModel: &kingpin.CmdGroupModel{},
				FlagGroupModel: &kingpin.FlagGroupModel{
					Flags: []*kingpin.FlagModel{
						{Name: "flag"},
					},
				},
			},
			expected: []string{"Usage: gowrap sc1 [<flags> ...] <name>"},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			lines := usageCommand(testCase.app, testCase.cmd)
			assert.Equal(t, testCase.expected, lines)
		})
	}
}

func Test_usageDescription(t *testing.T) {
	testCases := map[string]struct {
		app *kingpin.ApplicationModel
		cmd *kingpin.CmdModel

		expected []string
	}{
		"AppOnly": {
			app: &kingpin.ApplicationModel{
				Help: "app help",
			},
			expected: []string{"app help"},
		},
		"SelectedCommandHasPreferenceOverApp": {
			app: &kingpin.ApplicationModel{
				Help: "app help",
			},
			cmd: &kingpin.CmdModel{
				Help: "cmd help",
			},
			expected: []string{"cmd help"},
		},
		"NoHelpDefinedInBoth": {
			app:      &kingpin.ApplicationModel{},
			cmd:      &kingpin.CmdModel{},
			expected: nil,
		},
		"NoHelpDefinedInAppAndNoCommands": {
			app:      &kingpin.ApplicationModel{},
			expected: nil,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			lines := usageDescription(testCase.app, testCase.cmd)
			assert.Equal(t, testCase.expected, lines)
		})
	}
}

func Test_usage(t *testing.T) {
	testCases := map[string]struct {
		app *kingpin.ApplicationModel
		cmd *kingpin.CmdModel

		expected []string
	}{
		"EmptyApp": {
			app: &kingpin.ApplicationModel{
				Name:           "gowrap",
				ArgGroupModel:  &kingpin.ArgGroupModel{},
				CmdGroupModel:  &kingpin.CmdGroupModel{},
				FlagGroupModel: &kingpin.FlagGroupModel{},
			},
			expected: []string{"Usage: gowrap"},
		},
		"AppNoSubCommands": {
			app: &kingpin.ApplicationModel{
				Name: "gowrap",
				Help: "gowrap help",
				ArgGroupModel: &kingpin.ArgGroupModel{
					Args: []*kingpin.ArgModel{
						{Name: "name"},
					},
				},
				CmdGroupModel: &kingpin.CmdGroupModel{},
				FlagGroupModel: &kingpin.FlagGroupModel{
					Flags: []*kingpin.FlagModel{
						{Name: "flag", Help: "flag help"},
					},
				},
			},
			expected: []string{
				"gowrap help",
				"",
				"Usage: gowrap [<flags> ...] <name>",
				"",
				"Flags:",
				"    --flag FLAG   flag help",
			},
		},
		"AppWithSubCommands": {
			app: &kingpin.ApplicationModel{
				Name: "gowrap",
				Help: "gowrap help",
				ArgGroupModel: &kingpin.ArgGroupModel{
					Args: []*kingpin.ArgModel{
						{Name: "name"},
					},
				},
				CmdGroupModel: &kingpin.CmdGroupModel{
					Commands: []*kingpin.CmdModel{
						{Name: "sb1", Help: "sb1 help"},
					},
				},
				FlagGroupModel: &kingpin.FlagGroupModel{},
			},
			expected: []string{
				"gowrap help",
				"",
				"Usage: gowrap <command> [<flags> ...]",
				"",
				"Commands:",
				"    sb1   sb1 help",
			},
		},
		"EmptyCommand": {
			app: &kingpin.ApplicationModel{
				Name:           "gowrap",
				ArgGroupModel:  &kingpin.ArgGroupModel{},
				CmdGroupModel:  &kingpin.CmdGroupModel{},
				FlagGroupModel: &kingpin.FlagGroupModel{},
			},
			cmd: &kingpin.CmdModel{
				FullCommand:    "sc1",
				ArgGroupModel:  &kingpin.ArgGroupModel{},
				CmdGroupModel:  &kingpin.CmdGroupModel{},
				FlagGroupModel: &kingpin.FlagGroupModel{},
			},
			expected: []string{"Usage: gowrap sc1"},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			actual := usage(testCase.app, testCase.cmd)
			assert.Equal(t, strings.Join(testCase.expected, "\n")+"\n", actual)
		})
	}
}
