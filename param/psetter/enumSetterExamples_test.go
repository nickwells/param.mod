package psetter_test

import (
	"fmt"

	"github.com/nickwells/param.mod/v5/param/psetter"
)

// ExampleEnum_standard demonstrates the use of an Enum setter.
func ExampleEnum_standard() {
	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	s := "x"

	ps.Add("my-string",
		psetter.Enum{
			Value: &s,
			AllowedVals: psetter.AllowedVals{
				"x": "X",
				"y": "Y",
			},
		}, "help text")

	fmt.Println("Before parsing")
	fmt.Printf("\ts = %s\n", s)
	ps.Parse([]string{"-my-string", "y"})
	fmt.Println("After  parsing")
	fmt.Printf("\ts = %s\n", s)
	// Output:
	// Before parsing
	//	s = x
	// After  parsing
	//	s = y
}

// ExampleEnum_withBadVal demonstrates the behaviour when a value not given
// in the AllowedValues is passed as a parameter. Note that there is normally
// no need to examine the return from ps.Parse as the standard Helper will
// report any errors and abort the program.
func ExampleEnum_withBadVal() {
	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	s := "x"

	ps.Add("my-string",
		psetter.Enum{
			Value: &s,
			AllowedVals: psetter.AllowedVals{
				"x": "X",
				"y": "Y",
			},
		}, "help text")

	fmt.Println("Before parsing")
	fmt.Printf("\ts = %s\n", s)
	// Parse the arguments. We supply a value but note that it is not in the
	// list of allowed values.
	errMap := ps.Parse([]string{"-my-string", "z"})
	// We expect to see an error reported.
	logErrs(errMap)
	// The value is unchanged due to the error.
	fmt.Println("After  parsing")
	fmt.Printf("\ts = %s\n", s)
	// Output:
	// Before parsing
	//	s = x
	// Errors for: my-string
	//	: value not allowed: "z"
	// At: [command line]: Supplied Parameter:2: -my-string z
	// After  parsing
	//	s = x
}

// ExampleEnum_withNilValue demonstrates the behaviour of the package when an
// invalid setter is provided. In this case the Value to be set has not been
// initialised. Note that in production code you should not recover from the
// panic, instead you should fix the code that caused it.
func ExampleEnum_withNilValue() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	// we expect this to panic because the Value has not been initialised
	ps.Add("my-string",
		psetter.Enum{
			AllowedVals: psetter.AllowedVals{
				"x": "X",
				"y": "Y",
			},
		}, "help text")
	// Output:
	// panic
	// my-string: psetter.Enum Check failed: the Value to be set is nil
}

// ExampleEnum_withBadInitialValue demonstrates the behaviour of the package
// when the initial value is invalid (not in the list of allowed
// values). Note that in production code you should not recover from the
// panic, instead you should fix the code that caused it.
func ExampleEnum_withBadInitialValue() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	s := "z"

	// we expect this to panic because the Value has an invalid initial value
	ps.Add("my-string",
		psetter.Enum{
			Value: &s,
			AllowedVals: psetter.AllowedVals{
				"x": "X",
				"y": "Y",
			},
		}, "help text")
	// Output:
	// panic
	// my-string: psetter.Enum Check failed: the initial value (z) is not valid
}
