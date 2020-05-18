package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/xabierlaiseca/gowrap/cmd/common"
	"github.com/xabierlaiseca/gowrap/cmd/generic-cmd-wrapper/cli"
)

var wrappedCmd = "<not-set>"

func main() {
	gowrapHome, err := common.GetGowrapHome()
	exitOnError(err)

	wd, err := os.Getwd()
	exitOnError(err)

	subCommand, err := cli.GenerateSubCommand(gowrapHome, wd, wrappedCmd, os.Args[1:])
	exitOnError(err)

	binary := subCommand.Binary
	err = syscall.Exec(binary, subCommand.Args, os.Environ())
	exitOnError(err)
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
