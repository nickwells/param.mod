package param

import "github.com/nickwells/strdist.mod/strdist"

// SuggestionFunc
type SuggestionFunc func(ps *PSet, s string) []string

// SuggestParams finds those parameter names the shortest distance from the
// passed value and returns them
func SuggestParams(ps *PSet, s string) []string {
	names := make([]string, 0, len(ps.nameToParam))
	for n := range ps.nameToParam {
		names = append(names, n)
	}

	return strdist.CaseBlindCosineFinder.FindNStrLike(3, s, names...)
}

// SuggestGroups finds those group names the shortest distance from the
// passed value and returns them
func SuggestGroups(ps *PSet, s string) []string {
	names := make([]string, 0, len(ps.nameToParam))
	for n := range ps.groups {
		names = append(names, n)
	}

	return strdist.CaseBlindCosineFinder.FindNStrLike(3, s, names...)
}

// SuggestNotes finds those note names the shortest distance from the
// passed value and returns them
func SuggestNotes(ps *PSet, s string) []string {
	names := make([]string, 0, len(ps.nameToParam))
	for n := range ps.notes {
		names = append(names, n)
	}

	return strdist.CaseBlindCosineFinder.FindNStrLike(3, s, names...)
}
