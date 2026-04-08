package psetter

import (
	"fmt"
	"regexp"
	"strings"
)

// RegexpListAppender allows you to specify a parameter that can be used to add
// to a list (a slice) of regular expressions.
//
// The user of the program which has a parameter of this type can pass
// multiple parameters and each will add to the list of values rather than
// replacing it each time. Note that each regexp must be passed separately;
// there is no way to pass multiple regexps at the same time. Also note that
// there is no way to reset the list of regexps, if this feature is required
// another parameter could be set up that will do this.
type RegexpListAppender struct {
	ValueReqMandatory

	// You must set a Value, the program will panic if not. This is the slice
	// of regular expressions that the setter is appending to.
	Value *[]*regexp.Regexp
}

// SetWithVal (called when a value follows the parameter) takes the parameter
// value and runs the checks against it. If any check returns a non-nil error
// it will return the error. Otherwise it will apply the Editor (if there is
// one) to the parameter value. If the Editor returns a non-nil error then
// that is returned and the Value is left unchanged.  Finally, it will append
// the checked and possibly edited value to the slice of strings.
func (s RegexpListAppender) SetWithVal(_, paramVal string) error {
	v, err := regexp.Compile(paramVal)
	if err != nil {
		return fmt.Errorf("could not parse %q into a regular expression: %s",
			paramVal, err)
	}

	*s.Value = append(*s.Value, v)

	return nil
}

// AllowedValues returns a description of the allowed values. It includes the
// separator to be used
func (s RegexpListAppender) AllowedValues() string {
	return "any value that can be compiled as a regular expression" +
		" by the standard regexp package"
}

// CurrentValue returns the current setting of the parameter value
func (s RegexpListAppender) CurrentValue() string {
	var cv strings.Builder

	sep := ""

	for _, v := range *s.Value {
		cv.WriteString(sep)
		cv.WriteString(v.String())

		sep = "\n"
	}

	return cv.String()
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or if it has nil Checks.
func (s RegexpListAppender) CheckSetter(name string) {
	// Check the value is not nil
	if s.Value == nil {
		panic(NilValueMessage(name, fmt.Sprintf("%T", s)))
	}
}

// ValDescribe returns a string giving a summary of the values that can
// follow the parameter name.
func (s RegexpListAppender) ValDescribe() string { return "pattern" }
