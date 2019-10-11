package phelp

import (
	"github.com/nickwells/param.mod/v3/param"
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
	twc.Println() //nolint: errcheck
	twc.WrapPrefixed("Note: ",
		"the order in which groups are processed is indeterminate"+
			" but within each group the files are processed in the"+
			" order listed above.", textIndent)
	twc.Println() //nolint: errcheck
	twc.WrapPrefixed("Note: ",
		"parameters given in the group config files must be valid"+
			" parameters of the program and they must be members"+
			" of the parameter group.",
		textIndent)
}

// showEnvironmentVariables prints the config files specific to particular
// groups of parameters that can be used to configure the behaviour of the
// program
func showEnvironmentVariables(twc *twrap.TWConf, ep []string) {
	if len(ep) == 0 {
		return
	}
	twc.Print("\n  Environment Variables\n\n")

	twc.Wrap(
		"The program can also be configured through"+
			" environment variables prefixed with:\n"+altSrcEnvVars(ep),
		textIndent)
	twc.Println() //nolint: errcheck
	twc.WrapPrefixed("Note: ",
		"The prefix is stripped off and any underscores ('_') in"+
			" the environment variable name after the prefix will be replaced"+
			" with dashes ('-') when matching the parameter name."+
			"\n\n"+
			"For instance, if the environment variables prefixes include 'XX_'"+
			" an environment variable called 'XX_a_b' will match a parameter"+
			" called 'a-b'",
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
// files. If there were no alternative sources it will not print anything and
// will return false.
func (h StdHelp) showAltSources(twc *twrap.TWConf, ps *param.PSet) bool {
	cf := ps.ConfigFiles()
	gf := getGroupConfigFiles(ps)
	ep := ps.EnvPrefixes()

	if len(cf) == 0 && len(gf) == 0 && len(ep) == 0 {
		return false
	}

	twc.Println("\nAlternative Sources") //nolint: errcheck

	twc.Println() //nolint: errcheck
	twc.Wrap("This program may also be configured through the"+
		" following alternative sources. These are in addition"+
		" to the parameters supplied on the command line.\n",
		0)

	showGroupConfigFiles(twc, gf)
	showConfigFiles(twc, cf)
	showEnvironmentVariables(twc, ep)

	twc.Println() //nolint: errcheck

	twc.WrapPrefixed("Note: ",
		"these additional sources are processed in the order shown above"+
			" and then the command line parameters are processed."+
			" This means that a value given on the command line will"+
			" replace any other settings in configuration files or"+
			" environment variables (unless the parameter may only be set"+
			" once). Similarly, settings in sources higher up this page"+
			" can be replaced by settings in sources lower in the page",
		0)
	return true
}
