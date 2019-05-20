package phelp

import (
	"crypto/md5"

	"github.com/nickwells/param.mod/v3/param/paction"
)

// helpStyle records the style of help message to generate - it can be
// standard, full or short
type helpStyle byte

// helpStyle values can be one of:
//    Std meaning that the standard help message is generated
//    Full meaning that a more complete message is generated
//    Short meaning that a concise message is generated
//    GroupNamesOnly meaning that only group names should be shown
const (
	Std helpStyle = iota
	Short
	GroupNamesOnly
)

// StdHelp implements the Helper interface. It adds the standard arguments
// and processes them. This is the helper you are most likely to want and it
// is the one that is used by the paramset.New func.
type StdHelp struct {
	groupsToShow     map[string]bool
	groupsToExclude  map[string]bool
	groupListCounter paction.Counter
	includeGroups    bool
	excludeGroups    bool

	avalShownAlready map[[md5.Size]byte]string

	reportWhereParamsAreSet bool
	reportUnusedParams      bool
	reportParamSources      bool

	dontReportErrors bool
	dontExitOnErrors bool

	exitAfterParsing bool

	showHelp      bool
	showAllParams bool

	style helpStyle
}

// SH is the instance of the standard help type
var SH = StdHelp{
	groupsToShow:     make(map[string]bool),
	groupsToExclude:  make(map[string]bool),
	avalShownAlready: make(map[[md5.Size]byte]string),
}
