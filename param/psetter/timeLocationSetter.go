package psetter

import (
	"fmt"
	"strings"
	"time"

	"github.com/nickwells/check.mod/check"
)

// TimeLocation allows you to give a parameter that can be used to set a
// time.Location pointer. You can also supply check functions that will
// validate the Value.
type TimeLocation struct {
	ValueReqMandatory

	// You must set a Value, the program will panic if not. Note that this is
	// a pointer to the pointer to the Location, you should initialise it
	// with the address of the Location pointer.
	Value **time.Location
	// The Checks, if any, are applied to the supplied parameter value and
	// the new parameter will be applied only if they all return a nil error.
	Checks []check.TimeLocation
}

// CountChecks returns the number of check functions this setter has
func (s TimeLocation) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) checks that the
// value can be parsed to a Location (note that it will try to replace any
// spaces with underscores if the first attempt fails). If it cannot be
// parsed successfully it returns an error. The Checks, if any, will be
// applied and if any of them return a non-nil error the Value will not be
// updated and the error will be returned.
func (s TimeLocation) SetWithVal(_ string, paramVal string) error {
	v, err := time.LoadLocation(paramVal)
	if err != nil {
		convertedVal := strings.ReplaceAll(paramVal, " ", "_")
		var e2 error
		v, e2 = time.LoadLocation(convertedVal)
		if e2 != nil {
			return fmt.Errorf("could not find %q as a time location: %s",
				paramVal, err)
		}
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
func (s TimeLocation) AllowedValues() string {
	return "any value that represents a location" +
		HasChecks(s) +
		". Typically this will be a string of the form" +
		" Continent/City_Name, for instance, Europe/London" +
		" or America/New_York." +
		" Additionally some of the three-letter timezone" +
		" names are also allowed such as UTC or CET."
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
