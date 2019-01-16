package param

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/strdist.mod/strdist"
)

// Remainder returns any arguments that come after the terminal
// parameter. Note this may be a nil slice if all the parameters have been
// processed.
func (ps *ParamSet) Remainder() []string { return ps.remainingParams }

// findClosestMatch finds parameters with the name which is the shortest
// distance from the passed value and returns a string describing them
func (ps *ParamSet) findClosestMatch(badParam string) string {
	paramNames := make([]string, 0, len(ps.nameToParam))
	for p := range ps.nameToParam {
		paramNames = append(paramNames, p)
	}

	matches := strdist.CaseBlindCosineFinder.FindNStrLike(
		3, badParam, paramNames...)

	return strings.Join(matches, " or ")
}

func (ps *ParamSet) reportMissingParams(missingCount int) {
	var err error

	byPosMiniHelp := "The first"
	if len(ps.byPos) == 1 {
		byPosMiniHelp += " parameter should be: <" + ps.byPos[0].name + ">"
	} else {
		byPosMiniHelp +=
			fmt.Sprintf(" %d parameters should be: ", len(ps.byPos))
		sep := "<"
		for _, bp := range ps.byPos {
			byPosMiniHelp += sep + bp.name
			sep = ">, <"
		}
		byPosMiniHelp += ">"
	}

	if missingCount == 1 {
		err = errors.New("A parameter is missing," +
			" one more positional parameter is needed. " +
			byPosMiniHelp)
	} else {
		err = fmt.Errorf(
			"Some parameters are missing,"+
				" %d more positional parameters are needed. %s",
			missingCount, byPosMiniHelp)
	}
	ps.errors[""] = append(ps.errors[""], err)
}

func (ps *ParamSet) getParamsFromStringSlice(loc *location.L, params []string) {
	if len(ps.byPos) > 0 {
		missingCount := len(ps.byPos) - len(params)
		if missingCount > 0 {
			ps.reportMissingParams(missingCount)
			return
		}

		for i, pp := range ps.byPos {
			pStr := params[i]
			loc.Incr()
			loc.SetContent(pStr)

			pp.processParam(loc.Source(), loc, pStr)

			if pp.isTerminal {
				ps.remainingParams = params[i+1:]
				return
			}
		}
	}

	for i := len(ps.byPos); i < len(params); i++ {
		pStr := params[i]
		loc.Incr()
		loc.SetContent(pStr)

		if pStr == ps.terminalParam {
			ps.remainingParams = params[i+1:]
			return
		}

		paramParts := strings.SplitN(pStr, "=", 2)
		trimmedParam, err := trimParam(paramParts[0])
		if err != nil {
			ps.errors[trimmedParam] = append(ps.errors[trimmedParam],
				loc.Error(err.Error()))
			continue
		}

		if p, ok := ps.nameToParam[trimmedParam]; ok {
			if p.setter.ValueReq() == Mandatory &&
				len(paramParts) == 1 {
				if i < (len(params) - 1) {
					i++
					loc.Incr()
					paramParts = append(paramParts, params[i])
					loc.SetContent(pStr + " " + params[i])
				}
			}
			p.processParam(loc.Source(), loc, paramParts)
		} else {
			ps.recordUnexpectedParam(trimmedParam, loc)
		}
	}
}

// trimParam trims the parameter of any leading dashes
func trimParam(param string) (string, error) {
	trimmedParam := strings.TrimPrefix(param, "--")
	if trimmedParam != param {
		return trimmedParam, nil
	}
	trimmedParam = strings.TrimPrefix(param, "-")
	if trimmedParam != param {
		return trimmedParam, nil
	}
	return param, fmt.Errorf(
		"'%s' is a parameter but does not start with either '--' or '-'",
		param)
}
