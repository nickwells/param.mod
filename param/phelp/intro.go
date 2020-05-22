package phelp

import (
	"github.com/nickwells/param.mod/v4/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// showIntro prints the program name and, optionally, the description
func showIntro(h StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	if h.sectionsChosen[usageSection] && h.hideDescriptions {
		return false
	}

	twc.Print(ps.ProgName() + "\n")

	if h.hideDescriptions {
		return true
	}
	twc.Wrap(ps.ProgDesc()+"\n", 0)
	return true
}
