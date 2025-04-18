package pedit_test

import (
	"fmt"

	"github.com/nickwells/param.mod/v6/paramset"
	"github.com/nickwells/param.mod/v6/pedit"
	"github.com/nickwells/param.mod/v6/psetter"
)

// ExampleAddPrefix_Edit demonstrates the behaviour of the AddPrefix editor.
//
// Note that the paramset function used here is just to make the example more
// reliable. In production code you would be best to use
// paramset.NewOrPanic(...) which will set the standard helper and panic if
// there's any error. Similarly, in production code, you would call Parse
// with no arguments in which case it will use the arguments given to the
// program (os.Args).
func ExampleAddPrefix_Edit() {
	ps, _ := paramset.NewNoHelpNoExitNoErrRpt()

	var s string

	ps.Add("param",
		psetter.String[string]{
			Value:  &s,
			Editor: pedit.AddPrefix{Prefix: "MyPrefix-"},
		}, "help text",
	)

	ps.Parse([]string{"-param", "Abc"})
	fmt.Println(s)

	// Output:
	// MyPrefix-Abc
}
