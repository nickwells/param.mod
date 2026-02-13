package paction

import (
	"fmt"

	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v7/param"
)

// CaptureParamVal will record the value of the parameter
func CaptureParamVal(v *string) param.ActionFunc {
	if v == nil {
		panic("nil pointer passed to CaptureParamVal")
	}

	return func(loc location.L, p *param.BaseParam, paramVals []string) error {
		if p.Setter().ValueReq() != param.Mandatory ||
			len(paramVals) == 0 {
			panic(fmt.Errorf(
				"capturing a parameter value must be done "+
					"with a parameter with a mandatory value: "+
					"param: %q, at: %s", p.Name(), loc))
		}

		*v = paramVals[len(paramVals)-1]

		return nil
	}
}
