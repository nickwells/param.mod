//go:build !unix

package phelp

import (
	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
)

// getHelpWidthSetter returns a setter suitable for use as the setter of the
// screen width for the help text wrapper.
func getHelpWidthSetter(helpWidth *int, chks ...check.ValCk[int]) param.Setter {
	return psetter.Int[int]{
		Value:  helpWidth,
		Checks: chks,
	}
}
