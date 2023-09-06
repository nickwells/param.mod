package param_test

import (
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/paramset"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

const (
	groupNameTypicalExample = "my-groupname"
	paramNameExample1       = "example1"
	paramNameExample2       = "example2"
)

// progExampleTypicalUse holds configuration values and any other
// program-wide values.
type progExampleTypicalUse struct {
	tuExample1 bool

	tuExample2 int64
}

// tuAddParams1 takes a prog pointer and returns a function which will add
// the parameters to the PSet. This allows you to avoid having global
// variables for your parameter values and still allows you to group the
// parameter setup in a separate function.
func tuAddParams1(prog *progExampleTypicalUse) param.PSetOptFunc {
	return func(ps *param.PSet) error {
		// we must set the group description before we can use the group name.
		// Parameters which don't explicitly set the group name are put in the
		// pre-declared "cmd" group
		ps.SetGroupDescription(groupNameTypicalExample,
			"The parameters for my command")

		ps.Add(paramNameExample1,
			psetter.Bool{Value: &prog.tuExample1},
			"here is where you would describe the parameter",
			// Here we add alternative parameter names
			param.AltNames("e1"),
			// Here we set the name of the group of parameters. This grouping
			// is used to shape the help message delivered
			param.GroupName(groupNameTypicalExample),
			// Here we add a reference to the other parameter. This will
			// appear in the parameter help text. It is an error if the named
			// parameter does not exist.
			param.SeeAlso(paramNameExample2),
		)

		return nil
	}
}

// tuAddParams2 takes a prog pointer and returns a function which will add
// the parameters to the PSet.
func tuAddParams2(prog *progExampleTypicalUse) param.PSetOptFunc {
	return func(ps *param.PSet) error {
		// add the example2 parameter to the set. Note that we don't set any
		// groupname and so this will be in the default group for the command
		// ("cmd")
		ps.Add(paramNameExample2,
			psetter.Int[int64]{Value: &prog.tuExample2},
			"the description of the parameter",
			param.AltNames("e2"),
			param.SeeAlso(paramNameExample1),
		)

		return nil
	}
}

// Example_typicalUse shows how you would typically use the param
// package. Construct the PSet, adding any parameters through AddParam
// functions either from the main package or else package specific parameter
// setters. Set some description of the program. Then just call Parse with no
// parameters so that it will use the command line parameters
func Example_typicalUse() {
	prog := &progExampleTypicalUse{}
	ps := paramset.NewOrPanic(
		tuAddParams1(prog),
		tuAddParams2(prog),
		param.SetProgramDescription("what this program does"))
	ps.Parse()

	// the rest of your program goes here
}
