package psetter_test

import (
	"fmt"

	"github.com/nickwells/errutil.mod/errutil"
	"github.com/nickwells/param.mod/v6/param"
)

// panicSafeCheck runs the CheckSetter and catches any panic, returning true
// if a panic was caught
func panicSafeCheck(s param.Setter) (panicked bool, panicVal any) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()

	s.CheckSetter("test")

	return
}

// logErrs will report the errors (if any) to stdout
func logErrs(errMap errutil.ErrMap) {
	for k, errs := range errMap {
		fmt.Println("Errors for:", k)

		for _, err := range errs {
			fmt.Println("\t:", err)
		}
	}
}
