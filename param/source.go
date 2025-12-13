package param

import (
	"fmt"
	"strings"

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
	var s strings.Builder

	s.WriteString(pSrc.From)
	s.WriteString(" (at ")
	s.WriteString(pSrc.Loc.String())
	s.WriteString(")")

	sep := " ["
	for _, p := range pSrc.ParamVals {
		s.WriteString(sep)
		s.WriteString(p)

		sep = "="
	}

	s.WriteString("]")

	return s.String()
}

// Sources is a slice of Source
type Sources []Source

// String formats a slice of Sources into a String
func (pSrcs Sources) String() string {
	var s strings.Builder

	sep := ""

	for _, ps := range pSrcs {
		s.WriteString(sep)
		s.WriteString(ps.String())

		sep = ", "
	}

	return s.String()
}
