package psetter

import (
	"github.com/nickwells/param.mod/v3/param"
)

// Nil is used if no value is to be set. It can be useful if the only
// effect is to be through the PostAction.
type Nil struct {
	param.ValueReqNone
}

// Set does nothing.
func (s Nil) Set(_ string) error {
	return nil
}

// AllowedValues returns a description of the allowed values.
func (s Nil) AllowedValues() string {
	return "none"
}

// CurrentValue returns the current setting of the parameter value, in this
// case there is never any current value.
func (s Nil) CurrentValue() string {
	return "none"
}

// CheckSetter does nothing.
func (s Nil) CheckSetter(name string) {
}
