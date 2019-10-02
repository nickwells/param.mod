package phelp

import (
	"fmt"
	"os"

	"github.com/nickwells/param.mod/v3/param"
)

// ProcessArgs will process the values set after parsing is complete. This is
// where any StdHelp parameters (as added by the StdHelp AddParams method)
// will be processed.
func (h StdHelp) ProcessArgs(ps *param.PSet) {
	var shouldExit = h.exitAfterParsing
	separator := ""
	const dfltSep = "\n" + equals + "\n\n"

	if h.reportWhereParamsAreSet {
		fmt.Fprint(ps.StdWriter(), separator)
		separator = dfltSep
		shouldExit = true

		showWhereParamsAreSet(ps)
	}
	if h.reportUnusedParams {
		fmt.Fprint(ps.StdWriter(), separator)
		separator = dfltSep
		shouldExit = true

		showUnusedParams(ps)
	}
	if h.reportParamSources {
		fmt.Fprint(ps.StdWriter(), separator)
		separator = dfltSep
		shouldExit = true

		showParamSources(ps)
	}

	if h.showHelp {
		fmt.Fprint(ps.StdWriter(), separator)
		separator = dfltSep // nolint
		shouldExit = true

		h.Help(ps)
	}

	if shouldExit {
		os.Exit(0)
	}
}
