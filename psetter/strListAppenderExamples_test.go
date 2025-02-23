package psetter_test

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v6/psetter"
)

// ExampleStrListAppender_standard demonstrates the use of a StrListAppender
func ExampleStrListAppender_standard() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	ss := []string{"Hello"}

	ps.Add("next", psetter.StrListAppender[string]{Value: &ss}, "help text")

	fmt.Println("Before parsing")

	for i, v := range ss {
		fmt.Printf("\tss[%d] = %q\n", i, v)
	}

	ps.Parse([]string{"-next", "darkness", "-next", "my old friend"})
	fmt.Println("After  parsing")

	for i, v := range ss {
		fmt.Printf("\tss[%d] = %q\n", i, v)
	}
	// Output:
	// Before parsing
	//	ss[0] = "Hello"
	// After  parsing
	//	ss[0] = "Hello"
	//	ss[1] = "darkness"
	//	ss[2] = "my old friend"
}

// ExampleStrListAppender_withPassingChecks demonstrates how to add checks to be
// applied to the value.
func ExampleStrListAppender_withPassingChecks() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	ss := []string{"Hello"}

	ps.Add("next",
		psetter.StrListAppender[string]{
			Value: &ss,
			Checks: []check.String{
				check.StringLength[string](check.ValGT(5)),
			},
		},
		"help text")

	fmt.Println("Before parsing")

	for i, v := range ss {
		fmt.Printf("\tss[%d] = %q\n", i, v)
	}

	ps.Parse([]string{"-next", "darkness", "-next", "my old friend"})
	fmt.Println("After  parsing")

	for i, v := range ss {
		fmt.Printf("\tss[%d] = %q\n", i, v)
	}
	// Output:
	// Before parsing
	//	ss[0] = "Hello"
	// After  parsing
	//	ss[0] = "Hello"
	//	ss[1] = "darkness"
	//	ss[2] = "my old friend"
}

// ExampleStrListAppender_withFailingChecks demonstrates how to add checks to be
// applied to the value. Note that there is normally no need to examine the
// return from ps.Parse as the standard Helper will report any errors and
// abort the program.
func ExampleStrListAppender_withFailingChecks() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	ss := []string{"Hello"}

	ps.Add("next",
		psetter.StrListAppender[string]{
			Value: &ss,
			Checks: []check.String{
				check.StringLength[string](check.ValLT(10)),
			},
		},
		"help text")

	fmt.Println("Before parsing")

	for i, v := range ss {
		fmt.Printf("\tss[%d] = %q\n", i, v)
	}

	errMap := ps.Parse([]string{"-next", "darkness", "-next", "my old friend"})

	// We expect to see an error reported.
	logErrs(errMap)

	// The value does not include the second parameter due to the error.
	fmt.Println("After  parsing")

	for i, v := range ss {
		fmt.Printf("\tss[%d] = %q\n", i, v)
	}
	// Output:
	// Before parsing
	//	ss[0] = "Hello"
	// Errors for: next
	//	: the length of the string (13) is incorrect: the value (13) must be less than 10
	// At: [command line]: Supplied Parameter:4: "-next" "my old friend"
	// After  parsing
	//	ss[0] = "Hello"
	//	ss[1] = "darkness"
}

// ExampleStrListAppender_withNilValue demonstrates the behaviour of the package
// when an invalid setter is provided. In this case the Value to be set has
// not been initialised. Note that in production code you should not recover
// from the panic, instead you should fix the code that caused it.
func ExampleStrListAppender_withNilValue() {
	defer func() { // For test purposes only - do not recover in live code
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // use paramset.NewOrPanic()

	// we expect this to panic because the list Value has not been initialised
	ps.Add("my-list", psetter.StrListAppender[string]{}, "help text")
	// Output:
	// panic
	// my-list: psetter.StrListAppender[string] Check failed: the Value to be set is nil
}
