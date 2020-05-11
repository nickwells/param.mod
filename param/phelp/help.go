package phelp

import (
	"crypto/md5"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/nickwells/param.mod/v4/param"
	"github.com/nickwells/param.mod/v4/param/psetter"
	"github.com/nickwells/twrap.mod/twrap"
)

const paramIndent = 6
const paramLine2Indent = 9
const descriptionIndent = 12
const textIndent = 4

// printMajorSeparator prints the separator between major parts of the help
// text
func printMajorSeparator(twc *twrap.TWConf) {
	twc.Print("\n===============\n\n")
}

// printMinorSeparator prints the separator between minor parts of the help
// text
func printMinorSeparator(twc *twrap.TWConf) {
	twc.Print("---------------\n")
}

// printSetValNote prints an explanation of how optional and mandatory values
// may be set
func (h StdHelp) printSetValNote(twc *twrap.TWConf) {
	printMinorSeparator(twc)
	twc.Println() //nolint: errcheck

	twc.WrapPrefixed("Note: ",
		"Optional parameter values (where the name is followed by [=...])"+
			" must come after an '=' rather than as the next argument."+
			" As follows,"+
			"\n\n"+
			"-xxx=false not -xxx false"+
			"\n\n"+
			"For parameters which must have a value it may be given in"+
			" either way",
		textIndent)
}

// printHelpMessages prints the messages
func (h StdHelp) printHelpMessages(twc *twrap.TWConf, messages ...string) {
	for _, message := range messages {
		twc.Wrap(message, 0)
	}
	if len(messages) > 0 {
		printMajorSeparator(twc)
	}
}

// Help prints any messages and then a standardised usage message based on
// the parameters supplied to the param set. If it is called directly (that
// is if the help style is set to noHelp) then the output will be written to
// the param.PSet's error writer (by default stderr) rather than to its
// standard writer (stdout) and os.Exit will be called with an exit status of
// 1 to indicate an error.
func (h StdHelp) Help(ps *param.PSet, messages ...string) {
	var twc *twrap.TWConf
	var err error
	exitAfterHelp := false

	if h.style == noHelp {
		h.style = stdHelp
		twc, err = twrap.NewTWConf(twrap.SetWriter(ps.ErrWriter()))
		exitAfterHelp = h.exitAfterHelp
	} else {
		twc, err = twrap.NewTWConf(twrap.SetWriter(ps.StdWriter()))
	}

	if err != nil {
		fmt.Fprint(os.Stderr, "Couldn't build the text wrapper:", err)
		return
	}

	h.printHelpMessages(twc, messages...)

	helper, ok := helpersByStyle[h.style]
	if !ok {
		panic(
			fmt.Errorf(
				"Help(...) has been called with an unhandled help style (%v)",
				h.style))
	}

	helper(h, twc, ps)

	if exitAfterHelp {
		os.Exit(1)
	}
}

// printProgDesc prints the program description
func printProgDesc(h StdHelp, twc *twrap.TWConf, ps *param.PSet) {
	twc.Wrap(ps.ProgDesc()+"\n", 0)
}

// printStdUsage will print the standard usage message
func printStdUsage(h StdHelp, twc *twrap.TWConf, ps *param.PSet) {
	if h.showFullHelp {
		printProgDesc(h, twc, ps)
	}
	twc.Print("Usage: ", ps.ProgName())
	if !h.printPositionalParams(twc, ps) {
		if ps.TrailingParamsExpected() {
			twc.Println(" ... " + //nolint: errcheck
				ps.TerminalParam() + " " +
				ps.TrailingParamsName() + "...")
		} else {
			twc.Println(" ...") //nolint: errcheck
		}
	}
	printMinorSeparator(twc)
	h.printByNameParams(twc, ps)
	if h.showFullHelp {
		if ps.HasNotes() {
			printMinorSeparator(twc)
			notes := ps.Notes()
			for _, n := range notes {
				twc.Println(n.Headline) //nolint: errcheck
				twc.Wrap(n.Text, textIndent)
			}
		}

		if ps.HasAltSources() {
			printMinorSeparator(twc)
			twc.Wrap(
				"Some parameters may also be set from alternative sources"+
					" such as configuration files or environment variables."+
					" For more details use the "+helpAltSourcesArgName+
					" parameter.",
				0)
		}
		if ps.HasReferences() {
			printMinorSeparator(twc)
			twc.Wrap(
				"References to other sources of information are available."+
					" To see them use the "+helpShowRefsArgName+" parameter.",
				0)
		}
		if ps.HasExamples() {
			printMinorSeparator(twc)
			twc.Wrap(
				"Examples are available."+
					" To see them use the "+helpShowExamplesArgName+
					" parameter.",
				0)
		}
	}
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

// printParamUsage prints the per-parameter help text
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

	if !h.showFullHelp {
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

// printParamsByName will print just the named parameters
func printParamsByName(h StdHelp, twc *twrap.TWConf, ps *param.PSet) {
	var shown = map[string]bool{}
	var badParams = []string{}
	for _, pName := range h.paramsToShow {
		trimmedName := strings.TrimLeft(pName, "-")
		p, err := ps.GetParamByName(trimmedName)
		if err != nil {
			badParams = append(badParams, pName)
			continue
		}

		if shown[p.Name()] {
			continue
		}
		h.printParamUsage(twc, p)
		shown[p.Name()] = true
	}
	switch len(badParams) {
	case 0:
		return
	case 1:
		twc.Wrap("parameter: "+badParams[0]+
			" is not a parameter of this program",
			0)
	default:
		twc.Wrap("The following are not parameters of this program: "+
			strings.Join(badParams, ", "),
			0)
	}
}

// printPositionalParams will print the positional parameters and their
// descriptions. It will return true if any positional parameters are printed
// and false otherwise
func (h StdHelp) printPositionalParams(twc *twrap.TWConf, ps *param.PSet) bool {
	bppCount := ps.CountByPosParams()
	if bppCount == 0 {
		return false
	}

	var hasTerminalParam bool
	for i := 0; i < bppCount; i++ {
		bp, _ := ps.GetParamByPos(i)
		twc.Print(" <", bp.Name(), ">")
		if bp.IsTerminal() {
			hasTerminalParam = true
		}
	}

	if hasTerminalParam {
		if ps.TrailingParamsExpected() {
			twc.Println(" " + //nolint: errcheck
				ps.TrailingParamsName() + "...")
		}
	} else {
		twc.Println(" ...") //nolint: errcheck
	}

	if !h.showFullHelp {
		return true
	}
	twc.Println("where") //nolint: errcheck

	for i := 0; i < bppCount; i++ {
		bp, _ := ps.GetParamByPos(i)
		twc.Wrap(bp.Name(), paramIndent)
		twc.Wrap(bp.Description(), descriptionIndent)
		h.showAllowedValsByPos(twc, bp)
		twc.WrapPrefixed("Initial value: ", bp.InitialValue(),
			descriptionIndent)
	}
	return true
}

// printGroupDetails prints the group name etc
func (h StdHelp) printGroupDetails(twc *twrap.TWConf, pg *param.Group) {
	twc.Printf("%s [ ", pg.Name)
	if len(pg.Params) == 1 {
		twc.Print("1 parameter")
	} else {
		twc.Printf("%d parameters", len(pg.Params))
	}
	if pg.HiddenCount > 0 {
		if pg.AllParamsHidden() {
			twc.Print(", all hidden")
		} else {
			twc.Printf(", %d hidden", pg.HiddenCount)
		}
	}
	twc.Println(" ]") //nolint: errcheck
	if !h.showFullHelp {
		return
	}
	if pg.Desc != "" {
		twc.Wrap(pg.Desc, textIndent)
	}
	printGroupConfigFile(twc, pg)
	twc.Println() //nolint: errcheck
}

// printGroupParams prints the parameters in the group
func (h StdHelp) printGroupParams(twc *twrap.TWConf, pg *param.Group) {
	for _, p := range pg.Params {
		if p.AttrIsSet(param.DontShowInStdUsage) &&
			!h.paramsShowHidden {
			continue
		}
		h.printParamUsage(twc, p)
	}
}

func (h StdHelp) printByNameParams(twc *twrap.TWConf, ps *param.PSet) {
	paramGroups := ps.GetGroups()

	sep := false
	for _, pg := range paramGroups {
		if pg.AllParamsHidden() && !h.paramsShowHidden {
			continue
		}
		if sep {
			printMinorSeparator(twc)
		}
		sep = true

		h.printGroupDetails(twc, pg)
		h.printGroupParams(twc, pg)
	}
}

// hasBadGroupNames checks the selected group names for validity and reports
// any unknown names. It returns true if any bad names were found
func (h StdHelp) hasBadGroupNames(twc *twrap.TWConf, ps *param.PSet) bool {
	var badGroupNames = []string{}

	for name := range h.groupsSelected {
		if ps.GetGroupByName(name) == nil {
			badGroupNames = append(badGroupNames, name)
		}
	}

	if len(badGroupNames) == 0 {
		return false
	}
	if len(badGroupNames) == 1 {
		twc.Wrap(badGroupNames[0]+
			" : is not a valid parameter group name", 0)
	} else {
		sort.Strings(badGroupNames)
		twc.Println( //nolint: errcheck
			"The following are not valid parameter group names:")
		for _, bgn := range badGroupNames {
			twc.Wrap(bgn, textIndent)
		}
	}
	twc.Println("possible group names are:") //nolint: errcheck
	for _, pg := range ps.GetGroups() {
		if !h.groupsSelected[pg.Name] {
			twc.Wrap(pg.Name, textIndent)
		}
	}
	return true
}

// printParamsInGroups prints the help for a selected list of parameter
// groups
func printParamsInGroups(h StdHelp, twc *twrap.TWConf, ps *param.PSet) {
	if h.hasBadGroupNames(twc, ps) {
		return
	}

	paramGroups := ps.GetGroups()

	sep := false
	for _, pg := range paramGroups {
		if !h.groupsSelected[pg.Name] {
			continue
		}
		if sep {
			printMinorSeparator(twc)
		}
		sep = true

		h.printGroupDetails(twc, pg)
		h.printGroupParams(twc, pg)
	}
}

// printParamsNotInGroups will print the help for all parameter groups except
// those listed
func printParamsNotInGroups(h StdHelp, twc *twrap.TWConf, ps *param.PSet) {
	if h.hasBadGroupNames(twc, ps) {
		return
	}

	paramGroups := ps.GetGroups()

	sep := false
	for _, pg := range paramGroups {
		if h.groupsSelected[pg.Name] {
			continue
		}
		if sep {
			printMinorSeparator(twc)
		}
		sep = true

		h.printGroupDetails(twc, pg)
		h.printGroupParams(twc, pg)
	}
}

func printGroupConfigFile(twc *twrap.TWConf, pg *param.Group) {
	if len(pg.ConfigFiles) > 0 {
		twc.Wrap(
			"\nParameters in this group may also be set "+
				altSrcConfigFiles(pg.ConfigFiles),
			textIndent)
	}
}

// printGroups will print the group name and details for all the registered
// parameter groups
func printGroups(h StdHelp, twc *twrap.TWConf, ps *param.PSet) {
	twc.Println("\nParameter groups") //nolint: errcheck
	paramGroups := ps.GetGroups()

	sep := false
	for _, pg := range paramGroups {
		if sep {
			printMinorSeparator(twc)
		}
		sep = true

		h.printGroupDetails(twc, pg)
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
