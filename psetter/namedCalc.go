package psetter

import "errors"

// NamedCalc holds a function and associated text. The text can be used in
// parameter descriptions to describe what the Calc function does. The Calc
// func takes the parameter name and value and returns a value of the
// appropriate type and an error.
type NamedCalc[T any] struct {
	Name string
	Calc func(string, string) (T, error)
}

// Check returns a non-nil error if the Name field is the empty string or the
// Calc func is nil. If both tests pass a nil error is returned.
func (tc NamedCalc[T]) Check() error {
	if tc.Name == "" {
		return errors.New("the Name must not be empty")
	}
	if tc.Calc == nil {
		return errors.New("the Calc must not be nil")
	}
	return nil
}
