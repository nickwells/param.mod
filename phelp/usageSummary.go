package phelp

import (
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// showUsageSummary prints the program name and a parameter summary
func showUsageSummary(_ StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	twc.Print("Usage: ", ps.ProgName())

	bppCount := ps.CountByPosParams()

	var lastPosParamIsTerminal bool
	for i := 0; i < bppCount; i++ {
		bp, _ := ps.GetParamByPos(i)
		twc.Print(" <", bp.Name(), ">")
		if bp.IsTerminal() {
			lastPosParamIsTerminal = true
		}
	}

	if !lastPosParamIsTerminal {
		groups := ps.GetGroups()
		var hasOptionalParams bool
		for _, g := range groups {
			for _, bn := range g.Params {
				if bn.AttrIsSet(param.MustBeSet) {
					twc.Print(" -" + bn.Name() +
						valueNeededStr(bn))
				} else {
					hasOptionalParams = true
				}
			}
		}
		if hasOptionalParams {
			twc.Print(" ...")
		}

		if ps.TrailingParamsExpected() {
			twc.Print(" " + ps.TerminalParam())
		}
	}

	if ps.TrailingParamsExpected() {
		twc.Print(" " + ps.TrailingParamsName() + "...")
	}

	twc.Print("\n")
	return true
}
