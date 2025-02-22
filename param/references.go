package param

// Reference records a name and a description for the "See Also" section of the
// help message
type Reference struct {
	name string
	desc string
}

// AddReference adds a reference to the set of references on the PSet
func (ps *PSet) AddReference(name, desc string) {
	ps.references = append(ps.references,
		Reference{
			name: name,
			desc: desc,
		})
}

// Name returns the Reference name
func (r Reference) Name() string {
	return r.name
}

// Desc returns the text of the Reference
func (r Reference) Desc() string {
	return r.desc
}

// HasReferences returns true if the PSet has any references
func (ps *PSet) HasReferences() bool {
	return len(ps.references) > 0
}

// References returns a copy of the current set of references.
func (ps *PSet) References() []Reference {
	sa := make([]Reference, len(ps.references))
	copy(sa, ps.references)

	return sa
}
