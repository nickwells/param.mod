package param

import (
	"fmt"

	"github.com/nickwells/english.mod/english"
)

// reportUnexpectedTrailingParams adds an error if there are trailing
// parameters and the trailingParamsExpected flag was not set.
func (ps *PSet) reportUnexpectedTrailingParams() {
	if ps.trailingParamsExpected {
		return
	}

	remCount := len(ps.Remainder())

	if remCount == 0 {
		return
	}

	remStr := english.JoinQuoted(ps.Remainder(), " ", " ")
	etc := "..."

	const maxLen = 20

	if len(remStr) > maxLen {
		remStr = remStr[0:maxLen-len(etc)] + etc
	}

	if remCount == 1 {
		ps.AddErr("",
			fmt.Errorf("there was an unexpected extra parameter: %s", remStr))
	} else {
		ps.AddErr("",
			fmt.Errorf("there were %d unexpected extra parameters: %s",
				remCount, remStr))
	}
}
