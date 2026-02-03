package paramset

import (
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/phelp"
)

// addHelperToOpts will take the slice of param set option functions and a
// SetHelper function to the start.
func addHelperToOpts(psof []param.PSetOptFunc) []param.PSetOptFunc {
	opts := make([]param.PSetOptFunc, 0, len(psof)+1)
	opts = append(opts, param.SetHelper(phelp.NewStdHelp()))
	opts = append(opts, psof...)

	return opts
}

// New creates a new PSet with the standard helper set. This is a suitable
// choice in most cases.
func New(psof ...param.PSetOptFunc) *param.PSet {
	opts := addHelperToOpts(psof)
	return param.NewSet(opts...)
}
