package psetter

import (
	"errors"
	"fmt"
	"github.com/nickwells/param.mod/param"
	"strings"
)

// EnumMapSetter sets the entry in a map of strings. The values must be in
// the allowed values map
type EnumMapSetter struct {
	Value       *map[string]bool
	AllowedVals AValMap // map[allowedValue] => description
	StrListSeparator
}

// ValueReq returns param.Mandatory indicating that some value must follow
// the parameter
func (s EnumMapSetter) ValueReq() param.ValueReq { return param.Mandatory }

// Set (called when there is no following value) returns an error
func (s EnumMapSetter) Set(_ string) error {
	return errors.New("no value given (it should be followed by '=...')")
}

// SetWithVal (called when a value follows the parameter) splits the value
// using the list separator. It then checks all the values for validity and
// only if all the values are in the allowed values list does it set the entry
// in the map of strings pointed to by the Value. It returns a error for the
// first invalid value.
func (s EnumMapSetter) SetWithVal(_ string, paramVal string) error {
	sep := s.GetSeparator()
	values := strings.Split(paramVal, sep)
	for _, v := range values {
		if _, ok := s.AllowedVals[v]; !ok {
			return errors.New("invalid value: '" + v)
		}
	}
	for _, v := range values {
		(*s.Value)[v] = true
	}
	return nil
}

// AllowedValues returns a string listing the allowed values
func (s EnumMapSetter) AllowedValues() string {
	return "a list of string values separated by '" + s.GetSeparator() +
		"'. The values must be from the following:\n" +
		allowedValues(s.AllowedVals)
}

// CurrentValue returns the current setting of the parameter value
func (s EnumMapSetter) CurrentValue() string {
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
func (s EnumMapSetter) CheckSetter(name string) {
	if s.Value == nil {
		panic(name +
			": EnumMapSetter Check failed: the Value to be set is nil")
	}
	if *s.Value == nil {
		panic(name +
			": EnumMapSetter Check failed: the map has not been created")
	}
	if len(s.AllowedVals) == 0 {
		panic(name +
			": EnumMapSetter Check failed: there are no allowed values")
	}
}
