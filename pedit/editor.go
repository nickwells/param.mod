package pedit

// Editor defines an interface providing an Edit function. This is used to
// edit a parameter value before setting the value. The expected use by the
// Setter is for the Setter to pass the parameter name as the first value,
// the parameter value as the second and for Edit to return the modified
// value. It is expected that the Setter will check the error value and if it
// is not nil it will return it and abort the setting of the value.
type Editor interface {
	Edit(string, string) (string, error)
}
