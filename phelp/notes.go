package phelp

import (
	"sort"
	"strings"

	"github.com/nickwells/english.mod/english"
	"github.com/nickwells/param.mod/v7/param"
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

	return h.notesChosen[n.Headline()]
}

// showNotes produces the Notes section of the help message
func showNotes(h StdHelp, ps *param.PSet) bool {
	notes := ps.Notes()
	if len(notes) == 0 {
		return false
	}

	switch h.helpFormat {
	case helpFmtTypeMarkdown:
		return showNotesFmtMD(h, ps)
	default:
		return showNotesFmtStd(h, ps)
	}
}

// showNotesFmtStd produces the Notes section of the help message
func showNotesFmtStd(h StdHelp, ps *param.PSet) bool {
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
		h.twc.Printf("Notes [ %d notes ]\n", len(notes))
	} else if hiddenCount == len(notes) {
		h.twc.Printf("Notes [ %d notes, all hidden ]\n", len(notes))
	} else {
		h.twc.Printf("Notes [ %d notes, %d hidden ]\n", len(notes), hiddenCount)
	}

	for _, headline := range keys {
		n := notes[headline]
		if !noteCanBeShown(h, n) {
			continue
		}

		h.twc.Wrap(n.Headline(), paramIndent)

		if h.showSummary {
			continue
		}

		h.twc.Wrap(n.Text(), descriptionIndent)

		showNotesRefsFmtStd(h.twc, n.SeeParams(), "Parameter")
		showNotesRefsFmtStd(h.twc, n.SeeNotes(), "Note")
	}

	return true
}

// showNotesRefsFmtStd adds a section to the Notes section of the help
// message showing the named references (if any)
func showNotesRefsFmtStd(twc *twrap.TWConf, refs []string, name string) {
	if len(refs) == 0 {
		return
	}

	prefix := "See " + english.Plural(name, len(refs)) + ": "

	var refStr strings.Builder

	targetSpace := twc.TargetLineLen - descriptionIndent - len(prefix)

	refStr.WriteString(refs[0])

	spaceUsed := len(refs[0])

	for i := 1; i < len(refs); i++ {
		refStr.WriteString(",")

		spaceUsed++

		if spaceUsed+1+len(refs[i]) > targetSpace {
			refStr.WriteString("\n")

			spaceUsed = 0
		} else {
			refStr.WriteString(" ")

			spaceUsed++
		}

		refStr.WriteString(refs[i])

		spaceUsed += len(refs[i])
	}

	twc.WrapPrefixed(prefix, refStr.String(), descriptionIndent)
}

// showNotesFmtMD produces the Notes section of the help message in Markdown
// format
func showNotesFmtMD(h StdHelp, ps *param.PSet) bool {
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

	h.twc.Print("# Notes\n\n")

	sort.Strings(keys)

	for _, headline := range keys {
		h.twc.Print("## " + makeTextMarkdownSafe(headline) + "\n")

		if h.showSummary {
			continue
		}

		n := notes[headline]
		text := makeTextMarkdownSafe(n.Text())
		r := strings.NewReplacer("\n", "\n\n")
		text = r.Replace(text)
		h.twc.Wrap(text, 0)

		showNotesRefsFmtMD(h.twc, n.SeeParams(), "Parameter")
		showNotesRefsFmtMD(h.twc, n.SeeNotes(), "Note")

		h.twc.Print("\n\n")
	}

	return true
}

// showNotesRefsFmtMD adds a section to the Markdown file showing the named
// references (if any)
func showNotesRefsFmtMD(twc *twrap.TWConf, refs []string, name string) {
	if len(refs) > 0 {
		twc.Print("### See " + english.Plural(name, len(refs)) + "\n")

		for _, ref := range refs {
			twc.Print("* " + makeTextMarkdownSafe(ref) + "\n")
		}

		twc.Print("\n")
	}
}
