package param

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/fileparse.mod/fileparse"
	"github.com/nickwells/location.mod/location"
)

// ConfigFileDetails records the details of a configuration
// file. Specifically its name and details about whether or not it must exist
type ConfigFileDetails struct {
	Name         string
	CfConstraint filecheck.Exists
	eRule        existenceRule
}

// ParamsMustExist will return true if the existence rule for the file is set
// to paramMustExist, that is, if the file should only contain valid
// parameters
func (cfd ConfigFileDetails) ParamsMustExist() bool {
	return cfd.eRule == paramMustExist
}

// String returns a string describing the ConfigFileDetails
func (cfd ConfigFileDetails) String() string {
	s := cfd.Name
	if cfd.CfConstraint == filecheck.MustExist {
		s += " (must exist)"
	}
	return s
}

// groupParamLineParser is a type which satisfies the LineParser
// interface and is used to parse group-specific parameter files for the
// paramSet member
type groupParamLineParser struct {
	ps    *PSet
	gName string
}

// paramLineParser is a type which satisfies the LineParser interface and
// is used to parse parameter files for the paramSet member
type paramLineParser struct {
	ps    *PSet
	eRule existenceRule
}

// cmdLineParamLineParser is a type which satisfies the LineParser interface
// and is used to parse command-line supplied parameter files for the
// paramSet member
type cmdLineParamLineParser struct {
	ps *PSet
}

// splitParamName splits the parameter name into two parts around a
// slash. The intention is that the part before the slash is a program name,
// or a comma-separated list of program names and the part after the slash is
// a parameter name.  If there is no slash then the program names list will
// be empty and the paramName will be the whole string. In either case the
// names are stripped of any surrounding whitespace
func splitParamName(pName string) (progNames []string, paramName string) {
	parts := strings.Split(pName, "/")

	if len(parts) == 2 {
		progNames = strings.Split(parts[0], ",")
		for i, pn := range progNames {
			progNames[i] = strings.TrimSpace(pn)
		}
		paramName = strings.TrimSpace(parts[1])
	} else {
		paramName = strings.TrimSpace(parts[0])
	}
	paramName = strings.TrimLeft(paramName, "-")

	return
}

// sliceContains returns true if the slice contains the passed string, false
// otherwise
func sliceContains(slc []string, s string) bool {
	for _, e := range slc {
		if e == s {
			return true
		}
	}
	return false
}

// ParseLine processes the line for a command-line supplied config file.
//
// Firstly it splits the line into two parts around an equal sign, the two
// parts being the parameter specification and the parameter value. Then it
// checks that if the parameter specification has a program part then the
// program name matches the current program name. Finally it attempts to set
// the parameter value from the parameter name and the value string which has
// been stripped of any surrounding whitespace
func (lp cmdLineParamLineParser) ParseLine(line string, loc *location.L) error {
	paramParts := strings.SplitN(line, "=", 2)

	progNames, paramName := splitParamName(paramParts[0])
	if len(progNames) != 0 {
		if !sliceContains(progNames, lp.ps.progBaseName) {
			return nil
		}
	}

	paramParts[0] = paramName

	if len(paramParts) == 2 {
		paramParts[1] = strings.TrimSpace(paramParts[1])
	}

	if p, ok := lp.ps.nameToParam[paramName]; ok {
		p.processParam(loc, paramParts)
	} else {
		lp.ps.recordUnexpectedParam(paramName, loc)
	}

	return nil
}

// ParseLine processes the line for a config file.
//
// Firstly it splits the line into two parts around an equal sign, the two
// parts being the parameter specification and the parameter value. Then it
// checks that if the parameter specification has a program part then the
// program name matches the current program name. Finally it attempts to set
// the parameter value from the parameter name and the value string which has
// been stripped of any surrounding whitespace
func (pflp paramLineParser) ParseLine(line string, loc *location.L) error {
	paramParts := strings.SplitN(line, "=", 2)

	eRule := pflp.eRule
	progNames, paramName := splitParamName(paramParts[0])
	if len(progNames) != 0 {
		if !sliceContains(progNames, pflp.ps.progBaseName) {
			return nil
		}
		eRule = paramMustExist
	}

	paramParts[0] = paramName

	if len(paramParts) == 2 {
		paramParts[1] = strings.TrimSpace(paramParts[1])
	}

	pflp.ps.setValue(paramParts, loc, eRule, "")

	return nil
}

// ParseLine processes the line for a group config file.
//
// Firstly it splits the line into two parts around an equal sign, the two
// parts being the parameter specification and the parameter value. Then it
// checks that if the parameter specification has a program part then the
// program name matches the current program name. Finally it attempts to set
// the parameter value from the parameter name and the value string which has
// been stripped of any surrounding whitespace
func (gpflp groupParamLineParser) ParseLine(line string, loc *location.L) error {
	paramParts := strings.SplitN(line, "=", 2)

	progNames, paramName := splitParamName(paramParts[0])
	if len(progNames) != 0 {
		if !sliceContains(progNames, gpflp.ps.progBaseName) {
			return nil
		}
	}

	paramParts[0] = paramName

	if len(paramParts) == 2 {
		paramParts[1] = strings.TrimSpace(paramParts[1])
	}

	gpflp.ps.setValue(paramParts, loc, paramMustExist, gpflp.gName)

	return nil
}

// checkExistenceConstraint will panic if the existence constraint is
// MustNotExist
func checkExistenceConstraint(fName string, c filecheck.Exists) {
	if c == filecheck.MustNotExist {
		panic(fmt.Sprintf("config file %q: bad existence constraint.", fName))
	}
}

// SetConfigFile will set the list of config files from which to read
// parameter values to just the value given. If it is used with the
// AddConfigFile method below then it should be the first method called.
//
// The config file name may start with ~/ to refer to the home directory
// of the user.
//
// The config file should contain parameter names and values separated
// by an equals sign. Any surrounding space around the parameter name and
// value are stripped off. For instance the following lines will have the
// same effect of setting the value of the myParam attribute to 42:
//
//	myParam  = 42
//	myParam=42
//
// The parameter name can be preceded by a comma-separated list of program
// names and a slash in which case the parameter will only be applied when
// the config file is being parsed by one of the listed programs. The match
// is applied to the basename of the program (the part after the last
// pathname separator). This is particularly useful if there is a config file
// which is shared amongst a number of different programs. It could also be
// used to give different default behaviour when a given program has several
// different names (one binary with different names linked to it). As for the
// parameter name and value any surrounding whitespace is stripped from the
// program name before comparison. For instance:
//
//	myProg/myProgParam = 99
//
// Parameters which don't take a value should appear on a line on their own,
// without an equals character following. As with parameters which take a
// value any surrounding white space is removed and ignored.
//
// Since a parameter file might be shared between several programs, a
// parameter in a config file which is not found in the set of parameters for
// that program is not reported as an error as it might be targeted at a
// different program. This is not the case for parameters which are marked as
// being for specific programs by having a list of program names before the
// parameter name. Similarly for parameters in files which are for a
// particular parameter group, the parameter must be recognised or else it is
// reported as an error.
//
// The config file supports the features of a file parsed by the
// fileparse.FP such as comments and include files.
func (ps *PSet) SetConfigFile(fName string, c filecheck.Exists) {
	checkExistenceConstraint(fName, c)

	ps.configFiles = []ConfigFileDetails{
		{
			Name:         fName,
			CfConstraint: c,
			eRule:        paramNeedNotExist,
		},
	}
}

// SetConfigFileStrict behaves as for SetConfigFile except that parameters
// given in the file must exist for the given program. This is suitable for
// program-specific config files where it can be expected that any parameter
// given in the config file will exist and so it should be reported as an
// error if it does not.
func (ps *PSet) SetConfigFileStrict(fName string, c filecheck.Exists) {
	checkExistenceConstraint(fName, c)

	ps.configFiles = []ConfigFileDetails{
		{
			Name:         fName,
			CfConstraint: c,
			eRule:        paramMustExist,
		},
	}
}

// SetGroupConfigFile sets the config file for the named group. Group config
// files have several constraints: the parameters in the file must only be
// for the named group and it is an error if any parameter in the file is not
// recognised.
//
// Additionally, the param group must already exist.
func (ps *PSet) SetGroupConfigFile(gName, fName string, c filecheck.Exists) {
	if c == filecheck.MustNotExist {
		panic(fmt.Sprintf(
			"config file %q (group %q): bad existence constraint.",
			fName, gName))
	}

	g, ok := ps.groups[gName]
	if !ok {
		panic("param group '" + gName + "' has not been created.")
	}

	g.configFiles = []ConfigFileDetails{
		{
			Name:         fName,
			CfConstraint: c,
			eRule:        paramMustExist,
		},
	}
}

// AddConfigFile adds an additional config file which will also be checked for
// existence and read from. Files are processed in the order they are added.
//
// This can be used to set a system-wide config file and a per-user config
// file that can be used to provide personal preferences.
func (ps *PSet) AddConfigFile(fName string, c filecheck.Exists) {
	checkExistenceConstraint(fName, c)

	ps.configFiles = append(ps.configFiles,
		ConfigFileDetails{
			Name:         fName,
			CfConstraint: c,
			eRule:        paramNeedNotExist,
		})
}

// AddConfigFileStrict behaves as for AddConfigFile except that parameters
// given in the file must exist for the given program. This is suitable for
// program-specific config files where it can be expected that any parameter
// given in the config file will exist and so it should be reported as an
// error if it does not.
func (ps *PSet) AddConfigFileStrict(fName string, c filecheck.Exists) {
	checkExistenceConstraint(fName, c)

	ps.configFiles = append(ps.configFiles,
		ConfigFileDetails{
			Name:         fName,
			CfConstraint: c,
			eRule:        paramMustExist,
		})
}

// AddGroupConfigFile adds an additional config file for the named group.
func (ps *PSet) AddGroupConfigFile(gName, fName string, c filecheck.Exists) {
	if c == filecheck.MustNotExist {
		panic(fmt.Sprintf(
			"config file %q (group %q): bad existence constraint.",
			fName, gName))
	}

	g, ok := ps.groups[gName]
	if !ok {
		panic("param group '" + gName + "' has not been created.")
	}

	g.configFiles = append(g.configFiles,
		ConfigFileDetails{
			Name:         fName,
			CfConstraint: c,
			eRule:        paramMustExist,
		})
}

// ConfigFiles returns a copy of the current config file details.
func (ps *PSet) ConfigFiles() []ConfigFileDetails {
	cf := make([]ConfigFileDetails, len(ps.configFiles))
	copy(cf, ps.configFiles)
	return cf
}

// ConfigFilesForGroup returns a copy of the current config file details for
// the given group name.
func (ps *PSet) ConfigFilesForGroup(gName string) []ConfigFileDetails {
	cf := make([]ConfigFileDetails, len(ps.groups[gName].configFiles))
	copy(cf, ps.groups[gName].configFiles)
	return cf
}

// checkCFErrs will add the errors to the PSet if the error is not a missing
// optional file. Note that we only ignore the error if the file that cannot
// be opened is the config file itself, failure to open included files is
// reported.
func checkCFErrs(ps *PSet, errs []error, cf ConfigFileDetails, desc string) {
	if len(errs) == 0 {
		return
	}

	if len(errs) == 1 {
		err := errs[0]
		var perr *os.PathError
		if errors.As(err, &perr) {
			if errors.Is(err, os.ErrNotExist) &&
				cf.CfConstraint == filecheck.Optional &&
				perr.Path == filepath.Clean(cf.Name) {
				return
			}
		}
	}
	ps.AddErr(desc+": "+cf.Name, errs...)
}

// getParamsFromConfigFiles will construct a line parser and then parse the
// config files - the group-specific config files first and then the common
// files.
func (ps *PSet) getParamsFromConfigFiles() {
	for gName, g := range ps.groups {
		desc := SrcConfigFilePfx + " for " + gName
		fp := fileparse.New(desc, groupParamLineParser{ps: ps, gName: gName})
		for _, cf := range g.configFiles {
			errs := fp.Parse(cf.Name)

			checkCFErrs(ps, errs, cf, desc)
		}
	}

	for _, cf := range ps.configFiles {
		desc := SrcConfigFilePfx
		fp := fileparse.New(desc, paramLineParser{ps: ps, eRule: cf.eRule})
		errs := fp.Parse(cf.Name)
		checkCFErrs(ps, errs, cf, desc)
	}
}

// ConfigFileActionFunc can be called as an action func and will take the
// second entry in the paramVals (which is expected to exist) as the name of
// a config file from which to take parameters.
func ConfigFileActionFunc(_ location.L, p *ByName, paramVals []string) error {
	if len(paramVals) != 2 {
		return errors.New("no config file name parameter has been given")
	}
	name := paramVals[1]
	desc := "supplied config file"

	cf := ConfigFileDetails{Name: name, CfConstraint: filecheck.MustExist}
	fp := fileparse.New(desc, cmdLineParamLineParser{ps: p.ps})
	errs := fp.Parse(cf.Name)
	checkCFErrs(p.ps, errs, cf, desc)

	return nil
}
