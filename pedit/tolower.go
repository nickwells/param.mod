package pedit

import "strings"

type ToLower struct{}

func (ToLower) Edit(_, paramVal string) (string, error) {
	return strings.ToLower(paramVal), nil
}
