package psetter_test

import (
	"fmt"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

// ExampleStrListAppender_standard demonstrates the use of a StrListAppender
func ExampleStrListAppender_standard() {
	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	ss := []string{"Hello"}

	ps.Add("next", psetter.StrListAppender{Value: &ss}, "help text")

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
	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	ss := []string{"Hello"}

	ps.Add("next",
		psetter.StrListAppender{
			Value: &ss,
			Checks: []check.String{
				check.StringLenGT(5),
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
	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	ss := []string{"Hello"}

	ps.Add("next",
		psetter.StrListAppender{
			Value: &ss,
			Checks: []check.String{
				check.StringLenLT(10),
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
	//	: the length of the value (13) must be less than 10 (at supplied parameters:4: -next my old friend)
	// After  parsing
	//	ss[0] = "Hello"
	//	ss[1] = "darkness"
}

// ExampleStrListAppender_withNilValue demonstrates the behaviour of the package
// when an invalid setter is provided. In this case the Value to be set has
// not been initialised. Note that in production code you should not recover
// from the panic, instead you should fix the code that caused it.
func ExampleStrListAppender_withNilValue() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	// we expect this to panic because the list Value has not been initialised
	ps.Add("my-list", psetter.StrListAppender{}, "help text")
	// Output:
	// panic
	// my-list: psetter.StrListAppender Check failed: the Value to be set is nil
}
