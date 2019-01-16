package psetter

import (
	"errors"
	"fmt"
	"github.com/nickwells/param.mod/v2/param"
	"strings"
)

// MapSetter sets the entry in a map of strings. Each value from the
// parameter is used as a key in the map with the map entry set to true.
type MapSetter struct {
	Value *map[string]bool
	StrListSeparator
}

// ValueReq returns param.Mandatory indicating that some value must follow
// the parameter
func (s MapSetter) ValueReq() param.ValueReq { return param.Mandatory }

// Set (called when there is no following value) returns an error
func (s MapSetter) Set(_ string) error {
	return errors.New("no value given (it should be followed by '=...')")
}

// SetWithVal (called when a value follows the parameter) splits the value
// using the list separator.
func (s MapSetter) SetWithVal(_ string, paramVal string) error {
	sep := s.GetSeparator()
	values := strings.Split(paramVal, sep)
	for _, v := range values {
		(*s.Value)[v] = true
	}
	return nil
}

// AllowedValues returns a string listing the allowed values
func (s MapSetter) AllowedValues() string {
	return "a list of string values separated by '" +
		s.GetSeparator() + "'."
}

// CurrentValue returns the current setting of the parameter value
func (s MapSetter) CurrentValue() string {
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
func (s MapSetter) CheckSetter(name string) {
	if s.Value == nil {
		panic(name + ": MapSetter Check failed: the Value to be set is nil")
	}
	if *s.Value == nil {
		panic(name + ": MapSetter Check failed: the map has not been created")
	}
}
