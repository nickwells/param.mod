package phelp

import (
	"crypto/md5"
	"fmt"

	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v4/param"
	"github.com/nickwells/param.mod/v4/param/psetter"
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
	introSection         = "intro"
	usageSection         = "usage"
	posParamsSection     = "pos-params"
	groupsSection        = "groups"
	namedParamsSection   = "named-params"
	groupedParamsSection = "grouped-params"
	notesSection         = "notes"
	sourcesSection       = "sources"
	examplesSection      = "examples"
	refsSection          = "refs"
)

var helpSectionsInOrder = []helpSection{
	{
		name: introSection,
		desc: "the program name and" +
			" optionally the program description",
		displayFunc: showIntro,
	},
	{
		name: usageSection,
		desc: "the program name, a parameter summary," +
			" and any trailing parameters",
		displayFunc: showUsageSummary,
	},
	{
		name: posParamsSection,
		desc: "the positional parameters coming just after the" +
			" program name",
		displayFunc: showByPosParams,
	},
	{
		name:        groupsSection,
		desc:        "the parameter groups",
		displayFunc: showGroups,
	},
	{
		name:        namedParamsSection,
		desc:        "the named parameters (flags)",
		displayFunc: showParamsByName,
	},
	{
		name:        groupedParamsSection,
		desc:        "the named parameters by group name",
		displayFunc: showParamsByGroupName,
	},
	{
		name:        notesSection,
		desc:        "additional notes on the program behaviour",
		displayFunc: showNotes,
	},
	{
		name: sourcesSection,
		desc: "any additional sources of parameter" +
			" values such as environment variables" +
			" or configuration files",
		displayFunc: showAltSources,
	},
	{
		name: examplesSection,
		desc: "examples of correct program use" +
			" and suggestions of ways to use the" +
			" program",
		displayFunc: showExamples,
	},
	{
		name: refsSection,
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
	standardSections = "std"
	paramSections    = "params"
	allSections      = "all"
)

var sectionAliases = psetter.Aliases{
	paramSections: []string{
		posParamsSection, groupedParamsSection},
	standardSections: []string{introSection, usageSection,
		posParamsSection, groupedParamsSection},
	allSections: []string{introSection, usageSection,
		posParamsSection, groupedParamsSection,
		notesSection, sourcesSection, examplesSection, refsSection},
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
	sectionsChosen map[string]bool
	groupsChosen   map[string]bool
	paramsChosen   map[string]bool
	notesChosen    map[string]bool

	showHiddenItems  bool
	hideDescriptions bool
	helpRequested    bool

	avalShownAlready map[[md5.Size]byte]string

	// params-... values
	paramsShowWhereSet bool
	paramsShowUnused   bool
	reportErrors       bool
	exitOnErrors       bool
	exitAfterParsing   bool

	exitAfterHelp bool // this can only be set in test code

	// completions-... values
	zshCompletionsDir  string
	zshMakeCompletions string
}

// NewStdHelp returns a pointer to a well-constructed instance of the
// standard help type ready to be used as the helper for a new param.PSet
// (the standard paramset.New() function will use this)
func NewStdHelp() *StdHelp {
	return &StdHelp{
		groupsChosen:   make(map[string]bool),
		paramsChosen:   make(map[string]bool),
		sectionsChosen: make(map[string]bool),

		avalShownAlready: make(map[[md5.Size]byte]string),

		reportErrors:  true,
		exitOnErrors:  true,
		exitAfterHelp: true,

		zshMakeCompletions: zshCompGenNone,
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
