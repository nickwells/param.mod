package psetter

import (
	"fmt"
	"regexp"
)

// Regexp allows you to give a parameter that can be used to set an
// regexp value.
type Regexp struct {
	ValueReqMandatory

	// You must set a Value, the program will panic if not. Note that this is
	// a pointer to a pointer, you should initialise it with the address of
	// the Regexp pointer.
	Value **regexp.Regexp
}

// SetWithVal (called when a value follows the parameter) checks that the value
// can be parsed to regular expression, if it cannot be parsed successfully it
// returns an error. Only if the value is parsed successfully is the Value set.
func (s Regexp) SetWithVal(_ string, paramVal string) error {
	v, err := regexp.Compile(paramVal)
	if err != nil {
		return fmt.Errorf("could not parse %q into a regular expression: %s",
			paramVal, err)
	}

	*s.Value = v
	return nil
}

// AllowedValues returns a string describing the allowed values
func (s Regexp) AllowedValues() string {
	return "any value that can be compiled as a regular expression" +
		" by the standard regexp package"
}

// CurrentValue returns the current setting of the parameter value
func (s Regexp) CurrentValue() string {
	if s.Value == nil {
		return "Illegal value"
	}
	if *s.Value == nil {
		return ""
	}
	return (*s.Value).String()
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s Regexp) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, fmt.Sprintf("%T", s)))
	}
}
