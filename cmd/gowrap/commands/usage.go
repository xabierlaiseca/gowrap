package commands

import (
	"fmt"
	"strings"

	"github.com/alecthomas/kingpin"
)

func usage(application, cmd interface{}) string {
	app, _ := application.(*kingpin.ApplicationModel)
	selectedCommand, _ := cmd.(*kingpin.CmdModel)
	lines := make([]string, 0)

	usageFns := []func(*kingpin.ApplicationModel, *kingpin.CmdModel) []string{
		usageDescription,
		usageCommand,
		usageSubCommands,
		usageFlags,
	}

	for _, usageFn := range usageFns {
		usageFnLines := usageFn(app, selectedCommand)
		if len(usageFnLines) > 0 {
			if len(lines) > 0 {
				lines = append(lines, "")
			}
			lines = append(lines, usageFnLines...)
		}
	}

	return strings.Join(lines, "\n") + "\n"
}

func usageDescription(app *kingpin.ApplicationModel, selectedCommand *kingpin.CmdModel) []string {
	if selectedCommand != nil && len(selectedCommand.Help) > 0 {
		return []string{selectedCommand.Help}
	} else if selectedCommand == nil && len(app.Help) > 0 {
		return []string{app.Help}
	}

	return nil
}

func usageCommand(app *kingpin.ApplicationModel, selectedCommand *kingpin.CmdModel) []string {
	var builder strings.Builder
	builder.WriteString("Usage: ")
	builder.WriteString(app.Name)

	cmds := app.CmdGroupModel
	flags := app.Flags
	args := app.Args

	if selectedCommand != nil {
		builder.WriteString(" ")
		builder.WriteString(selectedCommand.FullCommand)

		cmds = selectedCommand.CmdGroupModel
		flags = selectedCommand.Flags
		args = selectedCommand.Args
	}

	hasCommands := len(cmds.Commands) > 0

	if hasCommands {
		builder.WriteString(" <command>")
	}

	if hasCommands || len(flags) > 0 {
		builder.WriteString(" [<flags> ...]")
	}

	if !hasCommands {
		for _, arg := range args {
			builder.WriteString(" <")
			builder.WriteString(firstNonEmpty(arg.PlaceHolder, arg.Name))
			builder.WriteString(">")
		}
	}

	return []string{builder.String()}
}

func usageSubCommands(app *kingpin.ApplicationModel, selectedCommand *kingpin.CmdModel) []string {
	commands := app.Commands
	if selectedCommand != nil {
		commands = selectedCommand.Commands
	}

	if len(commands) == 0 {
		return nil
	}

	rows := make([][]string, 0)
	for _, command := range commands {
		if command.Hidden {
			continue
		}

		row := []string{command.Name, command.Help}
		rows = append(rows, row)
	}

	return appendFormattedRows([]string{"Commands:"}, rows, []int{minSpacesBeforeFirstColumn, minSpacesBeforeHelp})
}

func usageFlags(app *kingpin.ApplicationModel, selectedCommand *kingpin.CmdModel) []string {
	flags := app.Flags
	if selectedCommand != nil {
		flags = selectedCommand.Flags
	}

	if len(flags) == 0 {
		return nil
	}

	rows := make([][]string, 0)
	for _, flag := range flags {
		if flag.Hidden {
			continue
		}

		shortCell := ""
		if flag.Short != 0 {
			shortCell = fmt.Sprintf("-%c,", flag.Short)
		}

		longCell := fmt.Sprintf("--%s", flag.Name)
		if !flag.IsBoolFlag() {
			longCell += " " + flag.FormatPlaceHolder()
		}
		helpCell := flag.HelpWithEnvar()

		row := []string{shortCell, longCell, helpCell}
		rows = append(rows, row)
	}

	return appendFormattedRows([]string{"Flags:"}, rows, []int{minSpacesBeforeFirstColumn, 0, minSpacesBeforeHelp})
}

const (
	minSpacesBeforeFirstColumn = 4
	minSpacesBeforeHelp        = 2
)

func appendFormattedRows(lines []string, rows [][]string, minSpacesBeforeRows []int) []string {
	maxColumnLengths := make([]int, len(rows[0]))
	for _, row := range rows {
		for i := range row {
			maxColumnLengths[i] = max(maxColumnLengths[i], len(row[i]))
		}
	}

	indentationPerRow := make([]string, len(minSpacesBeforeRows))
	for i, spaces := range minSpacesBeforeRows {
		indentationPerRow[i] = strings.Repeat(" ", spaces)
	}

	for _, cellsInRow := range rows {
		var builder strings.Builder

		for i := range cellsInRow {
			builder.WriteString(indentationPerRow[i])
			length := maxColumnLengths[i]
			if length > 0 {
				length++
			}

			builder.WriteString(cellsInRow[i])

			if i+1 < len(cellsInRow) {
				space := length - len(cellsInRow[i])
				builder.WriteString(strings.Repeat(" ", space))
			}
		}

		lines = append(lines, builder.String())
	}

	return lines
}

func firstNonEmpty(strs ...string) string {
	for _, str := range strs {
		if len(str) > 0 {
			return str
		}
	}

	return ""
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
