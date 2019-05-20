package phelp

import (
	"os"

	"github.com/nickwells/param.mod/v3/param"
)

// ProcessArgs will process the values set after parsing is complete. This is
// where any StdHelp parameters (as added by the StdHelp AddParams method)
// will be processed.
func (h StdHelp) ProcessArgs(ps *param.PSet) {
	var shouldExit = h.exitAfterParsing

	if h.reportWhereParamsAreSet {
		showWhereParamsAreSet(ps)
		shouldExit = true
	}
	if h.reportUnusedParams {
		showUnusedParams(ps)
		shouldExit = true
	}
	if h.reportParamSources {
		showParamSources(ps)
		shouldExit = true
	}

	if h.showHelp {
		h.Help(ps)
		shouldExit = true
	}

	if shouldExit {
		os.Exit(0)
	}
}
