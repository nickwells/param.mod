package phelp

import (
	"fmt"
	"os"

	"github.com/nickwells/param.mod/v2/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// showParamSources will print a usage message showing the alternative
// sources that can be used to set parameters: environment variables or
// configuration files.
func showParamSources(ps *param.PSet) {
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

		twc, err := twrap.NewTWConf(twrap.TWConfOptSetWriter(w))
		if err != nil {
			fmt.Fprint(os.Stderr, "Couldn't build the text wrapper:", err)
			return
		}
		twc.Wrap(altSrcEnvVars(ep), 4)
	}
}
