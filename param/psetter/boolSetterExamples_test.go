package psetter_test

import (
	"fmt"

	"github.com/nickwells/param.mod/v4/param/psetter"
)

// ExampleBool_standard demonstrates the use of a Bool setter.
func ExampleBool_standard() {
	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	var p1 bool

	ps.Add("do-this", psetter.Bool{Value: &p1}, "help text")

	fmt.Printf("Before parsing    p1: %v\n", p1)
	ps.Parse([]string{"-do-this"})
	fmt.Printf("After  parsing    p1: %v\n", p1)
	// Output:
	// Before parsing    p1: false
	// After  parsing    p1: true
}

// ExampleBool_inverted demonstrates the use of a Bool setter with the Invert
// flag set to true. The standard behaviour will set the value to true when
// no explicit value is given but with this flag set the value is set to
// false. Any value given is inverted. This is useful for turning off some
// default behaviour rather than turning it on as the standard action of this
// setter would do.
func ExampleBool_inverted() {
	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	var p1 = true

	ps.Add("dont-do-this", psetter.Bool{Value: &p1, Invert: true}, "help text")

	fmt.Printf("Before parsing    p1: %v\n", p1)
	ps.Parse([]string{"-dont-do-this"})
	fmt.Printf("After  parsing    p1: %v\n", p1)
	// Output:
	// Before parsing    p1: true
	// After  parsing    p1: false
}

// ExampleBool_withValue demonstrates the use of a Bool setter showing how
// the value of the flag can be set to an explicit value by passing the value
// required after the parameter (following an "=").
func ExampleBool_withValue() {
	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	var p1 = true

	ps.Add("do-this", psetter.Bool{Value: &p1}, "help text")

	fmt.Printf("Before parsing    p1: %v\n", p1)
	ps.Parse([]string{"-do-this=false"})
	fmt.Printf("After  parsing    p1: %v\n", p1)
	// Output:
	// Before parsing    p1: true
	// After  parsing    p1: false
}

// ExampleBool_withBadValue demonstrates the use of a Bool setter showing the
// behaviour when an argument is supplied that cannot be translated into a
// bool value. Note that there is normally no need to examine the return from
// ps.Parse as the standard Helper will report any errors and abort the
// program.
func ExampleBool_withBadValue() {
	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	var p1 = true

	ps.Add("do-this", psetter.Bool{Value: &p1}, "help text")

	fmt.Printf("Before parsing    p1: %v\n", p1)
	// Parse the arguments. Note that the value after the '=' cannot be
	// translated into a bool value.
	errMap := ps.Parse([]string{"-do-this=blah"})
	// We expect to see an error reported.
	logErrs(errMap)
	// There was an error with the parameter so the value will be unchanged
	fmt.Printf("After  parsing    p1: %v\n", p1)
	// Output:
	// Before parsing    p1: true
	// Errors for: do-this
	//	: cannot interpret 'blah' as either true or false (at supplied parameters:1: -do-this=blah)
	// After  parsing    p1: true
}

// ExampleBool_withNilValue demonstrates the behaviour of the package when an
// invalid setter is provided. In this case the Value to be set has not been
// initialised. Note that in production code you should not recover
// from the panic, instead you should fix the code that caused it.
func ExampleBool_withNilValue() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	// we expect this to panic because the bool Value has not been initialised
	ps.Add("do-this", psetter.Bool{}, "help text")
	// Output:
	// panic
	// do-this: psetter.Bool Check failed: the Value to be set is nil
}
