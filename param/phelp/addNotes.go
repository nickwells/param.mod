package phelp

import "github.com/nickwells/param.mod/v5/param"

// addNotes adds all the notes for the help message
func addNotes(ps *param.PSet) {
	ps.AddNote("Parameters - Optional",
		"Parameters which are not required are shown surrounded by square"+
			" brackets [like-this].",
		param.NoteAttrs(param.DontShowNoteInStdUsage))

	ps.AddNote("Parameters - Values",
		"A parameter which must take a value is shown with a"+
			" following '=...'. In this case the value can either"+
			" be supplied immediately after the parameter (with"+
			" an '=' in between) or else as the next argument to"+
			" the program. As follows:\n\n"+
			"-xxx=42 or -xxx 42"+
			"\n\n"+
			"If the following value is optional it is shown with"+
			" a following '[=...]' (note the brackets). In this"+
			" case the following value must come after an '='"+
			" rather than as the next argument. As follows:\n\n"+
			"-xxx=false",
		param.NoteAttrs(param.DontShowNoteInStdUsage))

	ps.AddNote("Parameters - Groups",
		"The parameters are arranged into named groups which can"+
			" be selected or suppressed through other help parameters."+
			" Within each group the parameters are displayed in"+
			" alphabetical order."+
			"\n\n"+
			"Groups where all the parameters are hidden will not be shown."+
			" To see all the available parameter groups use the "+
			helpGroupsArgName+" parameter.",
		param.NoteAttrs(param.DontShowNoteInStdUsage))

	ps.AddNote("Alternative Sources - Environment Variables",
		"If the program can be configured through environment variables"+
			" then a prefix will be given. Only those environment"+
			" variables having this prefix will be considered."+
			"\n\n"+
			"When matching environment variables to program paremeters the"+
			" prefix is stripped off and any underscores ('_') in"+
			" the environment variable name after the prefix will"+
			" be replaced with dashes ('-') when matching the"+
			" parameter name."+
			"\n\n"+
			"For instance, for the prefix 'XX_' an environment"+
			" variable called 'XX_a_b' will match a"+
			" parameter called 'a-b'",
		param.NoteAttrs(param.DontShowNoteInStdUsage))

	ps.AddNote("Alternative Sources - Priority",
		"If there are alternative sources of parameters (for instance"+
			" configuration files) these will be processed before the"+
			" command line parameters. The order in which alternative"+
			" sources are processed is as given on the Alternative"+
			" Sources help page."+
			"\n\n"+
			"Processing command line parameters last means that a"+
			" value given on the command line will replace any settings"+
			" in configuration files or environment variables (unless"+
			" the parameter may only be set once).",
		param.NoteAttrs(param.DontShowNoteInStdUsage))

	ps.AddNote("Alternative Sources - Useful Parameters",
		"When alternative sources are available it can be useful to"+
			" know where parameters have been set and to show any"+
			" invalid parameters. The following parameters can be"+
			" useful with these tasks: "+
			paramsShowWhereSetArgName+", "+
			paramsShowUnusedArgName,
		param.NoteAttrs(param.DontShowNoteInStdUsage))
}
