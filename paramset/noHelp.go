package paramset

import (
	"os"

	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/phelp"
	"github.com/nickwells/twrap.mod/twrap"
)

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

	twc := twrap.NewTWConfOrPanic(twrap.SetWriter(ps.ErrW()))

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
