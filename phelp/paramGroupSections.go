package phelp

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nickwells/english.mod/english"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
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

	if h.showSummary {
		return
	}
	twc.Wrap(bp.Description(), descriptionIndent)
	h.showAllowedVals(twc, bp.Name(), bp.Setter())
	showInitialValue(twc, bp.InitialValue(), bp.Setter().CurrentValue())
}

// getMaxGroupNameLen returns the length of the longest group name
func getMaxGroupNameLen(groups []*param.Group) int {
	maxNameLen := 0
	for _, g := range groups {
		if len(g.Name()) > maxNameLen {
			maxNameLen = len(g.Name())
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
		if !h.showSummary {
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
		for _, p := range g.Params() {
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
		for _, p := range g.Params() {
			if paramCanBeShown(h, p) || h.groupsChosen[g.Name()] {
				if printGroup {
					printSep = printSepIf(twc, printSep, minorSectionSeparator)
					h.printGroup(twc, g, maxNameLen)
					printGroup = false
				}
			}
			if paramCanBeShown(h, p) {
				h.printParamUsage(twc, p)
				count++
			}
		}
	}

	return count != 0
}

// printParamUsage prints the named parameter help text
func (h StdHelp) printParamUsage(twc *twrap.TWConf, p *param.ByName) {
	valNeededSuffix := valueNeededStr(p)
	paramNames := ""
	optSuffix := ""
	if !p.AttrIsSet(param.MustBeSet) {
		paramNames = "["
		optSuffix += "]"
	}

	sep := ""

	for _, altParamName := range p.AltNames() {
		paramNames += sep + "-" + altParamName + valNeededSuffix
		sep = ", "
	}
	paramNames += optSuffix
	twc.Wrap2Indent(paramNames, paramIndent, paramLine2Indent)

	if h.showSummary {
		return
	}

	twc.Wrap(p.Description(), descriptionIndent)
	printParamAttributes(twc, p)
	showSeeAlso(twc, p)
	showSeeNotes(twc, p)

	if p.Setter().ValueReq() == param.None {
		return
	}
	h.showAllowedVals(twc, p.Name(), p.Setter())
	showInitialValue(twc, p.InitialValue(), p.Setter().CurrentValue())
}

// showInitialValue shows the initial value of the ByName parameter. If the
// value has changed then the initial value is always shown, otherwise if the
// initial value is empty or zero or "false" then it is not shown. Lastly if
// the value has changed then the current value is also shown.
func showInitialValue(twc *twrap.TWConf, initialValue, currentValue string) {
	if currentValue == initialValue &&
		(initialValue == "" ||
			initialValue == "0" ||
			initialValue == "0.0" ||
			initialValue == "false") {
		return
	}

	twc.WrapPrefixed("Initial value: ", initialValue, descriptionIndent)

	if currentValue != initialValue {
		twc.WrapPrefixed("Current value: ", currentValue, descriptionIndent)
	}
}

// showSeeAlso shows the references for the ByName parameter (if any)
func showSeeAlso(twc *twrap.TWConf, p *param.ByName) {
	refs := p.SeeAlso()
	if len(refs) == 0 {
		return
	}

	prompt := "See also: "
	twc.WrapPrefixed(prompt, strings.Join(refs, ", "), descriptionIndent)
}

// showSeeNotes shows the references to notes for the ByName parameter (if any)
func showSeeNotes(twc *twrap.TWConf, p *param.ByName) {
	notes := p.SeeNotes()
	if len(notes) == 0 {
		return
	}

	prompt := "See " + english.Plural("note", len(notes)) + ": "
	twc.WrapPrefixed(prompt, strings.Join(notes, ", "), descriptionIndent)
}

// showAllowedVals prints the allowed values for a parameter. It will print
// a reference to an earlier parameter if the allowed value text has been
// seen already and the text is longer than 50 characters
func (h StdHelp) showAllowedVals(
	twc *twrap.TWConf, pName string, s param.Setter,
) {
	const prefix = "Allowed values: "
	const longString = 50

	var valueList string
	if sAVM, ok := s.(psetter.AllowedValuesMapper); ok {
		avm := sAVM.AllowedValuesMap()
		if len(avm) == 1 {
			valueList = "The value must be:\n" + avm.String()
		} else if len(avm) > 1 {
			valueList = "The value must be one of the following:\n" +
				avm.String()
		}
	}

	var aliases string
	if sAVM, ok := s.(psetter.AllowedValuesAliasMapper); ok {
		avam := sAVM.AllowedValuesAliasMap()
		if len(avam) == 1 {
			aliases = "The following alias is available:\n" + avam.String()
		} else if len(avam) > 1 {
			aliases = "The following aliases are available:\n" + avam.String()
		}
	}

	keyStr := s.AllowedValues() + valueList + aliases
	if len(keyStr) > longString {
		if name, alreadyShown := h.avalShownAlready[keyStr]; alreadyShown {
			twc.WrapPrefixed(prefix,
				"(see parameter: "+name+")",
				descriptionIndent)
			return
		}
		h.avalShownAlready[keyStr] = pName
	}

	indent := descriptionIndent + len(prefix)
	valDescIndent := indent + 6
	twc.WrapPrefixed(prefix, s.AllowedValues(), descriptionIndent)
	twc.Wrap2Indent(valueList, indent, valDescIndent)
	twc.Wrap2Indent(aliases, indent, valDescIndent)
}

// paramCanBeShown will return true if the param can be shown. It checks
// whether hidden items can be shown and if the param is in the list of
// explicitly chosen params
func paramCanBeShown(h StdHelp, p *param.ByName) bool {
	if h.paramsChosen.hasNothingChosen() {
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

// groupCanBeShown will return true if the group can be shown. If no groups
// are explicitly chosen then any group can be shown. Otherwise the group
// must be in the list of explicitly chosen groups
func groupCanBeShown(h StdHelp, g *param.Group) bool {
	if h.groupsChosen.hasNothingChosen() {
		return true
	}

	return h.groupsChosen[g.Name()]
}

// printGroup prints the group name etc
func (h StdHelp) printGroup(twc *twrap.TWConf, g *param.Group, maxLen int) {
	twc.Printf("%-*.*s [ ", maxLen, maxLen, g.Name())
	if len(g.Params()) == 1 {
		twc.Print("1 parameter")
	} else {
		twc.Printf("%d parameters", len(g.Params()))
	}
	if g.HiddenCount() > 0 && !h.showHiddenItems {
		if g.AllParamsHidden() {
			twc.Print(", all hidden")
		} else {
			twc.Printf(", %d hidden", g.HiddenCount())
		}
	}
	twc.Print(" ]\n")
	if h.showSummary {
		return
	}
	if desc := g.Desc(); desc != "" {
		twc.Wrap(desc, textIndent)
	}
	printGroupConfigFile(twc, g)
	twc.Print("\n")
}

// valTypeName returns a descriptive string for the type of the Setter
func valTypeName(p *param.ByName) string {
	if name := p.ValueName(); name != "" {
		return name
	}

	s := p.Setter()
	if sVD, ok := s.(psetter.ValDescriber); ok {
		return sVD.ValDescribe()
	}

	valType := fmt.Sprintf("%T", s)

	parts := strings.Split(valType, ".")
	valType = parts[len(parts)-1]
	valType = strings.TrimRight(valType, "0123456789")

	return valType
}

// valueNeededStr returns a descriptive string indicating whether a trailing
// argument is needed and if so of what type it should be.
func valueNeededStr(p *param.ByName) string {
	valReq := p.Setter().ValueReq()
	if valReq == param.Mandatory {
		return "=" + valTypeName(p)
	}
	if valReq == param.Optional {
		return "[=" + valTypeName(p) + "] "
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
	if p.AttrIsSet(param.IsTerminalParam) {
		twc.Wrap(
			"\nNo more command-line parameters will be handled after this"+
				" parameter. They will be handled separately.",
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
	if len(g.ConfigFiles()) > 0 {
		twc.Wrap(
			"\nParameters in this group may also be set "+
				altSrcConfigFiles(g.ConfigFiles()),
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
