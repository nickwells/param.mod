package phelp

import (
	"fmt"

	"github.com/nickwells/errutil.mod/errutil"
	"github.com/nickwells/param.mod/v7/param"
)

// ErrorHandler will, by default, check for errors and if there are any
// report them and exit. It will respect the flags for controlling reporting
// of errors and exiting; these flags can be set by means of standard
// arguments as added by the StdHelp AddParams method, see the standard help
// message for details. It is typically called from the Parse(...) method
// being passed the PSet error writer, the program name and the PSet
// error map
func (h StdHelp) ErrorHandler(ps *param.PSet) {
	errMap := ps.Errors()

	if len(errMap) == 0 {
		return
	}

	if h.reportErrors {
		errutil.ErrMap(errMap).Report(h.ErrW(), ps.ProgName())
		fmt.Fprintln(h.ErrW(), "\n"+suggestHelpParam(ps))
	}

	if h.exitOnErrors {
		ps.SetExitStatus(exitStatusErrorsFound)
	}
}

// suggestHelpParam returns a string suggesting the standard help parameter
func suggestHelpParam(ps *param.PSet) string {
	return fmt.Sprintf("For more information use the %q parameter",
		ps.ShortestPrefix()+helpArgName)
}
