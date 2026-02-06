package phelp

import (
	"github.com/nickwells/param.mod/v6/param"
)

// ProcessArgs will process the values set after parsing is complete. This is
// where any StdHelp parameters (as added by the StdHelp AddParams method)
// will be processed.
func (h StdHelp) ProcessArgs(ps *param.PSet) {
	if zshCompHasAction(h) {
		ps.SetExitStatus(zshCompletionHandler(h, nil, ps))
	}

	if h.paramsShowWhereSet {
		h.sectionsChosen[whereSetHelpSectionName] = true
	}

	if h.paramsShowUnused {
		h.sectionsChosen[unusedParamsHelpSectionName] = true
	}

	if h.exitAfterParsing {
		ps.ShouldExit()
	}

	if h.sectionsChosen.count() > 0 {
		ps.HelpRequired()
		if h.exitAfterHelp {
			ps.ShouldExit()
		}
	}
}
