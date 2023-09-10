package pedit

import "strings"

// ToLower implements the Editor interface
type ToLower struct{}

// Edit returns the paramVal with all Unicode letters mapped to their lower
// case
func (ToLower) Edit(_, paramVal string) (string, error) {
	return strings.ToLower(paramVal), nil
}
