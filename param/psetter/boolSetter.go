package psetter

import (
	"fmt"
	"strconv"

	"github.com/nickwells/param.mod/param"
)

// BoolSetter is used to set boolean flags
type BoolSetter struct {
	Value *bool
}

// ValueReq returns param.Optional indicating that the parameter may have
// but need not have a following value
func (s BoolSetter) ValueReq() param.ValueReq { return param.Optional }

// Set sets the parameter value to true
func (s BoolSetter) Set(_ string) error {
	*s.Value = true
	return nil
}

// SetWithVal should be called when a value is given for the parameter
func (s BoolSetter) SetWithVal(_, val string) error {
	b, err := strconv.ParseBool(val)
	if err != nil {
		return err
	}
	*s.Value = b
	return nil
}

// AllowedValues returns a description of the allowed values.
func (s BoolSetter) AllowedValues() string {
	return "none (which will be taken as 'true')" +
		" or some value that can be interpreted as true or false"
}

// CurrentValue returns the current setting of the parameter value
func (s BoolSetter) CurrentValue() string {
	return fmt.Sprintf("%v", *s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil
func (s BoolSetter) CheckSetter(name string) {
	if s.Value == nil {
		panic(name + ": BoolSetter Check failed: the Value to be set is nil")
	}
}
