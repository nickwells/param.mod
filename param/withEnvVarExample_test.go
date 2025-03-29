package param_test

import (
	"fmt"
	"os"

	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/paramset"
	"github.com/nickwells/param.mod/v6/psetter"
)

// Example_withEnvVar shows how to use the param package using environment
// variables to set the parameter values
func Example_withEnvVar() {
	wevExample1 = false
	wevExample2 = 0

	ps := paramset.NewOrPanic(
		wevAddParams1,
		wevAddParams2,
		param.SetProgramDescription("what this program does"))
	ps.SetEnvPrefix("GOLEM_PARAM_TEST_")
	ps.AddEnvPrefix("golem_param_test2_")

	_ = os.Setenv("GOLEM_PARAM_TEST_"+"example1", "")
	_ = os.Setenv("golem_param_test2_"+"example2", "3")

	fmt.Println("example1:", wevExample1)
	fmt.Println("example2:", wevExample2)
	// For the purposes of the example we are passing a slice of
	// strings. This is just to prevent the Parse func from setting any
	// values (and complaining about invalid parameters) from the command
	// line (the os.Args slice).
	ps.Parse([]string{})
	fmt.Println("example1:", wevExample1)
	fmt.Println("example2:", wevExample2)

	// Output: example1: false
	// example2: 0
	// example1: true
	// example2: 3
}

const wevExampleGroupName = "groupname"

var (
	wevExample1 bool

	wevExample2 int64
)

// wevAddParams1 will set the "example1" parameter in the PSet
func wevAddParams1(ps *param.PSet) error {
	ps.AddGroup(wevExampleGroupName,
		"The parameters for my command")

	ps.Add("example1",
		psetter.Bool{Value: &wevExample1},
		"here is where you would describe the parameter",
		param.AltNames("e1"),
		param.GroupName(wevExampleGroupName))

	return nil
}

// wevAddParams2 will set the "example2" parameter in the PSet
func wevAddParams2(ps *param.PSet) error {
	ps.Add("example2",
		psetter.Int[int64]{Value: &wevExample2},
		"the description of the parameter",
		param.AltNames("e2"),
		param.GroupName(wevExampleGroupName))

	return nil
}
