package pedit

import "strings"

// ToUpper implements the Editor interface
type ToUpper struct{}

// Edit returns the paramVal with all Unicode letters mapped to their upper
// case
func (ToUpper) Edit(_, paramVal string) (string, error) {
	return strings.ToUpper(paramVal), nil
}
