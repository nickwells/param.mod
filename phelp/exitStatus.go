package phelp

// the exit statuses are used by StdHelp. Note that they are not exported and
// so should not be relied upon as they can change without notice from
// release to release but they could be used to diagnose immediate problems.
const (
	exitStatusOK = iota
	exitStatusErrorsFound
	exitStatusCompletionGenFailure
)
