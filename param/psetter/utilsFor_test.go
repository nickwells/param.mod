package psetter_test

import (
	"fmt"

	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/paramset"
)

// panicSafeCheck runs the CheckSetter and catches any panic, returning true
// if a panic was caught
func panicSafeCheck(s param.Setter) (panicked bool, panicVal interface{}) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()
	s.CheckSetter("test")
	return
}

// newPSetForTesting returns a PSet suitable for use in a test (without all
// the standard parameters and help functions)
//
// Note that the paramset function used here is just to make the example
// more reliable. In production code you would be best to use
// paramset.NewOrDie(...) which will set the standard helper and exit if
// there's any error.
func newPSetForTesting() *param.PSet {
	ps, _ := paramset.NewNoHelpNoExitNoErrRpt()
	return ps
}

// logErrs will report the errors (if any) to stdout
func logErrs(errMap param.ErrMap) {
	for k, errs := range errMap {
		fmt.Println("Errors for:", k)
		for _, err := range errs {
			fmt.Println("\t:", err)
		}
	}
}
