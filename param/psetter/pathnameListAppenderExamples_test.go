package psetter_test

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

// ExamplePathnameListAppender_standard demonstrates the use of a PathnameListAppender
func ExamplePathnameListAppender_standard() {
	ps := newPSetForTesting() // use paramset.NewOrDie()

	ss := []string{"testdata/pathname/nonesuch.go"}

	ps.Add("next",
		psetter.PathnameListAppender{
			Value:       &ss,
			Expectation: filecheck.IsNew(),
			Checks: []check.String{
				check.StringHasSuffix[string](".go"),
			},
		},
		"help text")

	fmt.Println("Before parsing")
	for i, v := range ss {
		fmt.Printf("\tss[%d] = %q\n", i, v)
	}
	ps.Parse([]string{"-next", "testdata/pathname/nonesuch2.go"})
	fmt.Println("After  parsing")
	for i, v := range ss {
		fmt.Printf("\tss[%d] = %q\n", i, v)
	}
	// Output:
	// Before parsing
	//	ss[0] = "testdata/pathname/nonesuch.go"
	// After  parsing
	//	ss[0] = "testdata/pathname/nonesuch.go"
	//	ss[1] = "testdata/pathname/nonesuch2.go"
}

// ExamplePathnameListAppender_withNilValue demonstrates the behaviour of the
// package when an invalid setter is provided. In this case the Value to be
// set has not been initialised. Note that in production code you should not
// recover from the panic, instead you should fix the code that caused it.
func ExamplePathnameListAppender_withNilValue() {
	defer func() { // For test purposes only - do not recover in live code
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // use paramset.NewOrDie()

	// we expect this to panic because the list Value has not been initialised
	ps.Add("my-list", psetter.PathnameListAppender{}, "help text")
	// Output:
	// panic
	// my-list: psetter.PathnameListAppender Check failed: the Value to be set is nil
}
