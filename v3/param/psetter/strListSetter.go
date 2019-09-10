package psetter

import (
	"strings"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/v3/param"
)

// StrList allows you to specify a parameter that can be used to set an
// list (a slice) of strings. You can override the list separator by setting
// the Sep value.
//
// If you have a list of allowed values you should use EnumList
type StrList struct {
	param.ValueReqMandatory
	param.NilAVM

	Value *[]string
	StrListSeparator
	Checks []check.StringSlice
}

// CountChecks returns the number of check functions this setter has
func (s StrList) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) splits the value
// into a slice of strings and sets the Value accordingly. It will return
// an error if a check is breached.
func (s StrList) SetWithVal(_ string, paramVal string) error {
	sep := s.GetSeparator()
	v := strings.Split(paramVal, sep)

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

// AllowedValues returns a description of the allowed values. It includes the
// separator to be used
func (s StrList) AllowedValues() string {
	return s.ListValDesc("string values") + HasChecks(s)
}

// CurrentValue returns the current setting of the parameter value
func (s StrList) CurrentValue() string {
	cv := ""
	sep := ""

	for _, v := range *s.Value {
		cv += sep + v
		sep = s.GetSeparator()
	}

	return cv
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s StrList) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, "psetter.StrList"))
	}
}
