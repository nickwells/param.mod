package psetter

import (
	"fmt"
	"strings"

	"github.com/nickwells/param.mod/v2/param"
)

// EnumMap sets the entry in a map of strings. The values must be in
// the allowed values map
type EnumMap struct {
	param.ValueReqMandatory

	Value       *map[string]bool
	AllowedVals AValMap // map[allowedValue] => description
	StrListSeparator
}

// SetWithVal (called when a value follows the parameter) splits the value
// using the list separator. It then checks all the values for validity and
// only if all the values are in the allowed values list does it set the entry
// in the map of strings pointed to by the Value. It returns a error for the
// first invalid value.
func (s EnumMap) SetWithVal(_ string, paramVal string) error {
	sep := s.GetSeparator()
	values := strings.Split(paramVal, sep)
	for i, v := range values {
		if _, ok := s.AllowedVals[v]; !ok {
			return fmt.Errorf(
				"invalid value: '%s': part: %d (='%s') is not an allowed value",
				paramVal, i+1, v)
		}
	}
	for _, v := range values {
		(*s.Value)[v] = true
	}
	return nil
}

// AllowedValues returns a string listing the allowed values
func (s EnumMap) AllowedValues() string {
	return s.ListValDesc("string values") +
		". The values must be from the following:\n" +
		s.AllowedVals.String()
}

// CurrentValue returns the current setting of the parameter value
func (s EnumMap) CurrentValue() string {
	cv := ""
	sep := ""
	for k, v := range *s.Value {
		cv += sep + fmt.Sprintf("%s=%v", k, v)
		sep = s.GetSeparator()
	}

	return cv
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or the map has not been created yet or if there are no
// allowed values.
func (s EnumMap) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, "psetter.EnumMap"))
	}
	intro := name + ": psetter.EnumMap Check failed: "
	if *s.Value == nil {
		panic(intro + "the map has not been created")
	}
	if err := s.AllowedVals.OK(); err != nil {
		panic(intro + err.Error())
	}
	for k := range *s.Value {
		if _, ok := s.AllowedVals[k]; !ok {
			panic(fmt.Sprintf("%sthe map entry (%s) is invalid", intro, k))
		}
	}
}
