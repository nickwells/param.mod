package paramset

import (
	"fmt"
	"os"

	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/phelp"
)

// addHelperToOpts will take the slice of param set option functions and a
// SetHelper function to the start.
func addHelperToOpts(psof []param.PSetOptFunc) []param.PSetOptFunc {
	opts := make([]param.PSetOptFunc, 0, len(psof)+1)
	opts = append(opts, param.SetHelper(phelp.NewStdHelp()))
	opts = append(opts, psof...)

	return opts
}

// New creates a new PSet with the standard helper set. This is a suitable
// choice in most cases.
func New(psof ...param.PSetOptFunc) (*param.PSet, error) {
	opts := addHelperToOpts(psof)
	return param.NewSet(opts...)
}

// NewOrDie creates a new PSet with the standard helper set. It then checks
// the error returned and if it is not nil it will report the error on stderr
// and exit with a non-zero exit status. This is a suitable choice unless you
// want to perform any special error handling.
func NewOrDie(psof ...param.PSetOptFunc) *param.PSet {
	ps, err := New(psof...)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"The program parameter set can't be made: %s", err)
		os.Exit(1)
	}

	return ps
}

// NewOrPanic creates a new PSet with the standard helper set. It then checks
// the error returned and if it is not nil it will panic with the returned
// error wrapped in an explanatory message. This is a suitable choice if you
// want to test the parameter generation and still have a simple API in your
// program.
func NewOrPanic(psof ...param.PSetOptFunc) *param.PSet {
	ps, err := New(psof...)
	if err != nil {
		panic(fmt.Errorf("the program parameter set can't be made: %w", err))
	}

	return ps
}
