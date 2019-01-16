package psetter

import (
	"errors"
	"fmt"
	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/v2/param"
	"strconv"
)

// Int64Setter allows you to specify a parameter that can be used to set an
// int64 value. You can also supply a check function that will validate the
// Value. See the check package for some helper functions which will return
// functions that can perform a few common checks. For instance you can
// ensure that the value is positive by setting one of the Checks to the
// value returned by check.Int64GT(0)
type Int64Setter struct {
	Value  *int64
	Checks []check.Int64
}

// ValueReq returns param.Mandatory indicating that some value must follow
// the parameter
func (s Int64Setter) ValueReq() param.ValueReq { return param.Mandatory }

// Set (called when there is no following value) returns an error
func (s Int64Setter) Set(_ string) error {
	return errors.New("no number given (it should be followed by '=num')")
}

// SetWithVal (called when a value follows the parameter) checks that the
// value can be parsed to an integer, if it cannot be parsed successfully it
// returns an error. If there are checks and any check is violated it returns
// an error. Only if the value is parsed successfully and no checks are
// violated is the Value set.
func (s Int64Setter) SetWithVal(_ string, paramVal string) error {
	v, err := strconv.ParseInt(paramVal, 0, 0)
	if err != nil {
		return fmt.Errorf("could not parse '%s' as an integer value: %s",
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
func (s Int64Setter) AllowedValues() string {
	rval := "any value that can be read as a whole number"
	if len(s.Checks) != 0 {
		rval += " subject to checks"
	}
	return rval
}

// CurrentValue returns the current setting of the parameter value
func (s Int64Setter) CurrentValue() string {
	return fmt.Sprintf("%v", *s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s Int64Setter) CheckSetter(name string) {
	if s.Value == nil {
		panic(name + ": Int64Setter Check failed: the Value to be set is nil")
	}
}
