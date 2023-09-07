package psetter

import (
	"fmt"
	"strings"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v6/pedit"
)

// StrList allows you to specify a parameter that can be used to set a list
// (a slice) of strings.
//
// If only certain, predefined, values are allowed you might prefer to use
// EnumList
type StrList struct {
	ValueReqMandatory

	// You must set a Value, the program will panic if not. This is the
	// string that the setter is setting.
	Value *[]string
	StrListSeparator
	// The Checks, if any, are applied to the supplied parameter value and
	// the new parameter will be applied only if they all return a nil error.
	Checks []check.StringSlice
	// The Editor, if present, is applied to each of the listed parameter
	// values after any checks are applied and allows the programmer to
	// modify the value supplied before using it to set the Value.
	Editor pedit.Editor
}

// CountChecks returns the number of check functions this setter has
func (s StrList) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) splits the value
// into a slice of strings and sets the Value accordingly. The Checks, if
// any, will be applied and if any of them return an error the Value will not
// be updated and the error will be returned. If the Editor is non-nil then
// it is applied to each of the listed parameter values. If the Editor
// returns an error for any of the listed values then the Value will not be
// updated and the error is returned.
func (s StrList) SetWithVal(paramName, paramVal string) error {
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

	if s.Editor != nil {
		for i, sv := range v {
			newVal, err := s.Editor.Edit(paramName, sv)
			if err != nil {
				return err
			}
			v[i] = newVal
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
		panic(NilValueMessage(name, fmt.Sprintf("%T", s)))
	}
}
