package phelp

import (
	"fmt"

	"github.com/nickwells/param.mod/v2/param"
)

// showParamSources will print a usage message showing the alternative
// sources that can be used to set parameters: environment variables or
// configuration files.
func showParamSources(ps *param.PSet) {
	cf := ps.ConfigFiles()
	ep := ps.EnvPrefixes()

	w := ps.StdWriter()

	fmt.Fprintln(w, "\nAdditional Sources") // nolint: errcheck
	if len(cf) == 0 && len(ep) == 0 {
		fmt.Fprintln(w, "None") // nolint: errcheck
		return
	}

	if len(cf) != 0 {
		fmt.Fprintln(w, "  Configuration Files") // nolint: errcheck

		for _, f := range cf {
			fmt.Fprintln(w, "    ", f.String()) // nolint: errcheck
		}
	}

	if len(ep) != 0 {
		fmt.Fprintln(w, "  Environment Variables") // nolint: errcheck

		formatText(w, altSrcEnvVars(ep), 4, 4)
	}
}
