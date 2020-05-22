package phelp

import (
	"crypto/md5"
	"fmt"
	"sort"
	"strings"

	"github.com/nickwells/param.mod/v4/param"
	"github.com/nickwells/param.mod/v4/param/psetter"
	"github.com/nickwells/twrap.mod/twrap"
)

// showByPosParams will print the positional parameters and their
// descriptions.
func showByPosParams(h StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	bppCount := ps.CountByPosParams()
	if bppCount == 0 {
		return false
	}

	twc.Print("Positional Parameters\n\n")

	for i := 0; i < bppCount; i++ {
		printByPosParam(h, twc, ps, i)
	}
	return true
}

// printByPosParam prints the details of the i'th positional parameter
func printByPosParam(h StdHelp, twc *twrap.TWConf, ps *param.PSet, i int) {
	bp, _ := ps.GetParamByPos(i)
	twc.Wrap(fmt.Sprintf("%d) %s", i+1, bp.Name()), paramIndent)

	if h.hideDescriptions {
		return
	}
	twc.Wrap(bp.Description(), descriptionIndent)
	h.showAllowedValsByPos(twc, bp)
	twc.Wrap("Initial value: "+bp.InitialValue(), descriptionIndent)
}

// getMaxGroupNameLen returns the length of the longest group name
func getMaxGroupNameLen(groups []*param.Group) int {
	maxNameLen := 0
	for _, g := range groups {
		if len(g.Name) > maxNameLen {
			maxNameLen = len(g.Name)
		}
	}
	return maxNameLen
}

// showGroups will show the group name and details for all the registered
// parameter groups
func showGroups(h StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	twc.Print("\nParameter groups\n\n")
	groups := ps.GetGroups()

	maxNameLen := getMaxGroupNameLen(groups)
	printSep := false
	for _, g := range groups {
		if !groupCanBeShown(h, g) {
			continue
		}
		if !h.hideDescriptions {
			printSep = printSepIf(twc, printSep, minorSectionSeparator)
		}
		h.printGroup(twc, g, maxNameLen)
	}
	return true
}

// showParamsByName will show the named parameters in name order. If
// paramsChosen is not empty then only those parameters will be shown.
func showParamsByName(h StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	groups := ps.GetGroups()
	paramsToShow := make([]*param.ByName, 0)
	for _, g := range groups {
		for _, p := range g.Params {
			if paramCanBeShown(h, p) {
				paramsToShow = append(paramsToShow, p)
			}
		}
	}

	if len(paramsToShow) == 0 {
		return false
	}

	sort.Slice(paramsToShow,
		func(i, j int) bool {
			return paramsToShow[i].Name() < paramsToShow[j].Name()
		})

	for _, p := range paramsToShow {
		h.printParamUsage(twc, p)
	}
	return true
}

// showParamsByGroupName will show the named parameters in name order. If
// paramsChosen is not empty then only those parameters will be shown.
func showParamsByGroupName(h StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	groups := ps.GetGroups()

	maxNameLen := getMaxGroupNameLen(groups)
	printSep := false
	count := 0
	for _, g := range groups {
		if !groupCanBeShown(h, g) {
			continue
		}

		printGroup := true
		for _, p := range g.Params {
			if paramCanBeShown(h, p) {
				if printGroup {
					printSep = printSepIf(twc, printSep, minorSectionSeparator)
					h.printGroup(twc, g, maxNameLen)
					printGroup = false
				}
				h.printParamUsage(twc, p)
				count++
			}
		}
	}

	return count != 0
}

// printParamUsage prints the named parameter help text
func (h StdHelp) printParamUsage(twc *twrap.TWConf, p *param.ByName) {
	valueReq := p.Setter().ValueReq()
	prefix := "-"
	suffix := valueNeededStr(valueReq)
	if !p.AttrIsSet(param.MustBeSet) {
		prefix = "[" + prefix
		suffix += "]"
	}

	paramNames := ""
	sep := ""

	for _, altParamName := range p.AltNames() {
		paramNames += sep + prefix + altParamName + suffix
		sep = " or "
	}
	twc.Wrap2Indent(paramNames, paramIndent, paramLine2Indent)

	if h.hideDescriptions {
		return
	}

	twc.Wrap(p.Description(), descriptionIndent)
	printParamAttributes(twc, p)
	h.showAllowedValsByName(twc, p)
	if valueReq == param.None {
		return
	}
	twc.WrapPrefixed("Initial value: ", p.InitialValue(), descriptionIndent)
}

// showAllowedValsByName shows the allowed values for the ByName parameter
func (h StdHelp) showAllowedValsByName(twc *twrap.TWConf, p *param.ByName) {
	if p.Setter().ValueReq() == param.None {
		return
	}
	h.showAllowedValues(twc, p.Name(), p.Setter())
}

// showAllowedValsByPos shows the allowed values for the ByPos parameter
func (h StdHelp) showAllowedValsByPos(twc *twrap.TWConf, p *param.ByPos) {
	h.showAllowedValues(twc, p.Name(), p.Setter())
}

// showAllowedValues prints the allowed values for a parameter. It will print
// a reference to an earlier parameter if the allowed value text has been
// seen already
func (h StdHelp) showAllowedValues(twc *twrap.TWConf, pName string, s param.Setter) {
	const prefix = "Allowed values: "

	aval := s.AllowedValues()
	keyStr := aval

	var avm psetter.AllowedVals
	if sAVM, ok := s.(psetter.AllowedValuesMapper); ok {
		avm = sAVM.AllowedValuesMap()
	}
	if avm != nil {
		keyStr += avm.String()
	}

	var avam psetter.Aliases
	if sAVM, ok := s.(psetter.AllowedValuesAliasMapper); ok {
		avam = sAVM.AllowedValuesAliasMap()
	}
	if avam != nil {
		keyStr += avam.String()
	}

	key := md5.Sum([]byte(keyStr))

	if name, alreadyShown := h.avalShownAlready[key]; alreadyShown {
		twc.WrapPrefixed(prefix,
			"(see parameter: "+name+")",
			descriptionIndent)
		return
	}

	h.avalShownAlready[key] = pName
	twc.WrapPrefixed(prefix, aval, descriptionIndent)
	if avm != nil {
		indent := descriptionIndent + len(prefix)
		aval = "The value must be one of the following:\n" + avm.String()
		twc.Wrap2Indent(aval, indent, indent+4)
	}
	if avam != nil {
		indent := descriptionIndent + len(prefix)
		aval = "The following aliases are available:\n" + avam.String()
		twc.Wrap2Indent(aval, indent, indent+4)
	}
}

// paramCanBeShown will return true if the param can be shown. It checks
// whether hidden items can be shown and if the param is in the list of
// explicitly chosen params
func paramCanBeShown(h StdHelp, p *param.ByName) bool {
	if len(h.paramsChosen) == 0 {
		if h.showHiddenItems {
			return true
		}
		if p.AttrIsSet(param.DontShowInStdUsage) {
			return false
		}
		return true
	}
	for _, name := range p.AltNames() {
		if h.paramsChosen[name] {
			return true
		}
	}

	return h.paramsChosen[p.Name()]
}

// groupCanBeShown will return true if the group can be shown. It checks
// whether hidden items can be shown and if the group is in the list of
// explicitly chosen groups
func groupCanBeShown(h StdHelp, g *param.Group) bool {
	if len(h.groupsChosen) == 0 {
		return true
	}

	return h.groupsChosen[g.Name]
}

// printGroup prints the group name etc
func (h StdHelp) printGroup(twc *twrap.TWConf, g *param.Group, maxLen int) {
	twc.Printf("%-*.*s [ ", maxLen, maxLen, g.Name)
	if len(g.Params) == 1 {
		twc.Print("1 parameter")
	} else {
		twc.Printf("%d parameters", len(g.Params))
	}
	if g.HiddenCount > 0 {
		if g.AllParamsHidden() {
			twc.Print(", all hidden")
		} else {
			twc.Printf(", %d hidden", g.HiddenCount)
		}
	}
	twc.Print(" ]\n")
	if h.hideDescriptions {
		return
	}
	if g.Desc != "" {
		twc.Wrap(g.Desc, textIndent)
	}
	printGroupConfigFile(twc, g)
	twc.Print("\n")
}

func valueNeededStr(vr param.ValueReq) string {
	if vr == param.Mandatory {
		return "=..."
	}
	if vr == param.Optional {
		return "[=...]"
	}
	return ""
}

// printParamAttributes prints additional text according to the settings of
// the parameter's attributes
func printParamAttributes(twc *twrap.TWConf, p *param.ByName) {
	if p.AttrIsSet(param.SetOnlyOnce) {
		twc.Wrap(
			"\nThis parameter value may only be set once."+
				" Any appearances after the first will not be used",
			descriptionIndent)
	}
	if p.AttrIsSet(param.CommandLineOnly) && p.PSet().HasAltSources() {
		var sources []string
		if p.PSet().HasGlobalConfigFiles() {
			sources = append(sources, "in the configuration files")
		} else {
			grpCF := p.PSet().ConfigFilesForGroup(p.GroupName())
			if len(grpCF) > 0 {
				sources = append(sources,
					"in the configuration files for this group")
			}
		}

		if p.PSet().HasEnvPrefixes() {
			sources = append(sources, "as an environment variable")
		}

		if len(sources) > 0 {
			twc.Wrap(
				"\nThis parameter may only be given on the command line, not "+
					strings.Join(sources, " or "), descriptionIndent)
		}
	}
}

func printGroupConfigFile(twc *twrap.TWConf, g *param.Group) {
	if len(g.ConfigFiles) > 0 {
		twc.Wrap(
			"\nParameters in this group may also be set "+
				altSrcConfigFiles(g.ConfigFiles),
			textIndent)
	}
}

// altSrcConfigFiles generates the fragment of the help message that
// lists the Config file names and returns it. If there are no names it
// returns the empty string.
func altSrcConfigFiles(cf []param.ConfigFileDetails) string {
	s := ""
	switch len(cf) {
	case 0:
	case 1:
		s = "in the configuration file: "
	default:
		s = "in one of the configuration files:\n"
	}

	sep := ""
	for _, f := range cf {
		s += sep + f.String()
		sep = "\n"
	}
	return s
}

// altSrcEnvVars generates the fragment of the help message that lists
// the allowed environment variable prefixes and returns it. If there are no
// valid prefixes it returns the empty string
func altSrcEnvVars(ep []string) string {
	switch len(ep) {
	case 0:
		return ""
	case 1:
		return ep[0]
	default:
		return strings.Join(ep[:len(ep)-1], ",\n") + " or\n" + ep[len(ep)-1]
	}
}
