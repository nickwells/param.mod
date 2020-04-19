package phelp

import (
	"github.com/nickwells/param.mod/v4/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// printReferences prints the See Also section of the help message
func printReferences(h StdHelp, twc *twrap.TWConf, ps *param.PSet) {
	refs := ps.References()
	if len(refs) == 0 {
		twc.Wrap("There are no references", textIndent)
		return
	}

	twc.Println("See Also") //nolint: errcheck

	for _, r := range refs {
		twc.Wrap("\n"+r.Name+"\n", paramIndent)
		twc.Wrap(r.Desc, descriptionIndent)
	}
}
