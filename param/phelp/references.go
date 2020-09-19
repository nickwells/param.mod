package phelp

import (
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// showReferences produces the See Also section of the help message
func showReferences(h StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	refs := ps.References()
	if len(refs) == 0 {
		return false
	}

	switch h.helpFormat {
	case helpFmtTypeMD:
		return showReferencesFmtMD(h, twc, ps)
	default:
		return showReferencesFmtStd(h, twc, ps)
	}
}

// showReferencesFmtStd prints references in the standard help format
func showReferencesFmtStd(h StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	refs := ps.References()
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

// showReferencesFmtMD prints references in markdown format
func showReferencesFmtMD(h StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	refs := ps.References()
	twc.Print("# See Also\n\n")

	for _, r := range refs {
		twc.Print("```\n")
		twc.Wrap(r.Name, 0)
		twc.Print("```\n")
		if h.hideDescriptions {
			continue
		}
		desc := makeTextMarkdownSafe(r.Desc)
		twc.Wrap(desc, 0)
	}
	return true
}
