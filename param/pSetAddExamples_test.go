package param_test

import (
	"fmt"

	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/paction"
	"github.com/nickwells/param.mod/v5/param/paramset"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

// ExamplePSet_Add shows the usage of the Add method of the PSet. This is
// used to add new parameters into the set. This, most simple, use of the Add
// function will add a parameter to the parameter set in the default group,
// with no alternative names that can be used and the parameter will appear
// in the standard usage message.
func ExamplePSet_Add() {
	ps := paramset.NewOrDie()

	var f float64

	ps.Add("param-name", psetter.Float64{Value: &f},
		"a parameter description for the usage message")
}

// ExamplePSet_Add_withExtras shows the usage of the Add method of the
// PSet. This is used to add new parameters into the set. This example shows
// the use of the additional options that can be passed to the Add call.
func ExamplePSet_Add_withExtras() {
	ps := paramset.NewOrDie()

	var f float64
	var fHasBeenSet bool

	// We only capture the return value so that we can report the settings
	// below.
	p := ps.Add("param-name", psetter.Float64{Value: &f},
		"a parameter description for the usage message",
		// The following parameters are optional - a list of zero or more
		// functions that set some of the extra features of the parameter

		// This next line returns a function that sets the group name of the
		// parameter. All parameters with the same group name will be shown
		// together in the usage message.
		param.GroupName("test.group"),
		// This sets flags on the parameter:
		//
		// DontShowInStdUsage: the parameter will not be shown in the default
		//                     help message. This can be useful for less
		//                     commonly useful parameters.
		//
		// MustBeSet: it is an error that Parse will report if the parameter
		//            is not given when Parse is called
		param.Attrs(param.DontShowInStdUsage|param.MustBeSet),
		// This gives an alternative string that can be used to set the
		// parameter value. Each parameter name must be unique; the param
		// package will panic if any duplicate names are used.
		param.AltNames("pn"),
		// This provides a function to be called immediately after the
		// parameter has been set. It uses an action function from the
		// paction package that will set the fHasBeenSet variable to
		// true. This can also be found by the HasBeenSet method on the
		// ByName parameter.
		param.PostAction(paction.SetBool(&fHasBeenSet, true)),
	)

	fmt.Printf("param (f) value:   %3.1f\n", f)
	fmt.Printf("fHasBeenSet value: %v\n", fHasBeenSet)
	fmt.Printf("param HasBeenSet?: %v\n", p.HasBeenSet())
	fmt.Printf("group name: %s\n", p.GroupName())
	fmt.Printf("param name: %s\n", p.Name())
	fmt.Printf("CommandLineOnly:    %t\n", p.AttrIsSet(param.CommandLineOnly))
	fmt.Printf("MustBeSet:          %t\n", p.AttrIsSet(param.MustBeSet))
	fmt.Printf("SetOnlyOnce:        %t\n", p.AttrIsSet(param.SetOnlyOnce))
	fmt.Printf("DontShowInStdUsage: %t\n", p.AttrIsSet(param.DontShowInStdUsage))

	// Output: param (f) value:   0.0
	// fHasBeenSet value: false
	// param HasBeenSet?: false
	// group name: test.group
	// param name: param-name
	// CommandLineOnly:    false
	// MustBeSet:          true
	// SetOnlyOnce:        false
	// DontShowInStdUsage: true
}
