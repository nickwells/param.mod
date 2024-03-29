package paramset_test

import (
	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/paramset"
	"github.com/nickwells/param.mod/v6/psetter"
)

var (
	thingName string
	action    = "nothing"
)

// addParams will add parameters to the passed ParamSet
func addParams(ps *param.PSet) error {
	// This adds a parameter to the PSet that can be given with either of two
	// names: '-name' or '-n'. The parameter parsing will report an error if
	// the parameter is not given or if the value given is an empty string
	ps.Add("name", psetter.String[string]{
		Value:  &thingName,
		Checks: []check.String{check.StringLength[string](check.ValGT(0))},
	},
		"set the name of the thing to do that other thing to",
		param.AltNames("n"),
		param.Attrs(param.CommandLineOnly|param.MustBeSet),
	)

	// This adds another parameter to the PSet. This can be given with either
	// of '-action' or '-a'. The value given must be one of the allowed
	// values: 'delete', 'copy' or 'nothing'. An error will be reported if
	// the parameter is not seen
	ps.Add("action",
		psetter.Enum[string]{
			Value: &action,
			AllowedVals: psetter.AllowedVals[string]{
				"delete":  "delete the thing",
				"copy":    "copy the thing",
				"nothing": "do nothing",
			},
		},
		"give the action to perform on the thing",
		param.AltNames("a"),
		param.Attrs(param.MustBeSet),
	)

	return nil
}

func ExampleNewOrPanic_simple() {
	ps := paramset.NewOrPanic(addParams,
		param.SetProgramDescription(
			"a description of the purpose of the program"))
	ps.Parse()

	// the rest of your program goes here
}
