package psetter_test

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v6/psetter"
)

// ExampleIntList_standard demonstrates the use of a IntList setter.
func ExampleIntList_standard() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	il := []int64{42}

	ps.Add("my-ints",
		psetter.IntList[int64]{
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

// ExampleIntList_withPassingChecks demonstrates how to add checks to be
// applied to the value.
func ExampleIntList_withPassingChecks() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	il := []int64{42}

	ps.Add("my-ints",
		psetter.IntList[int64]{
			Value: &il,
			Checks: []check.ValCk[[]int64]{
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

// ExampleIntList_withFailingChecks demonstrates how to add checks to be
// applied to the value. Note that there is normally no need
// to examine the return from ps.Parse as the standard Helper will report any
// errors and abort the program.
func ExampleIntList_withFailingChecks() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	il := []int64{42}

	ps.Add("my-ints",
		psetter.IntList[int64]{
			Value: &il,
			Checks: []check.ValCk[[]int64]{
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

// ExampleIntList_withNilValue demonstrates the behaviour of the package when
// an invalid setter is provided. In this case the Value to be set has not
// been initialised. Note that in production code you should not recover from
// the panic, instead you should fix the code that caused it.
func ExampleIntList_withNilValue() {
	defer func() { // For test purposes only - do not recover in live code
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // use paramset.NewOrPanic()

	// we expect this to panic because the map Value has not been initialised
	ps.Add("my-ints", psetter.IntList[int64]{}, "help text")
	// Output:
	// panic
	// my-ints: psetter.IntList[int64] Check failed: the Value to be set is nil
}
