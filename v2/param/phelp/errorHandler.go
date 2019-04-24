package phelp

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v2/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// ErrorHandler will, by default, check for errors and if there are any
// report them and exit. It will respect the flags for controlling reporting
// of errors and exiting; these flags can be set by means of standard
// arguments as added by the StdHelp AddParams method, see the standard help
// message for details. It is typically called from the Parse(...) method
// being passed the PSet error writer, the program name and the PSet
// error map
func (h StdHelp) ErrorHandler(w io.Writer, name string, errMap param.ErrMap) {
	if len(errMap) == 0 {
		return
	}

	if !h.dontReportErrors {
		ReportErrors(w, name, errMap)
		twc, err := twrap.NewTWConf(twrap.SetWriter(w))
		if err != nil {
			fmt.Fprint(os.Stderr, "Couldn't build the text wrapper:", err)
			return
		}

		twc.Wrap("\nFor help with the correct use of the parameters"+
			" and to see which parameters are available please use the '-"+
			usageArgName+
			"' parameter which will print a usage message\n",
			textIndent)
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
	twc, err := twrap.NewTWConf(twrap.SetWriter(w))
	if err != nil {
		fmt.Fprint(os.Stderr, "Couldn't build the text wrapper:", err)
		return
	}

	paramNames := make([]string, 0, len(errMap))
	for paramName := range errMap {
		paramNames = append(paramNames, paramName)
	}
	sort.Strings(paramNames)

	fmt.Fprint(w, name, ": ")
	if len(errMap) == 1 {
		fmt.Fprint(w, "an error was")
	} else {
		fmt.Fprint(w, len(errMap), " errors were")
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
				twc.Wrap(e.Msg+
					"\nparameter set at: "+e.Loc.String(), indent2Len)
			default:
				twc.Wrap(e.Error(), indent2Len)
			}
		}
	}
}
