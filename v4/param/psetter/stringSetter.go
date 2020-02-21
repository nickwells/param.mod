package psetter

import (
	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/v4/param"
)

// String is the type for setting string values from
// parameters
type String struct {
	param.ValueReqMandatory

	Value  *string
	Checks []check.String
}

// CountChecks returns the number of check functions this setter has
func (s String) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal checks that the parameter value meets the checks if any. It
// returns an error if the check is not satisfied. Only if the check
// is not violated is the Value set.
func (s String) SetWithVal(_ string, paramVal string) error {
	for _, check := range s.Checks {
		if check == nil {
			continue
		}

		err := check(paramVal)
		if err != nil {
			return err
		}
	}

	*s.Value = paramVal
	return nil
}

// AllowedValues simply returns "any string" since String
// does not check its value
func (s String) AllowedValues() string {
	return "any string" + HasChecks(s)
}

// CurrentValue returns the current setting of the parameter value
func (s String) CurrentValue() string {
	return *s.Value
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s String) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, "psetter.String"))
	}
}
