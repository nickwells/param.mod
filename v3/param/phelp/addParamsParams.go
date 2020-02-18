package phelp

import (
	"github.com/nickwells/param.mod/v3/param"
	"github.com/nickwells/param.mod/v3/param/psetter"
)

const (
	paramsShowWhereSetArgName = "params-show-where-set"
	paramsShowUnusedArgName   = "params-show-unused"
)

const (
	exitAfterParamProcessing = "\n\nThe program will exit" +
		" after the parameters are processed."
)

// addParamHandlingParams will add the standard parameter-handling parameters
// into the parameter set
func (h *StdHelp) addParamHandlingParams(ps *param.PSet) {
	groupName := groupNamePfx + "-params"

	ps.AddGroup(groupName,
		"These are the parameter-handling parameters."+
			" There are parameters for showing where parameters"+
			" have been set and for the handling of parameter errors.")

	ps.Add(paramsShowWhereSetArgName,
		psetter.Bool{Value: &h.paramsShowWhereSet},
		"after all the parameters are set a message will be printed"+
			" showing where they were set. This can be useful for"+
			" debugging (especially if there are several config"+
			" files in use)."+
			exitAfterParamProcessing,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName))

	ps.Add(paramsShowUnusedArgName,
		psetter.Bool{Value: &h.paramsShowUnused},
		"after all the parameters are set a message will be printed"+
			" showing any parameters (including those from configuration"+
			" files or the environment) which were not recognised."+
			"\n\n"+
			"Parameters set in configuration files or through"+
			" environment variables may be intended for other programs"+
			" and so unused values are not classed as errors. Command"+
			" line options are obviously intended for this program and"+
			" so any command line parameter which is not recognised is"+
			" treated as an error. Setting this parameter will let you"+
			" check for spelling mistakes in parameters"+
			" that you've set in your alternative sources."+
			exitAfterParamProcessing,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName))

	ps.Add("params-dont-show-errors",
		psetter.Bool{Value: &h.dontReportErrors},
		"after all the parameters are set any errors detected will be"+
			" reported unless this flag is set to true",
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName))

	ps.Add("params-dont-exit-on-errors",
		psetter.Bool{Value: &h.dontExitOnErrors},
		"if errors are detected when processing the parameters the"+
			" program will exit unless this flag is set to true. Note"+
			" that the behaviour of the program cannot be guaranteed"+
			" if this option is chosen and it should only be used in"+
			" extreme circumstances",
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName))

	ps.Add("params-exit-after-parsing",
		psetter.Bool{Value: &h.exitAfterParsing},
		"exit after the parameters have been read and processed. This"+
			" lets you check the parameters are valid and see what"+
			" values get set without actually running the program."+
			"\n\n"+
			"Note that the program may perform some operations as the"+
			" parameters are processed and these will still take place"+
			" even if this parameter is set.",
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName))

	ps.Add("params-file",
		configFileSetter{seenBefore: make(map[string]bool)},
		"read in parameters from the given file. Note that the"+
			" parameter file will be read as a configuration file"+
			" with each parameter on a separate line. Comments,"+
			" white space etc. will be treated as in any other"+
			" configuration file",
		param.AltName("params-from"),
		param.PostAction(param.ConfigFileActionFunc),
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName))
}
