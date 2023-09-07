package phelp

import (
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// showConfigFiles prints the config files that can be used to configure the
// behaviour of the program
func showConfigFiles(h StdHelp, twc *twrap.TWConf, cf []param.ConfigFileDetails) {
	if len(cf) == 0 {
		return
	}
	if h.showSummary {
		for _, f := range cf {
			if f.ParamsMustExist() {
				twc.Println("config-file::" + f.String())
			} else {
				twc.Println("multi-program-config-file::" + f.String())
			}
		}
		return
	}

	twc.Print("\n  Common Configuration Files\n\n")
	var hasNonStrictFiles bool
	for _, f := range cf {
		prefix := " "
		if !f.ParamsMustExist() {
			hasNonStrictFiles = true
			prefix = "*"
		}
		twc.Printf("    %s %s\n", prefix, f.String())
	}
	if hasNonStrictFiles {
		twc.Println() //nolint: errcheck
		twc.WrapPrefixed("Note: ",
			"the files marked with a '*' are allowed to contain"+
				" parameters not valid for this program. Any such"+
				" parameters will be silently ignored. To detect"+
				" such parameters call the program with the"+
				" '"+paramNameShowUnused+"' parameter.",
			textIndent)
	}
}

type groupCF struct {
	groupName string
	cf        param.ConfigFileDetails
}

// showGroupConfigFiles prints the config files specific to particular groups
// of parameters that can be used to configure the behaviour of the program
func showGroupConfigFiles(h StdHelp, twc *twrap.TWConf, gf []groupCF) {
	if len(gf) == 0 {
		return
	}
	if h.showSummary {
		for _, f := range gf {
			twc.Println("group-config-file:" + f.groupName +
				":" + f.cf.String()) //nolint: errcheck
		}
		return
	}

	twc.Print("\n  Group Configuration Files\n\n")

	for _, f := range gf {
		twc.Println("    "+f.groupName+": ", f.cf.String()) //nolint: errcheck
	}

	if len(gf) > 1 {
		twc.Println() //nolint: errcheck
		twc.WrapPrefixed("Note: ",
			"the order in which groups are processed is indeterminate"+
				" but within each group the files are processed in the"+
				" order listed above.", textIndent)
	}

	twc.Println() //nolint: errcheck
	twc.WrapPrefixed("Note: ",
		"parameters given in group config files must be valid"+
			" parameters of the program and members"+
			" of the parameter group.",
		textIndent)
}

// showEnvPrefixes prints the config files specific to particular
// groups of parameters that can be used to configure the behaviour of the
// program
func showEnvPrefixes(h StdHelp, twc *twrap.TWConf, ep []string) {
	if len(ep) == 0 {
		return
	}
	if h.showSummary {
		for _, e := range ep {
			twc.Println("env-var-prefix::" + e)
		}
		return
	}
	twc.Print("\n  Environment Variables\n\n")

	twc.Wrap(
		"The program can also be configured through"+
			" environment variables prefixed with:\n"+altSrcEnvVars(ep),
		textIndent)
}

// getGroupConfigFiles this returns the collection of group config files
func getGroupConfigFiles(ps *param.PSet) []groupCF {
	gf := []groupCF{}

	groups := ps.GetGroups()
	for _, grp := range groups {
		for _, configFile := range grp.ConfigFiles {
			gf = append(gf, groupCF{
				groupName: grp.Name,
				cf:        configFile,
			})
		}
	}
	return gf
}

// showAltSources will print a usage message showing the alternative sources
// that can be used to set parameters: environment variables or configuration
// files. If there were no alternative sources it will not print saying that
// there are no alternative sources.
func showAltSources(h StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	gf := getGroupConfigFiles(ps)
	cf := ps.ConfigFiles()
	ep := ps.EnvPrefixes()

	if len(gf) == 0 && len(cf) == 0 && len(ep) == 0 {
		if h.showSummary {
			twc.Wrap("none", textIndent)
		} else {
			twc.Wrap("There are no alternative sources, parameters can only"+
				" be set through the command line",
				textIndent)
		}
		return true
	}

	if !h.showSummary {
		twc.Print("Alternative Sources\n\n")
		twc.Wrap("Program parameters may be set through the command line"+
			" but also through these additional sources.",
			0)
	}

	showGroupConfigFiles(h, twc, gf)
	showConfigFiles(h, twc, cf)
	showEnvPrefixes(h, twc, ep)

	twc.Print("\n")
	return true
}
