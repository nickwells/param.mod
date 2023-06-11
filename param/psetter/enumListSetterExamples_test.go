package psetter_test

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

// ExampleEnumList_standard demonstrates the use of an EnumList setter
func ExampleEnumList_standard() {
	ps := newPSetForTesting() // use paramset.NewOrDie()

	const (
		XOption = "x"
		YOption = "y"
	)

	var ss []string

	ps.Add("my-list",
		psetter.EnumList{
			Value: &ss,
			AllowedVals: psetter.AllowedVals{
				XOption: "a description of this option",
				YOption: "what this option means",
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

// ExampleEnumList_withBadVals demonstrates the behaviour when a value not
// given in the AllowedValues is passed. Note that there is normally no need
// to examine the return from ps.Parse as the standard Helper will report any
// errors and abort the program.
func ExampleEnumList_withBadVals() {
	ps := newPSetForTesting() // use paramset.NewOrDie()

	const (
		XOption = "x"
		YOption = "y"
	)

	var ss []string

	ps.Add("my-list",
		psetter.EnumList{
			Value: &ss,
			AllowedVals: psetter.AllowedVals{
				XOption: "a description of this option",
				YOption: "what this option means",
			},
		}, "help text")

	fmt.Println("Before parsing")
	for i, v := range ss {
		fmt.Printf("\tss[%d] = %q\n", i, v)
	}

	// Parse the arguments. We supply a list of strings but note that one of
	// them is not in the list of allowed values.
	errMap := ps.Parse([]string{"-my-list", "x,z"})

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
	//	: value is not allowed: "z"
	// At: [command line]: Supplied Parameter:2: "-my-list" "x,z"
	// After  parsing
}

// ExampleEnumList_withPassingChecks demonstrates how you can specify
// additional checks to be applied to the passed arguments before the value
// is set.
func ExampleEnumList_withPassingChecks() {
	ps := newPSetForTesting() // use paramset.NewOrDie()

	const (
		XOption = "x"
		YOption = "y"
	)

	var ss []string

	ps.Add("my-list",
		psetter.EnumList{
			Value: &ss,
			AllowedVals: psetter.AllowedVals{
				XOption: "a description of this option",
				YOption: "what this option means",
			},
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

// ExampleEnumList_withFailingChecks demonstrates the behaviour of the
// package when an invalid value is given. In this case the resulting list is
// not of the required length. It demonstrates the checks that can be
// supplied to ensure that the resulting list is as expected. Note that there
// is normally no need to examine the return from ps.Parse as the standard
// Helper will report any errors and abort the program.
func ExampleEnumList_withFailingChecks() {
	ps := newPSetForTesting() // use paramset.NewOrDie()

	const (
		XOption = "x"
		YOption = "y"
	)

	var ss []string

	ps.Add("my-list",
		psetter.EnumList{
			Value: &ss,
			AllowedVals: psetter.AllowedVals{
				XOption: "a description of this option",
				YOption: "what this option means",
			},
			Checks: []check.StringSlice{
				check.SliceLength[[]string](check.ValEQ(2)),
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

// ExampleEnumList_withNilValue demonstrates the behaviour of the package
// when an invalid setter is provided. In this case the Value to be set has
// not been initialised. Note that in production code you should not recover
// from the panic, instead you should fix the code that caused it.
func ExampleEnumList_withNilValue() {
	defer func() { // For test purposes only - do not recover in live code
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // use paramset.NewOrDie()

	const (
		XOption = "x"
		YOption = "y"
	)

	// we expect this to panic because the list Value has not been initialised
	ps.Add("my-list",
		psetter.EnumList{
			AllowedVals: psetter.AllowedVals{
				XOption: "a description of this option",
				YOption: "what this option means",
			},
		}, "help text")
	// Output:
	// panic
	// my-list: psetter.EnumList Check failed: the Value to be set is nil
}
