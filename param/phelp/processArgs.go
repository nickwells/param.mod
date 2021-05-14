package phelp

import (
	"os"

	"github.com/nickwells/pager.mod/pager"
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

// helpAction records an action to be performed while processing the help
// arguments. Note that the shouldRun member is a func as otherwise errors
// raised by earlier actions (zshCompletionHandler) cannot be detected when
// we come to report errors.
type helpAction struct {
	shouldRun  func() bool
	shouldExit bool
	action     func(StdHelp, *twrap.TWConf, *param.PSet) int
}

// ProcessArgs will process the values set after parsing is complete. This is
// where any StdHelp parameters (as added by the StdHelp AddParams method)
// will be processed.
func (h StdHelp) ProcessArgs(ps *param.PSet) {
	actions := []helpAction{
		{
			shouldRun:  func() bool { return zshCompHasAction(h) },
			shouldExit: zshCompHasAction(h),
			action:     zshCompletionHandler,
		},
		{
			shouldRun:  func() bool { return h.paramsShowWhereSet },
			shouldExit: h.exitAfterHelp,
			action:     showWhereParamsAreSet,
		},
		{
			shouldRun:  func() bool { return h.paramsShowUnused },
			shouldExit: h.exitAfterHelp,
			action:     showUnusedParams,
		},
		{
			shouldRun:  func() bool { return len(ps.Errors()) > 0 },
			shouldExit: h.exitOnErrors,
			action:     reportErrors,
		},
		{
			shouldRun:  func() bool { return !h.sectionsChosen.hasNothingChosen() },
			shouldExit: h.exitAfterHelp,
			action:     help,
		},
	}

	var twc *twrap.TWConf
	printSep := false
	var pgr *pager.Pager

	var shouldExit = h.exitAfterParsing
	var exitStatus = 0

	for _, a := range actions {
		if !a.shouldRun() {
			continue
		}

		if twc == nil {
			if h.pageOutput {
				pgr = pager.Start(ps)
			}
			twc = twrap.NewTWConfOrPanic(twrap.SetWriter(ps.StdWriter()))
		}

		if h.helpFormat != helpFmtTypeMarkdown {
			printSep = printSepIf(twc, printSep, majorSectionSeparator)
		}
		es := a.action(h, twc, ps)

		if es > exitStatus && a.shouldExit {
			exitStatus = es
		}

		shouldExit = shouldExit || a.shouldExit
	}

	pgr.Done()

	if shouldExit {
		os.Exit(exitStatus)
	}
}

func help(h StdHelp, _ *twrap.TWConf, ps *param.PSet) int {
	h.Help(ps)
	return 0
}
