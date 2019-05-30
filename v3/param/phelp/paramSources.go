package phelp

import (
	"fmt"
	"os"

	"github.com/nickwells/param.mod/v3/param"
	"github.com/nickwells/twrap.mod/twrap"
)

const (
	altSrcSectionNoteIndent = 4
	altSrcCommonNoteIndent  = 4
)

// showConfigFiles prints the config files that can be used to configure the
// behaviour of the program
func showConfigFiles(twc *twrap.TWConf, cf []param.ConfigFileDetails) {
	if len(cf) == 0 {
		return
	}
	fmt.Fprint(twc.W, "\n  Common Configuration Files\n\n")

	var hasStrictFiles bool
	for _, f := range cf {
		fmt.Fprint(twc.W, "    ")
		if f.ParamsMustExist() {
			hasStrictFiles = true
			fmt.Fprint(twc.W, "* ")
		}
		fmt.Fprintln(twc.W, f.String())
	}
	if hasStrictFiles {
		fmt.Fprintln(twc.W)
		twc.WrapPrefixed("Note: ",
			"the files marked with a '*' should only contain"+
				" parameters valid for this program.", altSrcSectionNoteIndent)
	}
	fmt.Fprintln(twc.W)
	twc.WrapPrefixed("Note: ",
		"the files are processed in the order listed above.",
		altSrcSectionNoteIndent)
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
	fmt.Fprint(twc.W, "\n  Group Configuration Files\n\n")

	for _, f := range gf {
		fmt.Fprintln(twc.W, "    "+f.groupName+": ", f.cf.String())
	}
	fmt.Fprintln(twc.W)
	twc.WrapPrefixed("Note: ",
		"the order in which groups are processed is indeterminate"+
			" but within each group the files are processed in the"+
			" order listed above.", altSrcSectionNoteIndent)
	fmt.Fprintln(twc.W)
	twc.WrapPrefixed("Note: ",
		"parameters given in the group config files must be valid"+
			" parameters and they must be members of the group.",
		altSrcSectionNoteIndent)
}

// showEnvironmentVariables prints the config files specific to particular
// groups of parameters that can be used to configure the behaviour of the
// program
func showEnvironmentVariables(twc *twrap.TWConf, ep []string) {
	if len(ep) == 0 {
		return
	}
	fmt.Fprint(twc.W, "\n  Environment Variables\n\n")

	twc.Wrap("The program can also be configured "+altSrcEnvVars(ep), 4)
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

// showParamSources will print a usage message showing the alternative
// sources that can be used to set parameters: environment variables or
// configuration files.
func showParamSources(ps *param.PSet) {
	twc, err := twrap.NewTWConf(twrap.TWConfOptSetWriter(ps.StdWriter()))
	if err != nil {
		fmt.Fprint(os.Stderr, "Couldn't build the text wrapper:", err)
		return
	}

	cf := ps.ConfigFiles()
	gf := getGroupConfigFiles(ps)
	ep := ps.EnvPrefixes()

	fmt.Fprintln(twc.W, "\nAdditional Sources")
	if len(cf) == 0 && len(gf) == 0 && len(ep) == 0 {
		fmt.Fprintln(twc.W, "None")
		return
	}

	showGroupConfigFiles(twc, gf)
	showConfigFiles(twc, cf)
	showEnvironmentVariables(twc, ep)

	fmt.Fprintln(twc.W)
	twc.WrapPrefixed("Note: ",
		"the alternative sources are processed in the order shown above"+
			" and then the command line parameters are processed."+
			" This means that a value given on the command line will"+
			" replace any other settings in configuration files or"+
			" environment variables. Similarly, settings in sources"+
			" higher up this page can be replaced by settings lower"+
			" in the page",
		altSrcCommonNoteIndent)
}
