package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/xabierlaiseca/gowrap/cmd/gowrap/commands"
)

func main() {
	parser := argparse.NewParser("gowrap", "gowrap util")
	versionsFileCmd, versionsFileCmdHandler := commands.NewVersionsFileCommand(&parser.Command)
	listCmd, listCmdHandler := commands.NewListCommand(&parser.Command)
	parser.Parse(os.Args)

	if versionsFileCmd.Happened() {
		exitOnError(versionsFileCmdHandler())
	} else if listCmd.Happened() {
		exitOnError(listCmdHandler())
	}
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
