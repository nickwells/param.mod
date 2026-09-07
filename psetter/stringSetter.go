package psetter

import (
	"fmt"
)

// String is the type for setting string values from
// parameters
type String[T ~string | []byte | []rune] struct {
	ValueReqMandatory
	ValueChecker[T]

	// You must set a Value, the program will panic if not. This is the
	// string that the setter is setting.
	Value *T
	// The Editor, if present, is applied to the parameter value after any
	// checks are applied and allows the programmer to modify the value
	// supplied before using it to set the Value.
	Editor Editor
}

// SetWithVal checks that the parameter value meets the checks if any. It
// returns an error if the check is not satisfied. Only if the check
// is not violated is the Value set.
func (s String[T]) SetWithVal(paramName, paramVal string) error {
	if s.Editor != nil {
		var err error

		paramVal, err = s.Editor.Edit(paramName, paramVal)
		if err != nil {
			return err
		}
	}

	if err := s.ApplyChecks(T(paramVal)); err != nil {
		return err
	}

	*s.Value = T(paramVal)

	return nil
}

// AllowedValues simply returns "any string" since String
// does not check its value
func (s String[T]) AllowedValues() string {
	return "any string" + HasChecks(s)
}

// ValDescribe returns text describing the value expected
func (s String[T]) ValDescribe() string {
	return "string"
}

// CurrentValue returns the current setting of the parameter value
func (s String[T]) CurrentValue() string {
	return string(*s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or if it has nil Checks.
func (s String[T]) CheckSetter(name string) {
	// Check the value is not nil
	if s.Value == nil {
		panic(NilValueMessage(name, fmt.Sprintf("%T", s)))
	}

	s.VerifyChecks(name, fmt.Sprintf("%T", s))
}
