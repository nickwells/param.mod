package paction

import (
	"strings"

	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v5/param"
)

type ParamTestFunc func(location.L, *param.ByName, []string) bool

// IsACommandLineParam returns true if the parameter has been set through
// the command line, false otherwise.
func IsACommandLineParam(loc location.L, _ *param.ByName, _ []string) bool {
	return loc.Note() == param.SrcCommandLine
}

// IsNotACommandLineParam returns true if the parameter has not been set through
// the command line, false otherwise.
func IsNotACommandLineParam(loc location.L, _ *param.ByName, _ []string) bool {
	return !IsACommandLineParam(loc, nil, nil)
}

// IsAConfigFileParam returns true if the parameter has been set through
// a line in a configuration file, false otherwise.
func IsAConfigFileParam(loc location.L, _ *param.ByName, _ []string) bool {
	return strings.HasPrefix(loc.Note(), param.SrcConfigFilePfx)
}

// IsNotAConfigFileParam returns true if the parameter has not been set through
// a line in a configuration file, false otherwise.
func IsNotAConfigFileParam(loc location.L, _ *param.ByName, _ []string) bool {
	return !IsAConfigFileParam(loc, nil, nil)
}

// IsAnEnvVarParam returns true if the parameter has been set through
// an environment variable, false otherwise.
func IsAnEnvVarParam(loc location.L, _ *param.ByName, _ []string) bool {
	return loc.Note() == param.SrcEnvironment
}

// IsNotAnEnvVarParam returns true if the parameter has not been set through
// an environment variable, false otherwise.
func IsNotAnEnvVarParam(loc location.L, _ *param.ByName, _ []string) bool {
	return loc.Note() == param.SrcEnvironment
}
