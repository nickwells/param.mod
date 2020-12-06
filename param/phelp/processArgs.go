package phelp

import (
	"os"

	"github.com/nickwells/param.mod/v5/param"
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

	actions := []struct {
		shouldRun  func() bool
		shouldExit bool
		action     func(StdHelp, *twrap.TWConf, *param.PSet) int
		exitStatus int // only used if action is nil
	}{
		{
			func() bool { return zshCompHasAction(h) },
			zshCompHasAction(h), zshCompletionHandler, 0,
		},
		{
			func() bool { return h.paramsShowWhereSet },
			h.exitAfterHelp, showWhereParamsAreSet, 0,
		},
		{
			func() bool { return h.paramsShowUnused },
			h.exitAfterHelp, showUnusedParams, 0,
		},
		{
			func() bool { return len(ps.Errors()) > 0 && h.reportErrors },
			h.exitOnErrors, reportErrors, 0,
		},
		{
			func() bool { return len(ps.Errors()) > 0 },
			h.exitOnErrors, nil, 1,
		},
		{
			func() bool { return !h.sectionsChosen.hasNothingChosen() },
			h.exitAfterHelp, help, 0,
		},
	}

	twc := twrap.NewTWConfOrPanic(twrap.SetWriter(ps.StdWriter()))

	printSep := false

	for _, a := range actions {
		if !a.shouldRun() {
			continue
		}

		var es int
		if a.action != nil {
			printSep = printSepIf(twc, printSep, majorSectionSeparator)
			es = a.action(h, twc, ps)
		} else {
			es = a.exitStatus
		}
		if es > exitStatus && a.shouldExit {
			exitStatus = es
		}

		shouldExit = shouldExit || a.shouldExit
	}

	if shouldExit {
		os.Exit(exitStatus)
	}
}

func help(h StdHelp, _ *twrap.TWConf, ps *param.PSet) int {
	h.Help(ps)
	return 0
}
