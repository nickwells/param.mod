package paramset

import (
	"github.com/nickwells/param.mod/param"
	"github.com/nickwells/param.mod/param/phelp"
	"os"
)

// New creates a new ParamSet with the standard helper set. This is the
// one you should use in most cases
func New(psof ...param.ParamSetOptFunc) (*param.ParamSet, error) {
	opts := make([]param.ParamSetOptFunc, 0, len(psof)+1)
	opts = append(opts, param.SetHelper(&phelp.SH))
	opts = append(opts, psof...)
	return param.NewSet(opts...)
}

// noHelp is a minimal implementation of the param.Helper interface. In
// particular there are no parameters added
type noHelp struct{}

func (nh noHelp) ProcessArgs(ps *param.ParamSet)       {}
func (nh noHelp) Help(ps *param.ParamSet, s ...string) {}
func (nh noHelp) AddParams(ps *param.ParamSet)         {}
func (nh noHelp) ErrorHandler(ps *param.ParamSet) {
	phelp.ReportErrors(ps)

	os.Exit(1)
}

var nh noHelp

// NewNoHelp creates a new ParamSet with the helper set to the noHelp
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
func NewNoHelp(psof ...param.ParamSetOptFunc) (*param.ParamSet, error) {
	return param.NewSet(append(psof, param.SetHelper(nh))...)
}
