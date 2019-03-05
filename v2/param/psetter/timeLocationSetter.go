package psetter

import (
	"fmt"
	"strings"
	"time"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/v2/param"
)

// TimeLocation allows you to specify a parameter that can be used to
// set a time.Location pointer. You can also supply check functions that will
// validate the Value.
type TimeLocation struct {
	param.ValueReqMandatory

	Value  **time.Location
	Checks []check.TimeLocation
}

// CountChecks returns the number of check functions this setter has
func (s TimeLocation) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) checks that the
// value can be parsed to a location, if it cannot be parsed successfully it
// returns an error. If there is a check and the check is violated it returns
// an error. Only if the value is parsed successfully and the check is not
// violated is the Value set. If the supplied value cannot be successfully
// translated into a time.Location then any embedded spaces will be converted
// to underscores and the value will be retried. If this also fails then the
// original error is returned.
func (s TimeLocation) SetWithVal(_ string, paramVal string) error {
	v, err := time.LoadLocation(paramVal)
	if err != nil {
		convertedVal := strings.Replace(paramVal, " ", "_", -1)
		var e2 error
		v, e2 = time.LoadLocation(convertedVal)
		if e2 != nil {
			return fmt.Errorf("could not parse '%s' as a location: %s",
				paramVal, err)
		}
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
func (s TimeLocation) AllowedValues() string {
	return "any value that represents a location" +
		HasChecks(s) +
		". Typically this will be a string of the form" +
		" Continent/City_Name, for instance, Europe/London or America/New_York"
}

// CurrentValue returns the current setting of the parameter value
func (s TimeLocation) CurrentValue() string {
	return (*s.Value).String()
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s TimeLocation) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, "psetter.TimeLocation"))
	}
}
