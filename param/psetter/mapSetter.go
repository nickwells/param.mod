package psetter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nickwells/check.mod/check"
)

// Map sets the entry in a map of strings. Each value from the
// parameter is used as a key in the map with the map entry set to true.
type Map struct {
	ValueReqMandatory

	// You must set a Value, the program will panic if not. This is the map
	// of strings to bool that the setter is setting
	Value *map[string]bool
	// The Checks, if any, are applied to the supplied parameter value and
	// the new parameter will be applied only if they all return a nil error
	Checks []check.MapStringBool
	StrListSeparator
}

// CountChecks returns the number of check functions this setter has
func (s Map) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) splits the value
// using the list separator.
func (s Map) SetWithVal(_ string, paramVal string) error {
	sep := s.GetSeparator()
	values := strings.Split(paramVal, sep)
	m := map[string]bool{}

	for i, v := range values {
		parts := strings.SplitN(v, "=", 2)
		switch len(parts) {
		case 1:
			m[parts[0]] = true
		case 2:
			// check that the bool can be parsed
			b, err := strconv.ParseBool(parts[1])
			if err != nil {
				return fmt.Errorf("bad value: %q: part: %d (%q) is invalid."+
					" The value (%q) cannot be interpreted"+
					" as true or false: %s",
					paramVal, i+1, v, parts[1], err)
			}
			m[parts[0]] = b
		}
	}

	for _, check := range s.Checks {
		if check == nil {
			continue
		}

		err := check(m)
		if err != nil {
			return err
		}
	}

	for k, b := range m {
		(*s.Value)[k] = b
	}
	return nil
}

// AllowedValues returns a string listing the allowed values
func (s Map) AllowedValues() string {
	return "a list of string values separated by '" +
		s.GetSeparator() + "'" + HasChecks(s)
}

// CurrentValue returns the current setting of the parameter value
func (s Map) CurrentValue() string {
	cv := ""
	sep := ""
	for k, v := range *s.Value {
		cv += sep + fmt.Sprintf("%s=%v", k, v)
		sep = s.GetSeparator()
	}

	return cv
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or the map has not been created yet.
func (s Map) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, "psetter.Map"))
	}
	if *s.Value == nil {
		*s.Value = make(map[string]bool)
	}
}
