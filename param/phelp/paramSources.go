package phelp

import (
	"github.com/nickwells/param.mod/v4/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// showConfigFiles prints the config files that can be used to configure the
// behaviour of the program
func showConfigFiles(twc *twrap.TWConf, cf []param.ConfigFileDetails) {
	if len(cf) == 0 {
		return
	}
	twc.Print("\n  Common Configuration Files\n\n")

	var hasStrictFiles bool
	for _, f := range cf {
		prefix := "      "
		if f.ParamsMustExist() {
			hasStrictFiles = true
			prefix = "    * "
		}
		twc.Printf("%s%s\n", prefix, f.String())
	}
	if hasStrictFiles {
		twc.Println() //nolint: errcheck
		twc.WrapPrefixed("Note: ",
			"the files marked with a '*' should only contain"+
				" parameters valid for this program.", textIndent)
	}
}

type groupCF struct {
	groupName string
	cf        param.ConfigFileDetails
}

// showGroupConfigFiles prints the config files specific to particular groups
// of parameters that can be used to configure the behaviour of the program
func showGroupConfigFiles(twc *twrap.TWConf, gf []groupCF) {
	if len(gf) == 0 {
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
		twc.Wrap("There are no alternative sources, parameters can only"+
			" be set through the command line",
			textIndent)
		return true
	}

	twc.Print("Alternative Sources\n\n")
	if !h.hideDescriptions {
		twc.Wrap("Program parameters may be set through the command line"+
			" but also through these additional sources.",
			0)
	}

	showGroupConfigFiles(twc, gf)
	showConfigFiles(twc, cf)
	showEnvPrefixes(h, twc, ep)

	twc.Print("\n")
	return true
}
