package phelp

import (
	"github.com/nickwells/param.mod/v7/param"
)

// showExamples prints the Examples section of the help message
func showExamples(h StdHelp, ps *param.PSet) bool {
	ex := ps.Examples()
	if len(ex) == 0 {
		return false
	}

	switch h.helpFormat {
	case helpFmtTypeMarkdown:
		return showExamplesFmtMD(h, ps)
	default:
		return showExamplesFmtStd(h, ps)
	}
}

// showExamplesFmtStd prints examples in the standard help format
func showExamplesFmtStd(h StdHelp, ps *param.PSet) bool {
	ex := ps.Examples()

	_, _ = h.twc.Print("Examples\n")

	for _, e := range ex {
		h.twc.Wrap("\n"+e.Ex()+"\n", paramIndent)

		if h.showSummary {
			continue
		}

		h.twc.Wrap(e.Desc(), descriptionIndent)
	}

	return true
}

// showExamplesFmtMD prints examples in markdown format
func showExamplesFmtMD(h StdHelp, ps *param.PSet) bool {
	ex := ps.Examples()

	h.twc.Print("# Examples\n\n")

	for _, e := range ex {
		h.twc.Print("```sh\n")
		h.twc.Print(e.Ex() + "\n")
		h.twc.Print("```\n")

		if h.showSummary {
			continue
		}

		desc := makeTextMarkdownSafe(e.Desc())

		h.twc.Wrap(desc, 0)
		h.twc.Print("\n")
	}

	return true
}
