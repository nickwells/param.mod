package psetter

import (
	"errors"
	"fmt"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/fileparse.mod/fileparse"
	"github.com/nickwells/param.mod/param"
)

// PathnameSetter allows you to specify a parameter that can be used to set an
// pathname value. You can specify some required attributes of the referenced
// file-system object by setting the AttributesNeeded member. You can also
// supply a check function that will validate the Value. For instance you
// could check that the value was in a particular directory
type PathnameSetter struct {
	Value       *string
	Expectation filecheck.Provisos
	Checks      []check.String
}

// ValueReq returns param.Mandatory indicating that some value must follow
// the parameter
func (s PathnameSetter) ValueReq() param.ValueReq { return param.Mandatory }

// Set (called when there is no following value) returns an error
func (s PathnameSetter) Set(_ string) error {
	return errors.New("no pathname given (it should be followed by '=...')")
}

// SetWithVal (called when a value follows the parameter) checks first that the
// value, if it cannot be parsed successfully it returns an
// error. If there is a check and the check is violated it returns an
// error. Only if the value is parsed successfully and the check is not
// violated is the Value set.
func (s PathnameSetter) SetWithVal(_ string, paramVal string) error {
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
func (s PathnameSetter) AllowedValues() string {
	rval := "a pathname"

	extras := s.Expectation.String()
	if extras != "" {
		rval += ". " + extras
	}

	return rval
}

// CurrentValue returns the current setting of the parameter value
func (s PathnameSetter) CurrentValue() string {
	return fmt.Sprintf("%v", *s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s PathnameSetter) CheckSetter(name string) {
	if s.Value == nil {
		panic(name +
			": PathnameSetter Check failed: the Value to be set is nil")
	}
}
