package psetter

import (
	"fmt"
	"regexp"

	"github.com/nickwells/param.mod/v2/param"
)

// Regexp allows you to specify a parameter that can be used to set an
// regexp value.
type Regexp struct {
	param.ValueReqMandatory

	Value **regexp.Regexp
}

// SetWithVal (called when a value follows the parameter) checks that the value
// can be parsed to regular expression, if it cannot be parsed successfully it
// returns an error. Only if the value is parsed successfully is the Value set.
func (s Regexp) SetWithVal(_ string, paramVal string) error {
	v, err := regexp.Compile(paramVal)
	if err != nil {
		return fmt.Errorf("could not parse '%s' to a regular expression: %s",
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
	return (*s.Value).String()
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s Regexp) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, "psetter.Regexp"))
	}
}
