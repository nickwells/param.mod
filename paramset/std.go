package paramset

import (
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/phelp"
)

// New creates a new PSet with the standard helper set. This is a suitable
// choice in most cases.
func New(psof ...param.PSetOptFunc) *param.PSet {
	return param.NewSet(phelp.NewStdHelp(), psof...)
}
