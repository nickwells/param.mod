package phelp

import (
	"github.com/nickwells/param.mod/v4/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// showReferences produces the See Also section of the help message
func showReferences(h StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	refs := ps.References()
	if len(refs) == 0 {
		return false
	}

	twc.Print("See Also\n")

	for _, r := range refs {
		twc.Wrap("\n"+r.Name+"\n", paramIndent)
		if h.hideDescriptions {
			continue
		}
		twc.Wrap(r.Desc, descriptionIndent)
	}
	return true
}
