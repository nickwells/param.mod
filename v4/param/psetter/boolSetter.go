package psetter

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/nickwells/param.mod/v3/param"
)

// Bool is used to set boolean flags
//
// The Invert flag is used to invert the normal meaning of a boolean parameter.
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
type Bool struct {
	param.ValueReqOptional

	Value *bool

	Invert bool
}

// Set sets the parameter value to true
func (s Bool) Set(_ string) error {
	if s.Invert {
		*s.Value = false
	} else {
		*s.Value = true
	}
	return nil
}

// SetWithVal should be called when a value is given for the parameter
func (s Bool) SetWithVal(_, val string) error {
	b, err := strconv.ParseBool(val)
	if err != nil {
		return errors.New(
			"cannot interpret '" +
				val +
				"' as either true or false")
	}
	if s.Invert {
		*s.Value = !b
	} else {
		*s.Value = b
	}
	return nil
}

// AllowedValues returns a description of the allowed values.
func (s Bool) AllowedValues() string {
	return "none (which will be taken as 'true')" +
		" or some value that can be interpreted as true or false." +
		" The value must be given after an '='," +
		" not as a following value, as this is optional"
}

// CurrentValue returns the current setting of the parameter value
func (s Bool) CurrentValue() string {
	if s.Invert {
		return fmt.Sprintf("%v", !*s.Value)
	}
	return fmt.Sprintf("%v", *s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil
func (s Bool) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, "psetter.Bool"))
	}
}
