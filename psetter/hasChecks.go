package psetter

// CheckCounter can be used if you want to report the number of checks that a
// type has
type CheckCounter interface {
	CountChecks() int
}

// HasChecks returns a string reporting whether or not the number of checks
// is zero. This is suitable for constructing a string to be returned by an
// AllowedValues function. If you have a Setter which may have checks then
// add the result of this to the end of the value returned by the
// AllowedValues method on the Setter.
func HasChecks(cc CheckCounter) string {
	if cc.CountChecks() != 0 {
		return " subject to checks"
	}

	return ""
}
