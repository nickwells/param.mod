package psetter

// NilValueMessage returns a standard message documenting that a Value is
// nil. This is used by the psetter CheckSetter methods and should be used to
// report nil Value fields by any implementation of the CheckSetter method so
// as to provide a consistent user experience.
func NilValueMessage(paramName, setterType string) string {
	return paramName + ": " + setterType +
		" Check failed: the Value to be set is nil"
}
