package phelp

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/paction"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

const (
	helpArgName           = "help"
	helpShowArgName       = "help-show"
	helpFullArgName       = "help-full"
	helpShowHiddenArgName = "help-all"
	helpSummaryArgName    = "help-summary"
	helpGroupsArgName     = "help-groups"
	helpParamsArgName     = "help-params"
	helpNotesArgName      = "help-notes"
	helpFormatArgName     = "help-format"
)

const (
	helpFmtTypeStd = "standard"
	helpFmtTypeMD  = "markdown"
)

const (
	exitAfterHelpMessage = "\n\nThe program will exit" +
		" after the help message is shown."
)

type trimDashes struct{}

// Edit will remove any leading dashes ("-") from a parameter name and return
// it.
func (trimDashes) Edit(_, val string) (string, error) {
	return strings.TrimLeft(val, "-"), nil
}

// addUsageParams will add the usage parameters into the parameter set
func (h *StdHelp) addUsageParams(ps *param.PSet) {
	groupName := groupNamePfx + "-help"

	ps.AddGroup(groupName,
		"These are parameters for printing a help message.")

	ps.Add(helpArgName, psetter.Nil{},
		"print this help message and exit."+
			"\n\n"+
			"To see hidden parameters use the -"+helpShowHiddenArgName+
			" parameter."+
			"\n"+
			"For a brief help message use the -"+helpSummaryArgName+
			" parameter",
		param.Attrs(param.CommandLineOnly),
		param.AltName("usage"),
		param.PostAction(setHelpSections(h, standardHelpSectionNames)),
		param.GroupName(groupName))

	ps.Add(helpFullArgName, psetter.Nil{},
		" show all parts of the help message and all"+
			" parameters, including hidden ones."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.PostAction(setHelpSections(h, allHelpSectionNames)),
		param.PostAction(paction.SetBool(&h.showHiddenItems, true)),
		param.GroupName(groupName))

	ps.Add(helpShowHiddenArgName, psetter.Nil{},
		" show all the parameters."+
			" Less commonly useful parameters are not shown in the"+
			" standard help message. This will reveal them."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-a"),
		param.PostAction(paction.SetBool(&h.showHiddenItems, true)),
		param.PostAction(paction.SetBool(&h.helpRequested, true)),
		param.GroupName(groupName))

	ps.Add(helpSummaryArgName, psetter.Nil{},
		"print a shorter help message. Only minimal details"+
			" are show, descriptions are not shown."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-s"),
		param.AltName("help-short"),
		param.PostAction(paction.SetBool(&h.hideDescriptions, true)),
		param.PostAction(paction.SetBool(&h.helpRequested, true)),
		param.GroupName(groupName))

	ps.Add("help-all-short", psetter.Nil{},
		"print a shorter help message but with all the"+
			" parameters shown. This is the equivalent"+
			" of giving both the "+helpShowHiddenArgName+
			" and the "+helpSummaryArgName+" parameters."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-as"),
		param.AltName("help-sa"),
		param.PostAction(paction.SetBool(&h.showHiddenItems, true)),
		param.PostAction(paction.SetBool(&h.hideDescriptions, true)),
		param.PostAction(paction.SetBool(&h.helpRequested, true)),
		param.GroupName(groupName))

	ps.Add(helpGroupsArgName,
		psetter.Map{
			Value: (*map[string]bool)(&h.groupsChosen),
			Checks: []check.MapStringBool{
				check.MapStringBoolTrueCountGT(0),
			},
		},
		"when printing the help message only show the listed groups."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-group"),
		param.AltName("help-g"),
		param.PostAction(checkGroups(h, ps)),
		param.GroupName(groupName))

	ps.Add(helpParamsArgName,
		psetter.Map{
			Value: (*map[string]bool)(&h.paramsChosen),
			Checks: []check.MapStringBool{
				check.MapStringBoolTrueCountGT(0),
			},
			Editor: trimDashes{},
		},
		"when printing the help message only show the listed parameters."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-param"),
		param.AltName("help-p"),
		param.PostAction(checkParams(h, ps)),
		param.GroupName(groupName))

	ps.Add(helpNotesArgName,
		psetter.Map{
			Value: (*map[string]bool)(&h.notesChosen),
			Checks: []check.MapStringBool{
				check.MapStringBoolTrueCountGT(0),
			},
		},
		"when printing the help message only show the listed notes."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-note"),
		param.AltName("help-n"),
		param.PostAction(checkNotes(h, ps)),
		param.PostAction(setHelpSections(h, notesHelpSectionName)),
		param.GroupName(groupName))

	ps.Add(helpShowArgName,
		psetter.EnumMap{
			Value:       (*map[string]bool)(&h.sectionsChosen),
			AllowedVals: makeSectionAllowedVals(),
			Aliases:     sectionAliases,
		},
		"specify the parts of the help message you wish to see",
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName),
	)

	ps.Add(helpFormatArgName,
		psetter.Enum{
			Value: &h.helpFormat,
			AllowedVals: psetter.AllowedVals{
				helpFmtTypeStd: "the standard format." +
					" This is almost certainly what you want",
				helpFmtTypeMD: "markdown format. This will have markdown" +
					" annotations applied. This can be useful to produce" +
					" online documentation",
			},
		},
		"specify how the help message should be produced. Only some parts"+
			" of the help message support this feature. They will mostly"+
			" produce Standard format regardless of this setting.",
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName),
	)

	// Final checks

	ps.AddFinalCheck(
		makeForceChosenSection(h, h.groupsChosen, groupedParamsHelpSectionName,
			groupsHelpSectionName, groupedParamsHelpSectionName))
	ps.AddFinalCheck(
		makeForceChosenSection(h, h.paramsChosen, namedParamsHelpSectionName,
			namedParamsHelpSectionName, groupedParamsHelpSectionName))
	ps.AddFinalCheck(
		makeForceChosenSection(h, h.notesChosen, notesHelpSectionName,
			notesHelpSectionName))
	ps.AddFinalCheck(makeForceDefaultSectionsFunc(h))
}

// makeForceChosenSection returns a FinalCheckFunction that checks to see if
// any choices have been made and if so whether any of the altSections are in
// the list of help sections. If not the forceSelection is applied.
func makeForceChosenSection(h *StdHelp, c choices, dflt string, alts ...string) param.FinalCheckFunc {
	return func() error {
		if c.hasNothingChosen() {
			return nil
		}
		for _, s := range alts {
			if h.sectionsChosen[s] {
				return nil
			}
		}
		return h.setHelpSections(dflt)
	}
}

// makeForceDefaultSectionsFunc returns a FinalCheckFunction that will set
// the help sections to a sensible default value if help has been implicitly
// requested but no help sections have been set
func makeForceDefaultSectionsFunc(h *StdHelp) param.FinalCheckFunc {
	return func() error {
		if !h.helpRequested {
			return nil
		}
		if h.sectionsChosen.hasNothingChosen() {
			return h.setHelpSections(standardHelpSectionNames)
		}
		return nil
	}
}

// checkGroups returns an ActionFunc which will check that the groupsChosen
// element of the StdHelp structure only contains valid group names
func checkGroups(h *StdHelp, ps *param.PSet) param.ActionFunc {
	return func(_ location.L, _ *param.ByName, _ []string) error {
		var badNames []string
		for gName := range h.groupsChosen {
			if !ps.HasGroupName(gName) {
				badNames = append(badNames, fmt.Sprintf("%q", gName))
			}
		}
		if len(badNames) == 0 {
			return nil
		}
		if len(badNames) == len(h.groupsChosen) {
			for k := range h.groupsChosen {
				delete(h.groupsChosen, k)
			}
			err := h.unsetHelpSections(
				groupsHelpSectionName, groupedParamsHelpSectionName)
			if err != nil {
				return err
			}
		}

		altNames := altNames(h, ps, badNames, param.SuggestGroups)
		return makeBadNameError(badNames, altNames, "group",
			" For a list of available group names try '-"+
				helpShowArgName+" "+groupsHelpSectionName+"'")
	}
}

// checkNotes returns an ActionFunc which will check that the notesChosen
// element of the StdHelp structure only contains valid note names
func checkNotes(h *StdHelp, ps *param.PSet) param.ActionFunc {
	return func(_ location.L, _ *param.ByName, _ []string) error {
		var badNames []string
		for n := range h.notesChosen {
			if _, err := ps.GetNote(n); err != nil {
				badNames = append(badNames, n)
			}
		}
		if len(badNames) == 0 {
			return nil
		}
		if len(badNames) == len(h.notesChosen) {
			for k := range h.notesChosen {
				delete(h.notesChosen, k)
			}
			err := h.unsetHelpSections(notesHelpSectionName)
			if err != nil {
				return err
			}
		}

		altNames := altNames(h, ps, badNames, param.SuggestNotes)
		return makeBadNameError(badNames, altNames, "note", "")
	}
}

// checkParams returns an ActionFunc which will check that the paramsChosen
// element of the StdHelp structure only contains valid parameter names
func checkParams(h *StdHelp, ps *param.PSet) param.ActionFunc {
	return func(_ location.L, _ *param.ByName, _ []string) error {
		var badNames []string
		for pName := range h.paramsChosen {
			trimmedName := strings.TrimLeft(pName, "-")
			if _, err := ps.GetParamByName(trimmedName); err != nil {
				badNames = append(badNames, fmt.Sprintf("%q", pName))
			}
		}
		if len(badNames) == 0 {
			return nil
		}
		if len(badNames) == len(h.paramsChosen) {
			for k := range h.paramsChosen {
				delete(h.paramsChosen, k)
			}
			err := h.unsetHelpSections(
				namedParamsHelpSectionName, groupedParamsHelpSectionName)
			if err != nil {
				return err
			}
		}

		altNames := altNames(h, ps, badNames, param.SuggestParams)
		return makeBadNameError(badNames, altNames, "parameter",
			" To see a list of parameter names try "+
				"'-"+helpShowArgName+" "+namedParamsHelpSectionName+"'."+
				" To see a full list, add '-"+helpShowHiddenArgName+"'."+
				" To see just the names, add '-"+helpSummaryArgName+"'.")
	}
}

// altNames returns a list of possible alternative names using the passed
// suggestion function
func altNames(h *StdHelp, ps *param.PSet, badNames []string, sf param.SuggestionFunc) []string {
	altNames := make([]string, 0)
	alreadyAdded := map[string]bool{}

	for _, n := range badNames {
		suggestions := sf(ps, n)
		for _, s := range suggestions {
			if !h.paramsChosen[s] && !alreadyAdded[s] {
				alreadyAdded[s] = true
				altNames = append(altNames, fmt.Sprintf("%q", s))
			}
		}
	}
	return altNames
}

// makeBadNameError will create an error formatting the bad and alternate names
func makeBadNameError(badNames, altNames []string, tName, extra string) error {
	if len(badNames) == 0 {
		return nil
	}
	badStr := ""
	switch len(badNames) {
	case 0:
		return nil
	case 1:
		badStr = fmt.Sprintf("Bad %s name: %s.",
			tName, badNames[0])
	default:
		sort.Strings(badNames)
		badStr = fmt.Sprintf("Bad %s names: %s.",
			tName, strings.Join(badNames, ", "))
	}
	alts := ""
	switch len(altNames) {
	case 0:
		alts = ""
	case 1:
		alts = " Did you mean " + altNames[0]
	default:
		sort.Strings(altNames)
		alts = " Did you mean " + strings.Join(altNames, " or ")
	}
	return errors.New(badStr + alts + extra)
}
