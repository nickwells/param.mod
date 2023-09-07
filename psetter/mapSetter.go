package psetter

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v6/pedit"
)

// Map sets the entry in a map of strings. Each value from the
// parameter is used as a key in the map with the map entry set to true.
type Map[T ~string] struct {
	ValueReqMandatory

	// You must set a Value, the program will panic if not. This is the map
	// of strings to bool that the setter is setting
	Value *map[T]bool
	// The Checks, if any, are applied to the supplied parameter value and
	// the new parameter will be applied only if they all return a nil error
	Checks []check.ValCk[map[T]bool]
	// The Editor, if present, is applied to the parameter value after any
	// checks are applied and allows the programmer to modify the value
	// supplied before using it to set the Value.
	Editor pedit.Editor
	// The StrListSeparator allows you to override the default separator
	// between list elements.
	StrListSeparator
}

// CountChecks returns the number of check functions this setter has
func (s Map[T]) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) splits the value
// using the list separator. For each of these values it will try to split it
// into two parts around an '='. The first part has the Editor (if any)
// applied to it and the name is replaced with the edited value. If the
// Editor returns a non-nil error then that is returned and the value is
// unchanged. If there is only one part the named map entry is set to true,
// otherwise it will try to parse the second part as a bool. If it can be so
// parsed then the named map entry will be set to that value, otherwise it
// will return the parsing error. It will run the checks (if any) against the
// map and if any check returns a non-nil error that is returned. Finally it
// will update the Value with the new entries.
//
// Note that the Value map is not replaced compmletely, just updated.
func (s Map[T]) SetWithVal(paramName string, paramVal string) error {
	var err error
	sep := s.GetSeparator()
	values := strings.Split(paramVal, sep)
	m := map[T]bool{}

	for i, v := range values {
		parts := strings.SplitN(v, "=", 2)
		key := parts[0]
		if s.Editor != nil {
			key, err = s.Editor.Edit(paramName, key)
			if err != nil {
				return err
			}
		}
		b := true
		if len(parts) == 2 {
			// check that the bool can be parsed
			b, err = strconv.ParseBool(parts[1])
			if err != nil {
				return fmt.Errorf("bad value: %q: part: %d (%q) is invalid."+
					" The value (%q) cannot be interpreted"+
					" as true or false: %s",
					paramVal, i+1, v, parts[1], err)
			}
		}
		m[T(key)] = b
	}

	for _, check := range s.Checks {
		err = check(m)
		if err != nil {
			return err
		}
	}

	for k, b := range m {
		(*s.Value)[T(k)] = b
	}
	return nil
}

// AllowedValues returns a string listing the allowed values
func (s Map[T]) AllowedValues() string {
	return "a list of string values separated by '" +
		s.GetSeparator() + "'" + HasChecks(s)
}

// CurrentValue returns the current setting of the parameter value
func (s Map[T]) CurrentValue() string {
	cv := ""

	keys := make([]string, 0, len(*s.Value))
	for k := range *s.Value {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)

	sep := ""
	for _, k := range keys {
		cv += sep + fmt.Sprintf("%s=%v", k, (*s.Value)[T(k)])
		sep = "\n"
	}

	return cv
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or or if it has nil Checks. If the the map has not been
// created yet it will be created here.
func (s Map[T]) CheckSetter(name string) {
	// Check the value is not nil
	if s.Value == nil {
		panic(NilValueMessage(name, fmt.Sprintf("%T", s)))
	}

	// Check there are no nil Check funcs
	for i, check := range s.Checks {
		if check == nil {
			panic(NilCheckMessage(name, fmt.Sprintf("%T", s), i))
		}
	}

	// make the map if it is nil
	if *s.Value == nil {
		*s.Value = make(map[T]bool)
	}
}
