package phelp

import (
	"github.com/nickwells/param.mod/v7/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// ProcessArgs will process the values set after parsing is complete. This is
// where any StdHelp parameters (as added by the StdHelp AddParams method)
// will be processed.
func (h *StdHelp) ProcessArgs(ps *param.PSet) {
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
		// when help has been requested, only report errors if there have
		// been errors with the help parameters
		if !hasHelpErrors(ps) {
			h.reportErrors = false
		}

		ps.HelpRequired()

		if h.exitAfterHelp {
			ps.ShouldExit()
		}
	}
}

// hasHelpErrors checks to see if any of the errors are help-related and
// returns true if there are any such errors.
func hasHelpErrors(ps *param.PSet) bool {
	errMap := ps.Errors()

	for k := range errMap {
		p, err := ps.GetParamByName(k)
		if err != nil {
			continue // it's not a named parameter error
		}

		if p.GroupName() == helpGroupName {
			return true
		}
	}

	return false
}
