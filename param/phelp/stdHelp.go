package phelp

import (
	"fmt"

	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/psetter"
	"github.com/nickwells/twrap.mod/twrap"
)

// helpSection records the information about a particular section of the help
// message. The displayFunc prints the appropriate section of the help
// message and returns false if there was nothing displayed, true otherwise
type helpSection struct {
	name        string
	desc        string
	displayFunc func(StdHelp, *twrap.TWConf, *param.PSet) bool
}

const (
	introHelpSectionName         = "intro"
	usageHelpSectionName         = "usage"
	groupsHelpSectionName        = "groups"
	posParamsHelpSectionName     = "params-pos"
	namedParamsHelpSectionName   = "params-named"
	groupedParamsHelpSectionName = "params-grouped"
	notesHelpSectionName         = "notes"
	sourcesHelpSectionName       = "sources"
	examplesHelpSectionName      = "examples"
	refsHelpSectionName          = "refs"
)

var helpSectionsInOrder = []helpSection{
	{
		name: introHelpSectionName,
		desc: "the program name and" +
			" optionally the program description",
		displayFunc: showIntro,
	},
	{
		name: usageHelpSectionName,
		desc: "the program name, a parameter summary," +
			" and any trailing parameters",
		displayFunc: showUsageSummary,
	},
	{
		name: posParamsHelpSectionName,
		desc: "the positional parameters coming just after the" +
			" program name",
		displayFunc: showByPosParams,
	},
	{
		name:        groupsHelpSectionName,
		desc:        "the parameter groups",
		displayFunc: showGroups,
	},
	{
		name:        namedParamsHelpSectionName,
		desc:        "the named parameters (flags)",
		displayFunc: showParamsByName,
	},
	{
		name:        groupedParamsHelpSectionName,
		desc:        "the named parameters by group name",
		displayFunc: showParamsByGroupName,
	},
	{
		name:        notesHelpSectionName,
		desc:        "additional notes on the program behaviour",
		displayFunc: showNotes,
	},
	{
		name: sourcesHelpSectionName,
		desc: "any additional sources of parameter" +
			" values such as environment variables" +
			" or configuration files",
		displayFunc: showAltSources,
	},
	{
		name: examplesHelpSectionName,
		desc: "examples of correct program use" +
			" and suggestions of ways to use the" +
			" program",
		displayFunc: showExamples,
	},
	{
		name: refsHelpSectionName,
		desc: "references to other programs or" +
			" further sources of information",
		displayFunc: showReferences,
	},
}

// makeSectionAllowedVals constructs an AllowedVals map from the
// helpSectionsInOrder slice
func makeSectionAllowedVals() psetter.AllowedVals {
	rval := psetter.AllowedVals{}
	for _, s := range helpSectionsInOrder {
		if _, duplicate := rval[s.name]; duplicate {
			panic(fmt.Sprintf("Bad help section: %q appears twice", s.name))
		}

		rval[s.name] = s.desc
	}
	return rval
}

// Alias names
const (
	standardHelpSectionAlias = "std"
	paramHelpSectionAlias    = "params"
	allHelpSectionAlias      = "all"

	groupHelpSectionAlias = "group"
	grpHelpSectionAlias   = "grp"

	exampleHelpSectionAlias = "example"
	egHelpSectionAlias      = "eg"

	refHelpSectionAlias     = "ref"
	seeAlsoHelpSectionAlias = "see-also"

	posParamsHelpSectionAlias     = "pos-params"
	namedParamsHelpSectionAlias   = "named-params"
	groupedParamsHelpSectionAlias = "grouped-params"
)

var sectionAliases = psetter.Aliases{
	paramHelpSectionAlias: []string{
		posParamsHelpSectionName, groupedParamsHelpSectionName,
	},
	standardHelpSectionAlias: []string{
		introHelpSectionName, usageHelpSectionName,
		posParamsHelpSectionName, groupedParamsHelpSectionName,
	},
	allHelpSectionAlias: []string{
		introHelpSectionName, usageHelpSectionName,
		posParamsHelpSectionName, groupedParamsHelpSectionName,
		notesHelpSectionName, sourcesHelpSectionName,
		examplesHelpSectionName, refsHelpSectionName,
	},

	groupHelpSectionAlias: []string{groupsHelpSectionName},
	grpHelpSectionAlias:   []string{groupsHelpSectionName},

	exampleHelpSectionAlias: []string{examplesHelpSectionName},
	egHelpSectionAlias:      []string{examplesHelpSectionName},

	refHelpSectionAlias:     []string{refsHelpSectionName},
	seeAlsoHelpSectionAlias: []string{refsHelpSectionName},

	posParamsHelpSectionAlias:     []string{posParamsHelpSectionName},
	namedParamsHelpSectionAlias:   []string{namedParamsHelpSectionName},
	groupedParamsHelpSectionAlias: []string{groupedParamsHelpSectionName},
}

type choices map[string]bool

// hasNothingChosen returns true if there is no entry in the choices set to
// true
func (c choices) hasNothingChosen() bool {
	for _, v := range c {
		if v {
			return false
		}
	}
	return true
}

// StdHelp implements the Helper interface. It records the parameter values
// set by the common parameters and proivides methods for generating the help
// message. This is the helper you should use unless you are testing the
// package or have some requirement to provide help other than at the command
// line.
//
// It will be used automatically if you create your param.PSet using the
// paramset.New function (recommended).
type StdHelp struct {
	// help-... values
	sectionsChosen choices
	groupsChosen   choices
	paramsChosen   choices
	notesChosen    choices

	showHiddenItems  bool
	hideDescriptions bool
	helpRequested    bool
	pageOutput       bool

	helpFormat string

	avalShownAlready map[string]string

	// params-... values
	paramsShowWhereSet bool
	paramsShowUnused   bool
	reportErrors       bool
	exitOnErrors       bool
	exitAfterParsing   bool

	exitAfterHelp bool // this can only be set in test code

	// completions-... values
	zshCompDir    string
	zshCompAction string
}

// NewStdHelp returns a pointer to a well-constructed instance of the
// standard help type ready to be used as the helper for a new param.PSet
// (the standard paramset.New() function will use this)
func NewStdHelp() *StdHelp {
	return &StdHelp{
		sectionsChosen: make(choices),
		groupsChosen:   make(choices),
		paramsChosen:   make(choices),
		notesChosen:    make(choices),

		avalShownAlready: make(map[string]string),

		pageOutput:    true,
		reportErrors:  true,
		exitOnErrors:  true,
		exitAfterHelp: true,

		zshCompAction: zshCompActionNone,

		helpFormat: helpFmtTypeStd,
	}
}

// setHelpSections returns an ActionFunc to set the sectionsChosen in the
// StdHelp instance
func setHelpSections(h *StdHelp, sections ...string) param.ActionFunc {
	return func(_ location.L, _ *param.ByName, _ []string) error {
		return h.setHelpSections(sections...)
	}
}

// setHelpSections sets the sectionsChosen in the StdHelp instance
func (h *StdHelp) setHelpSections(sections ...string) error {
	av := makeSectionAllowedVals()
	for _, s := range sections {
		if _, ok := av[s]; ok {
			h.sectionsChosen[s] = true
		} else if secList, ok := sectionAliases[s]; ok {
			for _, s := range secList {
				h.sectionsChosen[s] = true
			}
		} else {
			return fmt.Errorf("%q is not a valid section", s)
		}
	}
	return nil
}

// unsetHelpSections sets to false the sectionsChosen in the StdHelp instance
func (h *StdHelp) unsetHelpSections(sections ...string) error {
	av := makeSectionAllowedVals()
	for _, s := range sections {
		if _, ok := av[s]; ok {
			h.sectionsChosen[s] = false
		} else if secList, ok := sectionAliases[s]; ok {
			for _, s := range secList {
				h.sectionsChosen[s] = false
			}
		} else {
			return fmt.Errorf("%q is not a valid section", s)
		}
	}
	return nil
}
