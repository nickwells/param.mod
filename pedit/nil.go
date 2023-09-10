package pedit

// Nil implements the Editor interface
type Nil struct{}

// Edit returns the paramVal unchanged
func (e Nil) Edit(_, paramVal string) (string, error) {
	return paramVal, nil
}
