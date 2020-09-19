package phelp

import (
	"fmt"
	"os"
	"sort"

	"github.com/nickwells/location.mod/location"
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
func reportErrors(_ StdHelp, twc *twrap.TWConf, ps *param.PSet) int {
	twc.Print("Parameter errors detected\n")
	twcErr := twrap.NewTWConfOrPanic(twrap.SetWriter(ps.ErrWriter()))
	ReportErrors(twcErr, ps.ProgName(), ps.Errors())
	return 1
}

// ReportErrors reports the errors (if any) to the writer. It can be
// used by any Helper and is used by the StdHelp instance.
func ReportErrors(twc *twrap.TWConf, name string, errMap param.ErrMap) {
	if len(errMap) == 0 {
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
	twc.Wrap("\nTry the '-"+helpArgName+"' parameter for more information.", 0)
}

// reportErrorCount calculates the number of errors in errMap and reports it
func reportErrorCount(twc *twrap.TWConf, name string, errMap param.ErrMap) {
	totErrs := 0
	paramsWithErrs := len(errMap)
	for _, errs := range errMap {
		totErrs += len(errs)
	}
	const detected = "detected while setting the parameters:"
	var msg string
	if totErrs == 1 {
		msg = "an error was " + detected
	} else if paramsWithErrs == 1 {
		msg = fmt.Sprintf("%d errors were "+detected, totErrs)
	} else {
		msg = fmt.Sprintf("%d errors with %d parameters were "+detected,
			totErrs, paramsWithErrs)
	}
	twc.WrapPrefixed(name+": ", msg, 0)
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
			twc.WrapPrefixed(prefix, e.Msg+"\nAt: "+e.Loc.String(), descriptionIndent)
		default:
			twc.WrapPrefixed(prefix, e.Error(), descriptionIndent)
		}
	}
}
