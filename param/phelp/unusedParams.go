package phelp

import (
	"fmt"
	"github.com/nickwells/param.mod/param"
	"sort"
)

func showUnusedParams(ps *param.ParamSet) {
	up := ps.UnusedParams()
	unusedParamCount := len(up)
	if unusedParamCount == 0 {
		fmt.Fprint(ps.ErrWriter(),
			ps.ProgName(), ": there were no unused parameters\n")
		return
	}

	fmt.Fprint(ps.ErrWriter(),
		ps.ProgName(), ": ", unusedParamCount,
		" parameters were set but not used:\n")
	var paramsByName = make([]string, 0, unusedParamCount)
	for name := range up {
		paramsByName = append(paramsByName, name)
	}
	sort.Strings(paramsByName)
	for _, pn := range paramsByName {
		fmt.Fprintln(ps.ErrWriter(), "\t", pn)
		for _, loc := range up[pn] {
			fmt.Fprintln(ps.ErrWriter(), "\t\tat: ", loc)
		}
	}
}
