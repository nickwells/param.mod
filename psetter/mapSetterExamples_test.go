package psetter_test

import (
	"fmt"

	"github.com/nickwells/param.mod/v6/psetter"
)

// ExampleMap_standard demonstrates the use of an Map setter.
func ExampleMap_standard() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	var m map[string]bool

	keys := []string{"x", "y"}

	ps.Add("my-map", psetter.Map[string]{Value: &m}, "help text")

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

// ExampleMap_fixingInitialValue demonstrates how an initial value may be
// changed through the command line. That is, it is possible to change the
// value of a map entry to false as well as to true.
func ExampleMap_fixingInitialValue() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	m := map[string]bool{"x": true}
	keys := []string{"x", "y"}

	ps.Add("my-map", psetter.Map[string]{Value: &m}, "help text")

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

// ExampleMap_withNilValue demonstrates the behaviour of the package when
// an invalid setter is provided. In this case the Value to be set has not
// been initialised. Note that in production code you should not recover from
// the panic, instead you should fix the code that caused it.
func ExampleMap_withNilValue() {
	defer func() { // For test purposes only - do not recover in live code
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // use paramset.NewOrPanic()

	// we expect this to panic because the map Value has not been initialised
	ps.Add("my-map", psetter.Map[string]{}, "help text")
	// Output:
	// panic
	// my-map: psetter.Map[string] Check failed: the Value to be set is nil
}
