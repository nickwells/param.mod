package phelp

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v3/param"
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
			helpArgName+
			"' parameter which will print a usage message\n",
			0)
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

	totErrs := 0
	for _, paramName := range paramNames {
		totErrs += len(errMap[paramName])
	}
	twc.WrapPrefixed(name+": ",
		fmt.Sprintf("%d %s detected while setting the parameters",
			totErrs, errOrErrs(totErrs)),
		0)

	for _, pName := range paramNames {
		errCount := len(errMap[pName])
		twc.Wrap(
			fmt.Sprintf("%s [%d %s]", pName, errCount, errOrErrs(errCount)),
			paramIndent)
		for i, e := range errMap[pName] {
			prefix := ""
			if errCount > 1 {
				prefix = fmt.Sprintf("%d : ", i+1)
			}
			switch e := e.(type) {
			case location.Err:
				twc.WrapPrefixed(prefix, e.Msg+
					"\nparameter set at: "+e.Loc.String(), descriptionIndent)
			default:
				twc.WrapPrefixed(prefix, e.Error(), descriptionIndent)
			}
		}
	}
}

// errOrErrs returns "error" or "errors" depending on the value of n
func errOrErrs(n int) string {
	if n == 1 {
		return "error"
	}
	return "errors"
}
