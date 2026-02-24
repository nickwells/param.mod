package phelp

import (
	"os"

	"github.com/nickwells/param.mod/v7/param"
)

// NoHelp minimally implements the param.Helper interface.
type NoHelp struct{}

// ProcessArgs is a null implementation of the ProcessArgs method in the
// param.Helper interface
func (nh NoHelp) ProcessArgs(_ *param.PSet) {}

// Help is a null implementation of the Help method in the param.Helper
// interface
func (nh NoHelp) Help(_ *param.PSet, _ ...string) {}

// AddParams is a null implementation of the AddParams method in the
// param.Helper interface
func (nh NoHelp) AddParams(_ *param.PSet) {}

// ErrorHandler is a minimal implementation of the ErrorHandler method in the
// param.Helper interface. It will check to see if any errors are present as
// reported by the param.PSet Errors method and if so it will report them and
// set the param.PSet exit status.
func (nh NoHelp) ErrorHandler(ps *param.PSet) {
	errMap := ps.Errors()
	if len(errMap) != 0 {
		errMap.Report(os.Stderr, ps.ProgName())
		ps.SetExitStatus(1)
	}
}
