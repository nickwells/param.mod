package phelp

import (
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// showIntro prints the program name and, optionally, the description
func showIntro(h StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	if h.sectionsChosen[usageHelpSectionName] && h.hideDescriptions {
		return false
	}

	twc.Print(ps.ProgName() + "\n")

	if h.hideDescriptions {
		return true
	}
	twc.Wrap(ps.ProgDesc()+"\n", 0)
	return true
}
