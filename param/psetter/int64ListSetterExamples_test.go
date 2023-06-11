package psetter_test

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

// ExampleInt64List_standard demonstrates the use of a Int64List setter.
func ExampleInt64List_standard() {
	ps := newPSetForTesting() // use paramset.NewOrDie()

	il := []int64{42}

	ps.Add("my-ints",
		psetter.Int64List{
			Value: &il,
		}, "help text")

	fmt.Println("Before parsing")
	for i, v := range il {
		fmt.Printf("\til[%d] = %d\n", i, v)
	}
	ps.Parse([]string{"-my-ints", "1,23"})
	fmt.Println("After  parsing")
	for i, v := range il {
		fmt.Printf("\til[%d] = %d\n", i, v)
	}
	// Output:
	// Before parsing
	//	il[0] = 42
	// After  parsing
	//	il[0] = 1
	//	il[1] = 23
}

// ExampleInt64List_withPassingChecks demonstrates how to add checks to be
// applied to the value.
func ExampleInt64List_withPassingChecks() {
	ps := newPSetForTesting() // use paramset.NewOrDie()

	il := []int64{42}

	ps.Add("my-ints",
		psetter.Int64List{
			Value: &il,
			Checks: []check.Int64Slice{
				check.SliceAll[[]int64, int64](check.ValGT[int64](5)),
			},
		}, "help text")

	fmt.Println("Before parsing")
	for i, v := range il {
		fmt.Printf("\til[%d] = %d\n", i, v)
	}
	ps.Parse([]string{"-my-ints", "6,23"})
	fmt.Println("After  parsing")
	for i, v := range il {
		fmt.Printf("\til[%d] = %d\n", i, v)
	}
	// Output:
	// Before parsing
	//	il[0] = 42
	// After  parsing
	//	il[0] = 6
	//	il[1] = 23
}

// ExampleInt64List_withFailingChecks demonstrates how to add checks to be
// applied to the value. Note that there is normally no need
// to examine the return from ps.Parse as the standard Helper will report any
// errors and abort the program.
func ExampleInt64List_withFailingChecks() {
	ps := newPSetForTesting() // use paramset.NewOrDie()

	il := []int64{42}

	ps.Add("my-ints",
		psetter.Int64List{
			Value: &il,
			Checks: []check.Int64Slice{
				check.SliceAll[[]int64, int64](check.ValGT[int64](5)),
			},
		}, "help text")

	fmt.Println("Before parsing")
	for i, v := range il {
		fmt.Printf("\til[%d] = %d\n", i, v)
	}
	// Parse the arguments. We supply a float value but note that it does not
	// satisfy the check for this parameter.
	errMap := ps.Parse([]string{"-my-ints", "1,23"})
	// We expect to see an error reported.
	logErrs(errMap)
	// The float value is unchanged due to the error.
	fmt.Println("After  parsing")
	for i, v := range il {
		fmt.Printf("\til[%d] = %d\n", i, v)
	}
	// Output:
	// Before parsing
	//	il[0] = 42
	// Errors for: my-ints
	//	: list entry: 0 (1) does not pass the test: the value (1) must be greater than 5
	// At: [command line]: Supplied Parameter:2: "-my-ints" "1,23"
	// After  parsing
	//	il[0] = 42
}

// ExampleInt64List_withNilValue demonstrates the behaviour of the package when
// an invalid setter is provided. In this case the Value to be set has not
// been initialised. Note that in production code you should not recover from
// the panic, instead you should fix the code that caused it.
func ExampleInt64List_withNilValue() {
	defer func() { // For test purposes only - do not recover in live code
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // use paramset.NewOrDie()

	// we expect this to panic because the map Value has not been initialised
	ps.Add("my-ints", psetter.Int64List{}, "help text")
	// Output:
	// panic
	// my-ints: psetter.Int64List Check failed: the Value to be set is nil
}
