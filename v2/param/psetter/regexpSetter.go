package psetter

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/nickwells/param.mod/v2/param"
)

// RegexpSetter allows you to specify a parameter that can be used to set an
// regexp value.
type RegexpSetter struct {
	Value **regexp.Regexp
}

// ValueReq returns param.Mandatory indicating that some value must follow
// the parameter
func (s RegexpSetter) ValueReq() param.ValueReq { return param.Mandatory }

// Set (called when there is no following value) returns an error
func (s RegexpSetter) Set(_ string) error {
	return errors.New("no pattern given (it should be followed by '=pattern')")
}

// SetWithVal (called when a value follows the parameter) checks that the value
// can be parsed to regular expression, if it cannot be parsed successfully it
// returns an error. Only if the value is parsed successfully is the Value set.
func (s RegexpSetter) SetWithVal(_ string, paramVal string) error {
	v, err := regexp.Compile(paramVal)
	if err != nil {
		return fmt.Errorf("could not parse '%s' to a regular expression: %s",
			paramVal, err)
	}

	*s.Value = v
	return nil
}

// AllowedValues returns a string describing the allowed values
func (s RegexpSetter) AllowedValues() string {
	return "any value that can be compiled as a regular expression" +
		" by the standard regexp package"
}

// CurrentValue returns the current setting of the parameter value
func (s RegexpSetter) CurrentValue() string {
	if s.Value == nil {
		return "Illegal value"
	}
	return (*s.Value).String()
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s RegexpSetter) CheckSetter(name string) {
	if s.Value == nil {
		panic(name + ": RegexpSetter Check failed: the Value to be set is nil")
	}
}
