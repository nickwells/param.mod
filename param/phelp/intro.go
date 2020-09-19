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

	switch h.helpFormat {
	case helpFmtTypeMD:
		return showIntroFmtMD(h, twc, ps)
	default:
		return showIntroFmtStd(h, twc, ps)
	}
}

// showIntroFmtStd prints the intro section in the standard help format
func showIntroFmtStd(h StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	twc.Print(ps.ProgName() + "\n")

	if h.hideDescriptions {
		return true
	}
	twc.Wrap(ps.ProgDesc(), 0)
	return true
}

// showIntroFmtMD prints the intro section in markdown format
func showIntroFmtMD(h StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	twc.Print("# " + ps.ProgBaseName() + "\n\n")

	if h.hideDescriptions {
		return true
	}
	desc := makeTextMarkdownSafe(ps.ProgDesc())
	twc.Wrap(desc, 0)
	twc.Print("\n")
	return true
}
