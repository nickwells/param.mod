package psetter

import (
	"errors"
	"fmt"

	"github.com/nickwells/param.mod/v2/param"
)

// Enum allows you to specify a parameter that will only allow an
// enumerated range of values which are specified in the AllowedVals map
// which maps each allowed value to a description
type Enum struct {
	param.ValueReqMandatory

	Value       *string
	AllowedVals AValMap // map[allowedValue] => description
}

// SetWithVal (called when a value follows the parameter) checks the value
// for validity and only if it is in the allowed values list does it set the
// Value. It returns an error if the value is invalid.
func (s Enum) SetWithVal(_ string, paramVal string) error {
	if _, ok := s.AllowedVals[paramVal]; ok {
		*s.Value = paramVal
		return nil
	}
	return errors.New("invalid value: '" + paramVal + "'")
}

// AllowedValues returns a string listing the allowed values
func (s Enum) AllowedValues() string {
	return "one of\n" + s.AllowedVals.String()
}

// CurrentValue returns the current setting of the parameter value
func (s Enum) CurrentValue() string {
	return fmt.Sprintf("%v", *s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or there are no allowed values.
func (s Enum) CheckSetter(name string) {
	intro := name + ": Enum Check failed: "
	if s.Value == nil {
		panic(intro + "the Value to be set is nil")
	}
	if err := s.AllowedVals.OK(); err != nil {
		panic(intro + err.Error())
	}
	if _, ok := s.AllowedVals[*s.Value]; !ok {
		panic(fmt.Sprintf("%sthe initial value (%s) is not valid",
			intro, *s.Value))
	}
}
