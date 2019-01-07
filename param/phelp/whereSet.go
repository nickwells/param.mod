package phelp

import (
	"fmt"
	"github.com/nickwells/param.mod/param"
)

func showWhereParamsAreSet(ps *param.ParamSet) {
	paramGroups := ps.GetParamGroups()
	w := ps.StdWriter()

	for _, pg := range paramGroups {
		printGroupDetails(w, pg, Short)
		for _, p := range pg.Params {
			if p.HasBeenSet() {
				fmt.Fprint(w, "Set : ")
			} else {
				fmt.Fprint(w, "--- : ")
			}

			fmt.Fprint(w, p.Name())
			for _, altName := range p.AltNames() {
				if altName != p.Name() {
					fmt.Fprint(w, " or ", altName)
				}
			}
			fmt.Fprintln(w)

			intro := "          at: "
			whereSet := p.WhereSet()
			if len(whereSet) != 0 {
				for _, loc := range whereSet {
					fmt.Fprintln(w, intro, loc)
					intro = "         and: "
				}
			}
		}
	}
}
