package pedit

import "strings"

type ToUpper struct{}

func (ToUpper) Edit(_, paramVal string) (string, error) {
	return strings.ToUpper(paramVal), nil
}
