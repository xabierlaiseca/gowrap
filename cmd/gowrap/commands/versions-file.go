package commands

import (
	"github.com/alecthomas/kingpin"
	"github.com/xabierlaiseca/gowrap/pkg/versionsfile"
)

func NewVersionsFileCommand(parent *kingpin.Application) {
	cmd := parent.Command("versions-file", "").Hidden()
	newVersionsFileGenerateCommand(cmd)
}

func newVersionsFileGenerateCommand(parent *kingpin.CmdClause) {
	cmd := parent.Command("generate", "").Hidden()
	file := cmd.Flag("file", "output file").
		Short('f').
		Default("versions.json").
		String()

	cmd.Action(func(*kingpin.ParseContext) error {
		return versionsfile.Generate(*file)
	})
}
