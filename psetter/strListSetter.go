package psetter

import (
	"fmt"
	"strings"

	"github.com/nickwells/check.mod/v2/check"
)

// StrList allows you to specify a parameter that can be used to set a list
// (a slice) of strings.
//
// If only certain, predefined, values are allowed you might prefer to use
// EnumList
type StrList[T ~string] struct {
	ValueReqMandatory

	// You must set a Value, the program will panic if not. This is the
	// string that the setter is setting.
	Value *[]T
	StrListSeparator
	// The Checks, if any, are applied to the supplied parameter value and
	// the new parameter will be applied only if they all return a nil error.
	Checks []check.ValCk[[]T]
	// The Editor, if present, is applied to each of the listed parameter
	// values after any checks are applied and allows the programmer to
	// modify the value supplied before using it to set the Value.
	Editor Editor
}

// CountChecks returns the number of check functions this setter has
func (s StrList[T]) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) splits the value
// into a slice of strings and sets the Value accordingly. The Checks, if
// any, will be applied and if any of them return an error the Value will not
// be updated and the error will be returned. If the Editor is non-nil then
// it is applied to each of the listed parameter values. If the Editor
// returns an error for any of the listed values then the Value will not be
// updated and the error is returned.
func (s StrList[T]) SetWithVal(paramName, paramVal string) error {
	sep := s.GetSeparator()
	sv := strings.Split(paramVal, sep)

	if s.Editor != nil {
		for i, val := range sv {
			newVal, err := s.Editor.Edit(paramName, val)
			if err != nil {
				return err
			}

			sv[i] = newVal
		}
	}

	v := make([]T, 0, len(sv))
	for _, strVal := range sv {
		v = append(v, T(strVal))
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

// AllowedValues returns a description of the allowed values. It includes the
// separator to be used
func (s StrList[T]) AllowedValues() string {
	return s.ListValDesc("string values") + HasChecks(s)
}

// CurrentValue returns the current setting of the parameter value
func (s StrList[T]) CurrentValue() string {
	var cv strings.Builder

	sep := ""

	for _, v := range *s.Value {
		cv.WriteString(sep)
		cv.WriteString(string(v))

		sep = s.GetSeparator()
	}

	return cv.String()
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or if it has nil Checks.
func (s StrList[T]) CheckSetter(name string) {
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
