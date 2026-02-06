package param

// Helper is the interface that a helper object must implement. The Helper is
// responsible for generating any help messages and for reporting any errors
// detected while parsing the parameters
type Helper interface {
	// AddParams will add Helper-specific parameters to the PSet. It is
	// called by the NewSet function after the Helper is set on the
	// newly created PSet
	AddParams(ps *PSet)
	// ProcessArgs is called by PSet.Parse once argument parsing is
	// complete. It should perform any operations prompted by the parameters
	// added by the Helper. Note that it is called before any remaining
	// (non-positional, unnamed) parameters are processed.
	ProcessArgs(ps *PSet)
	// ErrorHandler is called by PSet.Parse after the Helper parameters have
	// been processed and before any remaining parameters are processed. Note
	// that it is only called if errors have been detected.
	ErrorHandler(ps *PSet)
	// Help should generate a help message. It can be called directly through
	// the PSet.Help method. It will be called by PSet.Parse if the
	// PSet.HelpRequired method has been called.
	Help(ps *PSet, messages ...string)
}
