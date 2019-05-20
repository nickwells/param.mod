package param

import "io"

// Helper is the interface that a helper object must implement. It should
// supply a set of default parameters to be added by the AddParams func and a
// func (called ProcessArgs) to be called after the parsing is complete which
// will operate on the default parameter values. There should be a Help func
// for reporting a help message and an error handler for reporting errors.
type Helper interface {
	ProcessArgs(ps *PSet)
	ErrorHandler(w io.Writer, name string, errMap ErrMap)
	Help(ps *PSet, messages ...string)
	AddParams(ps *PSet)
}
