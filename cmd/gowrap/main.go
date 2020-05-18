package main

import (
	"fmt"
	"os"

	"github.com/xabierlaiseca/gowrap/cmd/gowrap/commands"

	"github.com/xabierlaiseca/gowrap/cmd/common"
)

func main() {
	gowrapHome, err := common.GetGowrapHome()
	exitOnError(err)

	wd, err := os.Getwd()
	exitOnError(err)

	exitOnError(commands.RunCli(gowrapHome, wd, os.Args[1:]))
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
