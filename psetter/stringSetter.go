package psetter

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v6/pedit"
)

// String is the type for setting string values from
// parameters
type String[T ~string | []byte | []rune] struct {
	ValueReqMandatory

	// You must set a Value, the program will panic if not. This is the
	// string that the setter is setting.
	Value *T
	// The Checks, if any, are applied to the supplied parameter value and
	// the new parameter will be applied only if they all return a nil error.
	Checks []check.ValCk[T]
	// The Editor, if present, is applied to the parameter value after any
	// checks are applied and allows the programmer to modify the value
	// supplied before using it to set the Value.
	Editor pedit.Editor
}

// CountChecks returns the number of check functions this setter has
func (s String[T]) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal checks that the parameter value meets the checks if any. It
// returns an error if the check is not satisfied. Only if the check
// is not violated is the Value set.
func (s String[T]) SetWithVal(paramName, paramVal string) error {
	for _, check := range s.Checks {
		if check == nil {
			continue
		}

		err := check(T(paramVal))
		if err != nil {
			return err
		}
	}

	if s.Editor != nil {
		var err error
		paramVal, err = s.Editor.Edit(paramName, paramVal)
		if err != nil {
			return err
		}
	}

	*s.Value = T(paramVal)
	return nil
}

// AllowedValues simply returns "any string" since String
// does not check its value
func (s String[T]) AllowedValues() string {
	return "any string" + HasChecks(s)
}

// CurrentValue returns the current setting of the parameter value
func (s String[T]) CurrentValue() string {
	return string(*s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s String[T]) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, fmt.Sprintf("%T", s)))
	}
}
