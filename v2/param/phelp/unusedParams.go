package phelp

import (
	"fmt"
	"sort"

	"github.com/nickwells/param.mod/v2/param"
)

func showUnusedParams(ps *param.PSet) {
	up := ps.UnusedParams()
	unusedParamCount := len(up)
	if unusedParamCount == 0 {
		fmt.Fprint(ps.ErrWriter(), // nolint: errcheck
			ps.ProgName(), ": there were no unused parameters\n")
		return
	}

	fmt.Fprint(ps.ErrWriter(), // nolint: errcheck
		ps.ProgName(), ": ", unusedParamCount,
		" parameters were set but not used:\n")
	var paramsByName = make([]string, 0, unusedParamCount)
	for name := range up {
		paramsByName = append(paramsByName, name)
	}
	sort.Strings(paramsByName)
	for _, pn := range paramsByName {
		fmt.Fprintln(ps.ErrWriter(), "\t", pn) // nolint: errcheck
		for _, loc := range up[pn] {
			fmt.Fprintln(ps.ErrWriter(), "\t\tat: ", loc) // nolint: errcheck
		}
	}
}
