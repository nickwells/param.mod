package psetter

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/nickwells/mathutil.mod/v2/mathutil"
	"golang.org/x/exp/constraints"
)

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

// ConstCalc returns a NamedCalc of the correct type which will set the value
// to the supplied value
func ConstCalc[T any](name string, val T) NamedCalc[T] {
	return NamedCalc[T]{
		Name: name,
		Calc: func(_, _ string) (T, error) {
			return val, nil
		},
	}
}

// Env2IntCalc returns a NamedCalc of the correct type which will return the
// value of the given environment variable converted to the appropriate
// signed integer type
func Env2IntCalc[T constraints.Signed](name, envName string) NamedCalc[T] {
	return NamedCalc[T]{
		Name: name,
		Calc: func(_, _ string) (T, error) {
			evVal := os.Getenv(envName)
			if evVal == "" {
				return 0, fmt.Errorf(
					"the %s environment variable is not set",
					envName)
			}
			v, err := strconv.ParseInt(evVal, 0, mathutil.BitsInType(T(0)))
			if err != nil {
				return 0, fmt.Errorf(
					"bad %s environment variable value: %q: %w",
					envName, evVal, err)
			}
			return T(v), nil
		},
	}
}

// Env2UintCalc returns a NamedCalc of the correct type which will return the
// value of the given environment variable converted to the appropriate
// unsigned integer type
func Env2UintCalc[T constraints.Unsigned](name, envName string) NamedCalc[T] {
	return NamedCalc[T]{
		Name: name,
		Calc: func(_, _ string) (T, error) {
			evVal := os.Getenv(envName)
			if evVal == "" {
				return 0, fmt.Errorf(
					"the %s environment variable is not set",
					envName)
			}
			v, err := strconv.ParseUint(evVal, 0, mathutil.BitsInType(T(0)))
			if err != nil {
				return 0, fmt.Errorf(
					"bad %s environment variable value: %q: %w",
					envName, evVal, err)
			}
			return T(v), nil
		},
	}
}

// Val2IntCalc returns a NamedCalc of the correct type which will set the
// value to the value of the parameter value converted to the
// appropriate signed integer type
func Val2IntCalc[T constraints.Signed]() NamedCalc[T] {
	return NamedCalc[T]{
		Name: "some value that can be read as a whole number",
		Calc: func(_, paramVal string) (T, error) {
			v, err := strconv.ParseInt(paramVal, 0, mathutil.BitsInType(T(0)))
			if err != nil {
				return 0, fmt.Errorf("bad parameter: %q: %w", paramVal, err)
			}
			return T(v), nil
		},
	}
}

// Val2UintCalc returns a NamedCalc of the correct type which will return the
// parameter value converted to the appropriate unsigned integer type
func Val2UintCalc[T constraints.Unsigned]() NamedCalc[T] {
	return NamedCalc[T]{
		Name: "a parameter that can be converted to an unsigned integer",
		Calc: func(_, paramVal string) (T, error) {
			v, err := strconv.ParseUint(paramVal, 0, mathutil.BitsInType(T(0)))
			if err != nil {
				return 0, fmt.Errorf("bad parameter: %q: %w", paramVal, err)
			}
			return T(v), nil
		},
	}
}

// Env2StringCalc returns a NamedCalc of the correct type which will return the
// value of the given environment variable
func Env2StringCalc[T ~string](name, envName string) NamedCalc[T] {
	return NamedCalc[T]{
		Name: name,
		Calc: func(_, _ string) (T, error) {
			return T(os.Getenv(envName)), nil
		},
	}
}

// Val2StringCalc returns a NamedCalc of the correct type which will set the
// value to the value of the parameter value
func Val2StringCalc[T ~string]() NamedCalc[T] {
	return NamedCalc[T]{
		Name: "a parameter value",
		Calc: func(_, paramVal string) (T, error) {
			return T(paramVal), nil
		},
	}
}
