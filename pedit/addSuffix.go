package pedit

// AddSuffix implements the Editor interface
type AddSuffix struct {
	Suffix string
}

// Edit returns the paramVal with the Suffix appended to it
func (e AddSuffix) Edit(_, paramVal string) (string, error) {
	return paramVal + e.Suffix, nil
}
