package paction

import (
	"fmt"
	"os"

	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v6/param"
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

// ReportAndExit returns an ActionFunc that will print its argument to the
// standard writer of the PSet (as given by the StdWriter method). Having
// printed the message it will exit with status 0.
func ReportAndExit(msg string) param.ActionFunc {
	return func(_ location.L, p *param.ByName, _ []string) error {
		fmt.Fprint(p.StdWriter(), msg) // nolint: errcheck
		os.Exit(0)

		return nil
	}
}

// ErrReportAndExit returns an ActionFunc that will print its argument to the
// error writer of the PSet (as given by the ErrWriter method). Having
// printed the message it will exit with status 1.
func ErrReportAndExit(msg string) param.ActionFunc {
	return func(_ location.L, p *param.ByName, _ []string) error {
		fmt.Fprint(p.ErrWriter(), msg) // nolint: errcheck
		os.Exit(1)

		return nil
	}
}
