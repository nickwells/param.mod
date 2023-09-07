package paction

import (
	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v6/param"
)

// SetVal returns an ActionFunc that sets the val to the given setTo value.
//
// This can be used, for instance, if you want to set some value if a
// parameter has been set. For instance if you are setting some configuration
// for an operation it implies that you want the operation performed. This
// saves forcing the user to both specify the configuration and request the
// operation.
func SetVal[T any](val *T, setTo T) param.ActionFunc {
	return func(_ location.L, _ *param.ByName, _ []string) error {
		*val = setTo
		return nil
	}
}

// SetValIf returns an ActionFunc that sets the val to the given setTo value
// if the supplied test function returns true.
//
// This can be used, for instance, if you want to set some value if a
// parameter has been set. For instance if you are setting some configuration
// for an operation it implies that you want the operation performed. This
// saves forcing the user to both specify the configuration and request the
// operation.
func SetValIf[T any](val *T, setTo T, test ParamTestFunc) param.ActionFunc {
	return func(loc location.L, p *param.ByName, s []string) error {
		if test(loc, p, s) {
			*val = setTo
		}
		return nil
	}
}

// SetMapVal returns an ActionFunc that sets the keyed entry in the map to
// the given setTo value.
//
// It is advised that the key be supplied as a named const value to avoid
// mismatches.
//
// This can be used if you want to set some map entry if a parameter has been
// set. For instance if you are setting some configuration for an operation
// it implies that you want the operation performed. This saves forcing the
// user to both specify the configuration and request the operation.
func SetMapVal[K comparable, V any](m map[K]V, k K, b V) param.ActionFunc {
	return func(_ location.L, _ *param.ByName, _ []string) error {
		m[k] = b
		return nil
	}
}

// SetMapValIF returns an ActionFunc that sets the keyed entry in the map to
// the given setTo value if the supplied test function returns true.
//
// It is advised that the key be supplied as a named const value to avoid
// mismatches.
//
// This can be used if you want to set some map entry if a parameter has been
// set. For instance if you are setting some configuration for an operation
// it implies that you want the operation performed. This saves forcing the
// user to both specify the configuration and request the operation.
func SetMapValIf[K comparable, V any](m map[K]V, k K, b V, test ParamTestFunc,
) param.ActionFunc {
	return func(loc location.L, p *param.ByName, s []string) error {
		if test(loc, p, s) {
			m[k] = b
		}
		return nil
	}
}
