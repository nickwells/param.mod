package phelp

import (
	"fmt"
	"io"
	"os"

	"github.com/nickwells/param.mod/v2/param"
	"github.com/nickwells/twrap.mod/twrap"
)

const stdIndent = "    "

const dashes = "---------------"
const equals = "==============="

const paramIndent = 8
const descriptionIndent = 16
const textIndent = 4

// badGroups checks that all the groups are in the PSet and reports the
// error if not. It returns a count of the number of problems found
func badGroups(ps *param.PSet, groups map[string]bool, name string) bool {
	badGroups := 0
	msg := ""
	for g := range groups {
		if !ps.HasGroupName(g) {
			if badGroups == 0 {
				msg = "group: '" + g + "' in the list of " + name + "," +
					" is not the name of a parameter group." +
					" Please check the spelling."
			} else {
				msg += "\nalso: '" + g + "'"
			}
			badGroups++
		}
	}
	if badGroups > 0 {
		twc, err := twrap.NewTWConf(twrap.SetWriter(ps.ErrWriter()))
		if err != nil {
			fmt.Fprint(os.Stderr, "Couldn't build the text wrapper:", err)
			return true
		}

		twc.WrapPrefixed("Error: ", msg, 0)
	}
	return badGroups > 0
}

// printOptValNote prints an explanation of how optional values must be set
func (h StdHelp) printOptValNote(w io.Writer) {
	twc, err := twrap.NewTWConf(twrap.SetWriter(w))
	if err != nil {
		fmt.Fprint(os.Stderr, "Couldn't build the text wrapper:", err)
		return
	}
	fmt.Fprint(w, "\n"+equals+"\n\n")

	prefix := "Note: "
	twc.WrapPrefixed(prefix,
		"Optional values (those with a parameter name followed by [=...])"+
			" must be given with the parameter,"+
			" after an '=' rather than as a following argument."+
			" For instance,\n\n-xxx=...\nrather than\n-xxx ...",
		0)
}

// Help prints the messages and then a standardised usage message based on
// the parameters supplied to the param set. It then exits with an exit
// status of 1
func (h StdHelp) Help(ps *param.PSet, messages ...string) {
	w := ps.ErrWriter()
	twc, err := twrap.NewTWConf(twrap.SetWriter(w))
	if err != nil {
		fmt.Fprint(os.Stderr, "Couldn't build the text wrapper:", err)
		return
	}

	for _, message := range messages {
		twc.Wrap(message, 0)
	}
	if len(messages) > 0 {
		fmt.Fprint(w, "\n"+equals+"\n\n")
	}

	if h.style != Short &&
		h.style != GroupNamesOnly {
		twc.Wrap(ps.ProgDesc()+"\n", textIndent)
	}

	if h.includeGroups {
		if badGroups(ps, h.groupsToShow, "groups to show") {
			h.includeGroups = false
			h.style = GroupNamesOnly
		}
	}
	if h.excludeGroups {
		if badGroups(ps, h.groupsToExclude, "excluded groups") {
			h.excludeGroups = false
			h.style = GroupNamesOnly
		}
	}
	if h.groupListCounter.Count() > 1 {
		twc.WrapPrefixed("Error: ",
			"only include OR exclude parameter groups"+
				" not both at the same time."+
				" Excluded groups will be ignored."+
				" They have been set at:"+
				h.groupListCounter.SetBy(),
			0)
		h.excludeGroups = false
	}

	fmt.Fprint(w, "Usage: ", ps.ProgName())

	if h.style == GroupNamesOnly {
		fmt.Fprintln(w, "\nParameter groups")
		h.printGroups(w, ps)
	} else {
		h.printPositionalParams(w, ps)
		h.printParams(w, ps)
	}

	if h.style != Short {
		h.printAlternativeSources(ps)

		h.printOptValNote(w)
	}

	os.Exit(1)
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

func (h StdHelp) printParamUsage(w io.Writer, p *param.ByName) {
	twc, err := twrap.NewTWConf(twrap.SetWriter(w))
	if err != nil {
		fmt.Fprint(os.Stderr, "Couldn't build the text wrapper:", err)
		return
	}
	prefix := "-"
	suffix := valueNeededStr(p.ValueReq())
	if !p.AttrIsSet(param.MustBeSet) {
		prefix = "[ " + prefix
		suffix += " ]"
	}

	paramNames := ""
	sep := ""

	for _, altParamName := range p.AltNames() {
		paramNames += sep + prefix + altParamName + suffix
		sep = " or "
	}
	twc.Wrap(paramNames, paramIndent)

	if h.style == Short {
		return
	}

	twc.Wrap(p.Description(), descriptionIndent)
	twc.WrapPrefixed("Allowed values: ", p.AllowedValues(), descriptionIndent)
	twc.WrapPrefixed("Initial value: ", p.InitialValue(), descriptionIndent)
}

func (h StdHelp) printPositionalParams(w io.Writer, ps *param.PSet) {
	bppCount := ps.CountByPosParams()
	if bppCount == 0 {
		return
	}

	twc, err := twrap.NewTWConf(twrap.SetWriter(w))
	if err != nil {
		fmt.Fprint(os.Stderr, "Couldn't build the text wrapper:", err)
		return
	}
	for i := 0; i < bppCount; i++ {
		bp, _ := ps.GetParamByPos(i)
		fmt.Fprint(w, " <", bp.Name(), ">")
	}

	fmt.Fprintln(w)
	if h.style == Short {
		return
	}
	fmt.Fprintln(w, "where")

	for i := 0; i < bppCount; i++ {
		bp, _ := ps.GetParamByPos(i)
		twc.Wrap(bp.Name(), paramIndent)
		twc.Wrap(bp.Description(), descriptionIndent)
	}
}

// printGroupDetails prints the group name etc
func printGroupDetails(w io.Writer, pg *param.Group, style helpStyle) {
	fmt.Fprintln(w, "\n"+dashes)
	fmt.Fprintf(w, "%s [ ", pg.Name)
	if len(pg.Params) == 1 {
		fmt.Fprint(w, "1 parameter")
	} else {
		fmt.Fprintf(w, "%d parameters", len(pg.Params))
	}
	if pg.HiddenCount > 0 {
		if pg.AllParamsHidden() {
			fmt.Fprint(w, ", all hidden")
		} else {
			fmt.Fprintf(w, ", %d hidden", pg.HiddenCount)
		}
	}
	fmt.Fprintln(w, " ]")
	if style == Short {
		return
	}
	desc := pg.Desc
	if desc == "" {
		return
	}
	twc, err := twrap.NewTWConf(twrap.SetWriter(w))
	if err != nil {
		fmt.Fprint(os.Stderr, "Couldn't build the text wrapper:", err)
		return
	}
	twc.Wrap(desc, textIndent)
	fmt.Fprintln(w)
}

// showGroup will return true if the group should be reported and false
// otherwise
func (h StdHelp) showGroup(g string) bool {
	if h.includeGroups && !h.groupsToShow[g] {
		return false
	}
	if h.excludeGroups && h.groupsToExclude[g] {
		return false
	}
	return true
}

func (h StdHelp) printParams(w io.Writer, ps *param.PSet) {
	paramGroups := ps.GetGroups()

	for _, pg := range paramGroups {
		if !h.showGroup(pg.Name) {
			continue
		}

		if pg.AllParamsHidden() && !h.showAllParams {
			continue
		}

		printGroupDetails(w, pg, h.style)

		for _, p := range pg.Params {
			if p.AttrIsSet(param.DontShowInStdUsage) &&
				!h.showAllParams {
				continue
			}
			h.printParamUsage(w, p)
		}

		printGroupConfigFile(w, pg)
	}
}

func printGroupConfigFile(w io.Writer, pg *param.Group) {
	if len(pg.ConfigFiles) > 0 {
		msg := "\nParameters in this group may also be set "

		msg += altSrcConfigFiles(pg.ConfigFiles)

		twc, err := twrap.NewTWConf(twrap.SetWriter(w))
		if err != nil {
			fmt.Fprint(os.Stderr, "Couldn't build the text wrapper:", err)
			return
		}

		twc.Wrap(msg, textIndent)
	}
}

func (h StdHelp) printGroups(w io.Writer, ps *param.PSet) {
	paramGroups := ps.GetGroups()

	for _, pg := range paramGroups {
		if h.showGroup(pg.Name) {
			printGroupDetails(w, pg, h.style)
		}
	}
}

// printAlternativeSources prints the name(s) of the configuration file
// and any environment variable prefixes (if set)
func (h StdHelp) printAlternativeSources(ps *param.PSet) {
	ep := ps.EnvPrefixes()
	var hasEnvPrefixes bool
	if len(ep) > 0 {
		hasEnvPrefixes = true
	}

	cf := ps.ConfigFiles()
	var hasConfigFiles bool
	if len(cf) > 0 {
		hasConfigFiles = true
	}

	if hasConfigFiles || hasEnvPrefixes {
		message := "\n" + dashes + "\nAny of these parameters may also be set "

		message += altSrcConfigFiles(cf)

		if hasConfigFiles && hasEnvPrefixes {
			message += "    or "
		}

		message += altSrcEnvVars(ep) + "\n"

		twc, err := twrap.NewTWConf(twrap.SetWriter(ps.ErrWriter()))
		if err != nil {
			fmt.Fprint(os.Stderr, "Couldn't build the text wrapper:", err)
			return
		}

		twc.Wrap(message, 0)
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
	epLen := len(ep)
	if epLen == 0 {
		return ""
	}

	message := "through environment variables prefixed with"
	if epLen == 1 {
		message += ": "
	} else {
		message += " one of: "
	}
	sep := ""
	for i, pfx := range ep {
		message += sep + pfx
		sep = ", "
		if i == (epLen - 2) {
			sep = " or "
		}
	}
	message += `

The prefix is stripped off and any underscores ('_') in the environment variable name after the prefix will be replaced with dashes ('-') when matching the parameter name.

For instance, if the environment variables prefixes include 'XX_' an environment variable called 'XX_a_b' will match a parameter called 'a-b'
`

	return message
}
