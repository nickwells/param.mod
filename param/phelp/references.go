package phelp

import (
	"github.com/nickwells/param.mod/v4/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// showReferences prints the See Also section of the help message
func (h StdHelp) showReferences(twc *twrap.TWConf, ps *param.PSet) bool {
	refs := ps.References()
	if len(refs) == 0 {
		return false
	}

	twc.Println("See Also") //nolint: errcheck

	for _, r := range refs {
		twc.Wrap("\n"+r.Name+"\n", paramIndent)
		twc.Wrap(r.Desc, descriptionIndent)
	}

	return true
}
