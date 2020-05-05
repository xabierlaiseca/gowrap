package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/xabierlaiseca/gowrap/pkg/project"
	"github.com/xabierlaiseca/gowrap/pkg/versions"
)

func main() {
	wd, err := os.Getwd()
	exitOnError(err)

	version, err := project.Detect(wd)
	exitOnError(err)

	versionsDir, err := versions.GetVersionsDir()
	exitOnError(err)

	goBin := filepath.Join(versionsDir, version, "bin", "go")

	wrapperArgs := os.Args
	execArgs := make([]string, len(wrapperArgs))
	copy(execArgs, wrapperArgs)
	execArgs[0] = "go"
	exitOnError(syscall.Exec(goBin, execArgs, os.Environ()))
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
