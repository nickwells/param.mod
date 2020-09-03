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
// Before any further processing the helper's ProcessArgs method is
// called. This is expected to act on any helper parameters and to report any
// errors.
//
// Finally it will process any remaining parameters - these are any
// parameters following a positional parameter that has been marked as
// terminal or any parameters following the terminal parameter (which is "--"
// by default). If no trailing arguments are expected and no handler has been
// set for handling them then the default handler is called which will record
// an error and call the helper.ErrorHandler method.
//
// It will return a map of errors: mapping parameter name to a slice of all
// the errors seen for that parameter. In order to make sensible use of this
// the report-errors and exit-on-errors flags should be turned off - there
// are functions which allow the caller to do this (or they can be set
// through the StdHelp command-line flags) but they should be called before
// Parse is called. The default behaviour is to report any errors and
// exit. This means that you can sensibly ignore the return value unless you
// want to handle the errors yourself.
//
// It will panic if it is called twice.
func (ps *PSet) Parse(args ...[]string) ErrMap {
	if ps.parsed {
		panic(
			fmt.Sprintf("param.Parse has already been called,"+
				" previously from: %s now from: %s",
				ps.parseCalledFrom,
				caller()))
	}

	ps.parsed = true
	ps.parseCalledFrom = caller()

	ps.checkSeeAlsoRefs()

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
			ps.AddErr("Final Checks", err)
		}
	}

	ps.helper.ProcessArgs(ps)

	ps.remHandler.HandleRemainder(ps, loc)

	return ps.errors
}

// caller returns a string giving the filename and line number of the caller
// of the calling function. This is intended for providing useful debugging
// messages. Note that we ask for the second stack entry above this: 0 would
// give the location of the call to runtime.Caller, 1 would give the location
// of the call to caller() but we want to see where the parent function was
// called so we pass 2
func caller() string {
	if pc, file, line, ok := runtime.Caller(2); ok {
		f := runtime.FuncForPC(pc)
		funcName := "unknown"
		if f != nil {
			funcName = f.Name()
		}
		return fmt.Sprintf("%s:%d [%s]", file, line, funcName)
	}
	return "unknown-file:0 [unknown]"
}

func (ps *PSet) detectMandatoryParamsNotSet() {
	for _, p := range ps.byName {
		if p.AttrIsSet(MustBeSet) &&
			len(p.whereIsParamSet) == 0 {
			ps.AddErr(p.name,
				errors.New("this parameter must be set somewhere"))
		}
	}
}

// checkSeeAlsoRefs will make sure that every SeeAlso reference is to a valid
// parameter name and will panic if not
func (ps *PSet) checkSeeAlsoRefs() {
	for _, p := range ps.byName {
		refs := p.SeeAlso()
		for _, ref := range refs {
			if _, exists := ps.nameToParam[ref]; !exists {
				panic(
					fmt.Errorf(
						"Parameter %q has a SeeAlso reference to %q"+
							" but no such parameter exists. The bad"+
							" reference was added at: %s",
						p.Name(), ref, p.seeAlsoSource(ref)))
			}
		}
	}
}
