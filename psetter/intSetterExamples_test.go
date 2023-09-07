package psetter_test

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v6/psetter"
)

// ExampleInt64_standard demonstrates the use of a Int64 setter.
func ExampleInt64_standard() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	var i int64

	ps.Add("my-int",
		psetter.Int[int64]{
			Value: &i,
		}, "help text")

	fmt.Println("Before parsing")
	fmt.Printf("\ti = %d\n", i)
	ps.Parse([]string{"-my-int", "1"})
	fmt.Println("After  parsing")
	fmt.Printf("\ti = %d\n", i)
	// Output:
	// Before parsing
	//	i = 0
	// After  parsing
	//	i = 1
}

// ExampleInt64_withPassingChecks demonstrates how to add checks to be
// applied to the value.
func ExampleInt64_withPassingChecks() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	var i int64

	ps.Add("my-int",
		psetter.Int[int64]{
			Value: &i,
			Checks: []check.Int64{
				check.ValGT[int64](5),
			},
		}, "help text")

	fmt.Println("Before parsing")
	fmt.Printf("\ti = %d\n", i)
	ps.Parse([]string{"-my-int", "6"})
	fmt.Println("After  parsing")
	fmt.Printf("\ti = %d\n", i)
	// Output:
	// Before parsing
	//	i = 0
	// After  parsing
	//	i = 6
}

// ExampleInt64_withFailingChecks demonstrates how to add checks to be
// applied to the value. Note that there is normally no need to examine the
// return from ps.Parse as the standard Helper will report any errors and
// abort the program.
func ExampleInt64_withFailingChecks() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	var i int64

	ps.Add("my-int",
		psetter.Int[int64]{
			Value: &i,
			Checks: []check.Int64{
				check.ValGT[int64](5),
			},
		}, "help text")

	fmt.Println("Before parsing")
	fmt.Printf("\ti = %d\n", i)
	// Parse the arguments. We supply a int value but note that it does not
	// satisfy the check for this parameter.
	errMap := ps.Parse([]string{"-my-int", "1"})
	// We expect to see an error reported.
	logErrs(errMap)
	// The int value is unchanged due to the error.
	fmt.Println("After  parsing")
	fmt.Printf("\ti = %d\n", i)
	// Output:
	// Before parsing
	//	i = 0
	// Errors for: my-int
	//	: the value (1) must be greater than 5
	// At: [command line]: Supplied Parameter:2: "-my-int" "1"
	// After  parsing
	//	i = 0
}

// ExampleInt64_withNilValue demonstrates the behaviour of the package when
// an invalid setter is provided. In this case the Value to be set has not
// been initialised. Note that in production code you should not recover from
// the panic, instead you should fix the code that caused it.
func ExampleInt64_withNilValue() {
	defer func() { // For test purposes only - do not recover in live code
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // use paramset.NewOrPanic()

	// we expect this to panic because the map Value has not been initialised
	ps.Add("my-int", psetter.Int[int64]{}, "help text")
	// Output:
	// panic
	// my-int: psetter.Int[int64] Check failed: the Value to be set is nil
}
