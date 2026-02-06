package phelp

import (
	"github.com/nickwells/param.mod/v6/param"
)

// showUsageSummary prints the program name and a parameter summary. Only
// mandatory parameters are shown.
func showUsageSummary(h StdHelp, ps *param.PSet) bool {
	h.twc.Print("Usage: ", ps.ProgName())

	bppCount := ps.CountByPosParams()

	var lastPosParamIsTerminal bool

	for i := range bppCount {
		bp, _ := ps.GetParamByPos(i)

		h.twc.Print(" <", bp.Name(), ">")

		if bp.IsTerminal() {
			lastPosParamIsTerminal = true
		}
	}

	if !lastPosParamIsTerminal {
		var hasOptionalParams bool

		groups := ps.GetGroups()
		for _, g := range groups {
			for _, bn := range g.Params() {
				if bn.AttrIsSet(param.MustBeSet) {
					h.twc.Print(" " + ParamShortSummary(*bn))
				} else {
					hasOptionalParams = true
				}
			}
		}

		if hasOptionalParams {
			h.twc.Print(" ...")
		}

		if ps.TrailingParamsExpected() {
			h.twc.Print(" " + ps.TerminalParam())
		}
	}

	if ps.TrailingParamsExpected() {
		h.twc.Print(" " + ps.TrailingParamsName() + "...")
	}

	h.twc.Print("\n")

	return true
}
