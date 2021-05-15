package param_test

import (
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/paramset"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

// Example_typicalUse shows how you would typically use the param
// package. Construct the PSet, adding any parameters through AddParam
// functions either from the main package or else package specific parameter
// setters. Set some description of the program. Then just call Parse with no
// parameters so that it will use the command line parameters
func Example_typicalUse() {
	ps := paramset.NewOrDie(
		tuAddParams1,
		tuAddParams2,
		param.SetProgramDescription("what this program does"))
	ps.Parse()
}

const tuExampleGroupName = "my-groupname"

var (
	tuExample1 bool

	tuExample2 int64
)

// tuAddParams1 will add the "example1" parameter to the PSet
func tuAddParams1(ps *param.PSet) error {
	// we must set the group description before we can use the group name.
	// Parameters which don't explicitly set the group name are put in the
	// pre-declared "cmd" group
	ps.SetGroupDescription(tuExampleGroupName,
		"The parameters for my command")

	ps.Add("example1",
		psetter.Bool{Value: &tuExample1},
		"here is where you would describe the parameter",
		param.AltName("e1"),
		param.GroupName(tuExampleGroupName))

	return nil
}

// tuAddParams2 will add the "example2" parameter to the PSet
func tuAddParams2(ps *param.PSet) error {
	// add the example2 parameter to the set. Note that we don't set any
	// groupname and so this will be in the default group for the command
	// ("cmd")
	ps.Add("example2",
		psetter.Int64{Value: &tuExample2},
		"the description of the parameter",
		param.AltName("e2"))

	return nil
}
