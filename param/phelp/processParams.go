package phelp

// TODO: Rename this file to processArgs.go

import (
	"fmt"
	"os"

	"github.com/nickwells/param.mod/v4/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// printSepIf prints the separator line if sep is true. It always returns true
func printSepIf(twc *twrap.TWConf, printSep bool, sep string) bool {
	if printSep {
		twc.Print(sep)
	}
	return true
}

// ProcessArgs will process the values set after parsing is complete. This is
// where any StdHelp parameters (as added by the StdHelp AddParams method)
// will be processed.
func (h StdHelp) ProcessArgs(ps *param.PSet) {
	var shouldExit = h.exitAfterParsing
	var exitStatus = 0

	defer func() {
		if shouldExit {
			os.Exit(exitStatus)
		}
	}()

	twc := twrap.NewTWConfOrPanic(twrap.SetWriter(ps.StdWriter()))

	printSep := false

	if h.zshMakeCompletions != zshCompGenNone {
		printSep = printSepIf(twc, printSep, majorSectionSeparator)
		shouldExit = true

		err := h.zshMakeCompFile(twc, ps)
		if err != nil {
			fmt.Fprintln(ps.ErrWriter(),
				"Couldn't create the zsh completion file: ", err)
			exitStatus = 1
			return
		}
	}

	if h.paramsShowWhereSet {
		printSep = printSepIf(twc, printSep, majorSectionSeparator)
		shouldExit = shouldExit || h.exitAfterHelp

		h.showWhereParamsAreSet(twc, ps)
	}

	if h.paramsShowUnused {
		printSep = printSepIf(twc, printSep, majorSectionSeparator)
		shouldExit = shouldExit || h.exitAfterHelp

		showUnusedParams(twc, ps)
	}

	if len(ps.Errors()) > 0 && h.reportErrors {
		printSep = printSepIf(twc, printSep, majorSectionSeparator)
		twcErr := twrap.NewTWConfOrPanic(twrap.SetWriter(ps.ErrWriter()))
		ReportErrors(twcErr, ps.ProgName(), ps.Errors())
		shouldExit = shouldExit || h.exitOnErrors
		exitStatus = 1
	}

	if len(h.sectionsChosen) > 0 {
		printSep = printSepIf(twc, printSep, majorSectionSeparator)
		shouldExit = shouldExit || h.exitAfterHelp

		h.Help(ps)
	}

}
