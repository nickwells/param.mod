package psetter

import (
	"errors"
	"github.com/nickwells/param.mod/v2/param"
)

// NilSetter is used if no value is to be set. It can be useful if the only
// effect is to be through the PostAction.
type NilSetter struct {
}

// ValueReq returns param.None indicating that the parameter must not have
// a following value.
func (s NilSetter) ValueReq() param.ValueReq { return param.None }

// Set does nothing.
func (s NilSetter) Set(_ string) error {
	return nil
}

// SetWithVal should never be called - the parameter should not be passed a
// value.
func (s NilSetter) SetWithVal(_, _ string) error {
	return errors.New("no value should follow the parameter")
}

// AllowedValues returns a description of the allowed values.
func (s NilSetter) AllowedValues() string {
	return "none"
}

// CurrentValue returns the current setting of the parameter value, in this
// case there is never any current value.
func (s NilSetter) CurrentValue() string {
	return "none"
}

// CheckSetter does nothing.
func (s NilSetter) CheckSetter(name string) {
}
