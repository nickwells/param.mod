package psetter_test

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v6/psetter"
)

// ExampleStrList_standard demonstrates the use of a StrList setter
func ExampleStrList_standard() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	var ss []string

	ps.Add("my-list",
		psetter.StrList[string]{
			Value: &ss,
		}, "help text")

	fmt.Println("Before parsing")

	for i, v := range ss {
		fmt.Printf("\tss[%d] = %q\n", i, v)
	}

	ps.Parse([]string{"-my-list", "x,y"})
	fmt.Println("After  parsing")

	for i, v := range ss {
		fmt.Printf("\tss[%d] = %q\n", i, v)
	}
	// Output:
	// Before parsing
	// After  parsing
	//	ss[0] = "x"
	//	ss[1] = "y"
}

// ExampleStrList_withPassingChecks demonstrates how you can specify
// additional checks to be applied to the passed arguments before the value
// is set.
func ExampleStrList_withPassingChecks() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	var ss []string

	ps.Add("my-list",
		psetter.StrList[string]{
			Value: &ss,
			Checks: []check.StringSlice{
				check.SliceLength[[]string](check.ValEQ(2)),
			},
		}, "help text")

	fmt.Println("Before parsing")

	for i, v := range ss {
		fmt.Printf("\tss[%d] = %q\n", i, v)
	}

	ps.Parse([]string{"-my-list", "x,y"})
	fmt.Println("After  parsing")

	for i, v := range ss {
		fmt.Printf("\tss[%d] = %q\n", i, v)
	}
	// Output:
	// Before parsing
	// After  parsing
	//	ss[0] = "x"
	//	ss[1] = "y"
}

// ExampleStrList_withFailingChecks demonstrates the behaviour of the
// package when an invalid value is given. In this case the resulting list is
// not of the required length. It demonstrates the checks that can be
// supplied to ensure that the resulting list is as expected. Note that there
// is normally no need to examine the return from ps.Parse as the standard
// Helper will report any errors and abort the program.
func ExampleStrList_withFailingChecks() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	var ss []string

	ps.Add("my-list",
		psetter.StrList[string]{
			Value: &ss,
			Checks: []check.StringSlice{
				check.SliceLength[[]string, string](check.ValEQ(2)),
			},
		}, "help text")

	fmt.Println("Before parsing")

	for i, v := range ss {
		fmt.Printf("\tss[%d] = %q\n", i, v)
	}

	// Parse the arguments. We supply a list of strings, each of which is
	// allowed. The resulting slice is of the wrong length.
	errMap := ps.Parse([]string{"-my-list", "x"})

	// We expect to see an error reported.
	logErrs(errMap)

	// The slice of strings is unchanged due to the error.
	fmt.Println("After  parsing")

	for i, v := range ss {
		fmt.Printf("\tss[%d] = %q\n", i, v)
	}
	// Output:
	// Before parsing
	// Errors for: my-list
	//	: the length of the list (1) is incorrect: the value (1) must equal 2
	// At: [command line]: Supplied Parameter:2: "-my-list" "x"
	// After  parsing
}

// ExampleStrList_withNilValue demonstrates the behaviour of the package
// when an invalid setter is provided. In this case the Value to be set has
// not been initialised. Note that in production code you should not recover
// from the panic, instead you should fix the code that caused it.
func ExampleStrList_withNilValue() {
	defer func() { // For test purposes only - do not recover in live code
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // use paramset.NewOrPanic()

	// we expect this to panic because the list Value has not been initialised
	ps.Add("my-list", psetter.StrList[string]{}, "help text")
	// Output:
	// panic
	// my-list: psetter.StrList[string] Check failed: the Value to be set is nil
}
