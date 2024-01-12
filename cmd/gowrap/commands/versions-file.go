package commands

import (
	"github.com/alecthomas/kingpin"
	"github.com/xabierlaiseca/gowrap/pkg/versionsfile"
)

func newVersionsFileCommand(parent *kingpin.Application) {
	cmd := parent.Command("versions-file", "commands to manage versions file")
	newVersionsFileGenerateCommand(cmd)
	newVersionsFileDownloadCommand(cmd)
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

func newVersionsFileDownloadCommand(parent *kingpin.CmdClause) {
	cmd := parent.Command("download", "Downloads latest versions file")

	cmd.Action(func(*kingpin.ParseContext) error {
		_, err := versionsfile.DownloadToCache()
		return err
	})
}
