package phelp

// Code generated by mkparamfilefunc; DO NOT EDIT.
// with parameters set at:
//	[command line]: Argument:4: "-funcs" "personalOnly"
//	[command line]: Argument:2: "-group" "stdParams-cmpl"
//	[command line]: Argument:5: "-private"
import (
	"path/filepath"

	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/xdg.mod/xdg"
)

/*
setConfigFileForGroupStdParamsCmpl adds a config file to the set which the param
parser will process before checking the command line parameters.
*/
func setConfigFileForGroupStdParamsCmpl(ps *param.PSet) error {
	baseDir := xdg.ConfigHome()

	ps.AddGroupConfigFile("stdParams-cmpl",
		filepath.Join(baseDir,
			"github.com",
			"nickwells",
			"param.mod",
			"v6",
			"phelp",
			"group-stdParams-cmpl.cfg"),
		filecheck.Optional)
	return nil
}
