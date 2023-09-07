package psetter

import (
	"fmt"
	"strconv"

	"github.com/nickwells/check.mod/v2/check"
	"golang.org/x/exp/constraints"
)

// Float allows you to give a parameter that can be used to set a
// float value (64 or 32 bit).
type Float[T constraints.Float] struct {
	ValueReqMandatory

	// Value must be set, the program will panic if not. This is the value
	// being set
	Value *T
	// The Checks, if any, are applied to the supplied parameter value and
	// the new parameter will be applied only if they all return a nil error
	Checks []check.ValCk[T]
}

// CountChecks returns the number of check functions this setter has
func (s Float[T]) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) checks that the
// value can be parsed to a float, if it cannot be parsed successfully it
// returns an error. The Checks, if any, are called with the value to be
// applied and if any of them return a non-nil error the Value is not updated
// and the error is returned. Only if the parameter value is parsed
// successfully and no checks fail is the Value set.
func (s Float[T]) SetWithVal(_ string, paramVal string) error {
	v64, err := strconv.ParseFloat(paramVal, bitsInType(T(0)))
	if err != nil {
		return fmt.Errorf("could not interpret %q as a number: %s",
			paramVal, err)
	}

	v := T(v64)

	for _, check := range s.Checks {
		err := check(v)
		if err != nil {
			return err
		}
	}

	*s.Value = v
	return nil
}

// AllowedValues returns a string describing the allowed values
func (s Float[T]) AllowedValues() string {
	return "any value that can be read as a number with a decimal place" +
		HasChecks(s)
}

// CurrentValue returns the current setting of the parameter value
func (s Float[T]) CurrentValue() string {
	return fmt.Sprintf("%v", *s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or if it has nil Checks.
func (s Float[T]) CheckSetter(name string) {
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

// ValDescribe returns a name describing the values allowed
func (s Float[T]) ValDescribe() string {
	return "float"
}
