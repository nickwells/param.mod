package param

// Example records a sample usage and a description for the "Examples"
// section of the help message
type Example struct {
	ex   string
	desc string
}

// AddExample adds an example to the set of examples on the PSet. Note that
// there is no validation of the given example
func (ps *PSet) AddExample(ex, desc string) {
	ps.examples = append(ps.examples,
		Example{
			ex:   ex,
			desc: desc,
		})
}

// Ex returns the example text
func (ex Example) Ex() string {
	return ex.ex
}

// Desc returns the descriptive text for the example
func (ex Example) Desc() string {
	return ex.desc
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
