package paramset

import (
	"io"

	"github.com/nickwells/param.mod/v3/param"
	"github.com/nickwells/param.mod/v3/param/phelp"
)

type noHelpNoExit struct{}

func (nh noHelpNoExit) ProcessArgs(ps *param.PSet)       {}
func (nh noHelpNoExit) Help(ps *param.PSet, s ...string) {}
func (nh noHelpNoExit) AddParams(ps *param.PSet)         {}
func (nh noHelpNoExit) ErrorHandler(w io.Writer, name string, errs param.ErrMap) {
	phelp.ReportErrors(w, name, errs)
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

type noHelpNoExitNoErrRpt struct{}

func (nh noHelpNoExitNoErrRpt) ProcessArgs(ps *param.PSet)                               {}
func (nh noHelpNoExitNoErrRpt) Help(ps *param.PSet, s ...string)                         {}
func (nh noHelpNoExitNoErrRpt) AddParams(ps *param.PSet)                                 {}
func (nh noHelpNoExitNoErrRpt) ErrorHandler(w io.Writer, name string, errs param.ErrMap) {}

var nhnenr noHelpNoExitNoErrRpt

// NewNoHelpNoExitNoErrRpt returns a paramset and any errors encountered while
// creating it. It adds no parameters and doesn't provide a Usage message. It
// does report errors but doesn't exit if Parse errors are seen.
//
// This is only likely to be of any use for testing purposes
func NewNoHelpNoExitNoErrRpt(psof ...param.PSetOptFunc) (*param.PSet, error) {
	return param.NewSet(append(psof, param.SetHelper(nhnenr))...)
}
