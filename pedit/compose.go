package pedit

import "github.com/nickwells/param.mod/v6/psetter"

// Composite combines multiple editors which are applied in sequence with the
// results of each being passed to the next. It implements the Editor interface
type Composite struct {
	Editors []psetter.Editor
}

// Edit applies the Editors in sequence, passing the results of the first to
// the second and so on. Any error stops the editing and the error is
// returned.
func (c Composite) Edit(paramName, paramVal string) (string, error) {
	var err error
	for _, e := range c.Editors {
		paramVal, err = e.Edit(paramName, paramVal)
		if err != nil {
			break
		}
	}
	return paramVal, err
}
