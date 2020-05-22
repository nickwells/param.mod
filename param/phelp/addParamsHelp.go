package phelp

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v4/param"
	"github.com/nickwells/param.mod/v4/param/paction"
	"github.com/nickwells/param.mod/v4/param/psetter"
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
		param.PostAction(setHelpSections(h, standardSections)),
		param.GroupName(groupName))

	ps.Add(helpFullArgName, psetter.Nil{},
		" show all parts of the help message and all"+
			" parameters, including hidden ones."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.PostAction(setHelpSections(h, allSections)),
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
			Value: &h.groupsChosen,
			Checks: []check.MapStringBool{
				check.MapStringBoolTrueCountGT(0),
			},
		},
		"when printing the help message only show the listed groups."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-g"),
		param.PostAction(checkGroups(h, ps)),
		param.GroupName(groupName))

	ps.Add(helpParamsArgName,
		psetter.Map{
			Value: &h.paramsChosen,
			Checks: []check.MapStringBool{
				check.MapStringBoolTrueCountGT(0),
			},
			Editor: trimDashes{},
		},
		"when printing the help message only show the listed parameters."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-p"),
		param.PostAction(checkParams(h, ps)),
		param.GroupName(groupName))

	ps.Add(helpNotesArgName,
		psetter.Map{
			Value: &h.notesChosen,
			Checks: []check.MapStringBool{
				check.MapStringBoolTrueCountGT(0),
			},
		},
		"when printing the help message only show the listed notes."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-n"),
		param.PostAction(checkNotes(h, ps)),
		param.PostAction(setHelpSections(h, notesSection)),
		param.GroupName(groupName))

	ps.Add(helpShowArgName,
		psetter.EnumMap{
			Value:       &h.sectionsChosen,
			AllowedVals: makeSectionAllowedVals(),
			Aliases:     sectionAliases,
		},
		"specify the parts of the help message you wish to see",
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName),
	)

	// Final checks

	ps.AddFinalCheck(
		makeForceChosenSection(h, h.groupsChosen, groupsSection,
			groupsSection, groupedParamsSection))
	ps.AddFinalCheck(
		makeForceChosenSection(h, h.paramsChosen, namedParamsSection,
			namedParamsSection, groupedParamsSection))
	ps.AddFinalCheck(
		makeForceChosenSection(h, h.notesChosen, notesSection,
			notesSection))
	ps.AddFinalCheck(makeForceDefaultSectionsFunc(h))
}

// makeForceChosenSection returns a FinalCheckFunction that checks to see if
// any choices have been made and if so whether any of the altSections are in
// the list of help sections. If not the forceSelection is applied.
func makeForceChosenSection(h *StdHelp, choices map[string]bool, forceSection string, altSections ...string) param.FinalCheckFunc {
	return func() error {
		if len(choices) == 0 {
			return nil
		}
		for _, s := range altSections {
			if h.sectionsChosen[s] {
				return nil
			}
		}
		return h.setHelpSections(forceSection)
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
		if len(h.sectionsChosen) > 0 {
			return nil
		}
		return h.setHelpSections(standardSections)
	}
}

// checkGroups returns an ActionFunc which will check that the groupsChosen
// element of the StdHelp structure only contains valid group names
func checkGroups(h *StdHelp, ps *param.PSet) param.ActionFunc {
	// TODO: make this a FinalCheckFunc (and combine with makeForceChosenSection)
	// TODO: add a "did you mean..." suggestion - see strdist
	return func(_ location.L, _ *param.ByName, _ []string) error {
		var badNames []string
		var goodNames int
		for gName := range h.groupsChosen {
			if !ps.HasGroupName(gName) {
				badNames = append(badNames, fmt.Sprintf("%q", gName))
			} else {
				goodNames++
			}
		}
		if len(badNames) == 0 {
			return nil
		}
		if goodNames == 0 {
			h.unsetHelpSections(groupsSection, groupedParamsSection)
		}

		altNames := make([]string, 0)
		for _, g := range ps.GetGroups() {
			if !h.groupsChosen[g.Name] {
				altNames = append(altNames, fmt.Sprintf("%q", g.Name))
			}
		}
		return makeBadNameError(badNames, altNames, "group")
	}
}

// checkNotes returns an ActionFunc which will check that the notesChosen
// element of the StdHelp structure only contains valid note names
func checkNotes(h *StdHelp, ps *param.PSet) param.ActionFunc {
	// TODO: make this a FinalCheckFunc (and combine with makeForceChosenSection)
	// TODO: add a "did you mean..." suggestion - see strdist
	return func(_ location.L, _ *param.ByName, _ []string) error {
		var badNames []string
		var goodNames int
		for n := range h.notesChosen {
			if _, err := ps.GetNote(n); err != nil {
				badNames = append(badNames, n)
			} else {
				goodNames++
			}
		}
		if len(badNames) == 0 {
			return nil
		}
		if goodNames == 0 {
			h.unsetHelpSections(notesSection)
		}

		altNames := make([]string, 0)
		for _, n := range ps.Notes() {
			if !h.notesChosen[n.Headline] {
				altNames = append(altNames, n.Headline)
			}
		}
		return makeBadNameError(badNames, altNames, "note")
	}
}

// checkParams returns an ActionFunc which will check that the paramsChosen
// element of the StdHelp structure only contains valid parameter names
func checkParams(h *StdHelp, ps *param.PSet) param.ActionFunc {
	// TODO: make this a FinalCheckFunc (and combine with makeForceChosenSection)
	// TODO: add a "did you mean..." suggestion - see strdist
	return func(_ location.L, _ *param.ByName, _ []string) error {
		var badNames []string
		var goodNames int
		for pName := range h.paramsChosen {
			trimmedName := strings.TrimLeft(pName, "-")
			if _, err := ps.GetParamByName(trimmedName); err != nil {
				badNames = append(badNames, fmt.Sprintf("%q", pName))
			} else {
				goodNames++
			}
		}
		if len(badNames) == 0 {
			return nil
		}
		if goodNames == 0 {
			h.unsetHelpSections(namedParamsSection, groupedParamsSection)
		}

		altNames := make([]string, 0)
		alreadyAdded := map[string]bool{}
		for _, n := range badNames {
			suggestions := ps.FindClosestMatches(n)
			for _, s := range suggestions {
				if !h.paramsChosen[s] && !alreadyAdded[s] {
					alreadyAdded[s] = true
					altNames = append(altNames, s)
				}
			}
		}
		return makeBadNameError(badNames, altNames, "parameter")
	}
}

// makeBadNameError will create an error formatting the bad and alternate names
func makeBadNameError(badNames, altNames []string, nameType string) error {
	if len(badNames) == 0 {
		return nil
	}
	badStr := ""
	switch len(badNames) {
	case 0:
		return nil
	case 1:
		badStr = fmt.Sprintf("Bad %s name: %s.",
			nameType, badNames[0])
	default:
		sort.Strings(badNames)
		badStr = fmt.Sprintf("Bad %s names: %s.",
			nameType, strings.Join(badNames, ", "))
	}
	altStr := ""
	switch len(altNames) {
	case 0:
		altStr = ""
	case 1:
		altStr = fmt.Sprintf(" A possible %s name is %s.",
			nameType, altNames[0])
	default:
		sort.Strings(altNames)
		altStr = fmt.Sprintf(" Possible %s names are: %s.",
			nameType, strings.Join(altNames, ", "))
	}
	return errors.New(badStr + altStr)
}
