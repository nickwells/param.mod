package psetter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/v2/param"
)

// EnumListSetter sets the values in a slice of strings. The values must be in
// the allowed values map
type EnumListSetter struct {
	Value       *[]string
	AllowedVals AValMap // map[allowedValue] => description
	StrListSeparator
	Checks []check.StringSlice
}

// CountChecks returns the number of check functions this setter has
func (s EnumListSetter) CountChecks() int {
	return len(s.Checks)
}

// ValueReq returns param.Mandatory indicating that some value must follow
// the parameter
func (s EnumListSetter) ValueReq() param.ValueReq { return param.Mandatory }

// Set (called when there is no following value) returns an error
func (s EnumListSetter) Set(_ string) error {
	return errors.New("no value given (it should be followed by '=...')")
}

// SetWithVal (called when a value follows the parameter) splits the value
// using the list separator. It then checks all the values for validity and
// only if all the values are in the allowed values list does it add them
// to the slice of strings pointed to by the Value. It returns a error for
// the first invalid value or if a check is breached.
func (s EnumListSetter) SetWithVal(_ string, paramVal string) error {
	sep := s.GetSeparator()
	values := strings.Split(paramVal, sep)
	for _, v := range values {
		if _, ok := s.AllowedVals[v]; !ok {
			return errors.New("invalid value: '" + v + "'")
		}
	}

	if len(s.Checks) != 0 {
		for _, check := range s.Checks {
			if check == nil {
				continue
			}

			err := check(values)
			if err != nil {
				return err
			}
		}
	}
	*s.Value = values

	return nil
}

// AllowedValues returns a string listing the allowed values
func (s EnumListSetter) AllowedValues() string {
	return s.ListValDesc("string values") +
		HasChecks(s) +
		". The values must be from the following:\n" +
		s.AllowedVals.String()
}

// CurrentValue returns the current setting of the parameter value
func (s EnumListSetter) CurrentValue() string {
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
func (s EnumListSetter) CheckSetter(name string) {
	intro := name + ": EnumListSetter Check failed: "
	if s.Value == nil {
		panic(intro + "the Value to be set is nil")
	}
	if err := s.AllowedVals.OK(); err != nil {
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
