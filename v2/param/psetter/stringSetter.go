package psetter

import (
	"errors"
	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/v2/param"
)

// StringSetter is the type for setting string values from
// parameters
type StringSetter struct {
	Value  *string
	Checks []check.String
}

// ValueReq returns param.Mandatory indicating that some value must follow
// the parameter
func (s StringSetter) ValueReq() param.ValueReq { return param.Mandatory }

// Set (called when there is no following value) returns an error
func (s StringSetter) Set(_ string) error {
	return errors.New("no value given (it should be followed by '=...')")
}

// SetWithVal checks that the parameter value meets the checks if any. It
// returns an error if the check is not satisfied. Only if the check
// is not violated is the Value set.
func (s StringSetter) SetWithVal(_ string, paramVal string) error {

	if len(s.Checks) != 0 {
		for _, check := range s.Checks {
			if check == nil {
				continue
			}

			err := check(paramVal)
			if err != nil {
				return err
			}
		}
	}

	*s.Value = paramVal
	return nil
}

// AllowedValues simply returns "any string" since StringSetter
// does not check its value
func (s StringSetter) AllowedValues() string {
	rval := "any string"

	if len(s.Checks) != 0 {
		rval += " subject to checks"
	}
	return rval
}

// CurrentValue returns the current setting of the parameter value
func (s StringSetter) CurrentValue() string {
	return *s.Value
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s StringSetter) CheckSetter(name string) {
	if s.Value == nil {
		panic(name + ": StringSetter Check failed: the Value to be set is nil")
	}
}
