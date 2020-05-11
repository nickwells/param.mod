package phelp

import "github.com/nickwells/param.mod/v4/param"

// addNotes adds all the notes for the help message
func addNotes(ps *param.PSet) error {
	ps.AddNote("Optional Parameters",
		"Parameters which are not required are shown surrounded by square"+
			" brackets [like-this].")
	ps.AddNote("Parameter Values",
		"A parameter which takes a value is shown with a"+
			" following '=...'."+
			" If the following value is optional this is itself"+
			" bracketed [=...]. In this case the following value"+
			" must come after an '=' rather than as the next argument."+
			" As follows,"+
			"\n\n"+
			"-xxx=false not -xxx false"+
			"\n\n"+
			"For parameters which must have a value it may be given in"+
			" either way")
	ps.AddNote("Parameter Groups",
		"The parameters are arranged into named groups which can"+
			" be selected or suppressed through other help parameters."+
			" Within each group the parameters are displayed in"+
			" alphabetical order."+
			"\n\n"+
			" Groups where all the parameters are hidden will not be shown."+
			" To see all the available parameter groups use the "+
			helpGroupsArgName+" parameter.")

	return nil
}
