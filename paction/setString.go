package paction

import (
	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v6/param"
)

// SetString returns an ActionFunc that sets the val to the given string.
//
// This can be used if you want to set some string if a parameter has been
// set. For instance if you are setting some configuration for an operation
// it implies that you want the operation performed. This saves forcing the
// user to both specify the configuration and request the operation
//
// Deprecated: use SetVal[string] instead.
func SetString(val *string, setTo string) param.ActionFunc {
	return func(_ location.L, _ *param.ByName, _ []string) error {
		*val = setTo
		return nil
	}
}

// SetStringIf returns an ActionFunc that sets the val to the given string if
// the supplied test function returns true.
//
// This can be used if you want to set some string if a parameter has been
// set. For instance if you are setting some configuration for an operation
// it implies that you want the operation performed. This saves forcing the
// user to both specify the configuration and request the operation
//
// Deprecated: use SetValIf[string] instead.
func SetStringIf(val *string, setTo string, test ParamTestFunc,
) param.ActionFunc {
	return func(loc location.L, p *param.ByName, s []string) error {
		if test(loc, p, s) {
			*val = setTo
		}
		return nil
	}
}
