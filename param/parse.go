package param

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Parse will initialise the parameter values
//
// It will first look in the configuration files (if any filenames have been
// set using the SetConfigFile function).
//
// Next it will look in the environment (if any environment prefix strings
// have been set using the SetEnvPrefix function).
//
// Lastly it will process the command line arguments.
//
// It takes zero or more arguments each of which is a slice of strings. If no
// arguments are given then it uses the command line parameters (excluding the
// first which is used to set the program name). If any argument is passed
// then all the slices are concatenated together and that is parsed.
//
// It will return a map of errors: parameter name to a non-empty slice of
// error messages. In order to make sensible use of this the report-errors
// and exit-on-errors flags should be turned off - there are functions
// which allow the caller to do this (or they can be set through the
// command-line flags) but they should be called before Parse is called. The
// default behaviour is to report any errors and exit. This means that you
// can sensibly ignore the return value unless you want to handle the errors
// yourself
func (ps *ParamSet) Parse(args ...[]string) ErrMap {
	if ps.parsed {
		callStack := make([]byte, 10240)
		stackSize := runtime.Stack(callStack, false)

		ps.errors[""] = append(ps.errors[""],
			fmt.Errorf(
				"param.Parse has already been called, from: %s now from: %s",
				ps.parseCalledFrom,
				string(callStack[:stackSize])))
		ps.helper.ErrorHandler(ps)
		return ps.errors
	}

	if len(args) == 0 {
		ps.progName = os.Args[0]
		ps.progBaseName = filepath.Base(ps.progName)
	}

	ps.getParamsFromConfigFile()

	if len(ps.envPrefixes) != 0 {
		ps.getParamsFromEnvironment()
	}

	if len(args) == 0 {
		ps.getParamsFromStringSlice("command line", os.Args[1:])
	} else {
		var suppliedParams []string
		for _, sp := range args {
			suppliedParams = append(suppliedParams, sp...)
		}
		ps.getParamsFromStringSlice("supplied parameters", suppliedParams)
	}

	ps.detectMandatoryParamsNotSet()

	for _, fcf := range ps.finalChecks {
		err := fcf()
		if err != nil {
			ps.errors[""] = append(ps.errors[""], err)
		}
	}
	ps.parsed = true
	callStack := make([]byte, 10240)
	stackSize := runtime.Stack(callStack, false)
	ps.parseCalledFrom = string(callStack[:stackSize])

	ps.helper.ProcessArgs(ps)
	ps.helper.ErrorHandler(ps)

	return ps.errors
}

func makeParamLocDesc(source, location, uneditedParam string) string {
	if location != "" {
		return source + ": " + location + ": " + uneditedParam
	}
	return source + ": " + uneditedParam
}

func (ps *ParamSet) detectMandatoryParamsNotSet() {
	for _, p := range ps.byName {
		if p.attributes&MustBeSet == MustBeSet &&
			len(p.whereIsParamSet) == 0 {
			ps.errors[p.name] = append(ps.errors[p.name],
				errors.New("this parameter must be set somewhere"))
		}
	}
}
