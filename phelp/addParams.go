package phelp

import (
	"github.com/nickwells/param.mod/v6/param"
)

// groupNamePfx is a standard prefix applied to the names of the groups in
// which all the param package parameters are grouped. You should not give
// any of your parameter groups the same name (it'll be confusing)
const groupNamePfx = "stdParams"

// CommonParamsGroupNamePrefix returns the prefix used to start the names of
// all the common parameter groups
func CommonParamsGroupNamePrefix() string {
	return groupNamePfx
}

// AddParams will add the help parameters into the parameter set
func (h *StdHelp) AddParams(ps *param.PSet) {
	h.addParamHandlingParams(ps)
	h.addParamCompletionParams(ps)
	h.addUsageParams(ps)
	addNotes(ps)
}
