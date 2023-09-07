package pedit

type Composite struct {
	Editors []Editor
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
