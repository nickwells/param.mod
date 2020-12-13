package psetter_test

import (
	"fmt"
	"time"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

// ExampleTime_standard demonstrates the use of a Time setter
func ExampleTime_standard() {
	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	var t time.Time

	ps.Add("my-time", psetter.Time{Value: &t}, "help text")

	fmt.Printf("Before parsing: time = %s\n", t.Format("15:04"))
	ps.Parse([]string{"-my-time", "2000/Jan/01T15:00:00"})
	fmt.Printf("After  parsing: time = %s\n", t.Format("15:04"))
	// Output:
	// Before parsing: time = 00:00
	// After  parsing: time = 15:00
}

// ExampleTime_withFormat demonstrates the use of a Time setter with a
// non-default Format value
func ExampleTime_withFormat() {
	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	var t time.Time

	ps.Add("my-time",
		psetter.Time{
			Value:  &t,
			Format: psetter.TimeFmtHMS,
		},
		"help text")

	fmt.Printf("Before parsing: time = %s\n", t.Format("15:04:05"))
	ps.Parse([]string{"-my-time", "15:01:02"})
	fmt.Printf("After  parsing: time = %s\n", t.Format("15:04:05"))
	// Output:
	// Before parsing: time = 00:00:00
	// After  parsing: time = 15:01:02
}

// ExampleTime_withPassingChecks demonstrates how to add checks to be applied
// to the value.
func ExampleTime_withPassingChecks() {
	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	var t time.Time

	ps.Add("my-time",
		psetter.Time{
			Value: &t,
			Checks: []check.Time{
				check.TimeIsOnDOW(time.Friday),
			},
		},
		"help text")

	fmt.Printf("Before parsing: time = %s\n", t.Weekday())
	ps.Parse([]string{"-my-time", "2020/Apr/24T12:00:00"})
	fmt.Printf("After  parsing: time = %s\n", t.Weekday())
	// Output:
	// Before parsing: time = Monday
	// After  parsing: time = Friday
}

// ExampleTime_withFailingChecks demonstrates how to add checks to be applied
// to the value. Note that there is normally no need to examine the return
// from ps.Parse as the standard Helper will report any errors and abort the
// program.
func ExampleTime_withFailingChecks() {
	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	var t time.Time

	ps.Add("my-time",
		psetter.Time{
			Value: &t,
			Checks: []check.Time{
				check.TimeIsOnDOW(time.Friday),
			},
		},
		"help text")

	fmt.Printf("Before parsing: time = %s\n", t.Weekday())
	// Parse the arguments. Note that the time given is not on a Friday.
	errMap := ps.Parse([]string{"-my-time", "2020/Apr/25T12:00:00"})
	// We expect to see an error reported.
	logErrs(errMap)
	// The value is unchanged due to the error.
	fmt.Printf("After  parsing: time = %s\n", t.Weekday())
	// Output:
	// Before parsing: time = Monday
	// Errors for: my-time
	//	: the day of the week (Saturday) should be a Friday
	// At: [command line]: Supplied Parameter:2: -my-time 2020/Apr/25T12:00:00
	// After  parsing: time = Monday
}

// ExampleTime_withNilValue demonstrates the behaviour of the package when
// an invalid setter is provided. In this case the Value to be set has not
// been initialised. Note that in production code you should not recover from
// the panic, instead you should fix the code that caused it.
func ExampleTime_withNilValue() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("panic")
			fmt.Println(p)
		}
	}()

	ps := newPSetForTesting() // you would normally use paramset.NewOrDie()

	ps.Add("my-time", psetter.Time{}, "help text")

	// Output:
	// panic
	// my-time: psetter.Time Check failed: the Value to be set is nil
}
