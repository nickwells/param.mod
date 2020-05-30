package phelp

import (
	"sort"

	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// noteCanBeShown will return true if the note can be shown
func noteCanBeShown(h StdHelp, n *param.Note) bool {
	if h.notesChosen.hasNothingChosen() {
		if h.showHiddenItems {
			return true
		}
		if n.AttrIsSet(param.DontShowNoteInStdUsage) {
			return false
		}
		return true
	}
	return h.notesChosen[n.Headline]
}

// showNotes produces the Notes section of the help message
func showNotes(h StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	notes := ps.Notes()
	if len(notes) == 0 {
		return false
	}

	hiddenCount := 0

	keys := make([]string, 0, len(notes))
	for k, n := range notes {
		keys = append(keys, k)
		if n.AttrIsSet(param.DontShowNoteInStdUsage) {
			hiddenCount++
		}
	}
	sort.Strings(keys)
	if hiddenCount == len(notes) {
		twc.Printf("Notes [ %d notes, all hidden ]\n", len(notes))
	} else {
		twc.Printf("Notes [ %d notes, %d hidden ]\n", len(notes), hiddenCount)
	}

	for _, headline := range keys {
		n := notes[headline]
		if !noteCanBeShown(h, n) {
			continue
		}
		twc.Wrap(n.Headline, paramIndent)
		if h.hideDescriptions {
			continue
		}
		twc.Wrap(n.Text, descriptionIndent)
	}

	return true
}
