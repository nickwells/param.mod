package phelp

import (
	"github.com/nickwells/param.mod/v4/param"
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
func printWhereSetIntro(twc *twrap.TWConf, ps *param.PSet, p *param.ByName) {
	errCount := paramErrorCnt(ps, p)

	if errCount > 0 {
		if errCount < 10 {
			twc.Printf("Errs %d : ", errCount)
		} else {
			twc.Print(manyErrsIntro)
		}
	} else if p.HasBeenSet() {
		twc.Print(setIntro)
	} else {
		twc.Print(notSetIntro)
	}

}

func (h StdHelp) showWhereParamsAreSet(twc *twrap.TWConf, ps *param.PSet) {
	twc.Wrap("Parameter Summary\n\n"+
		"This shows a summary of all the parameters."+
		" If there are any errors with a parameter then that will be"+
		" indicated along with a count of the number of errors. If a"+
		" parameter has been set then that will be indicated along"+
		" with details of where it has been set.\n", 0)

	groups := ps.GetGroups()

	printSep := false
	for _, g := range groups {
		if printSep {
			twc.Println() //nolint: errcheck
			printMinorSeparator(twc)
		}
		printSep = true
		h.printGroupDetails(twc, g)
		for _, p := range g.Params {
			printWhereSetIntro(twc, ps, p)

			twc.Print(p.Name())
			for _, altName := range p.AltNames() {
				if altName != p.Name() {
					twc.Print(" or ", altName)
				}
			}
			twc.Println() //nolint: errcheck

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
