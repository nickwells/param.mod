package phelp

import (
	"fmt"
	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/param"
	"os"
	"sort"
)

// Help will process the values set after parsing is complete and process any
// errors
func (h StdHelp) ProcessArgs(ps *param.ParamSet) {
	var shouldExit = h.exitAfterParsing

	if h.reportWhereParamsAreSet {
		showWhereParamsAreSet(ps)
		shouldExit = true
	}
	if h.reportUnusedParams {
		showUnusedParams(ps)
		shouldExit = true
	}
	if h.reportParamSources {
		showParamSources(ps)
		shouldExit = true
	}

	if h.showHelp {
		h.Help(ps)
		shouldExit = true
	}

	if shouldExit {
		os.Exit(0)
	}
}

func (h StdHelp) ErrorHandler(ps *param.ParamSet) {
	if len(ps.Errors()) == 0 {
		return
	}

	if !h.dontReportErrors {
		showErrors(ps)
	}

	if h.dontExitOnErrors {
		return
	}
	os.Exit(1)
}

// ReportErrors reports the errors to the param set's error writer. It can be
// used by any Helper and is used by the StdHelp instance.
func ReportErrors(ps *param.ParamSet) {
	errMap := ps.Errors()
	if len(errMap) == 0 {
		return
	}

	paramNames := make([]string, 0, len(errMap))
	for paramName := range errMap {
		paramNames = append(paramNames, paramName)
	}
	sort.Strings(paramNames)

	fmt.Fprint(ps.ErrWriter(), ps.ProgName(), ": ", len(errMap))
	if len(errMap) == 1 {
		fmt.Fprint(ps.ErrWriter(), " error was")
	} else {
		fmt.Fprint(ps.ErrWriter(), " errors were")
	}
	fmt.Fprint(ps.ErrWriter(), " detected while setting the parameters:\n")

	firstLineIndent := stdIndent
	secondLineIndent := stdIndent + stdIndent
	indent2Len := len(secondLineIndent)

	paramSep := ""
	for _, paramName := range paramNames {
		fmt.Fprint(ps.ErrWriter(), paramSep)
		paramSep = firstLineIndent + "---\n"

		sep := ""
		fmt.Fprint(ps.ErrWriter(), firstLineIndent, paramName, "\n")
		for _, e := range errMap[paramName] {
			fmt.Fprint(ps.ErrWriter(), sep)
			sep = secondLineIndent + "---\n"
			switch e := e.(type) {
			case location.Err:
				formatText(ps.ErrWriter(), e.Msg,
					indent2Len, indent2Len)
				formatText(ps.ErrWriter(), "parameter set at: "+e.Loc.String(),
					indent2Len, indent2Len)
			default:
				formatText(ps.ErrWriter(), e.Error(),
					indent2Len, indent2Len)
			}
		}
	}
}

// showErrors reports the errors in the param set and prints advice on how to
// get help
func showErrors(ps *param.ParamSet) {
	ReportErrors(ps)

	formatText(ps.ErrWriter(),
		"\n"+`For help with the correct use of the parameters and to see which`+
			` parameters are available please use the '-`+
			usageArgName+
			`' parameter which will print a usage message`+"\n",
		textIndent, textIndent)
}
