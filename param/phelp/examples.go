package phelp

import (
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// showExamples prints the Examples section of the help message
func showExamples(h StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	ex := ps.Examples()
	if len(ex) == 0 {
		return false
	}

	twc.Print("Examples\n")

	for _, e := range ex {
		twc.Wrap("\n"+e.Ex+"\n", paramIndent)
		if h.hideDescriptions {
			continue
		}
		twc.Wrap(e.Desc, descriptionIndent)
	}
	return true
}
