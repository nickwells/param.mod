package psetter

import (
	"errors"
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
//
// The advantages of const values are:
//
// - typos become compilation errors rather than silently failing.
//
// - the name of the constant value can distinguish between the string value
// and it's meaning as a semantic element representing a flag used to choose
// program behaviour.
//
// - the name that you give the const value can distinguish between identical
// strings and show which of various flags with the same string value you
// actually mean.
type EnumMap[T ~string] struct {
	ValueReqMandatory
	// The AllowedVals must be set, the program will panic if not. These are
	// the allowed keys in the Values map
	AllowedVals[T]

	// The Aliases need not be given but if they are then each alias must not
	// be in AllowedVals and all of the resulting values must be in
	// AllowedVals.
	Aliases[T]

	// Value must be set, the program will panic if not. This is the map of
	// values that this setter is setting
	Value *map[T]bool
	// AllowHiddenMapEntries can be set to relax the checks on the initial
	// entries in the Values map
	AllowHiddenMapEntries bool
	// The StrListSeparator allows you to override the default separator
	// between list elements.
	StrListSeparator
}

// SetWithVal (called when a value follows the parameter) splits the value
// using the list separator. It then checks all the values for validity and
// only if all the values are in the allowed values list does it set the
// entries in the map of strings pointed to by the Value. It returns a error
// for the first invalid value.
func (s EnumMap[T]) SetWithVal(_ string, paramVal string) error {
	if paramVal == "" {
		return errors.New("empty value. Some value must be given")
	}

	values := strings.Split(paramVal, s.GetSeparator())

	err := s.checkValues(paramVal, values)
	if err != nil {
		return err
	}

	s.setValues(values)

	return nil
}

// checkValues checks that all the values given are OK and returns an error
// for any bad value
func (s EnumMap[T]) checkValues(paramVal string, values []string) error {
	for i, v := range values {
		namePart, boolPart, hasBoolPart := strings.Cut(v, "=")
		// check the name is an allowed value
		if !s.ValueAllowed(namePart) && !s.Aliases.IsAnAlias(namePart) {
			return fmt.Errorf("bad value: %q: part: %d (%q) is invalid."+
				" The name (%q) is not allowed",
				paramVal, i+1, v, namePart)
		}

		if hasBoolPart {
			// check that the bool can be parsed
			_, err := strconv.ParseBool(boolPart)
			if err != nil {
				return fmt.Errorf("bad value: %q:"+
					" part: %d (%q) is invalid."+
					" The value (%q) cannot be interpreted"+
					" as true or false: %w",
					paramVal, i+1, v, boolPart, err)
			}
		}
	}

	return nil
}

// setValues sets the values in the Value map from the strings in the slice
// which have already been checked for validity.
func (s EnumMap[T]) setValues(values []string) {
	for _, v := range values {
		namePart, boolPart, hasBoolPart := strings.Cut(v, "=")

		name := T(namePart)
		keys := []T{name}

		if s.Aliases.IsAnAlias(namePart) {
			keys = s.AliasVal(name)
		}

		b := true
		if hasBoolPart {
			b, _ = strconv.ParseBool(boolPart)
		}

		for _, k := range keys {
			(*s.Value)[k] = b
		}
	}
}

// AllowedValues returns a string listing the allowed values
func (s EnumMap[T]) AllowedValues() string {
	return s.ListValDesc("string values") +
		".\n\nEach value can be set to false by following the value" +
		" with '=false'; by default the value will be set to true."
}

// CurrentValue returns the current setting of the parameter value
func (s EnumMap[T]) CurrentValue() string {
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
// Value is nil or the map has not been created yet or if there are no
// allowed values.
func (s EnumMap[T]) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, fmt.Sprintf("%T", s)))
	}

	if *s.Value == nil {
		*s.Value = make(map[T]bool)
	}

	intro := fmt.Sprintf("%s: %T Check failed: ", name, s)

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
			panic(fmt.Sprintf("%sthe map entry with key %q is invalid"+
				" - it is not in the allowed values map",
				intro, k))
		}
	}
}

// ValDescribe returns a brief description of the allowed values suitable for
// appearing after the parameter name. Note that the full list of values is
// truncated if it gets too long.
func (s EnumMap[T]) ValDescribe() string {
	const maxValDescLen = 20

	var desc string

	avals, _ := s.AllowedVals.Keys()
	aliasKeys, _ := s.Aliases.Keys()
	avals = append(avals, aliasKeys...)

	sort.Strings(avals)

	sep := ""

	var incomplete bool

	optEqVal := [...]string{"", "=false", "=true"}

	for i, val := range avals {
		eqVal := optEqVal[i%len(optEqVal)]

		if len(desc)+len(val)+len(eqVal)+len(sep) > maxValDescLen {
			incomplete = true
			continue
		}

		desc += sep + val + eqVal
		sep = s.GetSeparator()
	}

	if incomplete {
		desc += "..."
	}

	return desc
}
