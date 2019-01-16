package psetter

import (
	"errors"
	"fmt"
	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/v2/param"
	"time"
)

// DurationSetter allows you to specify a parameter that can be used to set
// a time.Duration value. You can also supply a check function that will
// validate the Value. There are some helper functions given below (called
// DurationCheck...) which will return functions that can perform a few
// common checks. For instance you can ensure that the value is positive by
// setting one of the Checks to the value returned by
// DurationCheckGT(0)
type DurationSetter struct {
	Value  *time.Duration
	Checks []check.Duration
}

// ValueReq returns param.Mandatory indicating that some value must follow
// the parameter
func (s DurationSetter) ValueReq() param.ValueReq { return param.Mandatory }

// Set (called when there is no following value) returns an error
func (s DurationSetter) Set(_ string) error {
	return errors.New("no duration given (it should be followed by '=...')")
}

// SetWithVal (called when a value follows the parameter) checks that the
// value can be parsed to a duration, if it cannot be parsed successfully
// it returns an error. If there is a check and the check is violated it
// returns an error. Only if the value is parsed successfully and the check
// is not violated is the Value set.
func (s DurationSetter) SetWithVal(_ string, paramVal string) error {
	v, err := time.ParseDuration(paramVal)
	if err != nil {
		return fmt.Errorf("could not parse '%s' as a duration: %s",
			paramVal, err)
	}

	if len(s.Checks) != 0 {
		for _, check := range s.Checks {
			if check == nil {
				continue
			}

			err := check(v)
			if err != nil {
				return err
			}
		}
	}

	*s.Value = v
	return nil
}

// AllowedValues returns a string describing the allowed values
func (s DurationSetter) AllowedValues() string {
	rval := "any value that can be parsed to a duration"
	if len(s.Checks) != 0 {
		rval += " subject to checks"
	}
	return rval
}

// CurrentValue returns the current setting of the parameter value
func (s DurationSetter) CurrentValue() string {
	return fmt.Sprintf("%v", *s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil
func (s DurationSetter) CheckSetter(name string) {
	if s.Value == nil {
		panic(name +
			": DurationSetter Check failed: the Value to be set is nil")
	}
}
