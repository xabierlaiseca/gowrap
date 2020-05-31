package commands

import (
	"strings"

	"github.com/alecthomas/kingpin"
)

func asEnum(ac *kingpin.ArgClause, options ...string) *string {
	return ac.HintOptions(options...).
		PlaceHolder(strings.Join(options, "|")).
		Enum(options...)
}
