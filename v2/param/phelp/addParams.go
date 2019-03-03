package phelp

import (
	"github.com/nickwells/param.mod/v2/param"
	"github.com/nickwells/param.mod/v2/param/paction"
	"github.com/nickwells/param.mod/v2/param/psetter"
)

const (
	usageArgName        = "help"
	usageFullArgName    = "help-full"
	usageSummaryArgName = "help-summary"
)

// groupNamePfx is the name of the group in which all the param package
// parameters are grouped. You should not give any of your parameter groups
// the same name (it'll be confusing)
const groupNamePfx = "pkg.param"

// AddParams will add the help parameters into the parameter set
func (h *StdHelp) AddParams(ps *param.PSet) {
	h.addParamHandlingParams(ps)
	h.addUsageParams(ps)
}

// addParamHandlingParams will add the standard parameter-handling parameters
// into the parameter set
func (h *StdHelp) addParamHandlingParams(ps *param.PSet) {
	groupName := groupNamePfx + "-params"

	ps.AddGroup(groupName,
		`These are the parameter-handling parameters.

There are parameters for showing where parameters have been set and the handling of parameter errors.`)

	ps.Add("params-show-where-set",
		psetter.Bool{Value: &h.reportWhereParamsAreSet},
		`after all the parameters are set a message will be printed showing where they were set. This can be useful for debugging (especially if there are several config files in use).

The program will exit if this parameter is set`,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName))

	ps.Add("params-show-sources",
		psetter.Bool{Value: &h.reportParamSources},
		`after all the parameters are set a message will be printed showing all the places (other than the command line) that a parameter can be set. This will list any configuration files and environment prefixes.

The program will exit if this parameter is set`,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName))

	ps.Add("params-show-unused",
		psetter.Bool{Value: &h.reportUnusedParams},
		`after all the parameters are set a message will be printed showing any parameters (including those from config files or the environment) which were not recognised.

Parameters set in any config files or through environment variables may be intended for other programs and so unused values are not classed as errors. Command line options are obviously intended for this program and so any command line parameter which is not recognised is treated as an error. Setting this parameter will allow you to check for spelling mistakes or other typos.

The program will exit if this parameter is set`,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName))

	ps.Add("params-dont-show-errors",
		psetter.Bool{Value: &h.dontReportErrors},
		"after all the parameters are set any errors detected will be reported unless this flag is set to true",
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName))

	ps.Add("params-dont-exit-on-errors",
		psetter.Bool{Value: &h.dontExitOnErrors},
		"if errors are detected when processing the parameters the program will exit unless this flag is set to true",
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName))

	ps.Add("params-exit-after-parsing",
		psetter.Bool{Value: &h.exitAfterParsing},
		`exit after the parameters have been parsed. This can allow you to just check the parameters or investigate what would have been set and not actually run the program.

Note that the program may perform some operations as the parameters are processed and these will still take place even if this parameter is set.`,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName))
}

// addUsageParams will add the usage parameters into the parameter set
func (h *StdHelp) addUsageParams(ps *param.PSet) {
	groupListAF := (&h.groupListCounter).MakeActionFunc()
	groupName := groupNamePfx + "-help"

	ps.AddGroup(groupName,
		`These are the usage parameters. They can be used to show a usage message in various levels of detail.`)

	ps.Add(usageArgName, helpType{h: h},
		`print a message giving the valid params and exit.

By default any parameters which the application deems less commonly useful or advanced will not be shown. To see all the parameters use the `+usageFullArgName+` parameter.

For a more concise usage message use the `+usageSummaryArgName+` parameter`,
		param.Attrs(param.CommandLineOnly),
		param.AltName("usage"),
		param.GroupName(groupName))

	ps.Add(usageFullArgName, helpFull{h: h},
		"print the usage message giving all the valid params and exit. Parameters, including this one, which would normally be suppressed are also shown.",
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-all"),
		param.AltName("help-show-all"),
		param.GroupName(groupName))

	ps.Add(usageSummaryArgName,
		helpType{
			h:     h,
			style: Short,
		},
		"print a summarised usage message giving the valid params and exit. Whether the suppressed parameters are shown is determined by the "+usageFullArgName+" parameter. To see all the parameters (including suppressed ones) in a summarised form you will need to give both this parameter and the "+usageFullArgName+" parameter.",
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-short"),
		param.GroupName(groupName))

	ps.Add("help-show-groups",
		helpType{
			h:     h,
			style: GroupNamesOnly,
		},
		`print the parameter groups that have been set up and exit. Do not show the individual parameters in those groups.

This can allow you to see the sets of parameters in a concise way and allow you to choose which you want to select for closer examination.`,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-groups"),
		param.GroupName(groupName))

	ps.Add("help-groups-in-list",
		psetter.Map{
			Value: &h.groupsToShow,
		},
		`when printing the help message only show help for parameters in the listed groups.`,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName),
		param.PostAction(paction.SetBool(&h.includeGroups, true)),
		param.PostAction(paction.SetBool(&h.showHelp, true)),
		param.PostAction(groupListAF))

	ps.Add("help-groups-not-in-list",
		psetter.Map{
			Value: &h.groupsToExclude,
		},
		`when printing the help message don't show help for parameters in the listed groups.`,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.GroupName(groupName),
		param.PostAction(paction.SetBool(&h.excludeGroups, true)),
		param.PostAction(paction.SetBool(&h.showHelp, true)),
		param.PostAction(groupListAF))
}
