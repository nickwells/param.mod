package psetter

import (
	"sort"

	"github.com/nickwells/english.mod/english"
	"github.com/nickwells/strdist.mod/v2/strdist"
)

// SuggestionString returns a string suggesting the supplied values or the
// empty string if there are no values.
func SuggestionString(vals []string) string {
	const intro = ", did you mean "

	if len(vals) == 0 {
		return ""
	}

	if len(vals) == 1 {
		return intro + `"` + vals[0] + `"?`
	}

	sort.Strings(vals)

	return intro + english.JoinQuoted(vals, ", ", " or ", `"`, `"`) + "?"
}

// SuggestedVals returns a slice of suggested alternative values for
// the given value
func SuggestedVals(val string, alts []string) []string {
	const alternativeCount = 3

	finder := strdist.DefaultFinders[strdist.CaseBlindAlgoNameCosine]

	return finder.FindNStrLike(alternativeCount, val, alts...)
}
