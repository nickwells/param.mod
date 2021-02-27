package phelp

import (
	"sort"
	"strings"

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
	switch h.helpFormat {
	case helpFmtTypeMD:
		return showNotesFmtMD(h, twc, ps)
	default:
		return showNotesFmtStd(h, twc, ps)
	}
}

// showNotesFmtStd produces the Notes section of the help message
func showNotesFmtStd(h StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	notes := ps.Notes()
	hiddenCount := 0

	keys := make([]string, 0, len(notes))
	for k, n := range notes {
		keys = append(keys, k)
		if n.AttrIsSet(param.DontShowNoteInStdUsage) {
			hiddenCount++
		}
	}
	sort.Strings(keys)
	if h.showHiddenItems {
		twc.Printf("Notes [ %d notes ]\n", len(notes))
	} else if hiddenCount == len(notes) {
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

// showNotesFmtMD produces the Notes section of the help message
func showNotesFmtMD(h StdHelp, twc *twrap.TWConf, ps *param.PSet) bool {
	notes := ps.Notes()

	keys := make([]string, 0, len(notes))
	for k, n := range notes {
		if !noteCanBeShown(h, n) {
			continue
		}
		keys = append(keys, k)
	}

	if len(keys) == 0 {
		return false
	}
	twc.Print("# Notes\n\n")

	sort.Strings(keys)

	for _, headline := range keys {
		twc.Print("## " + makeTextMarkdownSafe(headline) + "\n")
		if h.hideDescriptions {
			continue
		}

		text := makeTextMarkdownSafe(notes[headline].Text)
		r := strings.NewReplacer("\n", "\n\n")
		text = r.Replace(text)
		twc.Wrap(text, 0)
		twc.Print("\n\n")
	}

	return true
}
