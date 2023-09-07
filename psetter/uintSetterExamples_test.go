package psetter_test

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v6/psetter"
)

// ExampleUint64_standard demonstrates the use of a Uint64 setter.
func ExampleUint64_standard() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	var i uint64

	ps.Add("my-uint",
		psetter.Uint[uint64]{
			Value: &i,
		}, "help text")

	fmt.Println("Before parsing")
	fmt.Printf("\ti = %d\n", i)
	ps.Parse([]string{"-my-uint", "1"})
	fmt.Println("After  parsing")
	fmt.Printf("\ti = %d\n", i)
	// Output:
	// Before parsing
	//	i = 0
	// After  parsing
	//	i = 1
}

// ExampleUint64_withPassingChecks demonstrates how to add checks to be
// applied to the value.
func ExampleUint64_withPassingChecks() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	var i uint64

	ps.Add("my-uint",
		psetter.Uint[uint64]{
			Value: &i,
			Checks: []check.ValCk[uint64]{
				check.ValGT[uint64](5),
			},
		}, "help text")

	fmt.Println("Before parsing")
	fmt.Printf("\ti = %d\n", i)
	ps.Parse([]string{"-my-uint", "6"})
	fmt.Println("After  parsing")
	fmt.Printf("\ti = %d\n", i)
	// Output:
	// Before parsing
	//	i = 0
	// After  parsing
	//	i = 6
}

// ExampleUint64_withFailingChecks demonstrates how to add checks to be
// applied to the value. Note that there is normally no need to examine the
// return from ps.Parse as the standard Helper will report any errors and
// abort the program.
func ExampleUint64_withFailingChecks() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	var i uint64

	ps.Add("my-uint",
		psetter.Uint[uint64]{
			Value: &i,
			Checks: []check.ValCk[uint64]{
				check.ValGT[uint64](5),
			},
		}, "help text")

	fmt.Println("Before parsing")
	fmt.Printf("\ti = %d\n", i)
	// Parse the arguments. We supply a uint value but note that it does not
	// satisfy the check for this parameter.
	errMap := ps.Parse([]string{"-my-uint", "1"})
	// We expect to see an error reported.
	logErrs(errMap)
	// The uint value is unchanged due to the error.
	fmt.Println("After  parsing")
	fmt.Printf("\ti = %d\n", i)
	// Output:
	// Before parsing
	//	i = 0
	// Errors for: my-uint
	//	: the value (1) must be greater than 5
	// At: [command line]: Supplied Parameter:2: "-my-uint" "1"
	// After  parsing
	//	i = 0
}

// ExampleUint64_withNilValue demonstrates the behaviour of the package when
// an invalid setter is provided. In this case the Value to be set has not
// been initialised. Note that in production code you should not recover from
// the panic, instead you should fix the code that caused it.
func ExampleUint64_withNilValue() {
	defer func() { // For test purposes only - do not recover in live code
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // use paramset.NewOrPanic()

	// we expect this to panic because the map Value has not been initialised
	ps.Add("my-uint", psetter.Uint[uint64]{}, "help text")
	// Output:
	// panic
	// my-uint: psetter.Uint[uint64] Check failed: the Value to be set is nil
}
