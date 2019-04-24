package phelp

import (
	"fmt"
	"sort"

	"github.com/nickwells/param.mod/v2/param"
)

func showUnusedParams(ps *param.PSet) {
	fmt.Fprintln(ps.ErrWriter(), dashes)
	up := ps.UnusedParams()
	unusedParamCount := len(up)
	if unusedParamCount == 0 {
		fmt.Fprint(ps.ErrWriter(),
			ps.ProgName(), ": there were no unused parameters\n")
		return
	}

	fmt.Fprint(ps.ErrWriter(),
		ps.ProgName(), ": ", unusedParamCount)
	if unusedParamCount == 1 {
		fmt.Fprint(ps.ErrWriter(), " parameter was")
	} else {
		fmt.Fprint(ps.ErrWriter(), " parameters were")
	}
	fmt.Fprintln(ps.ErrWriter(), " set but not used:")

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
