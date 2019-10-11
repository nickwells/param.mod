package phelp

import (
	"fmt"
	"os"

	"github.com/nickwells/param.mod/v3/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// ProcessArgs will process the values set after parsing is complete. This is
// where any StdHelp parameters (as added by the StdHelp AddParams method)
// will be processed.
func (h StdHelp) ProcessArgs(ps *param.PSet) {
	twc, err := twrap.NewTWConf(twrap.SetWriter(ps.StdWriter()))
	if err != nil {
		fmt.Fprint(os.Stderr, "Couldn't build the text wrapper:", err)
		return
	}

	var shouldExit = h.exitAfterParsing
	sep := false

	if h.paramsShowWhereSet {
		if sep {
			printMajorSeparator(twc)
		}
		sep = true
		shouldExit = h.exitAfterHelp

		h.showWhereParamsAreSet(twc, ps)
	}
	if h.paramsShowUnused {
		if sep {
			printMajorSeparator(twc)
		}
		sep = true
		shouldExit = h.exitAfterHelp

		showUnusedParams(twc, ps)
	}

	if h.style != noHelp {
		if sep {
			printMajorSeparator(twc)
		}
		shouldExit = h.exitAfterHelp

		h.Help(ps)
	}

	if shouldExit {
		os.Exit(0)
	}
}
