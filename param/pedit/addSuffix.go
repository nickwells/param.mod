package pedit

type AddSuffix struct {
	Suffix string
}

func (e AddSuffix) Edit(_, paramVal string) (string, error) {
	return paramVal + e.Suffix, nil
}
