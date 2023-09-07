package paramset

import (
	"fmt"

	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/phelp"
	"github.com/nickwells/twrap.mod/twrap"
)

type noHelpNoExit struct{}

func (nh noHelpNoExit) ProcessArgs(_ *param.PSet)       {}
func (nh noHelpNoExit) Help(_ *param.PSet, _ ...string) {}
func (nh noHelpNoExit) AddParams(_ *param.PSet)         {}
func (nh noHelpNoExit) ErrorHandler(ps *param.PSet, errs param.ErrMap) {
	twc, err := twrap.NewTWConf(twrap.SetWriter(ps.ErrW()))
	if err != nil {
		panic(fmt.Sprint("Couldn't build the text wrapper:", err))
	}
	phelp.ReportErrors(twc, ps.ProgName(), errs)
}

var nhne noHelpNoExit

// NewNoHelpNoExit returns a paramset and any errors encountered while
// creating it. It adds no parameters and doesn't provide a Usage message. It
// does report errors but doesn't exit if Parse errors are seen.
//
// This is only likely to be of any use for testing purposes
func NewNoHelpNoExit(psof ...param.PSetOptFunc) (*param.PSet, error) {
	return param.NewSet(append(psof, param.SetHelper(nhne))...)
}

// NewNoHelpNoExitOrPanic returns a paramset as per NewNoHelpNoExit but any
// error will cause the program to panic.
//
// This is only likely to be of any use for testing purposes
func NewNoHelpNoExitOrPanic(psof ...param.PSetOptFunc) *param.PSet {
	ps, err := param.NewSet(append(psof, param.SetHelper(nhne))...)
	if err != nil {
		panic(err)
	}
	return ps
}

type noHelpNoExitNoErrRpt struct{}

func (nh noHelpNoExitNoErrRpt) ProcessArgs(*param.PSet)                {}
func (nh noHelpNoExitNoErrRpt) Help(*param.PSet, ...string)            {}
func (nh noHelpNoExitNoErrRpt) AddParams(*param.PSet)                  {}
func (nh noHelpNoExitNoErrRpt) ErrorHandler(*param.PSet, param.ErrMap) {}

var nhnenr noHelpNoExitNoErrRpt

// NewNoHelpNoExitNoErrRpt returns a paramset and any errors encountered while
// creating it. It adds no parameters and doesn't provide a Usage message. It
// doesn't report errors and doesn't exit if Parse errors are seen.
//
// This is only likely to be of any use for testing purposes
func NewNoHelpNoExitNoErrRpt(psof ...param.PSetOptFunc) (*param.PSet, error) {
	return param.NewSet(append(psof, param.SetHelper(nhnenr))...)
}

// NewNoHelpNoExitNoErrRptOrPanic returns a paramset as per
// NewNoHelpNoExitNoErrRpt but any error will cause the program to panic.
//
// This is only likely to be of any use for testing purposes
func NewNoHelpNoExitNoErrRptOrPanic(psof ...param.PSetOptFunc) *param.PSet {
	ps, err := param.NewSet(append(psof, param.SetHelper(nhnenr))...)
	if err != nil {
		panic(err)
	}
	return ps
}
