package phelp

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v4/param"
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

		twc.Wrap("\nTry the '-"+helpArgName+
			"' parameter for more information.\n",
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

	reportErrorCount(twc, name, errMap)

	for _, pName := range paramNames {
		reportParamError(twc, pName, errMap[pName])
	}
}

// reportErrorCount calculates the number of errors in errMap and reports it
func reportErrorCount(twc *twrap.TWConf, name string, errMap param.ErrMap) {
	totErrs := 0
	paramsWithErrs := len(errMap)
	for _, errs := range errMap {
		totErrs += len(errs)
	}
	if totErrs == 1 {
		twc.WrapPrefixed(name+": ",
			"an error was detected while setting the parameters:",
			0)
	} else if paramsWithErrs == 1 {
		twc.WrapPrefixed(name+": ",
			fmt.Sprintf("%d errors were detected while setting the parameters:",
				totErrs),
			0)
	} else {
		twc.WrapPrefixed(name+": ",
			fmt.Sprintf(
				"%d errors with %d parameters were detected"+
					" while setting the parameters:",
				totErrs, paramsWithErrs),
			0)
	}
}

// reportParamError reports all the errors for an individual parameter
func reportParamError(twc *twrap.TWConf, pName string, errs []error) {
	errCount := len(errs)
	if errCount == 1 {
		twc.Wrap(pName, paramIndent)
	} else {
		twc.Wrap(fmt.Sprintf("%s - %d errors:", pName, errCount),
			paramIndent)
	}
	for i, e := range errs {
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
