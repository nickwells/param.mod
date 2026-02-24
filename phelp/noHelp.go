package phelp

import (
	"os"

	"github.com/nickwells/param.mod/v7/param"
)

// NoHelp is a minimal implementation of the param.Helper interface. In
// particular there are no parameters added
type NoHelp struct{}

func (nh NoHelp) ProcessArgs(_ *param.PSet)       {}
func (nh NoHelp) Help(_ *param.PSet, _ ...string) {}
func (nh NoHelp) AddParams(_ *param.PSet)         {}
func (nh NoHelp) ErrorHandler(ps *param.PSet) {
	errMap := ps.Errors()
	if len(errMap) != 0 {
		errMap.Report(os.Stderr, ps.ProgName())
		ps.SetExitStatus(1)
	}
}
