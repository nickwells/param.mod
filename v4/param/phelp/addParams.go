package phelp

import (
	"github.com/nickwells/param.mod/v4/param"
)

// groupNamePfx is the name of the group in which all the param package
// parameters are grouped. You should not give any of your parameter groups
// the same name (it'll be confusing)
const groupNamePfx = "common.params"

// AddParams will add the help parameters into the parameter set
func (h *StdHelp) AddParams(ps *param.PSet) {
	h.addParamHandlingParams(ps)
	h.addParamCompletionParams(ps)
	h.addUsageParams(ps)
}
