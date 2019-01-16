package paramset

import (
	"io"

	"github.com/nickwells/param.mod/v2/param"
	"github.com/nickwells/param.mod/v2/param/phelp"
)

type noHelpNoExit struct{}

func (nh noHelpNoExit) ProcessArgs(ps *param.ParamSet)       {}
func (nh noHelpNoExit) Help(ps *param.ParamSet, s ...string) {}
func (nh noHelpNoExit) AddParams(ps *param.ParamSet)         {}
func (nh noHelpNoExit) ErrorHandler(w io.Writer, name string, errs param.ErrMap) {
	phelp.ReportErrors(w, name, errs)
}

var nhne noHelpNoExit

// NewNoHelpNoExit returns a paramset and any errors encountered while
// creating it. It adds no parameters and doesn't provide a Usage message. It
// does report errors but doesn't exit if Parse errors are seen.
//
// This is only likely to be of any use for testing purposes
func NewNoHelpNoExit(psof ...param.ParamSetOptFunc) (*param.ParamSet, error) {
	return param.NewSet(append(psof, param.SetHelper(nhne))...)
}

type noHelpNoExitNoErrRpt struct{}

func (nh noHelpNoExitNoErrRpt) ProcessArgs(ps *param.ParamSet)                           {}
func (nh noHelpNoExitNoErrRpt) Help(ps *param.ParamSet, s ...string)                     {}
func (nh noHelpNoExitNoErrRpt) AddParams(ps *param.ParamSet)                             {}
func (nh noHelpNoExitNoErrRpt) ErrorHandler(w io.Writer, name string, errs param.ErrMap) {}

var nhnenr noHelpNoExitNoErrRpt

// NewNoHelpNoExitNoErrRpt returns a paramset and any errors encountered while
// creating it. It adds no parameters and doesn't provide a Usage message. It
// does report errors but doesn't exit if Parse errors are seen.
//
// This is only likely to be of any use for testing purposes
func NewNoHelpNoExitNoErrRpt(psof ...param.ParamSetOptFunc) (*param.ParamSet, error) {
	return param.NewSet(append(psof, param.SetHelper(nhnenr))...)
}
