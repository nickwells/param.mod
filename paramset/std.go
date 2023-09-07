package paramset

import (
	"fmt"
	"os"

	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/phelp"
	"github.com/nickwells/twrap.mod/twrap"
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
	opts := addHelperToOpts(psof)
	ps, err := param.NewSet(opts...)
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
	opts := addHelperToOpts(psof)
	ps, err := param.NewSet(opts...)
	if err != nil {
		panic(fmt.Errorf("The program parameter set can't be made: %w", err))
	}
	return ps
}

// noHelp is a minimal implementation of the param.Helper interface. In
// particular there are no parameters added
type noHelp struct{}

func (nh noHelp) ProcessArgs(_ *param.PSet)       {}
func (nh noHelp) Help(_ *param.PSet, _ ...string) {}
func (nh noHelp) AddParams(_ *param.PSet)         {}
func (nh noHelp) ErrorHandler(ps *param.PSet, errs param.ErrMap) {
	if len(errs) == 0 {
		return
	}
	twc, err := twrap.NewTWConf(twrap.SetWriter(ps.ErrW()))
	if err != nil {
		panic(fmt.Sprint("Couldn't build the text wrapper:", err))
	}

	phelp.ReportErrors(twc, ps.ProgName(), errs)

	os.Exit(1)
}

var nh noHelp

// NewNoHelp creates a new PSet with the helper set to the noHelp
// helper which does nothing. In particular it will not add any parameters
// and so it returns a suitable parameter set for the case where you want to
// add positional parameters the last of which is terminal. This style of
// interface is used if you have a positional parameter which will invoke a
// different command based on the value - see the 'git' or 'go' commands for
// examples of this CLI interface style. If you are choosing an interface
// like this you might want to consider having one of the possible parameter
// values being "help" so that the available options can be listed.
//
// If errors are detected then they will be reported and the program will
// exit.
func NewNoHelp(psof ...param.PSetOptFunc) (*param.PSet, error) {
	return param.NewSet(append(psof, param.SetHelper(nh))...)
}
