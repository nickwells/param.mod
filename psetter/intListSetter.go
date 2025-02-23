package psetter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/mathutil.mod/v2/mathutil"
	"golang.org/x/exp/constraints"
)

// IntList allows you to give a parameter that can be used to set a
// list (a slice) of ints's.
type IntList[T constraints.Signed] struct {
	ValueReqMandatory

	// Value must be set, the program will panic if not. This is the slice of
	// int64's that the setter is setting.
	Value *[]T
	// The StrListSeparator allows you to override the default separator
	// between list elements.
	StrListSeparator
	// The Checks, if any, are applied to the supplied parameter value and
	// the new parameter will be applied only if they all return a nil error
	Checks []check.ValCk[[]T]
}

// CountChecks returns the number of check functions this setter has
func (s IntList[T]) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) splits the value
// into a slice of int64's and sets the Value accordingly. The Checks, if
// any, are run against the new list of int64's and if any Check returns a
// non-nil error the Value is not updated and the error is returned.
func (s IntList[T]) SetWithVal(_ string, paramVal string) error {
	sep := s.GetSeparator()
	sv := strings.Split(paramVal, sep)
	v := make([]T, 0, len(sv))

	for i, strVal := range sv {
		i64, err := strconv.ParseInt(strVal, 0, mathutil.BitsInType(T(0)))
		if err != nil {
			return fmt.Errorf("bad value: %q:"+
				" part: %d (%s) cannot be interpreted as a whole number: %s",
				paramVal, i+1, strVal, err)
		}

		intVal := T(i64)

		v = append(v, intVal)
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
func (s IntList[T]) AllowedValues() string {
	return s.ListValDesc("whole numbers") + HasChecks(s)
}

// CurrentValue returns the current setting of the parameter value
func (s IntList[T]) CurrentValue() string {
	cv := ""
	sep := ""

	for _, v := range *s.Value {
		cv += sep + fmt.Sprintf("%v", v)
		sep = s.GetSeparator()
	}

	return cv
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or if it has nil Checks.
func (s IntList[T]) CheckSetter(name string) {
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
func (s IntList[T]) ValDescribe() string {
	return "list-of-ints"
}
