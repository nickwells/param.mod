package paction

import (
	"os"

	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v6/param"
)

// Exit returns an ActionFunc that will exit with the given exit status. This
// should always be the last ActionFunc as no subsequent ones will be called
func Exit(code int) param.ActionFunc {
	return func(_ location.L, _ *param.ByName, _ []string) error {
		os.Exit(code)
		return nil
	}
}
