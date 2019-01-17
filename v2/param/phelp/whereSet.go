package phelp

import (
	"fmt"

	"github.com/nickwells/param.mod/v2/param"
)

func showWhereParamsAreSet(ps *param.ParamSet) {
	paramGroups := ps.GetGroups()
	w := ps.StdWriter()

	for _, pg := range paramGroups {
		printGroupDetails(w, pg, Short)
		for _, p := range pg.Params {
			if p.HasBeenSet() {
				fmt.Fprint(w, "Set : ") // nolint: errcheck
			} else {
				fmt.Fprint(w, "--- : ") // nolint: errcheck
			}

			fmt.Fprint(w, p.Name()) // nolint: errcheck
			for _, altName := range p.AltNames() {
				if altName != p.Name() {
					fmt.Fprint(w, " or ", altName) // nolint: errcheck
				}
			}
			fmt.Fprintln(w) // nolint: errcheck

			intro := "          at: "
			whereSet := p.WhereSet()
			if len(whereSet) != 0 {
				for _, loc := range whereSet {
					fmt.Fprintln(w, intro, loc) // nolint: errcheck
					intro = "         and: "
				}
			}
		}
	}
}
