package phelp

import (
	"fmt"
	"github.com/nickwells/param.mod/v2/param"
)

// showParamSources will print a usage message showing the alternative
// sources that can be used to set parameters: environment variables or
// configuration files.
func showParamSources(ps *param.ParamSet) {
	cf := ps.ConfigFiles()
	ep := ps.EnvPrefixes()

	w := ps.StdWriter()

	fmt.Fprintln(w, "\nAdditional Sources")
	if len(cf) == 0 && len(ep) == 0 {
		fmt.Fprintln(w, "None")
		return
	}

	if len(cf) != 0 {
		fmt.Fprintln(w, "  Configuration Files")

		for _, f := range cf {
			fmt.Fprintln(w, "    ", f.String())
		}
	}

	if len(ep) != 0 {
		fmt.Fprintln(w, "  Environment Variables")

		formatText(w, altSrcEnvVars(ep), 4, 4)
	}
}
