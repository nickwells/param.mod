package param

// Note records additional text to be attached to the help message which does
// not sit under any of the other parts
type Note struct {
	Headline string
	Text     string
}

// AddNote adds an note to the set of notes on the PSet. Note that
// there is no validation of the given note
func (ps *PSet) AddNote(headline, text string) {
	n := Note{
		Headline: headline,
		Text:     text,
	}
	ps.notes = append(ps.notes, n)
}

// HasNotes returns true if the PSet has any entries in the set of notes
func (ps *PSet) HasNotes() bool {
	return len(ps.notes) > 0
}

// Notes returns a copy of the current set of notes.
func (ps *PSet) Notes() []Note {
	n := make([]Note, len(ps.notes))
	copy(n, ps.notes)
	return n
}
