package psetter

import (
	"github.com/nickwells/check.mod/v2/check"
)

// ValueChecker is intended as a mixin type for the param setters. It provides
// functions that will verify the supplied checks and apply them to the
// parsed parameter value.
type ValueChecker[T any] struct {
	// The Checks, if any, are applied to the parsed parameter value.
	Checks []check.ValCk[T]
}

// CountChecks returns the number of check functions.
func (c ValueChecker[T]) CountChecks() int {
	return len(c.Checks)
}

// VerifyChecks panics if any of the checks are nil. If this is called from
// within a CheckSetter method the paramName will have been supplied and the
// setterName can be obtained using fmt.Sprintf("%T", s).
func (c ValueChecker[T]) VerifyChecks(paramName, setterName string) {
	for i, check := range c.Checks {
		if check == nil {
			panic(NilCheckMessage(paramName, setterName, i))
		}
	}
}

// ApplyChecks applies each of the Checks in turn to the supplied value and
// will return the first non-nil error that is returned. If they all pass a
// nil error is returned. This should be called from the param.Setter's
// SetWithVal method.
func (c ValueChecker[T]) ApplyChecks(v T) error {
	for _, check := range c.Checks {
		err := check(v)
		if err != nil {
			return err
		}
	}

	return nil
}
