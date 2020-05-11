package phelp

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v4/param"
	"github.com/nickwells/param.mod/v4/param/paction"
	"github.com/nickwells/param.mod/v4/param/psetter"
)

const (
	helpArgName             = "help"
	helpFullArgName         = "help-full"
	helpSummaryArgName      = "help-summary"
	helpGroupsArgName       = "help-show-groups"
	helpAltSourcesArgName   = "help-show-sources"
	helpShowExamplesArgName = "help-show-examples"
	helpShowRefsArgName     = "help-show-references"
)

const (
	exitAfterHelpMessage = "\n\nThe program will exit" +
		" after the help message is shown."
)

// addUsageParams will add the usage parameters into the parameter set
func (h *StdHelp) addUsageParams(ps *param.PSet) {
	// TODO: Do we need the styleCounter
	var styleCounter paction.Counter
	styleCounterAF := (&styleCounter).MakeActionFunc()

	groupName := groupNamePfx + "-help"

	ps.AddGroup(groupName,
		"These are parameters for printing a help message.")

	ps.Add(helpArgName, psetter.Nil{},
		"print a help message explaining what the program does and"+
			" the available parameters."+
			"\n\n"+
			"Parameters which are less commonly useful will not be shown."+
			" To see these hidden parameters use the -"+helpFullArgName+
			" parameter."+
			"\n\n"+
			"For a shorter help message use the -"+helpSummaryArgName+
			" parameter"+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly),
		param.AltName("usage"),
		param.PostAction(setStyle(h, stdHelp)),
		param.PostAction(styleCounterAF),
		param.GroupName(groupName))

	ps.Add(helpFullArgName, psetter.Nil{},
		" show all the parameters when printing the help message."+
			" Parameters, including this one, which would not normally be"+
			" shown as part of the help message will be printed."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-a"),
		param.AltName("help-all"),
		param.AltName("help-show-all"),
		param.AltName("help-show-hidden"),
		param.PostAction(paction.SetBool(&h.styleNeedsSetting, true)),
		param.PostAction(paction.SetBool(&h.paramsShowHidden, true)),
		param.GroupName(groupName))

	ps.Add(helpSummaryArgName, psetter.Nil{},
		"print a summary of the help message."+
			" Whether the hidden parameters are shown is"+
			" determined by the "+helpFullArgName+" parameter."+
			" To see all the parameters (including hidden ones) in"+
			" a summarised form you will need to give both this parameter"+
			" and the "+helpFullArgName+" parameter."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-s"),
		param.AltName("help-short"),
		param.PostAction(paction.SetBool(&h.styleNeedsSetting, true)),
		param.PostAction(paction.SetBool(&h.showFullHelp, false)),
		param.GroupName(groupName))

	ps.Add("help-all-short", psetter.Nil{},
		"print a summarised help message giving all the valid parameters."+
			" This is the equivalent of giving both the "+helpFullArgName+
			" and the "+helpSummaryArgName+" parameters."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-as"),
		param.AltName("help-sa"),
		param.PostAction(paction.SetBool(&h.styleNeedsSetting, true)),
		param.PostAction(paction.SetBool(&h.paramsShowHidden, true)),
		param.PostAction(paction.SetBool(&h.showFullHelp, false)),
		param.GroupName(groupName))

	ps.Add(helpGroupsArgName, psetter.Nil{},
		"print the parameter groups that have been set up."+
			" Do not show the individual parameters in those groups."+
			"\n\n"+
			"This lets you see just the available groups of"+
			" parameters and choose"+
			" which you want to select for closer examination."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-groups"),
		param.PostAction(setStyle(h, groupNamesOnly)),
		param.PostAction(styleCounterAF),
		param.GroupName(groupName))

	ps.Add("help-groups-in-list",
		psetter.Map{
			Value: &h.groupsSelected,
			Checks: []check.MapStringBool{
				check.MapStringBoolTrueCountGT(0),
			},
		},
		"when printing the help message only show help for parameters"+
			" in the listed groups."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-groups-in"),
		param.PostAction(setStyle(h, paramsInGroups)),
		param.PostAction(styleCounterAF),
		param.PostAction(checkGroups(h, ps)),
		param.GroupName(groupName))

	ps.Add("help-groups-not-in-list",
		psetter.Map{
			Value: &h.groupsSelected,
			Checks: []check.MapStringBool{
				check.MapStringBoolTrueCountGT(0),
			},
		},
		"when printing the help message don't show help for parameters"+
			" in the listed groups."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-groups-not-in"),
		param.PostAction(setStyle(h, paramsNotInGroups)),
		param.PostAction(styleCounterAF),
		param.PostAction(checkGroups(h, ps)),
		param.GroupName(groupName))

	ps.Add("help-show-params",
		psetter.StrList{
			Value: &h.paramsToShow,
		},
		"when printing the help message only show help for the"+
			" listed parameters"+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-p"),
		param.AltName("help-params"),
		param.PostAction(setStyle(h, paramsByName)),
		param.PostAction(styleCounterAF),
		param.PostAction(checkParams(h, ps)),
		param.GroupName(groupName))

	ps.Add("help-show-desc", psetter.Nil{},
		"when printing the help message only show the program description"+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-show-prog-desc"),
		param.AltName("help-prog-desc"),
		param.AltName("help-program-description"),
		param.PostAction(setStyle(h, progDescOnly)),
		param.PostAction(styleCounterAF),
		param.GroupName(groupName))

	ps.Add(helpAltSourcesArgName, psetter.Nil{},
		"when printing the help message only show the"+
			" places (other than the command line) where"+
			" parameters may be set. This will list any"+
			" configuration files and environment prefixes."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-param-sources"),
		param.PostAction(setStyle(h, altSourcesOnly)),
		param.PostAction(styleCounterAF),
		param.GroupName(groupName))

	ps.Add(helpShowExamplesArgName, psetter.Nil{},
		"when printing the help message only show the"+
			" examples, if any."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-examples"),
		param.AltName("help-eg"),
		param.PostAction(setStyle(h, examplesOnly)),
		param.PostAction(styleCounterAF),
		param.GroupName(groupName))

	ps.Add(helpShowRefsArgName, psetter.Nil{},
		"when printing the help message only show the"+
			" references (the See Also section), if any."+
			exitAfterHelpMessage,
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
		param.AltName("help-show-refs"),
		param.AltName("help-refs"),
		param.AltName("help-see-also"),
		param.PostAction(setStyle(h, referencesOnly)),
		param.PostAction(styleCounterAF),
		param.GroupName(groupName))

	// Final checks

	ps.AddFinalCheck(func() error {
		if styleCounter.Count() > 1 {
			return fmt.Errorf(
				"you have chosen conflicting types of help: %s",
				styleCounter.SetBy())
		}
		return nil
	})

	ps.AddFinalCheck(func() error {
		if h.styleNeedsSetting && h.style == noHelp {
			h.style = stdHelp
		}
		return nil
	})
}

// checkGroups returns an ActionFunc which will check that the groupsSelected
// element of the StdHelp structure only contains valid group names
func checkGroups(h *StdHelp, ps *param.PSet) param.ActionFunc {
	return func(_ location.L, _ *param.ByName, _ []string) error {
		errCount := 0
		msg := ""
		groupNames := make([]string, 0, len(h.groupsSelected))
		for g := range h.groupsSelected {
			groupNames = append(groupNames, g)
		}
		sort.Strings(groupNames)
		for _, g := range groupNames {
			if !ps.HasGroupName(g) {
				if errCount == 0 {
					msg = "group: '" + g + "'," +
						" is not the name of a parameter group." +
						" Please check the spelling."
				} else {
					msg += "\nalso: '" + g + "'"
				}
				errCount++
			}
		}
		if errCount > 0 {
			return errors.New(msg)
		}
		return nil
	}
}

// checkParams returns an ActionFunc which will check that the paramsToShow
// element of the StdHelp structure only contains valid parameter names
func checkParams(h *StdHelp, ps *param.PSet) param.ActionFunc {
	return func(_ location.L, _ *param.ByName, _ []string) error {
		var badParams []string
		for _, pName := range h.paramsToShow {
			trimmedName := strings.TrimLeft(pName, "-")
			_, err := ps.GetParamByName(trimmedName)
			if err != nil {
				badParams = append(badParams, pName)
			}
		}
		switch len(badParams) {
		case 0:
			return nil
		case 1:
			return fmt.Errorf("%q is not a parameter of this program",
				badParams[0])
		default:
			return errors.New(
				`The following are not parameters of this program: "` +
					strings.Join(badParams, `", "`) + `"`)
		}
	}
}
