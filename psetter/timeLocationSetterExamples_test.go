package psetter_test

import (
	"fmt"
	"time"

	"github.com/nickwells/param.mod/v6/psetter"
)

// ExampleTimeLocation_standard demonstrates the use of a TimeLocation setter
func ExampleTimeLocation_standard() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	var loc *time.Location

	ps.Add("location", psetter.TimeLocation{Value: &loc}, "help text")

	fmt.Printf("Before parsing   location: %v\n", loc)
	ps.Parse([]string{"-location", "Europe/London"})
	fmt.Printf("After  parsing   location: %v\n", loc)
	// Output:
	// Before parsing   location: UTC
	// After  parsing   location: Europe/London
}

// ExampleTimeLocation_withNilValue demonstrates the behaviour of the package
// when an invalid setter is provided. In this case the Value to be set has
// not been initialised. Note that in production code you should not recover
// from the panic, instead you should fix the code that caused it.
func ExampleTimeLocation_withNilValue() {
	defer func() { // For test purposes only - do not recover in live code
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // use paramset.NewOrPanic()

	ps.Add("location", psetter.TimeLocation{}, "help text")

	// Output:
	// panic
	// location: psetter.TimeLocation Check failed: the Value to be set is nil
}
