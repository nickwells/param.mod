package psetter

// NilValueMessage ...
func NilValueMessage(paramName, setterType string) string {
	return paramName + ": " + setterType +
		" Check failed: the Value to be set is nil"
}
