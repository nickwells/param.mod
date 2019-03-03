package psetter

import (
	"fmt"
	"strconv"

	"github.com/nickwells/param.mod/v2/param"
)

// Bool is used to set boolean flags
type Bool struct {
	param.ValueReqOptional

	Value *bool
}

// Set sets the parameter value to true
func (s Bool) Set(_ string) error {
	*s.Value = true
	return nil
}

// SetWithVal should be called when a value is given for the parameter
func (s Bool) SetWithVal(_, val string) error {
	b, err := strconv.ParseBool(val)
	if err != nil {
		return err
	}
	*s.Value = b
	return nil
}

// AllowedValues returns a description of the allowed values.
func (s Bool) AllowedValues() string {
	return "none (which will be taken as 'true')" +
		" or some value that can be interpreted as true or false"
}

// CurrentValue returns the current setting of the parameter value
func (s Bool) CurrentValue() string {
	return fmt.Sprintf("%v", *s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil
func (s Bool) CheckSetter(name string) {
	if s.Value == nil {
		panic(name + ": Bool Check failed: the Value to be set is nil")
	}
}
