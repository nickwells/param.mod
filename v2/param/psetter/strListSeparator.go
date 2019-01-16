package psetter

// StrListDefaultSep is the default separator for a list of strings. it is set
// to a comma
const (
	StrListDefaultSep = ","
)

// StrListSeparator holds the separator value
type StrListSeparator struct {
	Sep string
}

// GetSeparator returns the separator or the default value (a comma) if it is
// an empty string
func (sls StrListSeparator) GetSeparator() string {
	sep := sls.Sep
	if sep == "" {
		sep = StrListDefaultSep
	}
	return sep
}
