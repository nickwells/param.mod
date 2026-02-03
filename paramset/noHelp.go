package paramset

import (
	"os"

	"github.com/nickwells/errutil.mod/errutil"
	"github.com/nickwells/param.mod/v6/param"
)

// noHelp is a minimal implementation of the param.Helper interface. In
// particular there are no parameters added
type noHelp struct{}

func (nh noHelp) ProcessArgs(_ *param.PSet)       {}
func (nh noHelp) Help(_ *param.PSet, _ ...string) {}
func (nh noHelp) AddParams(_ *param.PSet)         {}
func (nh noHelp) ErrorHandler(ps *param.PSet) {
	errMap := ps.Errors()
	if len(errMap) != 0 {
		errutil.ErrMap(errMap).Report(ps.ErrW(), ps.ProgName())
		os.Exit(1)
	}
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
func NewNoHelp(psof ...param.PSetOptFunc) *param.PSet {
	return param.NewSet(append(psof, param.SetHelper(nh))...)
}
