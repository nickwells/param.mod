package param

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
)

// BaseParam represents a parameter. It records the common information used
// by either a ByName or a ByPos parameter.
type BaseParam struct {
	ps *PSet

	name         string
	description  string
	valueName    string
	initialValue string

	setter     Setter
	postAction []ActionFunc

	seeAlso    map[string]string
	seeNote    map[string]string
	whereAdded string
}

// mkBaseParam constructs a BaseParam
func mkBaseParam(
	ps *PSet, name string, setter Setter, desc, whereAdded string,
) BaseParam {
	return BaseParam{
		ps:           ps,
		name:         name,
		setter:       setter,
		description:  desc,
		initialValue: setter.CurrentValue(),
		whereAdded:   whereAdded,
		seeAlso:      make(map[string]string),
		seeNote:      make(map[string]string),
	}
}

// PSet returns the parameter set to which the BaseParam parameter belongs
func (p BaseParam) PSet() *PSet { return p.ps }

// Name returns the name of the BaseParam parameter
func (p BaseParam) Name() string { return p.name }

// Description returns the description of the BaseParam parameter
func (p BaseParam) Description() string { return p.description }

// InitialValue returns the initialValue of the BaseParam parameter
func (p BaseParam) InitialValue() string { return p.initialValue }

// Setter returns the setter
func (p BaseParam) Setter() Setter { return p.setter }

// SeeAlso returns a sorted list of references to other parameters
func (p BaseParam) SeeAlso() []string {
	return slices.Sorted(maps.Keys(p.seeAlso))
}

// SeeNotes returns a sorted list of references to notes
func (p BaseParam) SeeNotes() []string {
	return slices.Sorted(maps.Keys(p.seeNote))
}

// ValueName returns the parameter's bespoke value name
func (p BaseParam) ValueName() string { return p.valueName }

// seeAlsoSource returns the string describing where the SeeAlso reference
// was added. This is suitable for reporting the location in code the mistake
// was made. If the reference is not found it will return an empty string
func (p BaseParam) seeAlsoSource(ref string) string {
	return p.seeAlso[ref]
}

// seeNoteSource returns the string describing where the SeeNote note was
// added. This is suitable for reporting the location in code the mistake was
// made. If the note is not found it will return an empty string
func (p BaseParam) seeNoteSource(note string) string {
	return p.seeNote[note]
}

// SetValueName sets the valueName on the BaseParam. The valueName is used as
// the short value name used in the parameter summary (it follows the "="
// after the parameter name). If this is not empty this will be used in
// preference to either the param setter's value description or the setter's
// type name. This allows a per-parameter value name to be given for a more
// helpful usage message. It will return an error if the vName is the empty
// string.
func (p *BaseParam) SetValueName(vName string) error {
	if vName == "" {
		return errors.New("some non-empty value name must be given")
	}

	p.valueName = vName

	return nil
}

// SetSeeAlsoRefs will add the names of parameters to the list of parameters
// to be referenced when showing the help message. They will be checked
// before the parameters are parsed to ensure that they are all valid
// names. Note that it is not possible to check the names as they are added
// since the referenced name might not have been added yet. It will return an
// error if the referenced name has already been used or if it appears in the
// list of values twice. A reference to the parameter itself will be ignored;
// this allows the same group of parameter names to be passed to each
// parameter in the group wihout self-reference.
//
// The whereAdded parameter should be some description of where this was
// called from (eg a stack trace).
func (p *BaseParam) SetSeeAlsoRefs(whereAdded string, refs ...string) error {
	dups := make(map[string]int)

	for i, ref := range refs {
		ref = strings.TrimSpace(ref)
		if ref == p.name { // don't add a see-also to yourself
			continue
		}

		if first, dupFound := dups[ref]; dupFound {
			return fmt.Errorf(
				"the SeeAlso reference %q appears in the list twice,"+
					" firstly at index %d and again at %d",
				ref, first, i)
		}

		dups[ref] = i

		if wherePreviouslyAdded, exists := p.seeAlso[ref]; exists {
			return fmt.Errorf(
				"the SeeAlso reference %q"+
					" has already been added, at %s"+
					" and now again at %s",
				ref, wherePreviouslyAdded, whereAdded)
		}

		p.seeAlso[ref] = whereAdded
	}

	return nil
}

// SetSeeNotes will add the names of notes to the list of notes to be
// referenced when showing the help message. They will be checked before the
// parameters are parsed to ensure that they are all valid note names. Note
// that it is not possible to check the note names as they are added since
// the referenced note might not have been added yet. It will return an error
// if the referenced note name has already been used or if it appears in the
// list of values twice.
//
// The whereAdded parameter should be some description of where this was
// called from (eg a stack trace).
func (p *BaseParam) SetSeeNotes(whereAdded string, notes ...string) error {
	dups := make(map[string]int)

	for i, note := range notes {
		note = strings.TrimSpace(note)

		if first, dupFound := dups[note]; dupFound {
			return fmt.Errorf(
				"the SeeNote reference %q appears in the list twice,"+
					" firstly at index %d and again at %d",
				note, first, i)
		}

		dups[note] = i

		if wherePreviouslyAdded, exists := p.seeNote[note]; exists {
			return fmt.Errorf(
				"the SeeNote reference %q"+
					" has already been added, at %s"+
					" and now again at %s",
				note, wherePreviouslyAdded, whereAdded)
		}

		p.seeNote[note] = whereAdded
	}

	return nil
}
