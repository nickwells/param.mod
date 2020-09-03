package psetter

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// EnumMap sets the entry in a map of strings. The values initially set in
// the map must be in the allowed values map unless AllowHiddenMapEntries is
// set to true. Only values with keys in the allowed values map can be
// set. If you allow hidden values then you can have entries in your map
// which cannot be set through this interface but this will still only allow
// values to be set which are in the allowed values map.
//
// It is recommended that you should use string constants for setting and
// accessing the map entries and for initialising the allowed values map to
// avoid possible errors.
type EnumMap struct {
	ValueReqMandatory
	// The AllowedVals must be set, the program will panic if not. These are
	// the allowed keys in the Values map
	AllowedVals

	// The Aliases need not be given but if they are then each alias must not
	// be in AllowedVals and all of the resulting values must be in
	// AllowedVals.
	Aliases

	// Value must be set, the program will panic if not. This is the map of
	// values that this setter is setting
	Value *map[string]bool
	// AllowHiddenMapEntries can be set to relax the checks on the initial
	// entries in the Values map
	AllowHiddenMapEntries bool
	StrListSeparator
}

// SetWithVal (called when a value follows the parameter) splits the value
// using the list separator. It then checks all the values for validity and
// only if all the values are in the allowed values list does it set the
// entries in the map of strings pointed to by the Value. It returns a error
// for the first invalid value.
func (s EnumMap) SetWithVal(_ string, paramVal string) error {
	sep := s.GetSeparator()
	values := strings.Split(paramVal, sep)

	for i, v := range values {
		parts := strings.SplitN(v, "=", 2)
		// check the name is an allowed value
		if !s.ValueAllowed(parts[0]) && !s.Aliases.IsAnAlias(parts[0]) {
			return fmt.Errorf("bad value: %q: part: %d (%q) is invalid."+
				" The name (%q) is not allowed",
				paramVal, i+1, v, parts[0])
		}
		if len(parts) == 2 {
			// check that the bool can be parsed
			_, err := strconv.ParseBool(parts[1])
			if err != nil {
				return fmt.Errorf("bad value: %q:"+
					" part: %d (%q) is invalid."+
					" The value (%q) cannot be interpreted"+
					" as true or false: %s",
					paramVal, i+1, v, parts[1], err)
			}
		}
	}
	for _, v := range values {
		parts := strings.SplitN(v, "=", 2)

		keys := []string{parts[0]}
		if s.Aliases.IsAnAlias(parts[0]) {
			keys = s.AliasVal(parts[0])
		}

		b := true
		if len(parts) == 2 {
			b, _ = strconv.ParseBool(parts[1])
		}

		for _, k := range keys {
			(*s.Value)[k] = b
		}
	}
	return nil
}

// AllowedValues returns a string listing the allowed values
func (s EnumMap) AllowedValues() string {
	return s.ListValDesc("string values") +
		".\n\nEach value can be set to false by following the value" +
		" with '=false'; by default the value will be set to true."
}

// CurrentValue returns the current setting of the parameter value
func (s EnumMap) CurrentValue() string {
	cv := ""

	keys := make([]string, 0, len(*s.Value))
	for k := range *s.Value {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sep := ""
	for _, k := range keys {
		cv += sep + fmt.Sprintf("%s=%v", k, (*s.Value)[k])
		sep = "\n"
	}

	return cv
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or the map has not been created yet or if there are no
// allowed values.
func (s EnumMap) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, "psetter.EnumMap"))
	}
	if *s.Value == nil {
		*s.Value = make(map[string]bool)
	}

	intro := name + ": psetter.EnumMap Check failed: "
	if err := s.AllowedVals.Check(); err != nil {
		panic(intro + err.Error())
	}
	if err := s.Aliases.Check(s.AllowedVals); err != nil {
		panic(intro + err.Error())
	}

	if s.AllowHiddenMapEntries {
		return
	}
	for k := range *s.Value {
		if _, ok := s.AllowedVals[k]; !ok {
			panic(fmt.Sprintf("%sthe map entry with key '%s' is invalid"+
				" - it is not in the allowed values map",
				intro, k))
		}
	}
}
