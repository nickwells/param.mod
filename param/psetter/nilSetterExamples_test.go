package psetter_test

import (
	"fmt"

	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/paction"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

// ExampleNil_standard demonstrates how you might use a Nil setter. Note that
// the Nil setter does nothing itself; any effect takes place through
// associated action functions
func ExampleNil_standard() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	var flag1 bool
	var flag2 bool

	ps.Add("my-param", psetter.Nil{}, "help text",
		param.PostAction(paction.SetVal(&flag1, true)),
		param.PostAction(paction.SetVal(&flag2, true)),
	)

	fmt.Println("Before parsing")
	fmt.Printf("\tflag1 = %v\n", flag1)
	fmt.Printf("\tflag2 = %v\n", flag2)
	ps.Parse([]string{"-my-param"})
	fmt.Println("After  parsing")
	fmt.Printf("\tflag1 = %v\n", flag1)
	fmt.Printf("\tflag2 = %v\n", flag2)
	// Output:
	// Before parsing
	//	flag1 = false
	//	flag2 = false
	// After  parsing
	//	flag1 = true
	//	flag2 = true
}
