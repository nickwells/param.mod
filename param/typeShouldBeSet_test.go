package param_test

// Code generated by mkfunccontrolparamtype; DO NOT EDIT.
// with parameters set at:
//	[command line]: Argument:2: "-d" "- this determines whether or not a ByName or ByPos param is expected to be set"
//	[command line]: Argument:3: "-for-testing"
//	[command line]: Argument:5: "-t" "ShouldBeSet"
//	[command line]: Argument:7: "-v" "ShouldNotBeSet"
//	[command line]: Argument:9: "-v" "ShouldBeSet"

/*
ShouldBeSetType - this determines whether or not a ByName or ByPos param is
expected to be set
*/
type ShouldBeSetType int

// These constants are the allowed values of ShouldBeSetType
const (
	ShouldNotBeSet ShouldBeSetType = iota
	ShouldBeSet
)

// IsValid is a method on the ShouldBeSetType type that can be used
// to check a received parameter for validity. It compares
// the value against the boundary values for the type
// and returns false if it is outside the valid range
func (v ShouldBeSetType) IsValid() bool {
	if v < ShouldNotBeSet {
		return false
	}

	if v > ShouldBeSet {
		return false
	}

	return true
}
