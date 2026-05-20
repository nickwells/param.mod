package phelp

import (
	"github.com/nickwells/param.mod/v7/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// ProcessArgs will process the values set after parsing is complete. This is
// where any StdHelp parameters (as added by the StdHelp AddParams method)
// will be processed.
func (h StdHelp) ProcessArgs(ps *param.PSet) {
	if zshCompHasAction(h) {
		twc := twrap.NewTWConfOrPanic()

		completionErrStatus := zshCompletionHandler(h, twc, ps)
		if completionErrStatus == 0 {
			h.reportErrors = false
		}

		ps.SetExitStatus(completionErrStatus)
		ps.ShouldExit()
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
		h.reportErrors = false // don't report errors if help requested

		ps.HelpRequired()

		if h.exitAfterHelp {
			ps.ShouldExit()
		}
	}
}
