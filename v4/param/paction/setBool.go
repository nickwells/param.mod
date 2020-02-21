package paction

import (
	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v3/param"
)

// SetBool returns an ActionFunc to set the val to the given bool.
//
// This can be used if you want to set some boolean if a parameter has been
// set. For instance if you are setting some configuration for an operation
// it implies that you want the operation performed. This saves forcing the
// user to both specify the configuration and request the operation
func SetBool(val *bool, setTo bool) param.ActionFunc {
	return func(_ location.L, _ *param.ByName, _ []string) error {
		*val = setTo
		return nil
	}
}
