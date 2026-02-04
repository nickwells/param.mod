package paramset

import (
	"github.com/nickwells/errutil.mod/errutil"
	"github.com/nickwells/param.mod/v6/param"
)

type noHelpNoExit struct{}

func (nh noHelpNoExit) ProcessArgs(_ *param.PSet)       {}
func (nh noHelpNoExit) Help(_ *param.PSet, _ ...string) {}
func (nh noHelpNoExit) AddParams(_ *param.PSet)         {}
func (nh noHelpNoExit) ErrorHandler(ps *param.PSet) {
	errMap := ps.Errors()
	if len(errMap) != 0 {
		errutil.ErrMap(errMap).Report(ps.ErrW(), ps.ProgName())
	}
}

var nhne noHelpNoExit

// NewNoHelpNoExit returns a paramset and any errors encountered while
// creating it. It adds no parameters and doesn't provide a Usage message. It
// does report errors but doesn't exit if Parse errors are seen.
//
// This is only likely to be of any use for testing purposes
func NewNoHelpNoExit(psof ...param.PSetOptFunc) *param.PSet {
	return param.NewSet(nhne, psof...)
}

type noHelpNoExitNoErrRpt struct{}

func (nh noHelpNoExitNoErrRpt) ProcessArgs(*param.PSet)     {}
func (nh noHelpNoExitNoErrRpt) Help(*param.PSet, ...string) {}
func (nh noHelpNoExitNoErrRpt) AddParams(*param.PSet)       {}
func (nh noHelpNoExitNoErrRpt) ErrorHandler(*param.PSet)    {}

var nhnenr noHelpNoExitNoErrRpt

// NewNoHelpNoExitNoErrRpt returns a paramset and any errors encountered while
// creating it. It adds no parameters and doesn't provide a Usage message. It
// doesn't report errors and doesn't exit if Parse errors are seen.
//
// This is only likely to be of any use for testing purposes
func NewNoHelpNoExitNoErrRpt(psof ...param.PSetOptFunc) *param.PSet {
	return param.NewSet(nhnenr, psof...)
}
