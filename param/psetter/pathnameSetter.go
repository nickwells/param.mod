package psetter

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/fileparse.mod/fileparse"
)

// Pathname allows you to give a parameter that can be used to set a pathname
// value.
type Pathname struct {
	ValueReqMandatory

	// You must set a Value, the program will panic if not. This is the
	// pathname that the setter is setting.
	Value *string
	// Expectation allows you to set some file-specific checks.
	Expectation filecheck.Provisos
	// The Checks, if any, are applied to the supplied parameter value and
	// the new parameter will be applied only if they all return a nil error.
	Checks []check.String
}

// CountChecks returns the number of check functions this setter has
func (s Pathname) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) checks first that
// the value can be converted into a pathname (a tilda at the start of the
// path is converted to the appropriate home directory). Then it confirms
// that the file conforms to the supplied provisos. The Checks, if any, are
// run and if any check returns a non-nil error the Value is not updated and
// the error is returned. Only if the value is converted successfully, the
// Expectations are all met and no checks fail is the Value set and a nil
// error is returned.
func (s Pathname) SetWithVal(_ string, paramVal string) error {
	pathname, err := fileparse.FixFileName(paramVal)
	if err != nil {
		return err
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

	*s.Value = pathname
	return nil
}

// AllowedValues returns a string describing the allowed values
func (s Pathname) AllowedValues() string {
	rval := "a pathname" + HasChecks(s)

	extras := s.Expectation.String()
	if extras != "" {
		rval += ". " + extras
	}

	return rval
}

// CurrentValue returns the current setting of the parameter value
func (s Pathname) CurrentValue() string {
	return fmt.Sprintf("%v", *s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s Pathname) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, "psetter.Pathname"))
	}
}
