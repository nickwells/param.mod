package psetter_test

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

// ExampleString_standard demonstrates the use of a String setter
func ExampleString_standard() {
	ps := newPSetForTesting() // use paramset.NewOrDie()

	var s string

	ps.Add("my-string", psetter.String{Value: &s}, "help text")

	fmt.Printf("Before parsing: s = %q\n", s)
	ps.Parse([]string{"-my-string", "Hello, World!"})
	fmt.Printf("After  parsing: s = %q\n", s)
	// Output:
	// Before parsing: s = ""
	// After  parsing: s = "Hello, World!"
}

// ExampleString_withPassingChecks demonstrates how to add checks to be
// applied to the value.
func ExampleString_withPassingChecks() {
	ps := newPSetForTesting() // use paramset.NewOrDie()

	var s string

	ps.Add("my-string",
		psetter.String{
			Value: &s,
			Checks: []check.String{
				check.StringLength[string](check.ValGT(5)),
			},
		}, "help text")

	fmt.Printf("Before parsing: s = %q\n", s)
	ps.Parse([]string{"-my-string", "Hello, World!"})
	fmt.Printf("After  parsing: s = %q\n", s)
	// Output:
	// Before parsing: s = ""
	// After  parsing: s = "Hello, World!"
}

// ExampleString_withFailingChecks demonstrates how to add checks to be
// applied to the value. Note that there is normally no need to examine the
// return from ps.Parse as the standard Helper will report any errors and
// abort the program.
func ExampleString_withFailingChecks() {
	ps := newPSetForTesting() // use paramset.NewOrDie()

	var s string

	ps.Add("my-string",
		psetter.String{
			Value: &s,
			Checks: []check.String{
				check.StringLength[string](check.ValGT(5)),
			},
		}, "help text")

	fmt.Printf("Before parsing: s = %q\n", s)
	// Parse the arguments. Note that the string supplied is too short to
	// satisfy the check for this parameter.
	errMap := ps.Parse([]string{"-my-string", "Hi!"})
	// We expect to see an error reported.
	logErrs(errMap)
	// The value is unchanged due to the error.
	fmt.Printf("After  parsing: s = %q\n", s)
	// Output:
	// Before parsing: s = ""
	// Errors for: my-string
	//	: the length of the string (3) is incorrect: the value (3) must be greater than 5
	// At: [command line]: Supplied Parameter:2: "-my-string" "Hi!"
	// After  parsing: s = ""
}

// ExampleString_withNilValue demonstrates the behaviour of the package when
// an invalid setter is provided. In this case the Value to be set has not
// been initialised. Note that in production code you should not recover from
// the panic, instead you should fix the code that caused it.
func ExampleString_withNilValue() {
	defer func() { // For test purposes only - do not recover in live code
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // use paramset.NewOrDie()

	ps.Add("my-string", psetter.String{}, "help text")

	// Output:
	// panic
	// my-string: psetter.String Check failed: the Value to be set is nil
}
