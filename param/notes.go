package param

import (
	"fmt"
	"runtime"
)

// NoteAttributes records various flags that can be set on a Note
type NoteAttributes int32

const (
	// DontShowNoteInStdUsage means that the note will be suppressed when the
	// usage message is printed unless an expanded usage message has been
	// requested
	DontShowNoteInStdUsage NoteAttributes = 1 << iota
)

// Note records additional text to be attached to the help message which does
// not sit under any of the other parts
type Note struct {
	Headline   string
	Text       string
	Attributes NoteAttributes
	addedAt    string
}

// AddNote adds an note to the set of notes on the PSet. The headline of the
// note must be unique
func (ps *PSet) AddNote(headline, text string, opts ...NoteOptFunc) *Note {
	if existingNote, alreadyExists := ps.notes[headline]; alreadyExists {
		panic(fmt.Sprintf(
			"a note with headline: %s has already been added\nat: %s",
			headline, existingNote.addedAt))
	}

	stk := make([]byte, 10000)
	stkSize := runtime.Stack(stk, false)

	n := &Note{
		Headline: headline,
		Text:     text,
		addedAt:  string(stk[:stkSize]),
	}
	for _, o := range opts {
		err := o(n)
		if err != nil {
			panic(err.Error())
		}
	}

	ps.notes[n.Headline] = n
	return n
}

// HasNotes returns true if the PSet has any entries in the set of notes
func (ps *PSet) HasNotes() bool {
	return len(ps.notes) > 0
}

// GetNote returns a copy of the named note, error will be non-nil if there
// is no such note.
func (ps *PSet) GetNote(headline string) (*Note, error) {
	n, ok := ps.notes[headline]
	if !ok {
		return nil, fmt.Errorf("There is no note with headline: %q", headline)
	}
	copyVal := *n
	return &copyVal, nil
}

// Notes returns a copy of the current set of notes.
func (ps *PSet) Notes() map[string]*Note {
	n := make(map[string]*Note, len(ps.notes))
	for k, v := range ps.notes {
		copyVal := *v
		n[k] = &copyVal
	}
	return n
}

// NoteOptFunc is the type of an option func used to set various flags etc on a
// note.
type NoteOptFunc func(n *Note) error

// NoteAttrs returns a NoteOptFunc which will set the attributes of the note to
// The passed value.
func NoteAttrs(attrs NoteAttributes) NoteOptFunc {
	return func(n *Note) error {
		n.Attributes = attrs
		return nil
	}
}

// AttrIsSet will return true if the supplied attribute is set on the
// note. Multiple attributes may be given in which case they must all be set
func (n Note) AttrIsSet(attr NoteAttributes) bool {
	return n.Attributes&attr == attr
}
