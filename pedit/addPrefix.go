package pedit

// AddPrefix implements the Editor interface
type AddPrefix struct {
	Prefix string
}

// Edit returns the paramVal with the Prefix prepended to it
func (e AddPrefix) Edit(_, paramVal string) (string, error) {
	return e.Prefix + paramVal, nil
}
