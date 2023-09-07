package psetter

import (
	"fmt"
	"path/filepath"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/fileparse.mod/fileparse"
)

// PathnameListAppender allows you to specify a parameter that can be used to
// add to a list (a slice) of pathnames.
//
// The user of the program which has a parameter of this type can pass
// multiple parameters and each will add to the list of values rather than
// replacing it each time. Note that each value must be passed separately;
// there is no way to pass multiple values at the same time. Also note that
// there is no way to reset the value, if this feature is required another
// parameter could be set up that will do this.
type PathnameListAppender struct {
	ValueReqMandatory

	// You must set a Value, the program will panic if not. This is the slice
	// of strings that the setter is appending to.
	Value *[]string
	// Expectation allows you to set some file-specific checks.
	Expectation filecheck.Provisos
	// The Checks, if any, are applied to the supplied parameter value and
	// the new parameter will be added to the list only if they all return a
	// nil error.
	Checks []check.String
	// Prepend will change the behaviour so that any new values are added at
	// the start of the list of pathnames rather than the end.
	Prepend bool
	// ForceAbsolute, if set, causes any pathname value to be passed
	// to filepath.Abs before setting the value.
	ForceAbsolute bool
}

// CountChecks returns the number of check functions this setter has
func (s PathnameListAppender) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) takes the parameter
// value and runs the checks against it. If any check returns a non-nil error
// it will return the error. Otherwise it will apply the Editor (if there is
// one) to the parameter value. If the Editor returns a non-nil error then
// that is returned and the Value is left unchanged.  Finally, it will append
// the checked and possibly edited value to the slice of strings.
func (s PathnameListAppender) SetWithVal(_, paramVal string) error {
	pathname, err := fileparse.FixFileName(paramVal)
	if err != nil {
		return err
	}

	if s.ForceAbsolute {
		pathname, err = filepath.Abs(pathname)
		if err != nil {
			return err
		}
	}

	err = s.Expectation.StatusCheck(pathname)
	if err != nil {
		return err
	}

	for _, check := range s.Checks {
		if check == nil {
			continue
		}

		err := check(pathname)
		if err != nil {
			return err
		}
	}

	if s.Prepend {
		*s.Value = append([]string{pathname}, *s.Value...)
		return nil
	}
	*s.Value = append(*s.Value, pathname)
	return nil
}

// AllowedValues returns a description of the allowed values. It includes the
// separator to be used
func (s PathnameListAppender) AllowedValues() string {
	const (
		intro = "a pathname that will be added to the"
		outro = " existing list of values"
	)
	prepend := ""
	if s.Prepend {
		prepend = " start of the"
	}
	return intro + prepend + outro + HasChecks(s)
}

// CurrentValue returns the current setting of the parameter value
func (s PathnameListAppender) CurrentValue() string {
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
func (s PathnameListAppender) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, fmt.Sprintf("%T", s)))
	}
}

// ValDescribe returns a brief description of the expected value
func (s PathnameListAppender) ValDescribe() string {
	return "pathname"
}
