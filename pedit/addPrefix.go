package pedit

type AddPrefix struct {
	Prefix string
}

func (e AddPrefix) Edit(_, paramVal string) (string, error) {
	return e.Prefix + paramVal, nil
}
