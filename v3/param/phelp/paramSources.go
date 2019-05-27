package phelp

import (
	"fmt"
	"io"
	"os"

	"github.com/nickwells/param.mod/v3/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// showConfigFiles prints the config files that can be used to configure the
// behaviour of the program
func showConfigFiles(w io.Writer, cf []param.ConfigFileDetails) {
	if len(cf) != 0 {
		fmt.Fprintln(w, "  Configuration Files")

		for _, f := range cf {
			fmt.Fprintln(w, "    ", f.String())
		}
	}
}

type groupCF struct {
	groupName string
	cf        param.ConfigFileDetails
}

// showGroupConfigFiles prints the config files specific to particular groups
// of parameters that can be used to configure the behaviour of the program
func showGroupConfigFiles(w io.Writer, gf []groupCF) {
	if len(gf) != 0 {
		fmt.Fprintln(w, "  Group Configuration Files")

		for _, f := range gf {
			fmt.Fprintln(w, "    "+f.groupName+": ", f.cf.String())
		}
	}
}

// showEnvironmentVariables prints the config files specific to particular groups
// of parameters that can be used to configure the behaviour of the program
func showEnvironmentVariables(w io.Writer, ep []string) {
	if len(ep) != 0 {
		fmt.Fprintln(w, "  Environment Variables")

		twc, err := twrap.NewTWConf(twrap.TWConfOptSetWriter(w))
		if err != nil {
			fmt.Fprint(os.Stderr, "Couldn't build the text wrapper:", err)
			return
		}
		twc.Wrap("The program can also be configured "+altSrcEnvVars(ep), 4)
	}
}

// showParamSources will print a usage message showing the alternative
// sources that can be used to set parameters: environment variables or
// configuration files.
func showParamSources(ps *param.PSet) {
	cf := ps.ConfigFiles()
	gf := []groupCF{}
	ep := ps.EnvPrefixes()

	groups := ps.GetGroups()
	for _, grp := range groups {
		for _, configFile := range grp.ConfigFiles {
			gf = append(gf, groupCF{
				groupName: grp.Name,
				cf:        configFile,
			})
		}
	}

	w := ps.StdWriter()

	fmt.Fprintln(w, "\nAdditional Sources")
	if len(cf) == 0 && len(gf) == 0 && len(ep) == 0 {
		fmt.Fprintln(w, "None")
		return
	}

	showConfigFiles(w, cf)
	showGroupConfigFiles(w, gf)
	showEnvironmentVariables(w, ep)
}
