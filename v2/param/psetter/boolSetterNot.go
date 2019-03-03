package psetter

import (
	"fmt"
	"strconv"

	"github.com/nickwells/param.mod/v2/param"
)

// BoolNot is used to invert the normal meaning of a boolean parameter.
// It is useful where you want to have a parameter of the form 'dont-xxx' but
// use it to set the value of a bool variable (default value: true) such as
// 'xxx' which you can then test by saying:
//
//      if xxx { doXXX() }
//
// rather than having to set the value of a variable which you would have to
// call dontXXX and then test by saying:
//
//      if !dontXXX { doXXX() }
//
// The benefit is that you can avoid the ugly double negative
type BoolNot struct {
	param.ValueReqOptional

	Value *bool
}

// Set sets the parameter value to false
func (s BoolNot) Set(_ string) error {
	*s.Value = false
	return nil
}

// SetWithVal should be called when a value is given for the parameter
func (s BoolNot) SetWithVal(_, val string) error {
	b, err := strconv.ParseBool(val)
	if err != nil {
		return err
	}
	*s.Value = !b
	return nil
}

// AllowedValues returns a description of the allowed values
func (s BoolNot) AllowedValues() string {
	return "none (which will be taken as 'true')" +
		" or some value that can be interpreted as true or false"
}

// CurrentValue returns the current setting of the parameter value
func (s BoolNot) CurrentValue() string {
	return fmt.Sprintf("%v", !*s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil
func (s BoolNot) CheckSetter(name string) {
	if s.Value == nil {
		panic(name + ": BoolNot Check failed: the Value to be set is nil")
	}
}
