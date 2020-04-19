package phelp

import (
	"github.com/nickwells/param.mod/v4/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// printExamples prints the Examples section of the help message
func printExamples(h StdHelp, twc *twrap.TWConf, ps *param.PSet) {
	ex := ps.Examples()
	if len(ex) == 0 {
		twc.Wrap("There are no examples", textIndent)
		return
	}

	twc.Println("Examples") //nolint: errcheck

	for _, e := range ex {
		twc.Wrap("\n"+e.Ex+"\n", paramIndent)
		twc.Wrap(e.Desc, descriptionIndent)
	}
}
