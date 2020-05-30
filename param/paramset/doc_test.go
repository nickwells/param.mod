package paramset_test

import (
	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/paramset"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

var thingName string
var action = "nothing"

// addParams will add parameters to the passed ParamSet
func addParams(ps *param.PSet) error {
	// This adds a parameter to the PSet that can be given with either of two
	// names: '-name' or '-n'. The parameter parsing will report an error if
	// the parameter is not given or if the value given is an empty string
	ps.Add("name", psetter.String{
		Value:  &thingName,
		Checks: []check.String{check.StringLenGT(0)},
	},
		"set the name of the thing to do that other thing to",
		param.AltName("n"),
		param.Attrs(param.CommandLineOnly|param.MustBeSet),
	)

	// This adds another parameter to the PSet. This can be given with either
	// of '-action' or '-a'. The value given must be one of the allowed
	// values: 'delete', 'copy' or 'nothing'. An error will be reported if
	// the parameter is not seen
	ps.Add("action",
		psetter.Enum{
			Value: &action,
			AllowedVals: psetter.AllowedVals{
				"delete":  "delete the thing",
				"copy":    "copy the thing",
				"nothing": "do nothing",
			},
		},
		"give the action to perform on the thing",
		param.AltName("a"),
		param.Attrs(param.MustBeSet),
	)

	return nil
}

func ExampleNewOrDie_simple() {
	ps := paramset.NewOrDie(addParams,
		param.SetProgramDescription(
			"a description of the purpose of the program"))
	ps.Parse()

	// the rest of your program goes here
}
