package param

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nickwells/english.mod/english"
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
		byPosMiniHelp += fmt.Sprintf(
			" %d parameters should be: ",
			len(ps.byPos))

		sep := "<"

		for _, bp := range ps.byPos {
			byPosMiniHelp += sep + bp.name
			sep = ">, <"
		}

		byPosMiniHelp += ">"
	}

	if missingCount == 1 {
		ps.AddErr("", errors.New("a parameter is missing,"+
			" one more positional parameter is needed. "+
			byPosMiniHelp))
	} else {
		ps.AddErr("", fmt.Errorf(
			"some parameters are missing,"+
				" %d more positional parameters are needed. %s",
			missingCount, byPosMiniHelp))
	}
}

type parsingStatus int

const (
	parsingFinished parsingStatus = iota
	parsingIncomplete
)

func (ps *PSet) handleParamsByPos(loc *location.L, params []string,
) parsingStatus {
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
				ps.trailingParams = params[i+1:]
				return parsingFinished
			}
		}
	}

	return parsingIncomplete
}

func (ps *PSet) handleParamsByName(loc *location.L, params []string) {
	var i int
	for i = len(ps.byPos); i < len(params); i++ {
		pStr := params[i]

		loc.Incr()
		loc.SetContent(fmt.Sprintf("%q", pStr))

		if pStr == ps.terminalParam {
			ps.terminalParamSeen = true
			break
		}

		paramName, paramVal, hasParamVal := strings.Cut(pStr, "=")

		trimmedParam, err := ps.trimParam(paramName)
		if err != nil {
			ps.AddErr(trimmedParam, loc.Error(err.Error()))
			continue
		}

		paramParts := append([]string{}, trimmedParam)

		p, ok := ps.nameToParam[trimmedParam]
		if !ok {
			ps.recordUnexpectedParam(trimmedParam, loc)
			continue
		}

		if hasParamVal {
			paramParts = append(paramParts, paramVal)
		} else if p.setter.ValueReq() == Mandatory {
			if i < (len(params) - 1) {
				i++

				loc.Incr()

				paramParts = append(paramParts, params[i])

				loc.SetContent(
					fmt.Sprintf("%q %q", paramName, paramParts[1]))
			}
		}

		p.processParam(loc, paramParts)

		if ps.terminalParamSeen {
			break
		}
	}

	if i < len(params) {
		ps.trailingParams = params[i+1:]
	}
}

// getParamsFromStringSlice processes first the positional parameters, if
// any, and then the named parameters
func (ps *PSet) getParamsFromStringSlice(loc *location.L, params []string) {
	if ps.handleParamsByPos(loc, params) == parsingFinished {
		return
	}

	ps.handleParamsByName(loc, params)
}

// TrimPrefixesFromParam goes through the list of allowed parameter prefixes
// and tries removing them in turn from the passed parameter name. As soon as
// a prefix is successfully removed the resulting shortened parameter name is
// returned. If the parameter name does not start with any of the prefixes
// then the parameter name is returned unchanged.
func (ps *PSet) TrimPrefixesFromParam(pName string) string {
	for _, pfx := range ps.paramPrefixes {
		trimmedParam := strings.TrimPrefix(pName, pfx)
		if trimmedParam != pName {
			return trimmedParam
		}
	}

	return pName
}

// trimParam trims the parameter of its prefix (if any). It returns an error
// if the parameter does not start with any of the given prefixes.
func (ps *PSet) trimParam(pName string) (string, error) {
	if len(ps.paramPrefixes) == 0 {
		return pName, nil
	}

	trimmedParam := ps.TrimPrefixesFromParam(pName)
	if trimmedParam != pName {
		return trimmedParam, nil
	}

	return pName, fmt.Errorf("parameter %q does not start with %s",
		pName,
		english.JoinQuoted(ps.paramPrefixes, ", ", " or "))
}
