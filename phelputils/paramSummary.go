package phelputils

import (
	"fmt"
	"strings"

	"github.com/nickwells/param.mod/v7/param"
	"github.com/nickwells/param.mod/v7/psetter"
)

// valTypeName returns a descriptive string for the type of the Setter
func valTypeName(p param.ByName) string {
	if name := p.ValueName(); name != "" {
		return name
	}

	s := p.Setter()
	if sVD, ok := s.(psetter.ValDescriber); ok {
		return sVD.ValDescribe()
	}

	valType := fmt.Sprintf("%T", s)

	parts := strings.Split(valType, ".")
	valType = parts[len(parts)-1]
	valType = strings.TrimRight(valType, "0123456789")

	return valType
}

// valueNeededStr returns a descriptive string indicating whether a trailing
// argument is needed and if so of what type it should be.
func valueNeededStr(p param.ByName) string {
	switch p.Setter().ValueReq() {
	case param.Mandatory:
		return "=" + valTypeName(p)
	case param.Optional:
		return "[=" + valTypeName(p) + "] "
	default:
		return ""
	}
}

// optionalWrapper wraps the string in "[" and "]" if the parameter is optional.
func optionalWrapper(s string, p param.ByName) string {
	if p.AttrIsSet(param.MustBeSet) {
		return s
	}

	return "[" + s + "]"
}

// ParamSummary returns a string summarising the usage of the parameter. The
// parameter name and all it's alternatives will be given and an indication
// of whether a following value is required. The value returned will be
// bracketted by '[' and ']' if it is not mandatory.
func ParamSummary(p param.ByName) string {
	var s strings.Builder

	sep := ""
	for _, altName := range p.AltNames() {
		s.WriteString(sep)
		sep = ", "

		s.WriteString(p.PSet().ShortestPrefix())
		s.WriteString(altName)
		s.WriteString(valueNeededStr(p))
	}

	return optionalWrapper(s.String(), p)
}

// ParamShortSummary returns a string summarising the usage of the
// parameter. Only the parameter name is shown (none of the alternatives will
// be given) and an indication of whether a following value is required. The
// value returned will be bracketted by '[' and ']' if it is not mandatory.
func ParamShortSummary(p param.ByName) string {
	var s strings.Builder

	s.WriteString(p.PSet().ShortestPrefix())
	s.WriteString(p.Name())
	s.WriteString(valueNeededStr(p))

	return optionalWrapper(s.String(), p)
}
