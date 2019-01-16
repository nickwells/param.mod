package psetter

import (
	"errors"
	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/v2/param"
	"strings"
)

// StrListSetter allows you to specify a parameter that can be used to set an
// list (a slice) of strings. You can override the list separator by setting
// the Sep value.
//
// If you have a list of allowed values you should use EnumListSetter
type StrListSetter struct {
	Value *[]string
	StrListSeparator
	Checks []check.StringSlice
}

// ValueReq returns param.Mandatory indicating that some value must follow
// the parameter
func (s StrListSetter) ValueReq() param.ValueReq { return param.Mandatory }

// Set (called when there is no following value) returns an error
func (s StrListSetter) Set(_ string) error {
	return errors.New("no value given (it should be followed by '=...')")
}

// SetWithVal (called when a value follows the parameter) splits the value
// into a slice of strings and sets the Value accordingly. It will return
// an error if a check is breached.
func (s StrListSetter) SetWithVal(_ string, paramVal string) error {
	sep := s.GetSeparator()
	v := strings.Split(paramVal, sep)

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

// AllowedValues returns a description of the allowed values. It includes the
// separator to be used
func (s StrListSetter) AllowedValues() string {
	rval := "a list of string values separated by '" +
		s.GetSeparator() + "'"

	if len(s.Checks) != 0 {
		rval += " subject to checks"
	}
	return rval
}

// CurrentValue returns the current setting of the parameter value
func (s StrListSetter) CurrentValue() string {
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
func (s StrListSetter) CheckSetter(name string) {
	if s.Value == nil {
		panic(name + ": StrListSetter Check failed: the Value to be set is nil")
	}
}
