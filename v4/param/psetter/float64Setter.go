package psetter

import (
	"fmt"
	"strconv"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/v4/param"
)

// Float64 allows you to specify a parameter that can be used to set an
// float64 value. You can also supply a check function that will validate
// the Value. There are some helper functions given below (called
// Float64Check...) which will return functions that can perform a few
// common checks. For instance you can ensure that the value is positive by
// setting one of the Checks to the value returned by
// Float64CheckGT(0)
type Float64 struct {
	param.ValueReqMandatory

	Value  *float64
	Checks []check.Float64
}

// CountChecks returns the number of check functions this setter has
func (s Float64) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) checks that the value
// can be parsed to a float, if it cannot be parsed successfully it returns an
// error. If there is a check and the check is violated it returns an
// error. Only if the value is parsed successfully and the check is not
// violated is the Value set.
func (s Float64) SetWithVal(_ string, paramVal string) error {
	v, err := strconv.ParseFloat(paramVal, 64)
	if err != nil {
		return fmt.Errorf("could not interpret %q as a number: %s",
			paramVal, err)
	}

	for _, check := range s.Checks {
		if check == nil {
			continue
		}

		err := check(v)
		if err != nil {
			return err
		}
	}

	*s.Value = v
	return nil
}

// AllowedValues returns a string describing the allowed values
func (s Float64) AllowedValues() string {
	return "any value that can be read as a number with a decimal place" +
		HasChecks(s)
}

// CurrentValue returns the current setting of the parameter value
func (s Float64) CurrentValue() string {
	return fmt.Sprintf("%v", *s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s Float64) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, "psetter.Float64"))
	}
}
