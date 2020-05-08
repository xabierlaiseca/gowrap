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
	installCmd, installCmdHandler := commands.NewInstallCommand(&parser.Command)
	exitOnError(parser.Parse(os.Args))

	switch {
	case versionsFileCmd.Happened():
		exitOnError(versionsFileCmdHandler())
	case listCmd.Happened():
		exitOnError(listCmdHandler())
	case installCmd.Happened():
		exitOnError(installCmdHandler())
	}
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
