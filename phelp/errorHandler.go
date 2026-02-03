package phelp

import (
	"os"

	"github.com/nickwells/errutil.mod/errutil"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// ErrorHandler will, by default, check for errors and if there are any
// report them and exit. It will respect the flags for controlling reporting
// of errors and exiting; these flags can be set by means of standard
// arguments as added by the StdHelp AddParams method, see the standard help
// message for details. It is typically called from the Parse(...) method
// being passed the PSet error writer, the program name and the PSet
// error map
func (h StdHelp) ErrorHandler(ps *param.PSet) {
	exitStatus := reportErrors(h, nil, ps)

	if exitStatus != 0 && h.exitOnErrors {
		os.Exit(exitStatus)
	}
}

// reportErrors reports the errors (if any) and returns a non-zero exit
// status if any errors were detected.
func reportErrors(h StdHelp, _ *twrap.TWConf, ps *param.PSet) int {
	errMap := ps.Errors()

	if len(errMap) == 0 {
		return 0
	}

	if h.reportErrors {
		errutil.ErrMap(errMap).Report(ps.ErrW(), ps.ProgName())

		twc := twrap.NewTWConfOrPanic(
			twrap.SetWriter(ps.ErrW()),
			twrap.SetTargetLineLen(h.helpLineLen))
		twc.Wrap("\n"+suggestHelpParam(ps), 0)
	}

	return 1
}

// suggestHelpParam returns a string suggesting the standard help parameter
func suggestHelpParam(ps *param.PSet) string {
	return "For more information use the '" +
		ps.ShortestPrefix() + helpArgName +
		"' parameter."
}
