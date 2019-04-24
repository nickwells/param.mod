package psetter

import (
	"fmt"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/fileparse.mod/fileparse"
	"github.com/nickwells/param.mod/v2/param"
)

// Pathname allows you to specify a parameter that can be used to set an
// pathname value. You can specify some required attributes of the referenced
// file-system object by setting the AttributesNeeded member. You can also
// supply a check function that will validate the Value. For instance you
// could check that the value was in a particular directory
type Pathname struct {
	param.ValueReqMandatory

	Value       *string
	Expectation filecheck.Provisos
	Checks      []check.String
}

// CountChecks returns the number of check functions this setter has
func (s Pathname) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) checks first that
// the value can be converted into a pathname. Then it confirms that the file
// conforms to the supplied provisos. If there are checks and any check is
// violated it returns an error. Only if the value is converted successfully
// and no checks are violated is the Value set and a nil error is returned.
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
