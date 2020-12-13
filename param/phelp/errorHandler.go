package phelp

import (
	"os"

	"github.com/nickwells/errutil.mod/errutil"
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// ErrorHandler will, by default, check for errors and if there are any
// report them and exit. It will respect the flags for controlling reporting
// of errors and exiting; these flags can be set by means of standard
// arguments as added by the StdHelp AddParams method, see the standard help
// message for details. It is typically called from the Parse(...) method
// being passed the PSet error writer, the program name and the PSet
// error map
func (h StdHelp) ErrorHandler(ps *param.PSet, errMap param.ErrMap) {
	if len(errMap) == 0 {
		return
	}

	if !h.reportErrors {
		return
	}
	twc := twrap.NewTWConfOrPanic(twrap.SetWriter(ps.ErrWriter()))

	ReportErrors(twc, ps.ProgName(), errMap)

	if !h.exitOnErrors {
		return
	}

	os.Exit(1)
}

// reportErrors will report that errors have been detected and then call
// ReportErrors but writing to the error writer
func reportErrors(_ StdHelp, _ *twrap.TWConf, ps *param.PSet) int {
	twc := twrap.NewTWConfOrPanic(twrap.SetWriter(ps.ErrWriter()))
	ReportErrors(twc, ps.ProgName(), ps.Errors())
	return 1
}

// ReportErrors reports the errors (if any) to the writer. It can be
// used by any Helper and is used by the StdHelp instance.
func ReportErrors(twc *twrap.TWConf, name string, errMap param.ErrMap) {
	if len(errMap) == 0 {
		return
	}

	errutil.ErrMap(errMap).Report(twc.W, name)

	twc.Wrap("\nTry the '-"+helpArgName+"' parameter for more information.", 0)
}
