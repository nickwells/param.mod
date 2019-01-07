package psetter

import (
	"errors"
	"fmt"
	"github.com/nickwells/param.mod/param"
)

// EnumSetter allows you to specify a parameter that will only allow an
// enumerated range of values which are specified in the AllowedVals map
// which maps each allowed value to a description
type EnumSetter struct {
	Value       *string
	AllowedVals AValMap // map[allowedValue] => description
}

// ValueReq returns param.Mandatory indicating that some value must follow
// the parameter
func (s EnumSetter) ValueReq() param.ValueReq { return param.Mandatory }

// Set (called when there is no following value) returns an error
func (s EnumSetter) Set(_ string) error {
	return errors.New("no value given (it should be followed by '=...')")
}

// SetWithVal (called when a value follows the parameter) checks the value
// for validity and only if it is in the allowed values list does it set the
// Value. It returns an error if the value is invalid.
func (s EnumSetter) SetWithVal(_ string, paramVal string) error {
	if _, ok := s.AllowedVals[paramVal]; ok {
		*s.Value = paramVal
		return nil
	}
	return errors.New("invalid value: '" + paramVal)
}

// AllowedValues returns a string listing the allowed values
func (s EnumSetter) AllowedValues() string {
	return "one of\n" + allowedValues(s.AllowedVals)
}

// CurrentValue returns the current setting of the parameter value
func (s EnumSetter) CurrentValue() string {
	return fmt.Sprintf("%v", *s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or there are no allowed values.
func (s EnumSetter) CheckSetter(name string) {
	if s.Value == nil {
		panic(name + ": EnumSetter Check failed: the Value to be set is nil")
	}
	if len(s.AllowedVals) == 0 {
		panic(name + ": EnumSetter Check failed: there are no allowed values")
	}
}
