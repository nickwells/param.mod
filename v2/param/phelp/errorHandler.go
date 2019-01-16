package phelp

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v2/param"
)

// ErrorHandler will, by default, check for errors and if there are any
// report them and exit. It will respect the flags for controlling reporting
// of errors and exiting; these flags can be set by means of standard
// arguments as added by the StdHelp AddParams method, see the standard help
// message for details. It is typically called from the Parse(...) method
// being passed the ParamSet error writer, the program name and the ParamSet
// error map
func (h StdHelp) ErrorHandler(w io.Writer, name string, errMap param.ErrMap) {
	if len(errMap) == 0 {
		return
	}

	if !h.dontReportErrors {
		ReportErrors(w, name, errMap)

		formatText(w,
			"\nFor help with the correct use of the parameters"+
				" and to see which parameters are available please use the '-"+
				usageArgName+
				"' parameter which will print a usage message\n",
			textIndent, textIndent)
	}

	if h.dontExitOnErrors {
		return
	}
	os.Exit(1)
}

// ReportErrors reports the errors (if any) to the writer. It can be
// used by any Helper and is used by the StdHelp instance.
func ReportErrors(w io.Writer, name string, errMap param.ErrMap) {
	if len(errMap) == 0 {
		return
	}

	paramNames := make([]string, 0, len(errMap))
	for paramName := range errMap {
		paramNames = append(paramNames, paramName)
	}
	sort.Strings(paramNames)

	fmt.Fprint(w, name, ": ", len(errMap))
	if len(errMap) == 1 {
		fmt.Fprint(w, " error was")
	} else {
		fmt.Fprint(w, " errors were")
	}
	fmt.Fprint(w, " detected while setting the parameters:\n")

	firstLineIndent := stdIndent
	secondLineIndent := stdIndent + stdIndent
	indent2Len := len(secondLineIndent)

	paramSep := ""
	for _, paramName := range paramNames {
		fmt.Fprint(w, paramSep)
		paramSep = firstLineIndent + "---\n"

		sep := ""
		fmt.Fprint(w, firstLineIndent, paramName, "\n")
		for _, e := range errMap[paramName] {
			fmt.Fprint(w, sep)
			sep = secondLineIndent + "---\n"
			switch e := e.(type) {
			case location.Err:
				formatText(w, e.Msg,
					indent2Len, indent2Len)
				formatText(w, "parameter set at: "+e.Loc.String(),
					indent2Len, indent2Len)
			default:
				formatText(w, e.Error(),
					indent2Len, indent2Len)
			}
		}
	}
}
