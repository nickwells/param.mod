package paction

import (
	"fmt"

	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v3/param"
)

// Report returns an ActionFunc that will print its argument to the standard
// writer of the PSet (as given by the StdWriter method).
func Report(msg string) param.ActionFunc {
	return func(_ location.L, p *param.ByName, _ []string) error {
		fmt.Fprint(p.StdWriter(), msg) // nolint: errcheck
		return nil
	}
}

// ErrReport returns an ActionFunc that will print its argument to the error
// writer of the PSet (as given by the ErrWriter method).
func ErrReport(msg string) param.ActionFunc {
	return func(_ location.L, p *param.ByName, _ []string) error {
		fmt.Fprint(p.ErrWriter(), msg) // nolint: errcheck
		return nil
	}
}
