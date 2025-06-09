package phelp

import (
	"github.com/nickwells/col.mod/v5/col"
	"github.com/nickwells/col.mod/v5/colfmt"
	"github.com/nickwells/english.mod/english"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/twrap.mod/twrap"
)

const (
	paramSetFmtStd   = "std"
	paramSetFmtShort = "short"
	paramSetFmtTable = "table"
)

const (
	setIntro      = "Set    : "
	notSetIntro   = "---    : "
	manyErrsIntro = "Errs 9+: "
)

const (
	atIndent = 4
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

// printWhereSetIntro prints the introduction to the parameter name
// indicating whether or not it has been set and if there are any errors
func printWhereSetIntro(twc *twrap.TWConf, p *param.ByName, errCount int) {
	const tooManyErrs = 10
	if errCount > 0 {
		if errCount < tooManyErrs {
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

func showWhereParamsAreSet(h StdHelp, twc *twrap.TWConf, ps *param.PSet) int {
	switch h.paramsSetFormat {
	case paramSetFmtStd:
		showWhereSetStd(h, twc, ps)
	case paramSetFmtShort:
		showWhereSetShort(h, twc, ps)
	case paramSetFmtTable:
		showWhereSetTable(h, twc, ps)
	default:
		panic("the Format of the report on where params are set is unknown")
	}

	return 0
}

func showWhereSetStd(h StdHelp, twc *twrap.TWConf, ps *param.PSet) {
	twc.Wrap("Parameter Summary\n\n"+
		"This shows a summary of all the parameters."+
		" If there are any errors with a parameter then that will be"+
		" indicated along with a count of the number of errors. If a"+
		" parameter has been set then that will be indicated along"+
		" with details of where it has been set.\n", 0)

	groups := ps.GetGroups()

	maxNameLen := getMaxGroupNameLen(groups)
	printSep := false

	for _, g := range groups {
		if printSep {
			twc.Print("\n")
			twc.Print(minorSectionSeparator)
		}

		printSep = true

		h.printGroup(twc, g, maxNameLen)

		for _, p := range g.Params() {
			printWhereSetIntro(twc, p, paramErrorCnt(ps, p))
			twc.Println(english.Join(p.AltNames(), ", ", " or "))

			intro := "at : "

			whereSet := p.WhereSet()
			for _, loc := range whereSet {
				twc.WrapPrefixed(intro, loc, len(notSetIntro)+atIndent)
				intro = "and: "
			}
		}
	}
}

func showWhereSetShort(_ StdHelp, twc *twrap.TWConf, ps *param.PSet) {
	groups := ps.GetGroups()

	for _, g := range groups {
		for _, p := range g.Params() {
			errCount := paramErrorCnt(ps, p)

			if errCount == 0 && !p.HasBeenSet() {
				continue
			}

			printWhereSetIntro(twc, p, errCount)
			twc.Println(english.Join(p.AltNames(), ", ", " or "))

			intro := "at : "

			whereSet := p.WhereSet()
			for _, loc := range whereSet {
				twc.WrapPrefixed(intro, loc, len(notSetIntro)+atIndent)
				intro = "and: "
			}
		}
	}
}

// skipWhereSetReport returns whether or not to report where the parameter has
// been set or if it has any errors associated. It returns a bool and the
// count of any errors found.
func skipWhereSetReport(p *param.ByName) (bool, int) {
	errCount := paramErrorCnt(p.PSet(), p)

	return errCount == 0 && !p.HasBeenSet(), errCount
}

// calcColumnWidths calculates the maximum parameter group name and parameter
// name lengths. It only gives the values for those which have been set or
// where an error was detected.
func calcColumnWidths(groups []*param.Group) (uint, uint) {
	maxGNLen, maxPNLen := 0, 0

	groupUsed := false

	for _, g := range groups {
		for _, p := range g.Params() {
			skip, _ := skipWhereSetReport(p)
			if skip {
				continue
			}

			groupUsed = true

			for _, pName := range p.AltNames() {
				if len(pName) > maxPNLen {
					maxPNLen = len(pName)
				}
			}
		}

		if groupUsed {
			if len(g.Name()) > maxGNLen {
				maxGNLen = len(g.Name())
			}
		}

		groupUsed = false
	}

	return uint(maxGNLen), uint(maxPNLen) //nolint:gosec
}

func showWhereSetTable(_ StdHelp, twc *twrap.TWConf, ps *param.PSet) {
	hdr, err := col.NewHeader()
	if err != nil {
		twc.Println("Cannot construct header for where-params-set table:", err)
		return
	}

	groups := ps.GetGroups()
	maxGroupNameLen, maxParamNameLen := calcColumnWidths(groups)

	rpt := col.NewReportOrPanic(hdr, twc.W,
		col.New(&colfmt.Int{HandleZeroes: true}, "errs"),
		col.New(&colfmt.String{W: maxGroupNameLen}, "parameter", "group"),
		col.New(&colfmt.WrappedString{W: maxParamNameLen},
			"parameter", "name"),
		col.New(&colfmt.String{}, "set at"),
	)

	for _, g := range groups {
		for _, p := range g.Params() {
			skip, errCount := skipWhereSetReport(p)

			if skip {
				continue
			}

			at := "-"
			sep := ""

			whereSet := p.WhereSet()
			if len(whereSet) != 0 {
				at = ""

				for _, loc := range whereSet {
					at += sep + loc
					sep = "\n"
				}
			}

			_ = rpt.PrintRow(errCount,
				g.Name(), english.Join(p.AltNames(), ", ", " or "),
				at)
		}
	}
}
