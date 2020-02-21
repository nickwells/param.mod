package psetter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/v4/param"
)

// Int64List allows you to specify a parameter that can be used to set a
// list (a slice) of int64's. You can override the list separator by setting
// the Sep value.
//
// If you have a list of allowed values you should use EnumList
type Int64List struct {
	param.ValueReqMandatory

	Value *[]int64
	StrListSeparator
	Checks []check.Int64Slice
}

// CountChecks returns the number of check functions this setter has
func (s Int64List) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) splits the value
// into a slice of int64's and sets the Value accordingly. It will return
// an error if a check is breached.
func (s Int64List) SetWithVal(_ string, paramVal string) error {
	sep := s.GetSeparator()
	sv := strings.Split(paramVal, sep)

	v := make([]int64, 0, len(sv))
	for i, strVal := range sv {
		intVal, err := strconv.ParseInt(strVal, 0, 0)
		if err != nil {
			return fmt.Errorf("bad value: %q:"+
				" part: %d (%s) cannot be interpreted as a whole number: %s",
				paramVal, i+1, strVal, err)
		}
		v = append(v, intVal)
	}

	for _, check := range s.Checks {
		if check == nil {
			continue
		}

		err := check(v)
		if err != nil {
			return err
		}
	}

	*s.Value = v

	return nil
}

// AllowedValues returns a description of the allowed values. It includes the
// separator to be used
func (s Int64List) AllowedValues() string {
	return s.ListValDesc("whole numbers") + HasChecks(s)
}

// CurrentValue returns the current setting of the parameter value
func (s Int64List) CurrentValue() string {
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
func (s Int64List) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, "psetter.Int64List"))
	}
}
