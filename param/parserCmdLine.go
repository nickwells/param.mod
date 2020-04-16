package param

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/strdist.mod/strdist"
)

// findClosestMatch finds parameters with the name which is the shortest
// distance from the passed value and returns a string describing them
func (ps *PSet) findClosestMatch(badParam string) string {
	paramNames := make([]string, 0, len(ps.nameToParam))
	for p := range ps.nameToParam {
		paramNames = append(paramNames, p)
	}

	matches := strdist.CaseBlindCosineFinder.FindNStrLike(
		3, badParam, paramNames...)

	return strings.Join(matches, " or ")
}

func (ps *PSet) reportMissingParams(missingCount int) {
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

type parsingStatus int

const (
	parsingFinished parsingStatus = iota
	parsingIncomplete
)

func (ps *PSet) handleParamsByPos(loc *location.L, params []string) parsingStatus {
	if len(ps.byPos) > 0 {
		missingCount := len(ps.byPos) - len(params)
		if missingCount > 0 {
			ps.reportMissingParams(missingCount)
			return parsingFinished
		}

		for i, pp := range ps.byPos {
			pStr := params[i]
			loc.Incr()
			loc.SetContent(pStr)

			pp.processParam(loc, pStr)

			if pp.isTerminal {
				ps.remainingParams = params[i+1:]
				return parsingFinished
			}
		}
	}
	return parsingIncomplete
}

func (ps *PSet) handleParamsByName(loc *location.L, params []string) {
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
			paramParts[0] = trimmedParam
			p.processParam(loc, paramParts)
		} else {
			ps.recordUnexpectedParam(trimmedParam, loc)
		}
	}
}

func (ps *PSet) getParamsFromStringSlice(loc *location.L, params []string) {
	status := ps.handleParamsByPos(loc, params)
	if status == parsingFinished {
		return
	}

	ps.handleParamsByName(loc, params)
}

// trimParam trims the parameter of any leading dashes
func trimParam(param string) (string, error) {
	trimmedParam := strings.TrimPrefix(param, "--")
	if trimmedParam != param {
		return trimmedParam, nil
	}
	trimmedParam = strings.TrimPrefix(param, "-")
	if trimmedParam != param {
		if trimmedParam == "" {
			return param, errors.New("the parameter name is blank")
		}
		return trimmedParam, nil
	}
	return param, fmt.Errorf(
		"parameter %q does not start with either '--' or '-'", param)
}
