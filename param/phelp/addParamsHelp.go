package phelp

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/english.mod/english"
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
	helpNoPageArgName     = "help-no-page"
)

const (
	helpFmtTypeStd      = "standard"
	helpFmtTypeMarkdown = "markdown"
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
			" parameter"+
			"\n"+
			"For the full help message use the -"+helpFullArgName+
			" parameter",
		param.Attrs(param.CommandLineOnly),
		param.AltNames("usage"),
		param.PostAction(setHelpSections(h, standardHelpSectionAlias)),
		param.GroupName(groupName))

	ps.Add(helpFullArgName, psetter.Nil{},
		"show all parts of the help message and all"+
			" parameters, including hidden ones."+
			exitAfterHelpMessage,
		param.AltNames("help-f"),
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.PostAction(setHelpSections(h, allHelpSectionAlias)),
		param.PostAction(paction.SetVal(&h.showHiddenItems, true)),
		param.GroupName(groupName))

	ps.Add(helpShowHiddenArgName, psetter.Nil{},
		"show all the parameters."+
			" Less commonly useful parameters are not shown in the"+
			" standard help message. This will reveal them."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltNames("help-a"),
		param.PostAction(paction.SetVal(&h.showHiddenItems, true)),
		param.PostAction(paction.SetVal(&h.helpRequested, true)),
		param.GroupName(groupName))

	ps.Add(helpSummaryArgName, psetter.Nil{},
		"print a shorter help message. Only minimal details"+
			" are shown, descriptions are not shown."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltNames("help-s", "help-short"),
		param.PostAction(paction.SetVal(&h.showSummary, true)),
		param.PostAction(paction.SetVal(&h.helpRequested, true)),
		param.GroupName(groupName))

	ps.Add("help-all-short", psetter.Nil{},
		"print a shorter help message but with all the"+
			" parameters shown. This is the equivalent"+
			" of giving both the "+helpShowHiddenArgName+
			" and the "+helpSummaryArgName+" parameters."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltNames("help-as", "help-sa"),
		param.PostAction(paction.SetVal(&h.showHiddenItems, true)),
		param.PostAction(paction.SetVal(&h.showSummary, true)),
		param.PostAction(paction.SetVal(&h.helpRequested, true)),
		param.GroupName(groupName))

	{
		boolCounter := check.NewCounter(check.ValEQ(true), check.ValGT(0))

		ps.Add(helpGroupsArgName,
			psetter.Map{
				Value: (*map[string]bool)(&h.groupsChosen),
				Checks: []check.MapStringBool{
					check.MapValAggregate[map[string]bool, string, bool](
						boolCounter),
				},
			},
			"when printing the help message only show the listed groups."+
				" This will also force hidden parameters to be shown."+
				exitAfterHelpMessage,
			param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
			param.AltNames("help-group", "help-g"),
			param.ValueName("group-name,..."),
			param.PostAction(checkGroups(h, ps)),
			param.PostAction(paction.SetVal(&h.showHiddenItems, true)),
			param.GroupName(groupName))
	}

	{
		boolCounter := check.NewCounter(check.ValEQ(true), check.ValGT(0))

		ps.Add(helpParamsArgName,
			psetter.Map{
				Value: (*map[string]bool)(&h.paramsChosen),
				Checks: []check.MapStringBool{
					check.MapValAggregate[map[string]bool, string, bool](
						boolCounter),
				},
				Editor: trimDashes{},
			},
			"when printing the help message only show the listed parameters."+
				exitAfterHelpMessage,
			param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
			param.AltNames("help-param", "help-p"),
			param.ValueName("param-name,..."),
			param.PostAction(checkParams(h, ps)),
			param.GroupName(groupName))
	}

	{
		boolCounter := check.NewCounter(check.ValEQ(true), check.ValGT(0))

		ps.Add(helpNotesArgName,
			psetter.Map{
				Value: (*map[string]bool)(&h.notesChosen),
				Checks: []check.MapStringBool{
					check.MapValAggregate[map[string]bool, string, bool](
						boolCounter),
				},
			},
			"when printing the help message only show the listed notes."+
				exitAfterHelpMessage,
			param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
			param.AltNames("help-note", "help-n"),
			param.ValueName("note-name,..."),
			param.PostAction(checkNotes(h, ps)),
			param.PostAction(setHelpSections(h, notesHelpSectionName)),
			param.GroupName(groupName))
	}

	ps.Add(helpShowArgName,
		psetter.EnumMap{
			Value:       (*map[string]bool)(&h.sectionsChosen),
			AllowedVals: makeSectionAllowedVals(),
			Aliases:     sectionAliases,
		},
		"specify the parts of the help message you wish to see",
		param.ValueName("part,..."),
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName),
	)

	ps.Add(helpFormatArgName,
		psetter.Enum{
			Value: &h.helpFormat,
			AllowedVals: psetter.AllowedVals{
				helpFmtTypeStd: "the standard format." +
					" This is almost certainly what you want",
				helpFmtTypeMarkdown: "markdown format. This will have" +
					" markdown annotations applied. This can be useful" +
					" to produce online documentation",
			},
		},
		"specify how the help message should be produced. Only some parts"+
			" of the help message support this feature. They will mostly"+
			" produce Standard format regardless of this setting.",
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName),
	)

	ps.Add(helpNoPageArgName,
		psetter.Bool{
			Value:  &h.pageOutput,
			Invert: true,
		},
		"show help but don't page the output. Without this parameter the"+
			" help message will be paged using the standard pager"+
			" (as given by the value of the 'PAGER' environment"+
			" variable or 'less' if 'PAGER' is not set or the command"+
			" it refers to cannot be found)",
		param.AltNames("help-dont-page", "help-no-pager"),
		param.GroupName(groupName),
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.PostAction(paction.SetVal(&h.helpRequested, true)),
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
			return h.setHelpSections(standardHelpSectionAlias)
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
				delete(h.groupsChosen, gName)
				matches, _ := ps.FindMatchingGroups(gName)
				if len(matches) > 0 {
					for _, gn := range matches {
						h.groupsChosen[gn] = true
					}
					continue
				}
				badNames = append(badNames, fmt.Sprintf("%q", gName))
			}
		}

		if len(badNames) == 0 {
			return nil
		}
		if len(h.groupsChosen) == 0 {
			err := h.unsetHelpSections(
				groupsHelpSectionName, groupedParamsHelpSectionName)
			if err != nil {
				return err
			}
		}

		altNames := altNames(h, ps, badNames, param.SuggestGroups)
		return makeBadNameError(badNames, altNames, "group",
			"\nFor a list of available group names try '-"+
				helpShowArgName+" "+groupsHelpSectionName+"'")
	}
}

// checkNotes returns an ActionFunc which will check that the notesChosen
// element of the StdHelp structure only contains valid note names
func checkNotes(h *StdHelp, ps *param.PSet) param.ActionFunc {
	return func(_ location.L, _ *param.ByName, _ []string) error {
		var badNames []string
		for nName := range h.notesChosen {
			if _, err := ps.GetNote(nName); err != nil {
				delete(h.notesChosen, nName)
				matches, _ := ps.FindMatchingNotes(nName)
				if len(matches) > 0 {
					for _, nn := range matches {
						h.notesChosen[nn] = true
					}
					continue
				}
				badNames = append(badNames, nName)
			}
		}
		if len(badNames) == 0 {
			return nil
		}
		if len(h.notesChosen) == 0 {
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
				delete(h.paramsChosen, pName)
				matches, _ := ps.FindMatchingNamedParams(trimmedName)
				if len(matches) > 0 {
					for _, pn := range matches {
						h.paramsChosen[pn] = true
					}
					continue
				}
				badNames = append(badNames, fmt.Sprintf("%q", pName))
			}
		}

		if len(badNames) == 0 {
			return nil
		}
		if len(h.paramsChosen) == 0 {
			err := h.unsetHelpSections(
				namedParamsHelpSectionName, groupedParamsHelpSectionName)
			if err != nil {
				return err
			}
		}

		altNames := altNames(h, ps, badNames, param.SuggestParams)
		return makeBadNameError(badNames, altNames, "parameter",
			"\nTo see a list of parameter names try "+
				"'-"+helpShowArgName+" "+namedParamsHelpSectionName+"'."+
				"\nTo see a full list, add '-"+helpShowHiddenArgName+"'."+
				"\nTo see just the names, add '-"+helpSummaryArgName+"'.")
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
	badStrPrefix := fmt.Sprintf("Bad %s %s: ",
		tName, english.Plural("name", len(badNames)))
	badStrJoin := "\n" + strings.Repeat(" ", len(badStrPrefix))
	sort.Strings(badNames)
	badStr := badStrPrefix + strings.Join(badNames, badStrJoin)

	alts := ""
	if len(altNames) > 0 {
		sort.Strings(altNames)
		alts = "\nDid you mean " + strings.Join(altNames, " or ")
	}
	return errors.New(badStr + alts + extra)
}
