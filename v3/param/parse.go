package param

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/nickwells/location.mod/location"
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
// yourself.
//
// Finally it will process any remaining parameters - these are any
// parameters following a positional parameter that has been marked as
// terminal or any parameters following the terminal parameter (which is "--"
// by default). If no trailing arguments are expected and no handler has been
// set for handling them then the default handler is called which will record
// an error and call the helper.ErrorHandler method.
func (ps *PSet) Parse(args ...[]string) ErrMap {
	if ps.parsed {
		errMap := ErrMap{
			"": []error{
				fmt.Errorf("param.Parse has already been called,"+
					" previously from: %s now from: %s",
					ps.parseCalledFrom,
					caller()),
			},
		}
		ps.helper.ErrorHandler(ps.ErrWriter(), ps.ProgName(), errMap)
		return errMap
	}

	if len(args) == 0 {
		ps.progName = os.Args[0]
		ps.progBaseName = filepath.Base(ps.progName)
	}

	ps.getParamsFromConfigFiles()
	ps.getParamsFromEnvironment()

	var loc *location.L
	if len(args) == 0 {
		loc = location.New("command line")
		ps.getParamsFromStringSlice(loc, os.Args[1:])
	} else {
		loc = location.New("supplied parameters")
		var suppliedParams []string
		for _, sp := range args {
			suppliedParams = append(suppliedParams, sp...)
		}
		ps.getParamsFromStringSlice(loc, suppliedParams)
	}

	ps.detectMandatoryParamsNotSet()

	for _, fcf := range ps.finalChecks {
		err := fcf()
		if err != nil {
			ps.errors[""] = append(ps.errors[""], err)
		}
	}
	ps.setParsed()

	ps.helper.ProcessArgs(ps)
	ps.helper.ErrorHandler(ps.ErrWriter(), ps.ProgName(), ps.errors)

	ps.remHandler.HandleRemainder(ps, loc)

	return ps.errors
}

// caller returns a string giving the filename and line number of the calling
// function. This is intended for providing useful debugging messages
func caller() string {
	if pc, file, line, ok := runtime.Caller(1); ok {
		f := runtime.FuncForPC(pc)
		funcName := "unknown"
		if f != nil {
			funcName = f.Name()
		}
		return fmt.Sprintf("%s:%d [%s]", file, line, funcName)
	}
	return "unknown-file:0 [unknown]"
}

// setParsed records that the parameters have already been parsed and where
// it has been called from
func (ps *PSet) setParsed() {
	ps.parsed = true
	ps.parseCalledFrom = caller()
}

func (ps *PSet) detectMandatoryParamsNotSet() {
	for _, p := range ps.byName {
		if p.AttrIsSet(MustBeSet) &&
			len(p.whereIsParamSet) == 0 {
			ps.errors[p.name] = append(ps.errors[p.name],
				errors.New("this parameter must be set somewhere"))
		}
	}
}
