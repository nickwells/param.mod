package psetter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nickwells/mathutil.mod/v2/mathutil"
	"golang.org/x/exp/constraints"
)

// IntList allows you to give a parameter that can be used to set a
// list (a slice) of ints's.
type IntList[T constraints.Signed] struct {
	ValueReqMandatory
	ValueChecker[[]T]

	// Value must be set, the program will panic if not. This is the slice of
	// int64's that the setter is setting.
	Value *[]T
	// The StrListSeparator allows you to override the default separator
	// between list elements.
	StrListSeparator
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

	if err := s.ApplyChecks(v); err != nil {
		return err
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
	var cv strings.Builder

	sep := ""

	for _, v := range *s.Value {
		cv.WriteString(sep)
		fmt.Fprintf(&cv, "%v", v)

		sep = s.GetSeparator()
	}

	return cv.String()
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or if it has nil Checks.
func (s IntList[T]) CheckSetter(name string) {
	// Check the value is not nil
	if s.Value == nil {
		panic(NilValueMessage(name, fmt.Sprintf("%T", s)))
	}

	s.VerifyChecks(name, fmt.Sprintf("%T", s))
}

// ValDescribe returns a name describing the values allowed
func (s IntList[T]) ValDescribe() string {
	return "number" + s.GetSeparator() + "number..."
}
