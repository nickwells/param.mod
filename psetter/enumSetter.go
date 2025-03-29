package psetter

import (
	"fmt"
	"maps"
	"slices"
)

// Enum allows you to give a parameter that will only allow one of an
// enumerated range of values which are specified in the AllowedVals map.
//
// It is recommended that you should use string constants for setting the
// value and for initialising the allowed values map to avoid possible
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
type Enum[T ~string] struct {
	ValueReqMandatory
	// The AllowedVals must be set, the program will panic if not. The Value
	// is guaranteed to take one of these values.
	AllowedVals[T]
	// Value must be set, the program will panic if not. This is the value
	// being set
	Value *T
	// AllowInvalidInitialValue can be set to relax the checks on the initial
	// Value. It can be set to allow, for instance, an empty initial value to
	// signify that no choice has yet been made.
	AllowInvalidInitialValue bool
}

// SetWithVal (called when a value follows the parameter) checks the value
// for validity and only if it is in the allowed values list does it set the
// Value. It returns an error if the value is invalid.
func (s Enum[T]) SetWithVal(_ string, paramVal string) error {
	if s.ValueAllowed(paramVal) {
		*s.Value = T(paramVal)
		return nil
	}

	return fmt.Errorf("value not allowed: %q", paramVal)
}

// AllowedValues returns a string listing the allowed values
func (s Enum[T]) AllowedValues() string {
	return "a string"
}

// CurrentValue returns the current setting of the parameter value
func (s Enum[T]) CurrentValue() string {
	return fmt.Sprintf("%v", *s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or there are no allowed values.
func (s Enum[T]) CheckSetter(name string) {
	if s.Value == nil {
		panic(NilValueMessage(name, fmt.Sprintf("%T", s)))
	}

	intro := fmt.Sprintf("%s: %T Check failed: ", name, s)

	if err := s.Check(); err != nil {
		panic(intro + err.Error())
	}

	if s.AllowInvalidInitialValue {
		return
	}

	if !s.ValueAllowed(string(*s.Value)) {
		panic(fmt.Sprintf("%sthe initial value (%s) is not valid",
			intro, *s.Value))
	}
}

// ValDescribe returns a brief description of the allowed values suitable for
// appearing after the parameter name. Note that the full list of values is
// truncated if it gets too long.
func (s Enum[T]) ValDescribe() string {
	const maxValDescLen = 20

	initialVal := string(*s.Value)

	var desc string

	if s.ValueAllowed(initialVal) {
		desc = initialVal
	}

	avals := slices.Sorted(maps.Keys(s.AllowedVals))

	for _, val := range avals {
		if string(val) == initialVal {
			continue
		}

		if len(desc) > 0 {
			desc += "|"
		}

		if len(desc)+len(val) > maxValDescLen {
			return desc + "..."
		}

		desc += string(val)
	}

	return desc
}
