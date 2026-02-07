package phelp

import (
	"github.com/nickwells/param.mod/v7/param"
)

// showIntro prints the program name and, optionally, the description
func showIntro(h StdHelp, ps *param.PSet) bool {
	if h.sectionsChosen[usageHelpSectionName] &&
		(h.showSummary || ps.ProgDesc() == "") {
		return false
	}

	switch h.helpFormat {
	case helpFmtTypeMarkdown:
		showIntroFmtMD(h, ps)
	default:
		showIntroFmtStd(h, ps)
	}

	return true
}

// showIntroFmtStd prints the intro section in the standard help format
func showIntroFmtStd(h StdHelp, ps *param.PSet) {
	h.twc.Print(ps.ProgName() + "\n")

	if h.showSummary {
		return
	}

	h.twc.Wrap(ps.ProgDesc(), 0)
}

// showIntroFmtMD prints the intro section in markdown format
func showIntroFmtMD(h StdHelp, ps *param.PSet) {
	h.twc.Print("# " + ps.ProgBaseName() + "\n\n")

	if h.showSummary {
		return
	}

	desc := makeTextMarkdownSafe(ps.ProgDesc())

	h.twc.Wrap(desc, 0)
	h.twc.Print("\n")
}
