package phelp

import (
	"fmt"
	"os"

	"github.com/nickwells/param.mod/v4/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// printSep prints the separator line if sep is true. It always returns true
func printSep(twc *twrap.TWConf, sep bool) bool {
	if sep {
		printMajorSeparator(twc)
	}
	return true
}

// ProcessArgs will process the values set after parsing is complete. This is
// where any StdHelp parameters (as added by the StdHelp AddParams method)
// will be processed.
func (h StdHelp) ProcessArgs(ps *param.PSet) {
	var shouldExit = h.exitAfterParsing

	defer func() {
		if shouldExit {
			os.Exit(0)
		}
	}()

	twc, err := twrap.NewTWConf(twrap.SetWriter(ps.StdWriter()))
	if err != nil {
		fmt.Fprintln(os.Stderr,
			"Couldn't build the text wrapper for handling parameters: ", err)
		return
	}

	sep := false

	if h.zshMakeCompletions != zshCompGenNone {
		sep = printSep(twc, sep)
		shouldExit = true

		err := h.zshMakeCompFile(twc, ps)
		if err != nil {
			fmt.Fprintln(ps.ErrWriter(),
				"Couldn't create the zsh completion file: ", err)
			return
		}
	}

	if h.paramsShowWhereSet {
		sep = printSep(twc, sep)
		shouldExit = h.exitAfterHelp

		h.showWhereParamsAreSet(twc, ps)
	}
	if h.paramsShowUnused {
		sep = printSep(twc, sep)
		shouldExit = h.exitAfterHelp

		showUnusedParams(twc, ps)
	}

	if h.style != noHelp {
		_ = printSep(twc, sep)
		shouldExit = h.exitAfterHelp

		h.Help(ps)
	}
}
