package psetter

import "fmt"

// NilValueMessage returns a standard message documenting that a Value is
// nil. This is used by the psetter CheckSetter methods and should be used to
// report nil Value fields by any implementation of the CheckSetter method so
// as to provide a consistent user experience.
func NilValueMessage(paramName, setterType string) string {
	return paramName + ": " + setterType +
		" Check failed: the Value to be set is nil"
}

// NilCheckMessage returns a standard message documenting that a Setter has
// one of its Checks set to nil. This is used by the psetter CheckSetter
// methods and should be used to report nil Check entries by any
// implementation of the CheckSetter method so as to provide a consistent
// user experience.
func NilCheckMessage(paramName, setterType string, i int) string {
	return fmt.Sprintf(
		"%s: %s Check failed: the Check func at index %d is nil",
		paramName, setterType, i)
}
