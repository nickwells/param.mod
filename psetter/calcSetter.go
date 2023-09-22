package psetter

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
)

// Calculated allows you to give a parameter where certain predetermined
// values are set by their associated functions and any other values are
// handled by the default func. All the NamedCalc values must have a
// non-empty Name and a non-nil Calc.
type Calculated[T any] struct {
	ValueReqMandatory

	// You must set a Value, the program will panic if not. This is a pointer
	// to the value that the setter is setting.
	Value *T

	// If a parameter value is given and it has an entry in the CalcMap then
	// the corresponding NamedCalc's Calc function is called and the
	// resulting value and error are used to set the Value or as the error
	// return. The CalcMap must have at least one entry. The map keys and the
	// Calc Name fields are used to generate the help text for the parameter.
	CalcMap map[string]NamedCalc[T]
	// If the given value does not match any entry in the CalcMap then the
	// Default NamedCalc is used (unless the NoDefault field is set to true)
	Default NamedCalc[T]
	// To disallow unmapped values set the NoDefault field to true. In this
	// case the CalcMap must have more than one entry
	NoDefault bool

	// The Checks, if any, are applied to the supplied parameter value and
	// the Value will only be update if they all return a nil error.
	Checks []check.ValCk[T]
}

// CountChecks returns the number of check functions this setter has
func (s Calculated[T]) CountChecks() int {
	return len(s.Checks)
}

// CurrentValue returns the current setting of the parameter value
func (s Calculated[T]) CurrentValue() string {
	return fmt.Sprintf("%v", *s.Value)
}

// SetWithVal (called when a value follows the parameter) checks the value
// for validity and only if it is allowed does it set the Value. It returns
// an error if the value is invalid.
func (s Calculated[T]) SetWithVal(paramName string, paramVal string) error {
	nc, ok := s.CalcMap[paramVal]
	if !ok {
		if s.NoDefault {
			return fmt.Errorf("bad value: %q is invalid", paramVal)
		}
		nc = s.Default
	}

	v, err := nc.Calc(paramName, paramVal)
	if err != nil {
		return err
	}

	for _, check := range s.Checks {
		err := check(v)
		if err != nil {
			return err
		}
	}

	*s.Value = v
	return nil
}

// AllowedValues returns the allowed values
func (s Calculated[T]) AllowedValues() string {
	if s.NoDefault {
		return ""
	}
	return s.Default.Name
}

// AllowedValuesMap returns the allowed values and their associated tags
func (s Calculated[T]) AllowedValuesMap() AllowedVals[string] {
	avm := AllowedVals[string]{}
	for k, nc := range s.CalcMap {
		avm[k] = nc.Name
	}
	return avm
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or if it has nil Checks. The value is forced to the default
// target line length if it is not greater than 0.
func (s Calculated[T]) CheckSetter(name string) {
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

	// Check the NamedCalc values
	if len(s.CalcMap) < 1 {
		panic("the CalcMap cannot be empty")
	}
	if s.NoDefault {
		if len(s.CalcMap) < 2 {
			panic("with no default value the CalcMap must have" +
				" at least 2 entries")
		}
	} else {
		err := s.Default.Check()
		if err != nil {
			panic(fmt.Sprintf("the default NamedCalc is invalid: %s", err))
		}
	}

	for k, nc := range s.CalcMap {
		err := nc.Check()
		if err != nil {
			panic(
				fmt.Sprintf("the CalcMap[%q] has an invalid NamedCalc: %s",
					k, err))
		}
	}
}

// ValDescribe returns a name describing the values allowed
func (s Calculated[T]) ValDescribe() string {
	return "..."
}
