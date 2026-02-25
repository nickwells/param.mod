package param

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/nickwells/location.mod/location"
)

const ePfxName = "environment variable prefix"

// checkEnvPrefix checks the environment prefix and returns a non-nil error
// if there are any problems.
func checkEnvPrefix(prefix string) error {
	if prefix == "" {
		return errors.New(ePfxName + " must not be empty")
	}

	return nil
}

// SetEnvPrefix will set the prefix for environment variables that are to be
// considered as potential parameters. This prefix is stripped from the name
// and any underscores are replaced with dashes before the environment
// variable name is passed on for matching against parameters.
//
// Any environment prefixes must be set before the parameters are parsed;
// this will panic otherwise.
func (ps *PSet) SetEnvPrefix(pfx string) {
	ps.panicIfAlreadyParsed(
		fmt.Sprintf("can't set the %s: %q", ePfxName, pfx))

	if err := checkEnvPrefix(pfx); err != nil {
		panic(err)
	}

	ps.envPrefixes = nil
	ps.envPrefixes = append(ps.envPrefixes, pfx)
}

// AddEnvPrefix adds a new prefix to the list of environment variable
// prefixes. If the new prefix is empty or a substring of any of the existing
// prefixes or vice versa then it panics.
//
// Any environment prefixes must be added before the parameters are parsed;
// this will panic otherwise.
func (ps *PSet) AddEnvPrefix(pfx string) {
	ps.panicIfAlreadyParsed(fmt.Sprintf("can't add the %s: %q", ePfxName, pfx))

	if err := checkEnvPrefix(pfx); err != nil {
		panic(err)
	}

	errMsgIntro := fmt.Sprintf("invalid %s: %q:", ePfxName, pfx)

	for _, ep := range ps.envPrefixes {
		if strings.HasPrefix(ep, pfx) {
			panic(fmt.Errorf("%s it's a prefix of the already added: %q",
				errMsgIntro, ep))
		}

		if strings.HasPrefix(pfx, ep) {
			panic(fmt.Errorf("%s the already added: %q is a prefix of it",
				errMsgIntro, ep))
		}
	}

	ps.envPrefixes = append(ps.envPrefixes, pfx)
}

// EnvPrefixes returns a copy of the current environment prefixes.
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
