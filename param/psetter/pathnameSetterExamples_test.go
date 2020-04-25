package psetter_test

import (
	"fmt"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/param.mod/v4/param/psetter"
)

// ExamplePathname_standard demonstrates the use of a Pathname setter.
func ExamplePathname_standard() {
	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	var pathname string

	ps.Add("my-pathname", psetter.Pathname{Value: &pathname}, "help text")

	fmt.Printf("Before parsing    pathname: %q\n", pathname)
	ps.Parse([]string{"-my-pathname", "testdata/noSuchFile.go"})
	fmt.Printf("After  parsing    pathname: %q\n", pathname)
	// Output:
	// Before parsing    pathname: ""
	// After  parsing    pathname: "testdata/noSuchFile.go"
}

// ExamplePathname_withPassingExpectation demonstrates the use of a Pathname
// setter which has the Expectation set.
func ExamplePathname_withPassingExpectation() {
	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	var pathname string

	ps.Add("my-pathname",
		psetter.Pathname{
			Value: &pathname,
			Expectation: filecheck.Provisos{
				Existence: filecheck.MustNotExist,
			},
		},
		"help text")

	fmt.Printf("Before parsing    pathname: %q\n", pathname)
	ps.Parse([]string{"-my-pathname", "testdata/noSuchFile.go"})
	fmt.Printf("After  parsing    pathname: %q\n", pathname)
	// Output:
	// Before parsing    pathname: ""
	// After  parsing    pathname: "testdata/noSuchFile.go"
}

// ExamplePathname_withFailingExpectation demonstrates the use of a Pathname
// setter which has the Expectation set. Note that there is normally no need
// to examine the return from ps.Parse as the standard Helper will report any
// errors and abort the program.
func ExamplePathname_withFailingExpectation() {
	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	var pathname string

	ps.Add("my-pathname",
		psetter.Pathname{
			Value: &pathname,
			Expectation: filecheck.Provisos{
				Existence: filecheck.MustExist,
			},
		},
		"help text")

	fmt.Printf("Before parsing    pathname: %q\n", pathname)
	errMap := ps.Parse([]string{"-my-pathname", "testdata/noSuchFile.go"})
	// We expect to see an error reported.
	logErrs(errMap)
	// There was an error with the parameter so the value will be unchanged
	fmt.Printf("After  parsing    pathname: %q\n", pathname)
	// Output:
	// Before parsing    pathname: ""
	// Errors for: my-pathname
	//	: error with parameter: path: "testdata/noSuchFile.go" should exist but doesn't (at supplied parameters:2: -my-pathname testdata/noSuchFile.go)
	// After  parsing    pathname: ""
}

// ExamplePathname_withPassingChecks demonstrates the use of a Pathname
// setter which has Checks. Note that it also has the Expectation set.
func ExamplePathname_withPassingChecks() {
	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	var pathname string

	ps.Add("my-pathname",
		psetter.Pathname{
			Value: &pathname,
			Expectation: filecheck.Provisos{
				Existence: filecheck.MustNotExist,
			},
			Checks: []check.String{
				check.StringHasSuffix(".go"),
			},
		},
		"help text")

	fmt.Printf("Before parsing    pathname: %q\n", pathname)
	ps.Parse([]string{"-my-pathname", "testdata/noSuchFile.go"})
	fmt.Printf("After  parsing    pathname: %q\n", pathname)
	// Output:
	// Before parsing    pathname: ""
	// After  parsing    pathname: "testdata/noSuchFile.go"
}

// ExamplePathname_withNilValue demonstrates the behaviour of the package when an
// invalid setter is provided. In this case the Value to be set has not been
// initialised. Note that in production code you should not recover
// from the panic, instead you should fix the code that caused it.
func ExamplePathname_withNilValue() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	// we expect this to panic because the bool Value has not been initialised
	ps.Add("do-this", psetter.Pathname{}, "help text")
	// Output:
	// panic
	// do-this: psetter.Pathname Check failed: the Value to be set is nil
}
