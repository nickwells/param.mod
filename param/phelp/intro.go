package phelp

import (
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// showIntro prints the program name and, optionally, the description
func showIntro(h StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	if h.sectionsChosen[usageHelpSectionName] &&
		(h.hideDescriptions || ps.ProgDesc() == "") {
		return false
	}

	switch h.helpFormat {
	case helpFmtTypeMD:
		showIntroFmtMD(h, twc, ps)
	default:
		showIntroFmtStd(h, twc, ps)
	}
	return true
}

// showIntroFmtStd prints the intro section in the standard help format
func showIntroFmtStd(h StdHelp, twc *twrap.TWConf, ps *param.PSet) {
	twc.Print(ps.ProgName() + "\n")

	if h.hideDescriptions {
		return
	}
	twc.Wrap(ps.ProgDesc(), 0)
}

// showIntroFmtMD prints the intro section in markdown format
func showIntroFmtMD(h StdHelp, twc *twrap.TWConf, ps *param.PSet) {
	twc.Print("# " + ps.ProgBaseName() + "\n\n")

	if h.hideDescriptions {
		return
	}
	desc := makeTextMarkdownSafe(ps.ProgDesc())
	twc.Wrap(desc, 0)
	twc.Print("\n")
}
