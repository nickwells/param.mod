package psetter

import (
	"fmt"
)

// Enum allows you to give a parameter that will only allow one of an
// enumerated range of values which are specified in the AllowedVals map.
type Enum struct {
	ValueReqMandatory
	// The AllowedVals must be set, the program will panic if not. The Value
	// is guaranteed to take one of these values.
	AllowedVals

	// Value must be set, the program will panic if not. This is the value
	// being set
	Value *string
}

// SetWithVal (called when a value follows the parameter) checks the value
// for validity and only if it is in the allowed values list does it set the
// Value. It returns an error if the value is invalid.
func (s Enum) SetWithVal(_ string, paramVal string) error {
	if s.ValueAllowed(paramVal) {
		*s.Value = paramVal
		return nil
	}
	return fmt.Errorf("value not allowed: %q", paramVal)
}

// AllowedValues returns a string listing the allowed values
func (s Enum) AllowedValues() string {
	return "a string"
}

// CurrentValue returns the current setting of the parameter value
func (s Enum) CurrentValue() string {
	return fmt.Sprintf("%v", *s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or there are no allowed values.
func (s Enum) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, "psetter.Enum"))
	}
	intro := name + ": psetter.Enum Check failed: "
	if err := s.ValueMapOK(); err != nil {
		panic(intro + err.Error())
	}
	if !s.ValueAllowed(*s.Value) {
		panic(fmt.Sprintf("%sthe initial value (%s) is not valid",
			intro, *s.Value))
	}
}
