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
func showByPosParams(h StdHelp, ps *param.PSet) bool {
	bppCount := ps.CountByPosParams()
	if bppCount == 0 {
		return false
	}

	h.twc.Print("Positional Parameters\n\n")

	for i := range bppCount {
		printByPosParam(h, ps, i)
	}

	return true
}

// printByPosParam prints the details of the i'th positional parameter
func printByPosParam(h StdHelp, ps *param.PSet, i int) {
	bp, err := ps.GetParamByPos(i)
	if err != nil {
		return
	}

	h.twc.Wrap(fmt.Sprintf("%d) %s", i+1, bp.Name()), paramIndent)

	if h.showSummary {
		return
	}

	h.twc.Wrap(bp.Description(), descriptionIndent)
	h.showAllowedVals(bp.Name(), bp.Setter())
	showInitialValue(h.twc, bp.InitialValue(), bp.Setter().CurrentValue())
}

// getMaxGroupNameLen returns the length of the longest group name
func getMaxGroupNameLen(groups []*param.Group) int {
	maxNameLen := 0
	for _, g := range groups {
		maxNameLen = max(len(g.Name()), maxNameLen)
	}

	return maxNameLen
}

// showGroups will show the group name and details for all the registered
// parameter groups
func showGroups(h StdHelp, ps *param.PSet) bool {
	h.twc.Print("\nParameter groups\n\n")

	groups := ps.GetGroups()
	maxNameLen := getMaxGroupNameLen(groups)
	sep := ""

	for _, g := range groups {
		if !groupCanBeShown(h, g) {
			continue
		}

		if !h.showSummary {
			h.twc.Print(sep)
			sep = minorSectionSeparator
		}

		h.printGroup(g, maxNameLen)
	}

	return true
}

// showParamsByName will show the named parameters in name order. If
// paramsChosen is not empty then only those parameters will be shown.
func showParamsByName(h StdHelp, ps *param.PSet) bool {
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
		h.printParamUsage(p)
	}

	return true
}

// showParamsByGroupName will show the named parameters in name order. If
// paramsChosen is not empty then only those parameters will be shown.
func showParamsByGroupName(h StdHelp, ps *param.PSet) bool {
	groups := ps.GetGroups()
	maxNameLen := getMaxGroupNameLen(groups)
	sep := ""
	count := 0

	for _, g := range groups {
		if !groupCanBeShown(h, g) {
			continue
		}

		printGroup := true

		for _, p := range g.Params() {
			if paramCanBeShown(h, p) || h.groupsChosen[g.Name()] {
				if printGroup {
					h.twc.Print(sep)
					sep = minorSectionSeparator
					h.printGroup(g, maxNameLen)
					printGroup = false
				}
			}

			if paramCanBeShown(h, p) {
				h.printParamUsage(p)

				count++
			}
		}
	}

	return count != 0
}

// printParamUsage prints the named parameter help text
func (h StdHelp) printParamUsage(p *param.ByName) {
	smy := ParamSummary(*p)
	h.twc.Wrap2Indent(smy, paramIndent, paramLine2Indent)

	if h.showSummary {
		return
	}

	h.twc.Wrap(p.Description(), descriptionIndent)
	printParamAttributes(h.twc, p)
	showSeeAlso(h.twc, p)
	showSeeNotes(h.twc, p)

	if p.Setter().ValueReq() == param.None {
		return
	}

	h.showAllowedVals(p.Name(), p.Setter())
	showInitialValue(h.twc, p.InitialValue(), p.Setter().CurrentValue())
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
func (h StdHelp) showAllowedVals(pName string, s param.Setter) {
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
			h.twc.WrapPrefixed(prefix,
				"(see parameter: "+name+")",
				descriptionIndent)

			return
		}

		h.avalShownAlready[keyStr] = pName
	}

	const extraIndent = 6

	indent := descriptionIndent + len(prefix)
	valDescIndent := indent + extraIndent

	h.twc.WrapPrefixed(prefix, s.AllowedValues(), descriptionIndent)
	h.twc.Wrap2Indent(valueList, indent, valDescIndent)
	h.twc.Wrap2Indent(aliases, indent, valDescIndent)
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
func (h StdHelp) printGroup(g *param.Group, maxLen int) {
	h.twc.Printf("%-*.*s [ ", maxLen, maxLen, g.Name())

	if len(g.Params()) == 1 {
		h.twc.Print("1 parameter")
	} else {
		h.twc.Printf("%d parameters", len(g.Params()))
	}

	if g.HiddenCount() > 0 && !h.showHiddenItems {
		if g.AllParamsHidden() {
			h.twc.Print(", all hidden")
		} else {
			h.twc.Printf(", %d hidden", g.HiddenCount())
		}
	}

	h.twc.Print(" ]\n")

	if h.showSummary {
		return
	}

	if desc := g.Desc(); desc != "" {
		h.twc.Wrap(desc, textIndent)
	}

	printGroupConfigFile(h.twc, g)
	h.twc.Print("\n")
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
		var sourcesNotAllowed []string

		if p.PSet().HasGlobalConfigFiles() {
			sourcesNotAllowed = append(sourcesNotAllowed,
				"in the configuration files")
		} else {
			grpCF := p.PSet().ConfigFilesForGroup(p.GroupName())
			if len(grpCF) > 0 {
				sourcesNotAllowed = append(sourcesNotAllowed,
					"in the configuration files for this group")
			}
		}

		if p.PSet().HasEnvPrefixes() {
			sourcesNotAllowed = append(sourcesNotAllowed,
				"as an environment variable")
		}

		if len(sourcesNotAllowed) > 0 {
			twc.Wrap(
				"\nThis parameter may only be given on the command line, not "+
					strings.Join(sourcesNotAllowed, " or "), descriptionIndent)
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
