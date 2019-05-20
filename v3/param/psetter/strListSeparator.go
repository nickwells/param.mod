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

// ListValDesc returns that fragment of the description of what values are
// allowed which explains how the list values are separated from one another.
func (sls StrListSeparator) ListValDesc(name string) string {
	return "a list of " + name + " separated by '" + sls.GetSeparator() + "'"
}
