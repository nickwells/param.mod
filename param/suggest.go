package param

import "github.com/nickwells/strdist.mod/v2/strdist"

// SuggestionFunc is the type of a function that returns a slice of suggested
// alternative names for the given string.
type SuggestionFunc func(ps *PSet, s string) []string

const alternativeCount = 3

// SuggestParams finds those parameter names the shortest distance from the
// passed value and returns them
func SuggestParams(ps *PSet, s string) []string {
	names := make([]string, 0, len(ps.nameToParam))
	for n := range ps.nameToParam {
		names = append(names, n)
	}

	finder := strdist.DefaultFinders[strdist.CaseBlindAlgoNameCosine]

	return finder.FindNStrLike(alternativeCount, s, names...)
}

// SuggestGroups finds those group names the shortest distance from the
// passed value and returns them
func SuggestGroups(ps *PSet, s string) []string {
	names := make([]string, 0, len(ps.groups))
	for n := range ps.groups {
		names = append(names, n)
	}

	finder := strdist.DefaultFinders[strdist.CaseBlindAlgoNameCosine]

	return finder.FindNStrLike(alternativeCount, s, names...)
}

// SuggestNotes finds those note names the shortest distance from the
// passed value and returns them
func SuggestNotes(ps *PSet, s string) []string {
	names := make([]string, 0, len(ps.notes))
	for n := range ps.notes {
		names = append(names, n)
	}

	finder := strdist.DefaultFinders[strdist.CaseBlindAlgoNameCosine]

	return finder.FindNStrLike(alternativeCount, s, names...)
}
