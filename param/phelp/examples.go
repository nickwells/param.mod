package phelp

import (
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// showExamples prints the Examples section of the help message
func showExamples(h StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	ex := ps.Examples()
	if len(ex) == 0 {
		return false
	}
	switch h.helpFormat {
	case helpFmtTypeMarkdown:
		return showExamplesFmtMD(h, twc, ps)
	default:
		return showExamplesFmtStd(h, twc, ps)
	}
}

// showExamplesFmtStd prints examples in the standard help format
func showExamplesFmtStd(h StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	ex := ps.Examples()
	twc.Print("Examples\n")

	for _, e := range ex {
		twc.Wrap("\n"+e.Ex+"\n", paramIndent)
		if h.hideDescriptions {
			continue
		}
		twc.Wrap(e.Desc, descriptionIndent)
	}
	return true
}

// showExamplesFmtMD prints examples in markdown format
func showExamplesFmtMD(h StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	ex := ps.Examples()
	twc.Print("# Examples\n\n")

	for _, e := range ex {
		twc.Print("```sh\n")
		twc.Print(e.Ex + "\n")
		twc.Print("```\n")
		if h.hideDescriptions {
			continue
		}
		desc := makeTextMarkdownSafe(e.Desc)
		twc.Wrap(desc, 0)
		twc.Print("\n")
	}
	return true
}
