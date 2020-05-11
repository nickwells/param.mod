package psetter

import (
	"fmt"
	"strings"

	"github.com/nickwells/check.mod/check"
)

// EnumList sets the values in a slice of strings. The values must be in
// the allowed values map.
//
// It is recommended that you should use string constants for setting and
// accessing the map entries and for initialising the allowed values map to
// avoid possible errors.
type EnumList struct {
	ValueReqMandatory
	// The AllowedVals must be set, the program will panic if not. These are
	// the only values that will be allowed in the slice of strings.
	AllowedVals

	// Value must be set, the program will panic if not. This is the slice of
	// values that this setter is setting
	Value *[]string
	StrListSeparator
	// The Checks, if any, are applied to the list of new values and the
	// Value will only be updated if they all return a nil error
	Checks []check.StringSlice
}

// CountChecks returns the number of check functions this setter has
func (s EnumList) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) splits the value
// using the list separator. It then checks all the values for validity and
// only if all the values are in the allowed values list does it add them
// to the slice of strings pointed to by the Value. It returns a error for
// the first invalid value or if a check is breached.
func (s EnumList) SetWithVal(_ string, paramVal string) error {
	sep := s.GetSeparator()
	values := strings.Split(paramVal, sep)
	for _, v := range values {
		if !s.ValueAllowed(v) {
			return fmt.Errorf("value is not allowed: %q", v)
		}
	}

	for _, check := range s.Checks {
		if check == nil {
			continue
		}

		err := check(values)
		if err != nil {
			return err
		}
	}
	*s.Value = values

	return nil
}

// AllowedValues returns a string listing the allowed values
func (s EnumList) AllowedValues() string {
	return s.ListValDesc("string values") +
		HasChecks(s) + "."
}

// CurrentValue returns the current setting of the parameter value
func (s EnumList) CurrentValue() string {
	str := ""
	sep := ""

	for _, v := range *s.Value {
		str += sep + fmt.Sprintf("%v", v)
		sep = s.GetSeparator()
	}

	return str
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or there are no allowed values or the initial value is not
// allowed.
func (s EnumList) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, "psetter.EnumList"))
	}
	intro := name + ": psetter.EnumList Check failed: "
	if err := s.AllowedVals.Check(); err != nil {
		panic(intro + err.Error())
	}
	for i, v := range *s.Value {
		if _, ok := s.AllowedVals[v]; !ok {
			panic(fmt.Sprintf(
				"%selement %d (%s) in the current list of entries is invalid",
				intro, i, v))
		}
	}
}
