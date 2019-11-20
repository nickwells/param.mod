package param

// Example records a sample usage and a description for the "Examples"
// section of the help message
type Example struct {
	Ex   string
	Desc string
}

// AddExample adds an example to the set of examples on the PSet. Note that
// there is no validation of the given example
func (ps *PSet) AddExample(ex, desc string) {
	ps.examples = append(ps.examples,
		Example{
			Ex:   ex,
			Desc: desc,
		})
}

// HasExamples returns true if the PSet has any entries in the set of examples
func (ps *PSet) HasExamples() bool {
	return len(ps.examples) > 0
}

// Examples returns a copy of the current set of examples.
func (ps *PSet) Examples() []Example {
	sa := make([]Example, len(ps.examples))
	copy(sa, ps.examples)
	return sa
}
