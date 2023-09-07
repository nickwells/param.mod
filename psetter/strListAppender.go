package psetter

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v6/pedit"
)

// StrListAppender allows you to specify a parameter that can be used to add
// to a list (a slice) of strings.
//
// The user of the program which has a parameter of this type can pass
// multiple parameters and each will add to the list of values rather than
// replacing it each time. Note that each value must be passed separately;
// there is no way to pass multiple values at the same time. Also note that
// there is no way to reset the value, if this feature is required another
// parameter could be set up that will do this.
type StrListAppender struct {
	ValueReqMandatory

	// You must set a Value, the program will panic if not. This is the slice
	// of strings that the setter is appending to.
	Value *[]string
	// The Checks, if any, are applied to the supplied parameter value and
	// the new parameter will be added to the list only if they all return a
	// nil error.
	Checks []check.String
	// The Editor, if present, is applied to the parameter value after any
	// checks are applied and allows the programmer to modify the value
	// supplied before using it to set the Value.
	Editor pedit.Editor
	// Prepend will change the behaviour so that any new values are added at
	// the start of the list of strings rather than the end.
	Prepend bool
}

// CountChecks returns the number of check functions this setter has
func (s StrListAppender) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) takes the parameter
// value and runs the checks against it. If any check returns a non-nil error
// it will return the error. Otherwise it will apply the Editor (if there is
// one) to the parameter value. If the Editor returns a non-nil error then
// that is returned and the Value is left unchanged.  Finally, it will append
// the checked and possibly edited value to the slice of strings.
func (s StrListAppender) SetWithVal(paramName, paramVal string) error {
	for _, check := range s.Checks {
		if check == nil {
			continue
		}

		err := check(paramVal)
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

	if s.Prepend {
		*s.Value = append([]string{paramVal}, *s.Value...)
		return nil
	}
	*s.Value = append(*s.Value, paramVal)
	return nil
}

// AllowedValues returns a description of the allowed values. It includes the
// separator to be used
func (s StrListAppender) AllowedValues() string {
	const (
		intro = "a string that will be added to the"
		outro = " existing list of values"
	)
	prepend := ""
	if s.Prepend {
		prepend = " start of the"
	}
	return intro + prepend + outro + HasChecks(s)
}

// CurrentValue returns the current setting of the parameter value
func (s StrListAppender) CurrentValue() string {
	cv := ""
	sep := ""

	for _, v := range *s.Value {
		cv += sep + v
		sep = "\n"
	}

	return cv
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s StrListAppender) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, fmt.Sprintf("%T", s)))
	}
}
