package ptypes

import "github.com/nickwells/twrap.mod/twrap"

// ExtraHelper is the interface to be satisfied by a type (typically a
// Setter) that wants to provide more help text.
type ExtraHelper interface {
	ExtraHelp(twc *twrap.TWConf, indent, extraIndent int)
}
