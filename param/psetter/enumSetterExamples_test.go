package psetter_test

import (
	"fmt"

	"github.com/nickwells/param.mod/v5/param/psetter"
)

// ExampleEnum_standard demonstrates the use of an Enum setter.
func ExampleEnum_standard() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	const (
		XOption = "x"
		YOption = "y"
	)

	s := XOption

	ps.Add("my-string",
		psetter.Enum{
			Value: &s,
			AllowedVals: psetter.AllowedVals{
				XOption: "a description of this option",
				YOption: "what this option means",
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
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	const (
		XOption = "x"
		YOption = "y"
	)

	s := XOption

	ps.Add("my-string",
		psetter.Enum{
			Value: &s,
			AllowedVals: psetter.AllowedVals{ // Note there's no 'z' value
				XOption: "a description of this option",
				YOption: "what this option means",
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
	// At: [command line]: Supplied Parameter:2: "-my-string" "z"
	// After  parsing
	//	s = x
}

// ExampleEnum_withNilValue demonstrates the behaviour of the package when an
// invalid setter is provided. In this case the Value to be set has not been
// initialised. Note that in production code you should not recover from the
// panic, instead you should fix the code that caused it.
func ExampleEnum_withNilValue() {
	defer func() { // For test purposes only - do not recover in live code
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // use paramset.NewOrPanic()

	const (
		XOption = "x"
		YOption = "y"
	)

	// we expect this to panic because the Value has not been initialised
	ps.Add("my-string",
		psetter.Enum{
			AllowedVals: psetter.AllowedVals{
				XOption: "a description of this option",
				YOption: "what this option means",
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
	defer func() { // For test purposes only - do not recover in live code
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // use paramset.NewOrPanic()

	const (
		XOption = "x"
		YOption = "y"
	)

	s := "z"

	// we expect this to panic because the Value has an invalid initial value
	ps.Add("my-string",
		psetter.Enum{
			Value: &s,
			AllowedVals: psetter.AllowedVals{ // Note there's no 'z' value
				XOption: "a description of this option",
				YOption: "what this option means",
			},
		}, "help text")
	// Output:
	// panic
	// my-string: psetter.Enum Check failed: the initial value (z) is not valid
}
