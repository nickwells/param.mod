package phelp

import (
	"github.com/nickwells/param.mod/v7/param"
)

// showReferences produces the See Also section of the help message
func showReferences(h StdHelp, ps *param.PSet) bool {
	refs := ps.References()
	if len(refs) == 0 {
		return false
	}

	switch h.helpFormat {
	case helpFmtTypeMarkdown:
		return showReferencesFmtMD(h, ps)
	default:
		return showReferencesFmtStd(h, ps)
	}
}

// showReferencesFmtStd prints references in the standard help format
func showReferencesFmtStd(h StdHelp, ps *param.PSet) bool {
	refs := ps.References()

	h.twc.Print("See Also\n")

	for _, r := range refs {
		h.twc.Wrap("\n"+r.Name()+"\n", paramIndent)

		if h.showSummary {
			continue
		}

		h.twc.Wrap(r.Desc(), descriptionIndent)
	}

	return true
}

// showReferencesFmtMD prints references in markdown format
func showReferencesFmtMD(h StdHelp, ps *param.PSet) bool {
	refs := ps.References()

	h.twc.Print("# See Also\n\n")

	for _, r := range refs {
		h.twc.Print("```\n")
		h.twc.Wrap(r.Name(), 0)
		h.twc.Print("```\n")

		if h.showSummary {
			continue
		}

		desc := makeTextMarkdownSafe(r.Desc())

		h.twc.Wrap(desc, 0)
	}

	return true
}
