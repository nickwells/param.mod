package paction

import (
	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v5/param"
)

// AppendStrings returns an ActionFunc that appends the given strings to the
// supplied slice.
//
// This can be used to add a string to a slice if a parameter has been set.
func AppendStrings(val *[]string, s ...string) param.ActionFunc {
	return func(_ location.L, _ *param.ByName, _ []string) error {
		*val = append(*val, s...)
		return nil
	}
}
