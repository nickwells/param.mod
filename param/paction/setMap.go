package paction

import (
	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v5/param"
)

// SetMap returns an ActionFunc that sets the keyed entry in the map to the
// given bool value.
//
// It is advised that the key be supplied as a named const value to avoid
// mismatches.
//
// This can be used if you want to set some map entry if a parameter has been
// set. For instance if you are setting some configuration for an operation
// it implies that you want the operation performed. This saves forcing the
// user to both specify the configuration and request the operation.
func SetMap(m map[string]bool, k string, b bool) param.ActionFunc {
	return func(_ location.L, _ *param.ByName, _ []string) error {
		m[k] = b
		return nil
	}
}

// SetMapIf returns an ActionFunc that sets the keyed entry in the map to the
// given bool value if the supplied test function returns true.
//
// It is advised that the key be supplied as a named const value to avoid
// mismatches.
//
// This can be used if you want to set some map entry if a parameter has been
// set. For instance if you are setting some configuration for an operation
// it implies that you want the operation performed. This saves forcing the
// user to both specify the configuration and request the operation.
func SetMapIf(m map[string]bool, k string, b bool, test ParamTestFunc) param.ActionFunc {
	return func(loc location.L, p *param.ByName, s []string) error {
		if test(loc, p, s) {
			m[k] = b
		}

		return nil
	}
}
