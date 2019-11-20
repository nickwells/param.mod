package phelp

import (
	"github.com/nickwells/param.mod/v3/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// showExamples prints the See Also section of the help message
func (h StdHelp) showExamples(twc *twrap.TWConf, ps *param.PSet) bool {
	ex := ps.Examples()
	if len(ex) == 0 {
		return false
	}

	twc.Println("Examples") //nolint: errcheck

	for _, e := range ex {
		twc.Wrap("\n"+e.Ex+"\n", paramIndent)
		twc.Wrap(e.Desc, descriptionIndent)
	}
	return true
}
