package phelp

import (
	"crypto/md5"
	"fmt"
	"os"
	"strings"

	"github.com/nickwells/param.mod/v3/param"
	"github.com/nickwells/twrap.mod/twrap"
)

const paramIndent = 6
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
func (h StdHelp) Help(ps *param.PSet, messages ...string) { //nolint: gocyclo
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

	switch h.style {
	case stdHelp:
		h.printStdUsage(twc, ps)

	case paramsByName:
		h.printParamsByName(twc, ps)

	case paramsInGroups:
		h.printParamsInGroups(twc, ps)

	case paramsNotInGroups:
		h.printParamsNotInGroups(twc, ps)

	case groupNamesOnly:
		h.printGroups(twc, ps)

	case progDescOnly:
		twc.Wrap(ps.ProgDesc()+"\n", 0)

	case altSourcesOnly:
		if !h.showAltSources(twc, ps) {
			twc.Wrap(
				"There are no alternative sources, parameters can only"+
					" be set through the command line",
				textIndent)
		}
	}

	if exitAfterHelp {
		os.Exit(1)
	}
}

// printStdUsage will print the standard usage message
func (h StdHelp) printStdUsage(twc *twrap.TWConf, ps *param.PSet) {
	if h.showFullHelp {
		twc.Wrap(ps.ProgDesc()+"\n", 0)
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
		h.printSetValNote(twc)
		if ps.HasAltSources() {
			printMinorSeparator(twc)
			h.showAltSources(twc, ps)
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

func (h StdHelp) printParamUsage(twc *twrap.TWConf, p *param.ByName) {
	prefix := "-"
	suffix := valueNeededStr(p.ValueReq())
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
	twc.Wrap(paramNames, paramIndent)

	if !h.showFullHelp {
		return
	}

	twc.Wrap(p.Description(), descriptionIndent)
	if p.AttrIsSet(param.SetOnlyOnce) {
		twc.Wrap(
			"\nThis parameter value may only be set once."+
				" Any appearances after the first will not be used",
			descriptionIndent)
	}
	if p.AttrIsSet(param.CommandLineOnly) && p.PSet().HasAltSources() {
		twc.Wrap(
			"\nThis parameter value may only be given on the command line"+
				" not in configuration files or as an environment variable.",
			descriptionIndent)
	}
	h.showAllowedValsByName(twc, p)
	if p.ValueReq() == param.None {
		return
	}
	twc.WrapPrefixed("Initial value: ", p.InitialValue(), descriptionIndent)
}

// showAllowedValsByName shows the allowed values for the ByName parameter
func (h StdHelp) showAllowedValsByName(twc *twrap.TWConf, p *param.ByName) {
	if p.ValueReq() == param.None {
		return
	}
	h.showAllowedValues(twc, p.Name(), p.AllowedValues(), p.AllowedValuesMap())
}

// showAllowedValsByPos shows the allowed values for the ByPos parameter
func (h StdHelp) showAllowedValsByPos(twc *twrap.TWConf, p *param.ByPos) {
	h.showAllowedValues(twc, p.Name(), p.AllowedValues(), p.AllowedValuesMap())
}

// showAllowedValues prints the allowed values for a parameter. It will print
// a reference to an earlier parameter if the allowed value text has been
// seen already
func (h StdHelp) showAllowedValues(twc *twrap.TWConf, pName, aval string, avalMap param.AValMap) {
	const prefix = "Allowed values: "

	keyStr := aval
	if avalMap != nil {
		keyStr += avalMap.String()
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
	if avalMap != nil {
		indent := descriptionIndent + len(prefix)
		aval = "The value must be one of the following:\n" + avalMap.String()
		twc.Wrap2Indent(aval, indent, indent+4)
	}
}

// printParamsByName will print just the named parameters
func (h StdHelp) printParamsByName(twc *twrap.TWConf, ps *param.PSet) {
	var shown = map[string]bool{}
	for _, pName := range h.paramsToShow {
		trimmedName := strings.TrimLeft(pName, "-")
		p, err := ps.GetParamByName(trimmedName)
		if err != nil {
			continue
		}

		if shown[p.Name()] {
			continue
		}
		h.printParamUsage(twc, p)
		shown[p.Name()] = true
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

func (h StdHelp) printParamsInGroups(twc *twrap.TWConf, ps *param.PSet) {
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

func (h StdHelp) printParamsNotInGroups(twc *twrap.TWConf, ps *param.PSet) {
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

func (h StdHelp) printGroups(twc *twrap.TWConf, ps *param.PSet) {
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

	for _, f := range cf {
		s += f.String() + "\n"
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
