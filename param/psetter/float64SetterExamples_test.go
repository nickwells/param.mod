package psetter_test

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

// ExampleFloat64_standard demonstrates the use of a Float64 setter.
func ExampleFloat64_standard() {
	ps := newPSetForTesting() // use paramset.NewOrDie()

	var f float64

	ps.Add("my-float",
		psetter.Float64{
			Value: &f,
		}, "help text")

	fmt.Println("Before parsing")
	fmt.Printf("\tf = %5.3f\n", f)
	ps.Parse([]string{"-my-float", "1.23"})
	fmt.Println("After  parsing")
	fmt.Printf("\tf = %5.3f\n", f)
	// Output:
	// Before parsing
	//	f = 0.000
	// After  parsing
	//	f = 1.230
}

// ExampleFloat64_withPassingChecks demonstrates how to add checks to be
// applied to the value.
func ExampleFloat64_withPassingChecks() {
	ps := newPSetForTesting() // use paramset.NewOrDie()

	var f float64

	ps.Add("my-float",
		psetter.Float64{
			Value: &f,
			Checks: []check.Float64{
				check.ValGT[float64](5.0),
			},
		}, "help text")

	fmt.Println("Before parsing")
	fmt.Printf("\tf = %5.3f\n", f)
	ps.Parse([]string{"-my-float", "6.23"})
	fmt.Println("After  parsing")
	fmt.Printf("\tf = %5.3f\n", f)
	// Output:
	// Before parsing
	//	f = 0.000
	// After  parsing
	//	f = 6.230
}

// ExampleFloat64_withFailingChecks demonstrates how to add checks to be
// applied to the value. Note that there is normally no need
// to examine the return from ps.Parse as the standard Helper will report any
// errors and abort the program.
func ExampleFloat64_withFailingChecks() {
	ps := newPSetForTesting() // use paramset.NewOrDie()

	var f float64

	ps.Add("my-float",
		psetter.Float64{
			Value: &f,
			Checks: []check.Float64{
				check.ValGT[float64](5.0),
			},
		}, "help text")

	fmt.Println("Before parsing")
	fmt.Printf("\tf = %5.3f\n", f)
	// Parse the arguments. We supply a float value but note that it does not
	// satisfy the check for this parameter.
	errMap := ps.Parse([]string{"-my-float", "1.23"})
	// We expect to see an error reported.
	logErrs(errMap)
	// The float value is unchanged due to the error.
	fmt.Println("After  parsing")
	fmt.Printf("\tf = %5.3f\n", f)
	// Output:
	// Before parsing
	//	f = 0.000
	// Errors for: my-float
	//	: the value (1.23) must be greater than 5
	// At: [command line]: Supplied Parameter:2: -my-float 1.23
	// After  parsing
	//	f = 0.000
}

// ExampleFloat64_withNilValue demonstrates the behaviour of the package when
// an invalid setter is provided. In this case the Value to be set has not
// been initialised. Note that in production code you should not recover from
// the panic, instead you should fix the code that caused it.
func ExampleFloat64_withNilValue() {
	defer func() { // For test purposes only - do not recover in live code
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // use paramset.NewOrDie()

	// we expect this to panic because the map Value has not been initialised
	ps.Add("my-float", psetter.Float64{}, "help text")
	// Output:
	// panic
	// my-float: psetter.Float64 Check failed: the Value to be set is nil
}
