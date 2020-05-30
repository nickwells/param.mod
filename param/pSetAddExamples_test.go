package param_test

import (
	"fmt"

	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/paramset"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

// ExamplePSet_Add shows the usage of the Add method of the PSet. This is
// used to add new parameters into the set.
func ExamplePSet_Add() {
	ps, _ := paramset.New()

	// we declare f here for the purposes of the example but typically it
	// would be declared in package scope somewhere or in the main() func
	var f float64

	p := ps.Add(
		"param-name",
		psetter.Float64{Value: &f},
		"a parameter description",
		param.GroupName("test.group"),         // Optional parameter
		param.Attrs(param.DontShowInStdUsage), // Optional parameter
	)

	fmt.Printf("%3.1f\n", f)
	fmt.Printf("group name: %s\n", p.GroupName())
	fmt.Printf("param name: %s\n", p.Name())
	fmt.Printf("CommandLineOnly: %t\n", p.AttrIsSet(param.CommandLineOnly))
	fmt.Printf("MustBeSet: %t\n", p.AttrIsSet(param.MustBeSet))
	fmt.Printf("SetOnlyOnce: %t\n", p.AttrIsSet(param.SetOnlyOnce))
	fmt.Printf("DontShowInStdUsage: %t\n", p.AttrIsSet(param.DontShowInStdUsage))

	// Output: 0.0
	// group name: test.group
	// param name: param-name
	// CommandLineOnly: false
	// MustBeSet: false
	// SetOnlyOnce: false
	// DontShowInStdUsage: true
}
