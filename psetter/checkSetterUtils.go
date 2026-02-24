package psetter

import "fmt"

// NilValueMessage returns a standard message documenting that a Setter has a
// nil Value. This can be used by any param.Setter CheckSetter methods and
// should be used to report nil Value fields by any implementation of the
// CheckSetter method so as to provide a consistent user experience.
func NilValueMessage(paramName, setterType string) string {
	return paramName + ": " + setterType +
		" Check failed: the Value to be set is nil"
}

// NilCheckMessage returns a standard message documenting that a Setter has
// one of its Checks set to nil. This can be used by any param.Setter
// CheckSetter methods and should be used to report nil Check entries by any
// implementation of the CheckSetter method so as to provide a consistent
// user experience.
func NilCheckMessage(paramName, setterType string, i int) string {
	return fmt.Sprintf(
		"%s: %s Check failed: the Check func at index %d is nil",
		paramName, setterType, i)
}

// BadSetterMessage returns a standard message documenting that a Setter has
// been improperly constructed. This can be used by any param.Setter
// CheckSetter methods and should be used to report improper construction of
// a Setter by any implementation of the CheckSetter method so as to provide
// a consistent user experience.
func BadSetterMessage(paramName, setterType, message string) string {
	return paramName + ": " + setterType +
		" Check failed: the Setter is improperly constructed: " +
		message
}

// BadValueMessage returns a standard message documenting that a Setter has a
// bad initial value. This can be used by any param.Setter CheckSetter
// methods and should be used to report the problem by any implementation of
// the CheckSetter method so as to provide a consistent user experience.
func BadValueMessage(paramName, setterType, message string) string {
	return paramName + ": " + setterType +
		" Check failed: the Value to be set is currently invalid: " +
		message
}
