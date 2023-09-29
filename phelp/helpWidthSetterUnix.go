//go:build unix

package phelp

import (
	"os"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/twrap.mod/twrap"
	"golang.org/x/sys/unix"
)

// getHelpWidthSetter returns a setter suitable for use as the setter of the
// screen width for the help text wrapper. If the ioctl call returns an error
// (typically because stdout is not attached to a tty device) it will return
// the default target line length from the twrap package.
func getHelpWidthSetter(helpWidth *int, chks ...check.ValCk[int]) param.Setter {
	return psetter.Calculated[int]{
		Value: helpWidth,

		CalcMap: map[string]psetter.NamedCalc[int]{
			"auto": {
				Name: "use the terminal width as the help width." +
					" If the help output is not to a terminal," +
					" the default width is used.",
				Calc: func(_, _ string) (int, error) {
					ws, err := unix.IoctlGetWinsize(
						int(os.Stdout.Fd()), unix.TIOCGWINSZ)
					if err == nil {
						return int(ws.Col), nil
					}
					return twrap.DfltTargetLineLen, nil
				},
			},
		},
		Default: psetter.Val2IntCalc[int](),
		Checks:  chks,
	}
}
