package paction

import (
	"fmt"
	"io"
	"os"

	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v7/param"
)

// ReportTo returns an ActionFunc that will print its msg argument to the
// given writer.
func ReportTo(w io.Writer, msg string) param.ActionFunc {
	return func(_ location.L, _ *param.BaseParam, _ []string) error {
		fmt.Fprint(w, msg)

		return nil
	}
}

// Report returns an ActionFunc that will print its argument to standard out.
func Report(msg string) param.ActionFunc {
	return ReportTo(os.Stdout, msg)
}

// ErrReport returns an ActionFunc that will print its argument to standard
// error.
func ErrReport(msg string) param.ActionFunc {
	return ReportTo(os.Stderr, msg)
}

// ReportToAndExit returns an ActionFunc that will print its argument to the
// given Writer. Having printed the message it will exit with the given
// status.
func ReportToAndExit(w io.Writer, exitStatus int, msg string) param.ActionFunc {
	return func(_ location.L, _ *param.BaseParam, _ []string) error {
		fmt.Fprint(w, msg)
		os.Exit(exitStatus)

		return nil
	}
}

// ReportAndExit returns an ActionFunc that will print its argument to
// standard out. Having printed the message it will exit with status 0.
func ReportAndExit(msg string) param.ActionFunc {
	return ReportToAndExit(os.Stdout, 0, msg)
}

// ErrReportAndExit returns an ActionFunc that will print its argument to
// standard error. Having printed the message it will exit with status 1.
func ErrReportAndExit(msg string) param.ActionFunc {
	return ReportToAndExit(os.Stderr, 1, msg)
}
