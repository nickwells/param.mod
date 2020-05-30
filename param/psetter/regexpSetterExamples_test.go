package psetter_test

import (
	"fmt"
	"regexp"

	"github.com/nickwells/param.mod/v5/param/psetter"
)

// ExampleRegexp_standard demonstrates the use of a Regexp setter.
func ExampleRegexp_standard() {
	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	var re *regexp.Regexp

	ps.Add("my-re", psetter.Regexp{Value: &re}, "help text")

	fmt.Printf("Before parsing    re: ")
	if re == nil {
		fmt.Printf(" nil\n")
	} else {
		fmt.Printf(" non-nil [%s]\n", re.String())
	}
	ps.Parse([]string{"-my-re", `.*\.go`})
	fmt.Printf("After  parsing    re: ")
	if re == nil {
		fmt.Printf(" nil\n")
	} else {
		fmt.Printf(" non-nil [%s]\n", re.String())
	}
	// Output:
	// Before parsing    re:  nil
	// After  parsing    re:  non-nil [.*\.go]
}

// ExampleRegexp_withNilValue demonstrates the behaviour of the package when
// an invalid setter is provided. In this case the Value to be set has not
// been initialised. Note that in production code you should not recover from
// the panic, instead you should fix the code that caused it.
func ExampleRegexp_withNilValue() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	// we expect this to panic because the regexp pointer Value has not been
	// initialised
	ps.Add("do-this", psetter.Regexp{}, "help text")
	// Output:
	// panic
	// do-this: psetter.Regexp Check failed: the Value to be set is nil
}
