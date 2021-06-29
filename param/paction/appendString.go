package paction

import (
	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v5/param"
)

// AppendStrings returns an ActionFunc that appends the given strings to the
// supplied slice.
//
// This can be used to add a string to a slice if a parameter has been set.
//
// Note that the value of the strings is passed by value and so the value
// used will be that at the time this is called. This is therefore useful
// only for setting fixed values; use the AppendStringVal function if you
// want to apply the value of a string variable at the time the ActionFunc is
// called.
func AppendStrings(val *[]string, s ...string) param.ActionFunc {
	return func(_ location.L, _ *param.ByName, _ []string) error {
		*val = append(*val, s...)
		return nil
	}
}

// AppendStringVal returns an ActionFunc that appends the given string to the
// supplied slice.
//
// This can be used to add a string to a slice if a parameter has been set.
//
// This should be used if you want to apply the value of a string variable at
// the time the ActionFunc is called rather than when it is first set.

func AppendStringVal(val *[]string, s *string) param.ActionFunc {
	return func(_ location.L, _ *param.ByName, _ []string) error {
		*val = append(*val, *s)
		return nil
	}
}
