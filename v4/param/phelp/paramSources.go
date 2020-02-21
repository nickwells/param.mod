package phelp

import (
	"fmt"

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
	if h.showFullHelp {
		twc.Println() //nolint: errcheck
		twc.WrapPrefixed("Note: ",
			"The prefix is stripped off and any underscores ('_') in"+
				" the environment variable name after the prefix will"+
				" be replaced with dashes ('-') when matching the"+
				" parameter name.\n\n"+
				"For instance, if the environment variables prefixes"+
				" include 'XX_' an environment variable called"+
				" 'XX_a_b' will match a parameter called 'a-b'",
			textIndent)
	}
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
	gf := getGroupConfigFiles(ps)
	cf := ps.ConfigFiles()
	ep := ps.EnvPrefixes()

	if len(gf) == 0 && len(cf) == 0 && len(ep) == 0 {
		return false
	}

	twc.Println("Alternative Sources") //nolint: errcheck

	twc.Println() //nolint: errcheck
	twc.Wrap("Program parameters may be set through the command line"+
		" but also through these additional sources.",
		0)

	showGroupConfigFiles(twc, gf)
	showConfigFiles(twc, cf)
	showEnvPrefixes(h, twc, ep)

	twc.Println() //nolint: errcheck

	if h.showFullHelp {
		twc.WrapPrefixed("Note: ",
			"additional sources are processed in the order shown above"+
				" and then the command line parameters are processed."+
				" This means that a value given on the command line will"+
				" replace any settings in configuration files or"+
				" environment variables (unless the parameter may only be set"+
				" once). Similarly, settings in sources higher up this page"+
				" can be replaced by settings in sources lower down the page",
			0)

		twc.Wrap("\nThe following parameters may be useful when working"+
			" with these alternative sources:",
			0)
		maxLen := 0
		for _, p := range []string{
			paramsShowWhereSetArgName, paramsShowUnusedArgName,
		} {
			if len(p) > maxLen {
				maxLen = len(p)
			}
		}
		twc.WrapPrefixed(
			fmt.Sprintf("%*s : ", -maxLen, paramsShowWhereSetArgName),
			"will show if values have been set from"+
				" any of the alternative sources.",
			textIndent)
		twc.WrapPrefixed(
			fmt.Sprintf("%*s : ", -maxLen, paramsShowUnusedArgName),
			"can be useful to check that parameters set"+
				" in alternative sources are all correct. Since some"+
				" shared config files can contain parameters intended"+
				" for other programs misspelled parameters may be"+
				" silently ignored. With this parameter you can see"+
				" all the potential parameters and check them.",
			textIndent)
	}
	return true
}
