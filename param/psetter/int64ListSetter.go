package psetter

import (
	"errors"
	"fmt"
	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/param"
	"strconv"
	"strings"
)

// Int64ListSetter allows you to specify a parameter that can be used to set a
// list (a slice) of int64's. You can override the list separator by setting
// the Sep value.
//
// If you have a list of allowed values you should use EnumListSetter
type Int64ListSetter struct {
	Value *[]int64
	StrListSeparator
	Checks []check.Int64Slice
}

// ValueReq returns param.Mandatory indicating that some value must follow
// the parameter
func (s Int64ListSetter) ValueReq() param.ValueReq { return param.Mandatory }

// Set (called when there is no following value) returns an error
func (s Int64ListSetter) Set(_ string) error {
	return errors.New("no value given (it should be followed by '=...')")
}

// SetWithVal (called when a value follows the parameter) splits the value
// into a slice of int64's and sets the Value accordingly. It will return
// an error if a check is breached.
func (s Int64ListSetter) SetWithVal(_ string, paramVal string) error {
	sep := s.GetSeparator()
	sv := strings.Split(paramVal, sep)

	v := make([]int64, 0, len(sv))
	for i, strVal := range sv {
		intVal, err := strconv.ParseInt(strVal, 0, 0)
		if err != nil {
			return fmt.Errorf(
				"list entry: %d (%s)"+
					" could not be parsed as an integer value: %s",
				i, strVal, err)
		}
		v = append(v, intVal)
	}

	if len(s.Checks) != 0 {
		for _, check := range s.Checks {
			if check == nil {
				continue
			}

			err := check(v)
			if err != nil {
				return err
			}
		}
	}
	*s.Value = v

	return nil
}

// AllowedValues returns a description of the allowed values. It includes the
// separator to be used
func (s Int64ListSetter) AllowedValues() string {
	rval := "a list of whole numbers separated by '" +
		s.GetSeparator() + "'"

	if len(s.Checks) != 0 {
		rval += " subject to checks"
	}
	return rval
}

// CurrentValue returns the current setting of the parameter value
func (s Int64ListSetter) CurrentValue() string {
	cv := ""
	sep := ""

	for _, v := range *s.Value {
		cv += sep + fmt.Sprintf("%v", v)
		sep = s.GetSeparator()
	}

	return cv
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s Int64ListSetter) CheckSetter(name string) {
	if s.Value == nil {
		panic(name +
			": Int64ListSetter Check failed: the Value to be set is nil")
	}
}
