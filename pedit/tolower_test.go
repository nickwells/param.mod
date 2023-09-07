package pedit_test

import (
	"fmt"

	"github.com/nickwells/param.mod/v6/paramset"
	"github.com/nickwells/param.mod/v6/pedit"
	"github.com/nickwells/param.mod/v6/psetter"
)

// ExampleToLower_Edit demonstrates the behaviour of the ToLower editor.
//
// Note that the paramset function used here is just to make the example more
// reliable. In production code you would be best to use
// paramset.NewOrPanic(...) which will set the standard helper and panic if
// there's any error. Similarly, in production code, you would call Parse
// with no arguments in which case it will use the arguments given to the
// program (os.Args).
func ExampleToLower_Edit() {
	ps, _ := paramset.NewNoHelpNoExitNoErrRpt()
	var s string

	ps.Add("param",
		psetter.String[string]{
			Value:  &s,
			Editor: pedit.ToLower{},
		}, "help text",
	)

	ps.Parse([]string{"-param", "Abc"})
	fmt.Println(s)

	// Output:
	// abc
}
