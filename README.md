# gowrap
`gowrap` is a Go version management tool written in Go.

`gowrap` tool consist of the following commands:
* `gowrap`: the command to manage installed versions
* `go`: a wrapper for the Go's `go` command
* `gofmt`: a wrapper for the Go's `gofmt` command

## gowrap command
This is the main Go version management command. It will allow you installing,
uninstalling and configuring some preferences (such as setting up a default Go
version). For more help, run `gowrap help`.

## Wrapper commands
As a user of `gowrap` tool, you should use wrapper commands (`go` and `gofmt`)
provided by this tool instead of directly executing specific versions of Go's
commands.

Wrapper commands will detect the most suitable Go version to use and they will
call the desired specific version of the executed command.

In order to decide which version to use, wrapper commands will follow these
rules:
1. If current directory is part of a Go project:
   1. If `.go-version` exists in project root, it will select the version
      defined in that file as candidate. Otherwise, it will select the version
      defined in `go.mod`
   1. If no matching version installed for selected version, it will suggest
      the user to install latest compatible version and use it
   1. If compatible versions are installed for selected version, it will use
      latest compatible version
1. If not in go project:
   1. If default version configured, it will use that version
   1. If no versions installed, it will suggest to install latest Go version
      and it will use it
   1. Otherwise, it will use latest installed Go version
