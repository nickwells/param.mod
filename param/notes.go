package param

import (
	"fmt"
	"runtime"
	"strings"

	"golang.org/x/exp/slices"
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
	headline     string
	text         string
	attributes   NoteAttributes
	addedAt      string
	seeAlsoNote  map[string]string
	seeAlsoParam map[string]string
}

// Headline returns the note's headline text
func (n Note) Headline() string {
	return n.headline
}

// Text returns the note's text
func (n Note) Text() string {
	return n.text
}

// AddNote adds an note to the set of notes on the PSet. The headline of the
// note must be unique
func (ps *PSet) AddNote(headline, text string, opts ...NoteOptFunc) *Note {
	if existingNote, alreadyExists := ps.notes[headline]; alreadyExists {
		panic(fmt.Sprintf(
			"a note with headline: %s has already been added\nat: %s",
			headline, existingNote.addedAt))
	}

	const stackDumpBufSz = 10_000

	stk := make([]byte, stackDumpBufSz)
	stkSize := runtime.Stack(stk, false)
	addedAt := string(stk[:stkSize])

	if stkSize == stackDumpBufSz {
		addedAt += " ..."
	}

	n := &Note{
		headline:     headline,
		text:         text,
		addedAt:      addedAt,
		seeAlsoNote:  make(map[string]string),
		seeAlsoParam: make(map[string]string),
	}

	for _, o := range opts {
		err := o(n)
		if err != nil {
			panic(err.Error())
		}
	}

	ps.notes[n.headline] = n

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
		n.attributes = attrs

		return nil
	}
}

// NoteSeeNote will add the names of notes to the list of notes to be
// referenced when showing the help message. They will be checked before the
// parameters are parsed to ensure that they are all valid names. Note that
// it is not possible to check the names as they are added since the
// referenced name might not have been added yet. It will return an error if
// the referenced name has already been used.
func NoteSeeNote(notes ...string) NoteOptFunc {
	source := caller()

	return func(n *Note) error {
		for _, note := range notes {
			note = strings.TrimSpace(note)
			if note == n.headline {
				continue
			}

			if whereAdded, exists := n.seeAlsoNote[note]; exists {
				return fmt.Errorf(
					"The NoteSeeAlso note %q has already been added, at %s",
					note, whereAdded)
			}

			n.seeAlsoNote[note] = source
		}

		return nil
	}
}

// NoteSeeParam will add the names of parameters to the list of parameters to be
// referenced when showing the help message. They will be checked before the
// parameters are parsed to ensure that they are all valid names. Note that
// it is not possible to check the names as they are added since the
// referenced name might not have been added yet. It will return an error if
// the referenced name has already been used.
func NoteSeeParam(params ...string) NoteOptFunc {
	source := caller()

	return func(n *Note) error {
		for _, param := range params {
			param = strings.TrimSpace(param)

			if whereAdded, exists := n.seeAlsoParam[param]; exists {
				return fmt.Errorf(
					"The NoteSeeParam parameter %q has already"+
						" been added, at %s",
					param, whereAdded)
			}

			n.seeAlsoParam[param] = source
		}

		return nil
	}
}

// AttrIsSet will return true if the supplied attribute is set on the
// note. Multiple attributes may be given in which case they must all be set
func (n Note) AttrIsSet(attr NoteAttributes) bool {
	return n.attributes&attr == attr
}

// SeeNotes returns a sorted list of the notes referenced by this note
func (n Note) SeeNotes() []string {
	notes := make([]string, 0, len(n.seeAlsoNote))
	for note := range n.seeAlsoNote {
		notes = append(notes, note)
	}

	slices.Sort(notes)

	return notes
}

// SeeParams returns a sorted list of the params referenced by this note
func (n Note) SeeParams() []string {
	params := make([]string, 0, len(n.seeAlsoParam))
	for param := range n.seeAlsoParam {
		params = append(params, param)
	}

	slices.Sort(params)

	return params
}
