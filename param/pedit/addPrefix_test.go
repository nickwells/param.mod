package pedit_test

import (
	"fmt"

	"github.com/nickwells/param.mod/v5/param/paramset"
	"github.com/nickwells/param.mod/v5/param/pedit"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

// ExampleAddPrefix_Edit demonstrates the behaviour of the AddPrefix editor.
//
// Note that the paramset function used here is just to make the example more
// reliable. In production code you would be best to use
// paramset.NewOrDie(...) which will set the standard helper and exit if
// there's any error. Similarly, in production code, you would call Parse
// with no arguments in which case it will use the arguments given to the
// program (os.Args).
func ExampleAddPrefix_Edit() {
	ps, _ := paramset.NewNoHelpNoExitNoErrRpt()
	var s string

	ps.Add("param",
		psetter.String{
			Value:  &s,
			Editor: pedit.AddPrefix{Prefix: "MyPrefix-"},
		}, "help text",
	)

	ps.Parse([]string{"-param", "Abc"})
	fmt.Println(s)

	// Output:
	// MyPrefix-Abc
}
