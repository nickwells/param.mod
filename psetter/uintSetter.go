package psetter

import (
	"fmt"
	"strconv"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/mathutil.mod/v2/mathutil"
	"golang.org/x/exp/constraints"
)

// Uint allows you to give a parameter that can be used to set an
// int64 value.
type Uint[T constraints.Unsigned] struct {
	ValueReqMandatory

	// You must set a Value, the program will panic if not. This is a pointer
	// to the int64 value that the setter is setting.
	Value *T
	// The Checks, if any, are applied to the supplied parameter value and
	// the Value will only be update if they all return a nil error.
	Checks []check.ValCk[T]
}

// CountChecks returns the number of check functions this setter has
func (s Uint[T]) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) checks that the
// value can be parsed to an integer, if it cannot be parsed successfully it
// returns an error. If there are checks and any check is violated it returns
// an error. Only if the value is parsed successfully and no checks are
// violated is the Value set.
func (s Uint[T]) SetWithVal(_ string, paramVal string) error {
	v64, err := strconv.ParseUint(paramVal, 0, mathutil.BitsInType(T(0)))
	if err != nil {
		return fmt.Errorf(
			"could not interpret %q as a positive whole number: %s",
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
func (s Uint[T]) AllowedValues() string {
	return "any value that can be read as a positive whole number" +
		HasChecks(s)
}

// CurrentValue returns the current setting of the parameter value
func (s Uint[T]) CurrentValue() string {
	return fmt.Sprintf("%v", *s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or if it has nil Checks.
func (s Uint[T]) CheckSetter(name string) {
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
func (s Uint[T]) ValDescribe() string {
	return "+ve int"
}
