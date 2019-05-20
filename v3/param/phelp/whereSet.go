package phelp

import (
	"fmt"
	"os"

	"github.com/nickwells/param.mod/v3/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// paramErrorCnt returns the number of errors that have been seen
func paramErrorCnt(ps *param.PSet, p *param.ByName) int {
	emap := ps.Errors()
	var errCount int

	for _, name := range p.AltNames() {
		errCount += len(emap[name])
	}

	return errCount
}

const (
	setIntro      = "Set    : "
	notSetIntro   = "---    : "
	manyErrsIntro = "Errs 9+: "
)

// printWhereSetIntro prints the introduction to the parameter name
// indicating whether or not it has been set and if there are any errors
func printWhereSetIntro(ps *param.PSet, p *param.ByName) {
	w := ps.StdWriter()

	errCount := paramErrorCnt(ps, p)

	if errCount > 0 {
		if errCount < 10 {
			fmt.Fprintf(w, "Errs %d : ", errCount)
		} else {
			fmt.Fprint(w, manyErrsIntro)
		}
	} else if p.HasBeenSet() {
		fmt.Fprint(w, setIntro)
	} else {
		fmt.Fprint(w, notSetIntro)
	}

}

func showWhereParamsAreSet(ps *param.PSet) {
	paramGroups := ps.GetGroups()
	w := ps.StdWriter()
	twc, err := twrap.NewTWConf(twrap.SetWriter(w))
	if err != nil {
		fmt.Fprint(os.Stderr, "Couldn't build the text wrapper:", err)
		return
	}

	for _, pg := range paramGroups {
		printGroupDetails(w, pg, Short)
		for _, p := range pg.Params {
			printWhereSetIntro(ps, p)

			fmt.Fprint(w, p.Name())
			for _, altName := range p.AltNames() {
				if altName != p.Name() {
					fmt.Fprint(w, " or ", altName)
				}
			}
			fmt.Fprintln(w)

			intro := "at : "
			whereSet := p.WhereSet()
			if len(whereSet) != 0 {
				for _, loc := range whereSet {
					twc.WrapPrefixed(intro, loc, len(notSetIntro)+4)
					intro = "and: "
				}
			}
		}
	}
}
