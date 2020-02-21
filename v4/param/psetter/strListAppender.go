package psetter

import (
	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/v4/param"
)

// StrListAppender allows you to specify a parameter that can be used to add
// to a list (a slice) of strings.
//
// The user of the program which has a parameter of this type can pass
// multiple parameters and each will add to the list of values rather than
// replacing it each time. Note that each value must be passed separately -
// there is no way to pass multiple values at the same time. Also note that
// there is no way to reset the value, if y=this feature is required another
// parameter could be set up that will do this.
//
// If you have a list of allowed values you should use EnumList
type StrListAppender struct {
	param.ValueReqMandatory

	Value *[]string
	StrListSeparator
	Checks []check.String
}

// CountChecks returns the number of check functions this setter has
func (s StrListAppender) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) splits the value
// into a slice of strings and sets the Value accordingly. It will return
// an error if a check is breached.
func (s StrListAppender) SetWithVal(_ string, paramVal string) error {
	for _, check := range s.Checks {
		if check == nil {
			continue
		}

		err := check(paramVal)
		if err != nil {
			return err
		}
	}

	*s.Value = append(*s.Value, paramVal)

	return nil
}

// AllowedValues returns a description of the allowed values. It includes the
// separator to be used
func (s StrListAppender) AllowedValues() string {
	return "a string value that will be added to the existing list of values" + HasChecks(s)
}

// CurrentValue returns the current setting of the parameter value
func (s StrListAppender) CurrentValue() string {
	cv := ""
	sep := ""

	for _, v := range *s.Value {
		cv += sep + v
		sep = s.GetSeparator()
	}

	return cv
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s StrListAppender) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, "psetter.StrListAppender"))
	}
}
