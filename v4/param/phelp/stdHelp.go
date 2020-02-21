package phelp

import (
	"crypto/md5"

	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v3/param"
)

// helpStyle records the style of help message to generate - it can be
// standard, full or short
type helpStyle byte

// helpStyle values control what help message is generated
const (
	noHelp helpStyle = iota
	stdHelp
	paramsByName
	paramsInGroups
	paramsNotInGroups
	groupNamesOnly
	progDescOnly
	altSourcesOnly
	examplesOnly
	referencesOnly
)

// StdHelp implements the Helper interface. It adds the standard arguments
// and processes them. This is the helper you are most likely to want and it
// is the one that is used by the paramset.New func.
type StdHelp struct {
	groupsSelected map[string]bool
	paramsToShow   []string

	avalShownAlready map[[md5.Size]byte]string

	paramsShowWhereSet bool
	paramsShowUnused   bool

	dontReportErrors bool
	dontExitOnErrors bool
	exitAfterHelp    bool // this can only be set in the test code
	exitAfterParsing bool

	paramsShowHidden  bool
	showFullHelp      bool
	styleNeedsSetting bool // if either of the previous is set then this is set

	style helpStyle

	zshCompletionsDir  string
	zshMakeCompletions string
}

// setStyle returns an ActionFunc to set the style element of the StdHelp
// structure to the given value.
func setStyle(h *StdHelp, setTo helpStyle) param.ActionFunc {
	return func(_ location.L, _ *param.ByName, _ []string) error {
		h.style = setTo
		return nil
	}
}

// NewStdHelp returns a pointer to a well-constructed instance of the
// standard help type ready to be used as the helper for a new param.PSet
// (the standard paramset.New() function will use this)
func NewStdHelp() *StdHelp {
	return &StdHelp{
		groupsSelected:     make(map[string]bool),
		avalShownAlready:   make(map[[md5.Size]byte]string),
		showFullHelp:       true,
		exitAfterHelp:      true,
		zshMakeCompletions: zshCompGenNone,
	}
}
