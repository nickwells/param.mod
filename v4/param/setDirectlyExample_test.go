package param_test

import (
	"fmt"

	"github.com/nickwells/param.mod/v3/param"
	"github.com/nickwells/param.mod/v3/param/paramset"
	"github.com/nickwells/param.mod/v3/param/psetter"
)

// Example_setDirectly shows how to use the param package. It is generally
// advisable to have the parameter setting grouped into separate functions
// (see the typicalUse example) but in order to show the use of the param
// funcs we have added the new parameters in line after constructing the new
// PSet.
//
// Note that the parameter names are given without any leading dashes. This
// is because they can be passed on the command line or through parameter
// files or environment variables where the leading dash is not used.
func Example_setDirectly() {
	var example1 bool
	var example2 int64

	ps := paramset.NewOrDie()
	ps.SetProgramDescription("what this program does")

	ps.Add("example1",
		psetter.Bool{Value: &example1},
		"here is where you would describe the parameter",
		// optional additional settings
		param.AltName("e1"))

	ps.Add("example2",
		psetter.Int64{Value: &example2},
		"the description of the parameter",
		// optional additional settings
		param.AltName("e2"))

	fmt.Println("example1:", example1)
	fmt.Println("example2:", example2)

	// For the purposes of the example we are passing the parameters in as a
	// slice. In practice you would almost always pass nothing in which case
	// Parse will use the command line arguments.
	ps.Parse([]string{
		"-e1",            // this uses the alternative name for the parameter
		"-example2", "3", // this parameter expects a following value
	})
	fmt.Println("example1:", example1)
	fmt.Println("example2:", example2)

	// Output: example1: false
	// example2: 0
	// example1: true
	// example2: 3
}
