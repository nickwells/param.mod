package psetter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nickwells/param.mod/v2/param"
)

// EnumMap sets the entry in a map of strings. The values initially set in
// the map must be in the allowed values map unless AllowHiddenMapEntries is
// set to true. Only values with keys in the allowed values map can be
// set. If you allow hidden values then you can have entries in your map
// which cannot be set through this interface but this will still only allow
// values to be set which are in the allowed values map.
//
// It is recommended that you should use string constants for setting and
// accessing the map entries and for initialising the allowed values map to
// avoid possible errors.
type EnumMap struct {
	param.ValueReqMandatory

	Value                 *map[string]bool
	AllowedVals           AValMap // map[allowedValue] => description
	AllowHiddenMapEntries bool
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
		parts := strings.Split(v, "=")
		if _, ok := s.AllowedVals[parts[0]]; !ok {
			return fmt.Errorf(
				"invalid value: '%s': part: %d (='%s') is not an allowed value",
				paramVal, i+1, v)
		}
		switch len(parts) {
		case 1:
			continue
		case 2:
			_, err := strconv.ParseBool(parts[1])
			if err != nil {
				return fmt.Errorf(
					"invalid value: '%s':"+
						" part: %d (='%s') is not an allowed value."+
						" The value after the '=' cannot be interpretted"+
						" as true or false: %s",
					paramVal, i+1, v, err)
			}
		default:
			return fmt.Errorf(
				"invalid value: '%s':"+
					" part: %d (='%s') is not an allowed value."+
					" There must be at most one '='",
				paramVal, i+1, v)
		}
	}
	for _, v := range values {
		parts := strings.Split(v, "=")
		switch len(parts) {
		case 1:
			(*s.Value)[v] = true
		case 2:
			b, _ := strconv.ParseBool(parts[1])
			(*s.Value)[parts[0]] = b
		}
	}
	return nil
}

// AllowedValues returns a string listing the allowed values
func (s EnumMap) AllowedValues() string {
	return s.ListValDesc("string values") +
		". The values must be from the following:\n" +
		s.AllowedVals.String() +
		"\n\nEach value can be set to false by following the value" +
		" with '=false'; by default the value will be set to true."
}

// CurrentValue returns the current setting of the parameter value
func (s EnumMap) CurrentValue() string {
	cv := ""
	sep := ""
	for k, v := range *s.Value {
		cv += sep + fmt.Sprintf("%s=%v", k, v)
		sep = "\n"
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
	if s.AllowHiddenMapEntries {
		return
	}
	for k := range *s.Value {
		if _, ok := s.AllowedVals[k]; !ok {
			panic(fmt.Sprintf("%sthe map entry with key '%s' is invalid"+
				" - it is not in the allowed values map",
				intro, k))
		}
	}
}
