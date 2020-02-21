package phelp

import (
	"fmt"
	"strings"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/param.mod/v4/param"
	"github.com/nickwells/param.mod/v4/param/psetter"
)

// addParamCompletionParams will add the standard parameters for specifying
// and creating shell completion functions into the parameter set
func (h *StdHelp) addParamCompletionParams(ps *param.PSet) {
	groupName := groupNamePfx + "-completion"

	ps.AddGroup(groupName,
		"These are the parameters concerned with specifying and"+
			" creating shell completion functions."+
			" There are parameters for specifying which directory"+
			" the completion files should be written to, for"+
			" triggering the generation of the files and for"+
			" specifying whether they should be overwritten.")

	SetGroupConfigFile_common_params_completion(ps) //nolint: errcheck

	zshDirParam := ps.Add("completions-zsh-dir",
		psetter.Pathname{
			Value: &h.zshCompletionsDir,
			Expectation: filecheck.Provisos{
				Existence: filecheck.MustExist,
				Checks: []check.FileInfo{
					check.FileInfoIsDir,
				},
			},
		},
		"which directory should a zsh completions function for this"+
			" program be written to."+
			" The directory should be in the list of directories"+
			" given in the fpath shell variable."+
			" See the zsh manual for more details.",
		param.GroupName(groupName),
		param.Attrs(param.DontShowInStdUsage),
	)

	const needZshDir = " The zsh completions directory name must be specified."
	zshMakeCompletionsParam := ps.Add("completions-zsh-make",
		psetter.Enum{
			AllowedVals: param.AllowedVals{
				zshCompGenRepl: "any existing zsh completions" +
					" file for the program will be overwritten or a" +
					" new file will be generated." +
					needZshDir,
				zshCompGenNew: "only generate the zsh completions file" +
					" if it doesn't already exist. Any pre-existing" +
					" file is protected and an error will be reported." +
					needZshDir,
				zshCompGenShow: "don't generate the zsh completions file." +
					" The file that would have been generated is" +
					" instead printed to standard output.",
				zshCompGenNone: "do nothing.",
			},
			Value: &h.zshMakeCompletions,
		},
		"how to create the zsh completions file."+
			" This specifies whether or if the file should be created."+
			" If it is set to any value other than '"+zshCompGenNone+
			"' then the program will exit after the parameters are processed.",
		param.GroupName(groupName),
		param.Attrs(param.CommandLineOnly|param.DontShowInStdUsage),
	)

	// Final checks

	ps.AddFinalCheck(func() error {
		if h.zshMakeCompletions == zshCompGenShow {
			return nil
		}
		if h.zshMakeCompletions == zshCompGenNone {
			return nil
		}

		if !zshDirParam.HasBeenSet() {
			// These parameters are processed before errors are reported so
			// we should abort the creation of the completion file
			h.zshMakeCompletions = zshCompGenNone

			return fmt.Errorf(
				"the %q parameter has been set (at: %s)"+
					" but the %q parameter has not",
				zshMakeCompletionsParam.Name(),
				strings.Join(
					zshMakeCompletionsParam.WhereSet(),
					" and at "),
				zshDirParam.Name())
		}
		return nil
	})
}
