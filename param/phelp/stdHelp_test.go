package phelp

// SetExitAfterHelp sets the exitAfterHelp flag to the value passed. Note
// that this is only available to the test functions as this function is
// declared in a file with name ending "_test" and so this will only be
// compiled when tests are being run.
func (h *StdHelp) SetExitAfterHelp(b bool) {
	h.exitAfterHelp = b
}

// SetDontExitOnErrors sets the exitOnErrors flag to the NOT of the value
// passed. Note that this is only available to the test functions as this
// function is declared in a file with name ending "_test" and so this will
// only be compiled when tests are being run.
func (h *StdHelp) SetDontExitOnErrors(b bool) {
	h.exitOnErrors = !b
}
