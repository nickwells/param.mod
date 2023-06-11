package psetter_test

import (
	"fmt"
	"time"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

// ExampleDuration_basic demonstrates the use of a Duration setter.
func ExampleDuration_basic() {
	ps := newPSetForTesting() // use paramset.NewOrDie()

	var v time.Duration

	ps.Add("how-long", psetter.Duration{Value: &v}, "help text")

	fmt.Printf("Before parsing    v: %v\n", v)
	ps.Parse([]string{"-how-long", "1h"})
	fmt.Printf("After  parsing    v: %v\n", v)
	// Output:
	// Before parsing    v: 0s
	// After  parsing    v: 1h0m0s
}

// ExampleDuration_withPassingChecks demonstrates how to specify additional
// checks for a Duration value.
func ExampleDuration_withPassingChecks() {
	ps := newPSetForTesting() // use paramset.NewOrDie()

	var v time.Duration

	ps.Add("how-long",
		psetter.Duration{
			Value: &v,
			Checks: []check.Duration{
				check.ValBetween[time.Duration](time.Duration(0), 2*time.Hour),
			},
		},
		"help text")

	fmt.Printf("Before parsing    v: %v\n", v)
	ps.Parse([]string{"-how-long", "1h"})
	fmt.Printf("After  parsing    v: %v\n", v)
	// Output:
	// Before parsing    v: 0s
	// After  parsing    v: 1h0m0s
}

// ExampleDuration_withFailingChecks demonstrates how to specify additional
// checks for a Duration value and shows the error that you can expect to see
// if a value is supplied which fails any of the checks. Note that there is
// normally no need to examine the return from ps.Parse as the standard
// Helper will report any errors and abort the program.
func ExampleDuration_withFailingChecks() {
	ps := newPSetForTesting() // use paramset.NewOrDie()

	var v time.Duration

	ps.Add("how-long",
		psetter.Duration{
			Value: &v,
			Checks: []check.Duration{
				check.ValGT[time.Duration](2 * time.Hour),
			},
		},
		"help text")

	fmt.Printf("Before parsing    v: %v\n", v)
	// Parse the arguments. Note that the duration given (1 hour) is less
	// than the minimum value given by the check function.
	errMap := ps.Parse([]string{"-how-long", "1h"})
	// The check will fail so we expect to see errors reported
	logErrs(errMap)
	// There was an error with the parameter so the value will be unchanged
	fmt.Printf("After  parsing    v: %v\n", v)
	// Output:
	// Before parsing    v: 0s
	// Errors for: how-long
	//	: the value (1h0m0s) must be greater than 2h0m0s
	// At: [command line]: Supplied Parameter:2: "-how-long" "1h"
	// After  parsing    v: 0s
}

// ExampleDuration_withNilValue demonstrates the behaviour of the package
// when an invalid setter is provided. In this case the Value to be set has
// not been initialised. Note that in production code you should not recover
// from the panic, instead you should fix the code that caused it.
func ExampleDuration_withNilValue() {
	defer func() { // For test purposes only - do not recover in live code
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // use paramset.NewOrDie()

	// we expect this to panic because the Duration Value has not been
	// initialised
	ps.Add("how-long", psetter.Duration{}, "help text")
	// Output:
	// panic
	// how-long: psetter.Duration Check failed: the Value to be set is nil
}
