package psetter

import (
	"fmt"
	"strings"

	"github.com/nickwells/param.mod/v2/param"
)

// Map sets the entry in a map of strings. Each value from the
// parameter is used as a key in the map with the map entry set to true.
type Map struct {
	param.ValueReqMandatory

	Value *map[string]bool
	StrListSeparator
}

// SetWithVal (called when a value follows the parameter) splits the value
// using the list separator.
func (s Map) SetWithVal(_ string, paramVal string) error {
	sep := s.GetSeparator()
	values := strings.Split(paramVal, sep)
	for _, v := range values {
		(*s.Value)[v] = true
	}
	return nil
}

// AllowedValues returns a string listing the allowed values
func (s Map) AllowedValues() string {
	return "a list of string values separated by '" +
		s.GetSeparator() + "'."
}

// CurrentValue returns the current setting of the parameter value
func (s Map) CurrentValue() string {
	cv := ""
	sep := ""
	for k, v := range *s.Value {
		cv += sep + fmt.Sprintf("%s=%v", k, v)
		sep = s.GetSeparator()
	}

	return cv
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or the map has not been created yet.
func (s Map) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, "psetter.Map"))
	}
	if *s.Value == nil {
		panic(name + ": psetter.Map Check failed: the map has not been created")
	}
}
