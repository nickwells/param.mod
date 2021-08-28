package psetter_test

import (
	"fmt"

	"github.com/nickwells/param.mod/v5/param/psetter"
)

// ExampleEnumMap_standard demonstrates the use of an EnumMap setter.
func ExampleEnumMap_standard() {
	ps := newPSetForTesting() // use paramset.NewOrDie()

	const (
		XOption = "x"
		YOption = "y"
	)

	var m map[string]bool
	keys := []string{XOption, YOption}

	ps.Add("my-map",
		psetter.EnumMap{
			Value: &m,
			AllowedVals: psetter.AllowedVals{
				XOption: "a description of this option",
				YOption: "what this option means",
			},
		}, "help text")

	fmt.Println("Before parsing")
	for _, k := range keys {
		if v, ok := m[k]; ok {
			fmt.Printf("\tm[%s] = %v\n", k, v)
		}
	}

	ps.Parse([]string{"-my-map", "x"})

	fmt.Println("After  parsing")
	for _, k := range keys {
		if v, ok := m[k]; ok {
			fmt.Printf("\tm[%s] = %v\n", k, v)
		}
	}
	// Output:
	// Before parsing
	// After  parsing
	//	m[x] = true
}

// ExampleEnumMap_fixingInitialValue demonstrates how an initial value may be
// changed through the command line. That is, it is possible to change the
// value of a map entry to false as well as to true.
func ExampleEnumMap_fixingInitialValue() {
	ps := newPSetForTesting() // use paramset.NewOrDie()

	const (
		XOption = "x"
		YOption = "y"
	)

	m := map[string]bool{"x": true}
	keys := []string{XOption, YOption}

	ps.Add("my-map",
		psetter.EnumMap{
			Value: &m,
			AllowedVals: psetter.AllowedVals{
				XOption: "a description of this option",
				YOption: "what this option means",
			},
		}, "help text")

	fmt.Println("Before parsing")
	for _, k := range keys {
		if v, ok := m[k]; ok {
			fmt.Printf("\tm[%s] = %v\n", k, v)
		}
	}

	ps.Parse([]string{"-my-map", "x=false,y"})

	fmt.Println("After  parsing")
	for _, k := range keys {
		if v, ok := m[k]; ok {
			fmt.Printf("\tm[%s] = %v\n", k, v)
		}
	}
	// Output:
	// Before parsing
	//	m[x] = true
	// After  parsing
	//	m[x] = false
	//	m[y] = true
}

// ExampleEnumMap_hiddenMapEntries demonstrates the behaviour of the package
// when the AllowHiddenMapEntries flag is set. In this case the Value to be
// set has an entry with a key not in the allowed values but no error is
// reported. Note that it is not possible to set such a map value as the key
// will be rejected as invalid.
func ExampleEnumMap_hiddenMapEntries() {
	ps := newPSetForTesting() // use paramset.NewOrDie()

	const (
		XOption = "x"
		YOption = "y"
	)

	m := map[string]bool{"z": true}
	keys := []string{XOption, YOption, "z"}

	ps.Add("my-map",
		psetter.EnumMap{
			Value: &m,
			AllowedVals: psetter.AllowedVals{ // Note there's no 'z' value
				XOption: "a description of this option",
				YOption: "what this option means",
			},

			// Setting AllowHiddenMapEntries to true prevents the 'z' entry
			// causing an error
			AllowHiddenMapEntries: true,
		}, "help text")

	fmt.Println("Before parsing")
	for _, k := range keys {
		if v, ok := m[k]; ok {
			fmt.Printf("\tm[%s] = %v\n", k, v)
		}
	}

	ps.Parse([]string{"-my-map", "y"})

	fmt.Println("After  parsing")
	for _, k := range keys {
		if v, ok := m[k]; ok {
			fmt.Printf("\tm[%s] = %v\n", k, v)
		}
	}
	// Output:
	// Before parsing
	//	m[z] = true
	// After  parsing
	//	m[y] = true
	//	m[z] = true
}

// ExampleEnumMap_withBadKey demonstrates the behaviour of the package when
// an invalid setter is provided. In this case the Value to be set has an
// entry with a key not in the allowed values. Note that in production code
// you should not recover from the panic, instead you should fix the code
// that caused it.
func ExampleEnumMap_withBadKey() {
	defer func() { // For test purposes only - do not recover in live code
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // use paramset.NewOrDie()

	const (
		XOption = "x"
		YOption = "y"
	)

	m := map[string]bool{"z": true}

	// we expect this to panic because the map has an entry which is not in
	// the allowed values
	ps.Add("my-map",
		psetter.EnumMap{
			Value: &m,
			AllowedVals: psetter.AllowedVals{ // Note there's no 'z' value
				XOption: "a description of this option",
				YOption: "what this option means",
			},
		}, "help text")
	// Output:
	// panic
	// my-map: psetter.EnumMap Check failed: the map entry with key "z" is invalid - it is not in the allowed values map
}

// ExampleEnumMap_withNilValue demonstrates the behaviour of the package when
// an invalid setter is provided. In this case the Value to be set has not
// been initialised. Note that in production code you should not recover from
// the panic, instead you should fix the code that caused it.
func ExampleEnumMap_withNilValue() {
	defer func() { // For test purposes only - do not recover in live code
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // use paramset.NewOrDie()

	const (
		XOption = "x"
		YOption = "y"
	)

	// we expect this to panic because the map Value has not been initialised
	ps.Add("my-map",
		psetter.EnumMap{
			AllowedVals: psetter.AllowedVals{
				XOption: "a description of this option",
				YOption: "what this option means",
			},
		}, "help text")
	// Output:
	// panic
	// my-map: psetter.EnumMap Check failed: the Value to be set is nil
}
