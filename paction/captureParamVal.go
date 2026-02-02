package paction

import (
	"fmt"

	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v6/param"
)

// CaptureParamVal will record the value of the parameter
func CaptureParamVal(v *string) param.ActionFunc {
	if v == nil {
		panic("nil pointer passed to CaptureParamVal")
	}

	return func(loc location.L, p *param.ByName, paramValues []string) error {
		const hasValue = 2
		if p.Setter().ValueReq() != param.Mandatory ||
			len(paramValues) != hasValue {
			panic(fmt.Errorf(
				"coding error: "+
					"capturing a parameter value must be done "+
					"with a parameter with a mandatory following value: "+
					"param: %q, at: %s", p.Name(), loc))
		}

		*v = paramValues[1]

		return nil
	}
}
