package param

import (
	"fmt"
	"os"
	"strings"

	"github.com/nickwells/location.mod/location"
)

// SetEnvPrefix will set the prefix for environment variables that are to be
// considered as potential parameters. This prefix is stripped from the name
// and any underscores are replaced with dashes before the environment
// variable name is passed on for matching against parameters
func (ps *PSet) SetEnvPrefix(prefix string) {
	panicMsgIntro := fmt.Sprintf(
		"Can't set %q as an environment variable prefix.", prefix)
	if prefix == "" {
		panic(panicMsgIntro + " The environment prefix must not be empty")
	}

	ps.envPrefixes = nil
	ps.envPrefixes = append(ps.envPrefixes, prefix)
}

// AddEnvPrefix adds a new prefix to the list of environment variable
// prefixes. If the new prefix is empty or a substring of any of the existing
// prefixes or vice versa then it panics
func (ps *PSet) AddEnvPrefix(prefix string) {
	panicMsgIntro := fmt.Sprintf(
		"Can't add %q as an environment variable prefix.", prefix)
	if prefix == "" {
		panic(panicMsgIntro + " The environment prefix must not be empty")
	}

	for _, ep := range ps.envPrefixes {
		if strings.HasPrefix(ep, prefix) {
			panic(fmt.Sprintf("%s It's a prefix of the already added: %q",
				panicMsgIntro, ep))
		}

		if strings.HasPrefix(prefix, ep) {
			panic(fmt.Sprintf("%s The already added: %q is a prefix of it",
				panicMsgIntro, ep))
		}
	}

	ps.envPrefixes = append(ps.envPrefixes, prefix)
}

// EnvPrefixes returns a copy of the current environment prefixes
func (ps *PSet) EnvPrefixes() []string {
	ep := make([]string, len(ps.envPrefixes))
	copy(ep, ps.envPrefixes)

	return ep
}

// ConvertParamNameToEnvVarName converts a parameter name to a valid
// environment variable name. Note that in order to be recognised it will
// need to be prefixed by a recognised environment variable prefix as
// added by AddEnvPrefix. It should have the opposite effect to the
// ConvertEnvVarNameToParamName function
func ConvertParamNameToEnvVarName(name string) string {
	return strings.ReplaceAll(name, "-", "_")
}

// ConvertEnvVarNameToParamName converts an environment variable name to a
// parameter name. Any environment variable prefix (as added by AddEnvPrefix)
// should have been stripped off first. It should have the opposite effect to
// the ConvertParamNameToEnvVarName function
func ConvertEnvVarNameToParamName(name string) string {
	return strings.ReplaceAll(name, "_", "-")
}

func (ps *PSet) getParamsFromEnvironment() {
	if len(ps.envPrefixes) == 0 {
		return
	}

	loc := location.New("")
	loc.SetNote(SrcEnvironment)

	for _, param := range os.Environ() {
		paramName, paramVal, hasParamVal := strings.Cut(param, "=")
		for _, envPrefix := range ps.envPrefixes {
			trimmedParam := strings.TrimPrefix(paramName, envPrefix)
			// We only process those env vars that start with the
			// envPrefix (so trimming the prefix will change the name)
			if trimmedParam != paramName {
				paramParts := append([]string{},
					ConvertEnvVarNameToParamName(trimmedParam))
				if hasParamVal {
					paramParts = append(paramParts, paramVal)
				}

				loc.SetContent(param)
				ps.setValue(paramParts, loc, paramNeedNotExist, "")

				break // we've found a match so stop looking
			}
		}
	}
}
