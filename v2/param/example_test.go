package param_test

import (
	"fmt"
	"github.com/nickwells/param.mod/v2/param"
	"github.com/nickwells/param.mod/v2/param/paramset"
	"github.com/nickwells/param.mod/v2/param/psetter"
	"os"
)

// Example_setDirectly shows how to use the param package. It is generally
// advisable to have the parameter setting grouped into separate functions
// (see the setByFunc example) but in order to show the use of the param
// funcs we have added the new parameters in line after constructing the new
// PSet.
//
// Note that the parameter names are given without any leading dashes. This
// is because they can be passed on the command line or through parameter
// files or environment variables where the leading dash is not used.
func Example_setDirectly() {
	var example1 bool
	var example2 int64

	ps, _ := paramset.New()
	ps.SetProgramDescription("what this program does")

	ps.Add("example1",
		psetter.BoolSetter{Value: &example1},
		"here is where you would describe the parameter",
		// optional additional settings
		param.AltName("e1"))

	ps.Add("example2",
		psetter.Int64Setter{Value: &example2},
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

// Example_setByFunc shows how to use the param package with param setting
// functions passed to the paramset.New method. These set the same parameters
// as in the setDirectly example. Because the parameters being set are
// external to the Example function we have defined them outside the function
// so that they are visible to the AddParams... funcs.
func Example_setByFunc() {
	ValExample1 = false
	ValExample2 = 0

	ps, _ := paramset.New(
		AddParams1,
		AddParams2,
		param.SetProgramDescription("what this program does"))

	fmt.Println("example1:", ValExample1)
	fmt.Println("example2:", ValExample2)
	// For the purposes of the example we are passing the parameters in as a
	// slice. In practice you would almost always pass nothing in which case
	// Parse will use the command line arguments.
	ps.Parse(
		[]string{
			"-e1", // this uses the alternative name for the parameter
		},
		[]string{
			"-example2", "3", // this parameter expects a following value
		},
	)
	fmt.Println("example1:", ValExample1)
	fmt.Println("example2:", ValExample2)

	// Output: example1: false
	// example2: 0
	// example1: true
	// example2: 3
}

// Example_withEnvVar shows how to use the param package using environment
// variables to set the parameter values
func Example_withEnvVar() {
	ValExample1 = false
	ValExample2 = 0

	ps, _ := paramset.New(
		AddParams1,
		AddParams2,
		param.SetProgramDescription("what this program does"))
	ps.SetEnvPrefix("GOLEM_PARAM_TEST_")
	ps.AddEnvPrefix("golem_param_test2_")

	os.Setenv("GOLEM_PARAM_TEST_"+"example1", "")
	os.Setenv("golem_param_test2_"+"example2", "3")

	fmt.Println("example1:", ValExample1)
	fmt.Println("example2:", ValExample2)
	// For the purposes of the example we are passing a slice of
	// strings. This is just to prevent the Parse func from setting any
	// values (and complaining about invalid parameters) from the command
	// line..
	ps.Parse([]string{})
	fmt.Println("example1:", ValExample1)
	fmt.Println("example2:", ValExample2)

	// Output: example1: false
	// example2: 0
	// example1: true
	// example2: 3
}

// Example_typicalUse shows how you would typically use the param
// package. Construct the PSet, adding any parameters through AddParam
// functions either from the main package or else package specific parameter
// setters. Set some description of the program. Then just call Parse with no
// parameters so that it will use the command line parameters
func Example_typicalUse() {
	ps, _ := paramset.New(
		AddParams1,
		AddParams2,
		param.SetProgramDescription("what this program does"))
	ps.Parse()
}

const exampleGroupName = "groupname"

var ValExample1 bool
var ValExample2 int64

// AddParams1 will set the "example1" parameter in the PSet
func AddParams1(ps *param.PSet) error {
	ps.SetGroupDescription(exampleGroupName,
		"The parameters for my command")

	ps.Add("example1",
		psetter.BoolSetter{Value: &ValExample1},
		"here is where you would describe the parameter",
		param.AltName("e1"),
		param.GroupName(exampleGroupName))

	return nil
}

// AddParams2 will set the "example2" parameter in the PSet
func AddParams2(ps *param.PSet) error {
	ps.Add("example2",
		psetter.Int64Setter{Value: &ValExample2},
		"the description of the parameter",
		param.AltName("e2"),
		param.GroupName(exampleGroupName))

	return nil
}
