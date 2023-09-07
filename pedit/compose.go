package pedit

import "github.com/nickwells/param.mod/v6/psetter"

type Composite struct {
	Editors []psetter.Editor
}

func (c Composite) Edit(paramName, paramVal string) (string, error) {
	pv := paramVal
	var err error
	for _, e := range c.Editors {
		pv, err = e.Edit(paramName, pv)
		if err != nil {
			break
		}
	}
	return pv, err
}
