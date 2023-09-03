package phelp

import (
	"github.com/nickwells/fileparse.mod/fileparse"
	"github.com/nickwells/param.mod/v5/param"
)

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
			" To see all the available parameter groups use the '"+
			helpShowArgName+" "+groupsHelpSectionName+"' parameter.",
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

	ps.AddNote("Alternative Sources - Configuration Files",
		"It is possible for a program to read parameters from files."+
			"\n\n"+
			"Parameters in such files are given one-per-line, the"+
			" leading dash is not required. If"+
			" any parameter value is required it is given after the"+
			" parameter name, separated with an '='. White space"+
			" at the start or end of the line is ignored as is any"+
			" around the '='."+
			" Parameters can be restricted to only be recognised for"+
			" specific programs by giving a comma-separated list of"+
			" program names followed by a '/'  before the parameter name."+
			" This can be useful for group or shared parameter files"+
			" (see below) and also allows you to configure the"+
			" behaviour of a program by creating multiple linked"+
			" copies under different names."+
			"\n\n"+
			"Blank lines in parameter files are ignored as is any text"+
			" following a '"+fileparse.DefaultCommentIntro+"'."+
			" Other files may be included by adding a"+
			" line to the file starting"+
			" with '"+fileparse.DefaultInclKeyWord+"'; any text"+
			" following this keyword has surrounding whitespace"+
			" removed and the remainder used as a filename to be processed."+
			"\n\n"+
			"Parameter files are either pre-declared and will be listed"+
			" in the sources section of the manual"+
			" (see '"+helpShowArgName+"' '"+sourcesHelpSectionName+"')"+
			" or are provided on the command line with"+
			" the '"+paramNameFile+"' parameter."+
			"\n\n"+
			"There is an additional distinction within the"+
			" pre-declared configuration files: some configuration"+
			" files are specific to a parameter group. Parameter"+
			" groups are means of organisinmg parameters into"+
			" logically-related collections. These groups of"+
			" parameters can each have their own group-specific"+
			" configuration files."+
			" (see '"+helpShowArgName+"' '"+groupsHelpSectionName+"')."+
			"\n\n"+
			"The parameters in these various types of configuration"+
			" file are handled slightly differently."+
			"\n"+
			"- Any valid"+
			" parameters of the program can be set in a file given"+
			" through the command-line"+
			" parameter '"+paramNameFile+"'. They are treated"+
			" as if they were given at the command line."+
			"\n"+
			"- Parameters given in pre-declared configuration"+
			" files have an additional restriction that prevents"+
			" parameters which are marked as 'command-line-only'"+
			" from being set. An error will be raised if one is"+
			" found in the file."+
			"\n"+
			"- Parameters given in a parameter-group"+
			" configuration file must also be members of the"+
			" parameter group."+
			"\n"+
			"- Additionally a configuration file may be shared"+
			" between multiple programs in which case the parameters"+
			" given in the file need not be parameters of the"+
			" program. Such parameters will be silently ignored."+
			" Such files, if any, will be highlighted in the list of"+
			" sources. To detect such ignored parameters use"+
			" the '"+paramNameShowUnused+"' parameter.",
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
			paramNameShowWhereSet+", "+
			paramNameShowUnused,
		param.NoteAttrs(param.DontShowNoteInStdUsage))
}
