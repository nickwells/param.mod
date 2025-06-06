package psetter_test

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/param.mod/v6/psetter"
)

// ExamplePathname_standard demonstrates the use of a Pathname setter.
func ExamplePathname_standard() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

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
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	var pathname string

	ps.Add("my-pathname",
		psetter.Pathname{
			Value:       &pathname,
			Expectation: filecheck.IsNew(),
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
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	var pathname string

	ps.Add("my-pathname",
		psetter.Pathname{
			Value:       &pathname,
			Expectation: filecheck.FileExists(),
		},
		"help text")

	fmt.Printf("Before parsing    pathname: %q\n", pathname)

	errMap := ps.Parse([]string{"-my-pathname", "testdata/noSuchFile.go"})
	// We expect to see an error reported. Note that the Pathname setter
	// suggests an alternative file from the same directory in the error
	// message.
	logErrs(errMap)
	// There was an error with the parameter so the value will be unchanged
	fmt.Printf("After  parsing    pathname: %q\n", pathname)
	// Output:
	// Before parsing    pathname: ""
	// Errors for: my-pathname
	//	: path: "testdata/noSuchFile.go": should exist but does not; "testdata" exists but "noSuchFile.go" does not, did you mean "testdata/SuchFile.go"?
	// At: [command line]: Supplied Parameter:2: "-my-pathname" "testdata/noSuchFile.go"
	// After  parsing    pathname: ""
}

// ExamplePathname_withPassingChecks demonstrates the use of a Pathname
// setter which has Checks. Note that it also has the Expectation set.
func ExamplePathname_withPassingChecks() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	var pathname string

	ps.Add("my-pathname",
		psetter.Pathname{
			Value:       &pathname,
			Expectation: filecheck.IsNew(),
			Checks: []check.String{
				check.StringHasSuffix[string](".go"),
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

// ExamplePathname_withNilValue demonstrates the behaviour of the package
// when an invalid setter is provided. In this case the Value to be set has
// not been initialised. Note that in production code you should not recover
// from the panic, instead you should fix the code that caused it.
func ExamplePathname_withNilValue() {
	defer func() { // For test purposes only - do not recover in live code
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // use paramset.NewOrPanic()

	// we expect this to panic because the bool Value has not been initialised
	ps.Add("do-this", psetter.Pathname{}, "help text")
	// Output:
	// panic
	// do-this: psetter.Pathname Check failed: the Value to be set is nil
}
