package param

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nickwells/location.mod/location"
)

// reportMissingParams adds an error to the param set error map reporting any
// missing positional parameters.
func (ps *PSet) reportMissingParams(missingCount int) {
	if missingCount == 0 {
		return
	}

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
		ps.AddErr("", errors.New("A parameter is missing,"+
			" one more positional parameter is needed. "+
			byPosMiniHelp))
	} else {
		ps.AddErr("", fmt.Errorf(
			"Some parameters are missing,"+
				" %d more positional parameters are needed. %s",
			missingCount, byPosMiniHelp))
	}
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
			ps.AddErr(trimmedParam, loc.Error(err.Error()))
			continue
		}

		if p, ok := ps.nameToParam[trimmedParam]; ok {
			if p.setter.ValueReq() == Mandatory &&
				len(paramParts) == 1 {
				if i < (len(params) - 1) {
					i++
					loc.Incr()
					paramParts = append(paramParts, params[i])
					loc.SetContent(strings.Join(paramParts, " "))
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
	if ps.handleParamsByPos(loc, params) == parsingFinished {
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
