package psetter

import (
	"fmt"
	"strings"
	"time"

	"github.com/nickwells/check.mod/v2/check"
)

// Duration allows you to specify a parameter that can be used to set a
// time.Duration value. You can also supply check functions that will
// validate the Value. See the check package for some common pre-defined
// checks.
type Duration struct {
	ValueReqMandatory

	// Value must be set, the program will panic if not. This is the
	// time.Duration that the setter is setting.
	Value *time.Duration
	// The Checks, if any, are applied to the new Duration and the Value will
	// only be updated if they all return a nil error.
	Checks []check.Duration
}

// CountChecks returns the number of check functions this setter has
func (s Duration) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) checks that the
// value can be parsed to a duration, if it cannot be parsed successfully it
// returns an error. If there is a check and the check is violated it returns
// an error. Only if the value is parsed successfully and the checks are not
// violated is the Value set.
func (s Duration) SetWithVal(_ string, paramVal string) error {
	v, err := time.ParseDuration(paramVal)
	if err != nil {
		return fmt.Errorf("could not parse %q as a duration: %s",
			paramVal, err)
	}

	for _, check := range s.Checks {
		err := check(v)
		if err != nil {
			return err
		}
	}

	*s.Value = v

	return nil
}

// AllowedValues returns a string describing the allowed values.
func (s Duration) AllowedValues() string {
	unitStrings := []string{"ns", "us", "µs", "ms", "s", "m", "h"}
	aval := "any value that can be parsed as a duration.\n" +
		"A duration string is a sequence of numbers with an optional" +
		" fraction and a unit. The allowed unit names are "
	aval += strings.Join(unitStrings, ", ")
	aval += ". The whole sequence can be signed and must not contain" +
		" any spaces.\n\n" +
		"For example: '300ms', '-1.5h' or '2h45m'" + HasChecks(s) + "\n"

	return aval
}

// CurrentValue returns the current setting of the parameter value.
func (s Duration) CurrentValue() string {
	return fmt.Sprintf("%v", *s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or if it has nil Checks.
func (s Duration) CheckSetter(name string) {
	// Check the value is not nil
	if s.Value == nil {
		panic(NilValueMessage(name, fmt.Sprintf("%T", s)))
	}

	// Check there are no nil Check funcs
	for i, check := range s.Checks {
		if check == nil {
			panic(NilCheckMessage(name, fmt.Sprintf("%T", s), i))
		}
	}
}
