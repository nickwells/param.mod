package phelp

import (
	"fmt"
	"strings"

	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
)

const (
	completionsQuiet          = "completions-quiet"
	completionsZshDirArgName  = "completions-zsh-dir"
	completionsZshMakeArgName = "completions-zsh-make"
)

// addParamCompletionParams will add the standard parameters for specifying
// and creating shell completion functions into the parameter set
func (h *StdHelp) addParamCompletionParams(ps *param.PSet) {
	groupName := groupNamePfx + "-completion"

	ps.AddGroup(groupName,
		"These are the parameters for"+
			" creating shell completion functions."+
			" You can specify where"+
			" the completion files should be written,"+
			" trigger the generation of the files and"+
			" control whether they should be overwritten.")

	setConfigFileForGroupCommonParamsCompletion(ps) //nolint: errcheck

	zshDirParam := ps.Add(completionsZshDirArgName,
		psetter.Pathname{
			Value:       &h.zshCompDir,
			Expectation: filecheck.DirExists(),
		},
		"which directory should a zsh completions function for this"+
			" program be written to."+
			" The directory should be in the list of directories"+
			" given in the fpath shell variable."+
			" See the zsh manual for more details.",
		param.GroupName(groupName),
		param.Attrs(param.DontShowInStdUsage),
	)

	ps.Add(completionsQuiet,
		psetter.Bool{
			Value: &h.completionsQuiet,
		},
		"suppress any messages produced after generating"+
			" or updating the completions file.",
		param.GroupName(groupName),
		param.Attrs(param.DontShowInStdUsage),
	)

	const needZshDir = " The zsh completions directory name must be specified."
	zshMakeCompletionsParam := ps.Add(completionsZshMakeArgName,
		psetter.Enum[string]{
			AllowedVals: psetter.AllowedVals{
				zshCompActionRepl: "any existing zsh completions" +
					" file for the program will be overwritten or a" +
					" new file will be generated." +
					needZshDir,
				zshCompActionNew: "only generate the zsh completions file" +
					" if it doesn't already exist. Any pre-existing" +
					" file is protected and an error will be reported." +
					needZshDir,
				zshCompActionShow: "don't generate the zsh completions file." +
					" The file that would have been generated is" +
					" instead printed to standard output.",
				zshCompActionNone: "do nothing.",
			},
			Value: &h.zshCompAction,
		},
		"how to create the zsh completions file."+
			" This specifies whether or if the file should be created."+
			" If it is set to any value other than '"+zshCompActionNone+
			"' then the program will exit after the parameters are processed.",
		param.SeeAlso(completionsZshDirArgName),
		param.GroupName(groupName),
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
	)

	// Final checks

	ps.AddFinalCheck(func() error {
		if h.zshCompAction == zshCompActionShow {
			return nil
		}
		if h.zshCompAction == zshCompActionNone {
			return nil
		}

		if !zshDirParam.HasBeenSet() {
			// These parameters are processed before errors are reported so
			// we should abort the creation of the completion file
			h.zshCompAction = zshCompActionNone

			return fmt.Errorf(
				"the %q parameter has been set (at: %s)"+
					" but the %q parameter has not",
				zshMakeCompletionsParam.Name(),
				strings.Join(zshMakeCompletionsParam.WhereSet(), " and at "),
				zshDirParam.Name())
		}
		return nil
	})
}
