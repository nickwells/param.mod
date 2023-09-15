package paramtest_test

import (
	"errors"
	"fmt"

	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/paramset"
	"github.com/nickwells/param.mod/v6/paramtest"
	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

// MyConfig is an example of a struct that holds values to be set by the
// param package.
type MyConfig struct {
	I64 int64
	S   string
	B   bool
}

// cmpMyConfigStruct compares the value with the expected value and returns
// an error if they differ.
//
// Note that the values are passed as any values and converted
// locally to MyConfig values with tests made that the conversions are
// successful.
func cmpMyConfigStruct(iVal, iExpVal any) error {
	val, ok := iVal.(*MyConfig)
	if !ok {
		return errors.New("Bad value: not a pointer to MyConfig")
	}
	expVal, ok := iExpVal.(*MyConfig)
	if !ok {
		return errors.New("Bad expected value: not a pointer to MyConfig")
	}

	if val.I64 != expVal.I64 {
		return fmt.Errorf("The I64 values differ: %d != %d",
			val.I64, expVal.I64)
	}

	if val.S != expVal.S {
		return fmt.Errorf("The S values differ: %q != %q",
			val.S, expVal.S)
	}

	if val.B != expVal.B {
		return fmt.Errorf("The B values differ: %t != %t",
			val.B, expVal.B)
	}

	return nil
}

// AddMCParams returns a function that can be used to add parameters to a PSet.
func AddMCParams(mc *MyConfig) func(ps *param.PSet) error {
	return func(ps *param.PSet) error {
		ps.Add("b", psetter.Bool{Value: &mc.B},
			"set the boolean flag",
			param.AltNames("bool"),
		)
		ps.Add("s", psetter.String[string]{Value: &mc.S},
			"set the string value",
			param.AltNames("string"),
		)
		ps.Add("i", psetter.Int[int64]{Value: &mc.I64},
			"set the int value",
			param.AltNames("int"),
		)

		return nil
	}
}

// ExampleParser demonstrates how you should use the Parser type to test your
// parameter sets.
//
// Note that each test case has its own PSet instance and the use of the
// paramset.NewNoHelpNoExitNoErrRptOrPanic func to create it. This is so that
// you don't get any confounding behaviour from the standard help package.
//
// Note also how each test case has a separate MyConfig value so that we
// don't get changes from previous tests confusing the results.
func ExampleParser() {
	var mc1 MyConfig
	var mc2 MyConfig

	testCases := []paramtest.Parser{
		{
			ID: testhelper.MkID("set I64 option"),
			Ps: paramset.NewNoHelpNoExitNoErrRptOrPanic(
				AddMCParams(&mc1)),
			Val:       &mc1,
			ExpVal:    &MyConfig{I64: 42},
			CheckFunc: cmpMyConfigStruct,
			Args:      []string{"-i", "42"},
		},
		{
			ID: testhelper.MkID("set B option"),
			Ps: paramset.NewNoHelpNoExitNoErrRptOrPanic(
				AddMCParams(&mc2)),
			Val:       &mc2,
			ExpVal:    &MyConfig{B: true},
			CheckFunc: cmpMyConfigStruct,
			Args:      []string{"-b"},
		},
	}

	for _, tc := range testCases {
		fmt.Printf("Here is where you would call tc.Test(t) for %q\n",
			tc.IDStr())
	}
}
