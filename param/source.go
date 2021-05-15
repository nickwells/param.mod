package param

import (
	"fmt"

	"github.com/nickwells/location.mod/location"
)

// Source records where a parameter has been set
type Source struct {
	From      string
	Loc       location.L
	ParamVals []string
	Param     *ByName
}

// String formats a Source into a string
func (pSrc Source) String() string {
	return fmt.Sprintf("Param: %s (at %s)", pSrc.Param.Name(), pSrc.Loc)
}

// Desc describes where the param was set
func (pSrc Source) Desc() string {
	s := pSrc.From + " (at " + pSrc.Loc.String() + ")"

	sep := " ["
	for _, p := range pSrc.ParamVals {
		s += sep + p
		sep = "="
	}
	s += "]"

	return s
}

// Sources is a slice of Source
type Sources []Source

// String formats a slice of Sources into a String
func (pSrcs Sources) String() string {
	var s string
	sep := ""

	for _, ps := range pSrcs {
		s += sep + ps.String()
		sep = ", "
	}

	return s
}
