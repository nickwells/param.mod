package psetter

import (
	"fmt"
	"strconv"

	"github.com/nickwells/check.mod/check"
)

// Float64 allows you to give a parameter that can be used to set a
// float64 value.
type Float64 struct {
	ValueReqMandatory

	// Value must be set, the program will panic if not. This is the value
	// being set
	Value *float64
	// The Checks, if any, are applied to the supplied parameter value and
	// the new parameter will be applied only if they all return a nil error
	Checks []check.Float64
}

// CountChecks returns the number of check functions this setter has
func (s Float64) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) checks that the
// value can be parsed to a float, if it cannot be parsed successfully it
// returns an error. The Checks, if any, are called with the value to be
// applied and if any of them return a non-nil error the Value is not updated
// and the error is returned. Only if the parameter value is parsed
// successfully and no checks fail is the Value set.
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
