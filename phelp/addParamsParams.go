package phelp

import (
	"github.com/nickwells/param.mod/v6/paction"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
)

const (
	paramNameWhereSetFormat   = "params-where-set-fmt"
	paramNameShowWhereSet     = "params-show-where-set"
	paramNameShowUnused       = "params-show-unused"
	paramNameDontShowErrors   = "params-dont-show-errors"
	paramNameDontExitOnErrors = "params-dont-exit-on-errors"
	paramNameExitAfterParsing = "params-exit-after-parsing"
	paramNameFile             = "params-file"
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

	ps.Add(paramNameWhereSetFormat,
		psetter.Enum[string]{
			Value: &h.paramsSetFormat,
			AllowedVals: psetter.AllowedVals[string]{
				paramSetFmtStd: "the standard format for showing" +
					" where and if parameters are set",
				paramSetFmtShort: "a short form of the information" +
					" and only showing values that have been set",
				paramSetFmtTable: "the information on where parameters" +
					" are set in a tabular format. Only values that" +
					" have been set are shown",
			},
		},
		"after all the parameters are set a message will be printed"+
			" showing where they were set. This parameter controls"+
			" how this information is shown."+
			exitAfterParamProcessing,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.PostAction(paction.SetVal[bool](&h.paramsShowWhereSet, true)),
		param.GroupName(groupName),
		param.SeeAlso(paramNameShowWhereSet),
	)

	ps.Add(paramNameShowWhereSet,
		psetter.Bool{Value: &h.paramsShowWhereSet},
		"after all the parameters are set a message will be printed"+
			" showing where they were set. This can be useful for"+
			" debugging (especially if there are several config"+
			" files in use)."+
			exitAfterParamProcessing,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName),
		param.SeeAlso(paramNameWhereSetFormat),
	)

	ps.Add(paramNameShowUnused,
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

	ps.Add(paramNameDontShowErrors,
		psetter.Bool{
			Value:  &h.reportErrors,
			Invert: true,
		},
		"after all the parameters are set any errors detected will be"+
			" reported unless this flag is set",
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName))

	ps.Add(paramNameDontExitOnErrors,
		psetter.Bool{
			Value:  &h.exitOnErrors,
			Invert: true,
		},
		"if errors are detected when processing the parameters the"+
			" program will exit unless this flag is set to true. Note"+
			" that the behaviour of the program cannot be guaranteed"+
			" if this option is chosen and it should only be used in"+
			" emergencies",
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName))

	ps.Add(paramNameExitAfterParsing,
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

	ps.Add(paramNameFile,
		&configFileSetter{},
		"read in parameters from the given file. Note that the"+
			" parameter file will be read as a configuration file"+
			" with each parameter on a separate line. Comments,"+
			" white space etc. will be"+
			" treated as in any other configuration file",
		param.AltNames("params-from", "params-f"),
		param.PostAction(param.ConfigFileActionFunc),
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName))
}
