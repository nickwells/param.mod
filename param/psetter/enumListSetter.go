package psetter

import (
	"fmt"
	"strings"

	"github.com/nickwells/check.mod/v2/check"
)

// EnumList sets the values in a slice of strings. The values must be in
// the allowed values map.
//
// It is recommended that you should use string constants for setting the
// list entries and for initialising the allowed values map to avoid possible
// errors.
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
type EnumList struct {
	ValueReqMandatory
	// The AllowedVals must be set, the program will panic if not. These are
	// the only values that will be allowed in the slice of strings.
	AllowedVals

	// The Aliases need not be given but if they are then each alias must not
	// be in AllowedVals and all of the resulting values must be in
	// AllowedVals.
	Aliases

	// Value must be set, the program will panic if not. This is the slice of
	// values that this setter is setting.
	Value *[]string
	// The StrListSeparator allows you to override the default separator
	// between list elements.
	StrListSeparator
	// The Checks, if any, are applied to the list of new values and the
	// Value will only be updated if they all return a nil error.
	Checks []check.StringSlice
}

// CountChecks returns the number of check functions this setter has
func (s EnumList) CountChecks() int {
	return len(s.Checks)
}

// SetWithVal (called when a value follows the parameter) splits the value
// using the list separator. It then checks all the values for validity and
// only if all the values are in the allowed values list does it add them
// to the slice of strings pointed to by the Value. It returns a error for
// the first invalid value or if a check is breached.
func (s EnumList) SetWithVal(_ string, paramVal string) error {
	sep := s.GetSeparator()
	vals := strings.Split(paramVal, sep)
	aliasedVals := []string{}
	for _, v := range vals {
		if !s.ValueAllowed(v) {
			if !s.Aliases.IsAnAlias(v) {
				return fmt.Errorf("value is not allowed: %q", v)
			}
			aliasedVals = append(aliasedVals, s.Aliases.AliasVal(v)...)
			continue
		}
		aliasedVals = append(aliasedVals, v)
	}

	for _, check := range s.Checks {
		if check == nil {
			continue
		}

		err := check(aliasedVals)
		if err != nil {
			return err
		}
	}
	*s.Value = aliasedVals

	return nil
}

// AllowedValues returns a string listing the allowed values
func (s EnumList) AllowedValues() string {
	return s.ListValDesc("string values") +
		HasChecks(s) + "."
}

// CurrentValue returns the current setting of the parameter value
func (s EnumList) CurrentValue() string {
	str := ""
	sep := ""

	for _, v := range *s.Value {
		str += sep + fmt.Sprintf("%v", v)
		sep = s.GetSeparator()
	}

	return str
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or there are no allowed values or the initial value is not
// allowed.
func (s EnumList[T]) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, fmt.Sprintf("%T", s)))
	}
	intro := fmt.Sprintf("%s: %T Check failed: ", name, s)
	if err := s.AllowedVals.Check(); err != nil {
		panic(intro + err.Error())
	}
	for i, v := range *s.Value {
		if _, ok := s.AllowedVals[string(v)]; !ok {
			panic(fmt.Sprintf(
				"%selement %d (%s) in the current list of entries is invalid",
				intro, i, v))
		}
	}
}
