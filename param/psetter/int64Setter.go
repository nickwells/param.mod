package psetter

import (
	"fmt"
	"strconv"

	"github.com/nickwells/check.mod/check"
)

// Int64 allows you to give a parameter that can be used to set an
// int64 value.
type Int64 struct {
	ValueReqMandatory

	// You must set a Value, the program will panic if not. This is a pointer
	// to the int64 value that the setter is setting.
	Value *int64
	// The Checks, if any, are applied to the supplied parameter value and
	// the Value will only be update if they all return a nil error.
	Checks []check.Int64
}

// CountChecks returns the number of check functions this setter has
func (s Int64) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) checks that the
// value can be parsed to an integer, if it cannot be parsed successfully it
// returns an error. If there are checks and any check is violated it returns
// an error. Only if the value is parsed successfully and no checks are
// violated is the Value set.
func (s Int64) SetWithVal(_ string, paramVal string) error {
	v, err := strconv.ParseInt(paramVal, 0, 0)
	if err != nil {
		return fmt.Errorf("could not interpret %q as a whole number: %s",
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
func (s Int64) AllowedValues() string {
	return "any value that can be read as a whole number" + HasChecks(s)
}

// CurrentValue returns the current setting of the parameter value
func (s Int64) CurrentValue() string {
	return fmt.Sprintf("%v", *s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s Int64) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, "psetter.Int64"))
	}
}
