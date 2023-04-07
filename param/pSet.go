package param

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/nickwells/errutil.mod/errutil"
	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/pager.mod/pager"
)

// DfltTerminalParam is the default value of the parameter that will stop
// command-line parameters from being processed. Any parameters found after
// this value will be available through the Remainder() func. This default
// value can be overridden through the SetTerminalParam func
const DfltTerminalParam = "--"

// DfltProgName is the program name that will be returned if Parse has not
// yet been called
const DfltProgName = "PROGRAM NAME UNKNOWN"

// ErrMap is the type used to store the errors recorded when parsing the
// parameters. Each map entry represents a parameter for which errors were
// detected; all the errors for that parameter are stored in a slice. Errors
// not related to any individual parameter are stored in the map entry with a
// key of an empty string.
type ErrMap errutil.ErrMap

// FinalCheckFunc is the type of a function to be called after all the
// parameters have been set
type FinalCheckFunc func() error

// PSet represents a collection of parameters to be parsed together. A
// program will typically only have one PSet but having this makes it
// easier to retain control. You would create the PSet in main using the
// paramset.New func which will automatically set the standard help
// member. This lets you know precisely which parameters have been enabled
// before calling Parse
type PSet struct {
	pager.Writers
	parseCalledFrom string

	progName     string
	progBaseName string
	progDesc     string

	byPos        []*ByPos
	byName       []*ByName
	nameToParam  map[string]*ByName
	groups       map[string]*Group
	unusedParams map[string][]string
	errors       ErrMap
	errorCount   int
	finalChecks  []FinalCheckFunc
	envPrefixes  []string
	configFiles  []ConfigFileDetails
	examples     []Example
	references   []Reference
	notes        map[string]*Note

	remainingParams        []string
	terminalParam          string
	terminalParamSeen      bool
	remHandler             RemHandler
	trailingParamsExpected bool
	trailingParamsName     string

	helper Helper

	exitOnParamSetupErr bool
	parsed              bool
}

// PSetOptFunc is the type of a function that can be passed to
// NewSet. These functions can be used to set optional behaviour on the
// parameter set.
type PSetOptFunc func(ps *PSet) error

// SetHelper returns a PSetOptFunc which can be passed to NewSet. This
// sets the value of the helper to be used by the parameter set and adds the
// parameters that the helper needs. Note that the helper must be set and so
// you must pass some such function. To avoid lots of duplicate code there
// are sensible defaults which can be used in the param/paramset package.
//
// This can only be set once and this will return an error if the helper has
// already been set
func SetHelper(h Helper) PSetOptFunc {
	return func(ps *PSet) error {
		if ps.helper != nil {
			return errors.New("The helper has already been set")
		}
		ps.helper = h
		h.AddParams(ps)
		return nil
	}
}

// SetRemHandler sets the value of the remainder handler to be used by the
// parameter set. Note that the handler must be set and so you cannot pass
// nil. The default behaviour is for an error to be reported if there are any
// unprocessed parameters. If you expect additional arguments after either a
// terminal positional parameter or after an explicit end-of-parameters
// parameter (see the TerminalParam method) then you have two choices. You
// can set the remainder handler to the NullRemHandler and process the
// remainder yourself in the body of the program. Alternatively you can pass
// a RemHandler that will handle the remainder in the Parse method.
//
// If this is not set the default behaviour is to report any extra parameters
// after the TerminalParam as errors
func (ps *PSet) SetRemHandler(rh RemHandler) error {
	if rh == nil {
		return errors.New("The remainder handler must not be nil")
	}
	if ps.parsed {
		return errors.New("Parsing is already complete" +
			" - you must set the RemHandler before calling Parse")
	}
	ps.remHandler = rh
	ps.trailingParamsExpected = true
	return nil
}

// SetNamedRemHandler calls SetRemHandler to set the remainder handler and if
// that succeeds it will set the text to be used for the remaining arguments.
// This name will be used in the help message
func (ps *PSet) SetNamedRemHandler(rh RemHandler, name string) error {
	err := ps.SetRemHandler(rh)
	if err != nil {
		return err
	}
	ps.trailingParamsName = name
	return nil
}

// TrailingParamsExpected returns true if a remainder handler has been set
// successfully and false otherwise
func (ps *PSet) TrailingParamsExpected() bool { return ps.trailingParamsExpected }

// TrailingParamsName returns the name that has been given to the trailing
// parameters (if any)
func (ps *PSet) TrailingParamsName() string { return ps.trailingParamsName }

// DontExitOnParamSetupErr turns off the standard behaviour of exiting if an
// error is detected while initialising the param set. The error is reported
// in either case
func DontExitOnParamSetupErr(ps *PSet) error {
	ps.exitOnParamSetupErr = false
	return nil
}

// SetErrWriter returns a PSetOptFunc which can be passed to NewSet. It
// sets the Writer to which error messages are written
func SetErrWriter(w io.Writer) PSetOptFunc {
	return func(ps *PSet) error {
		if w == nil {
			return fmt.Errorf("param.SetErrWriter cannot take a nil value")
		}
		ps.SetErrW(w)
		return nil
	}
}

// SetStdWriter returns a PSetOptFunc which can be passed to NewSet. It
// sets the Writer to which standard messages are written
func SetStdWriter(w io.Writer) PSetOptFunc {
	return func(ps *PSet) error {
		if w == nil {
			return fmt.Errorf("param.SetStdWriter cannot take a nil value")
		}
		ps.SetStdW(w)
		return nil
	}
}

// SetProgramDescription returns a PSetOptFunc which can be passed to
// NewSet. It will set the program description
func SetProgramDescription(desc string) PSetOptFunc {
	return func(ps *PSet) error {
		ps.progDesc = desc
		return nil
	}
}

// SetProgramDescription sets the program description
func (ps *PSet) SetProgramDescription(desc string) {
	ps.progDesc = desc
}

// ProgDesc returns the program description
func (ps *PSet) ProgDesc() string {
	return ps.progDesc
}

// NewSet creates a new PSet with the various maps and slices
// initialised. Generally you would be better off creating a PSet through the
// paramset.New function which will automatically set the default helper
func NewSet(psof ...PSetOptFunc) (*PSet, error) {
	ps := &PSet{
		parseCalledFrom: "Parse() not yet called",
		progName:        DfltProgName,
		progBaseName:    DfltProgName,
		nameToParam:     make(map[string]*ByName),
		groups:          make(map[string]*Group),
		notes:           make(map[string]*Note),
		unusedParams:    make(map[string][]string),
		errors:          ErrMap(*(errutil.NewErrMap())),
		finalChecks:     make([]FinalCheckFunc, 0),

		envPrefixes: make([]string, 0, 1),
		configFiles: make([]ConfigFileDetails, 0, 1),

		terminalParam: DfltTerminalParam,

		Writers: pager.W(),

		exitOnParamSetupErr: true,
	}

	for _, f := range psof {
		err := f(ps)
		if err != nil {
			fmt.Fprintf(ps.ErrW(),
				"An error was detected while creating the PSet: %s\n",
				err)
			if ps.exitOnParamSetupErr {
				os.Exit(1)
			}

			return nil, err
		}
	}
	if ps.remHandler == nil {
		ps.remHandler = dfltRemHandler{}
		ps.trailingParamsExpected = false
	}

	if ps.helper == nil {
		err := errors.New("A helper must be passed when creating a PSet")
		fmt.Fprintln(ps.ErrW(), err)
		if ps.exitOnParamSetupErr {
			os.Exit(1)
		}

		return nil, err
	}

	return ps, nil
}

// Remainder returns any arguments that come after the terminal parameter.
func (ps *PSet) Remainder() []string { return ps.remainingParams }

// Errors returns the map of errors for the param set
func (ps PSet) Errors() ErrMap { return ps.errors }

// AddErr adds the errors to the named entry in the Error Map
func (ps *PSet) AddErr(name string, err ...error) {
	ps.errorCount += len(err)
	ps.errors[name] = append(ps.errors[name], err...)
}

// Help will call the helper's Help function
func (ps *PSet) Help(message ...string) {
	ps.helper.Help(ps, message...)
}

// ProgName returns the name of the program - the value of the zeroth
// argument. Note that this should be called only after the arguments are
// already parsed - before that it will only give the default value
func (ps *PSet) ProgName() string { return ps.progName }

// ProgBaseName returns the base name of the program - the program name with
// any leading directories stripped off. Note that this should be called only
// after the arguments are already parsed - before that it will only give the
// default value
func (ps *PSet) ProgBaseName() string { return ps.progBaseName }

// AreSet will return true if Parse has been called or false otherwise
func (ps *PSet) AreSet() bool { return ps.parsed }

// UnusedParams returns a copy of the map of unused parameter names. The map
// associates a parameter name with a slice of strings which records where
// the parameter has been set. Unused parameters are always from config files
// or environment variables; unrecognised parameters given on the command
// line are reported as errors.
func (ps *PSet) UnusedParams() map[string][]string {
	up := make(map[string][]string, len(ps.unusedParams))
	for pName := range ps.unusedParams {
		up[pName] = make([]string, len(ps.unusedParams[pName]))
		copy(up[pName], ps.unusedParams[pName])
	}
	return up
}

// markAsUnused will add the named parameter to the list of unused parameters
func (ps *PSet) markAsUnused(name string, loc *location.L) {
	ps.unusedParams[name] = append(ps.unusedParams[name], loc.String())
}

// recordCmdLineOnlyErr records as an error the attempt to set a command-line
// only parameter from a non-command line source
func (ps *PSet) recordCmdLineOnlyErr(paramName string, loc *location.L) {
	ps.AddErr(paramName,
		loc.Error("The parameter can only be set on the command line"))
}

// cleanParamParts removes unwanted parts of the paramParts
func cleanParamParts(p *ByName, paramParts []string) []string {
	if len(paramParts) > 1 &&
		p.setter.ValueReq() != Mandatory &&
		paramParts[1] == "" {
		paramParts = paramParts[:1]
	}
	return paramParts
}

// recordUnexpectedParam records that the named parameter is not a parameter
// of this program and if a close match is found it will suggest that
// alternative in the error message
func (ps *PSet) recordUnexpectedParam(paramName string, loc *location.L) {
	msg := "this is not a parameter of this program."

	altNames := SuggestParams(ps, paramName)
	if len(altNames) != 0 {
		msg += "\n\nDid you mean: " + strings.Join(altNames, " or ") + " ?"
	}

	ps.AddErr(paramName, loc.Error(msg))
}

type existenceRule int

const (
	paramMustExist existenceRule = iota
	paramNeedNotExist
)

func (ps *PSet) setValue(paramParts []string, loc *location.L, eRule existenceRule, gName string) {
	paramName := paramParts[0]
	p, exists := ps.nameToParam[paramName]

	if !exists {
		if eRule == paramMustExist {
			ps.recordUnexpectedParam(paramName, loc)
		} else {
			ps.markAsUnused(paramName, loc)
		}
		return
	}

	if gName != "" && p.groupName != gName {
		ps.errors[paramName] = append(ps.errors[paramName],
			loc.Error("this parameter is not a member of group: "+gName))
		return
	}

	if p.AttrIsSet(CommandLineOnly) {
		ps.recordCmdLineOnlyErr(paramName, loc)
		return
	}

	paramParts = cleanParamParts(p, paramParts)

	p.processParam(loc, paramParts)
}

// GetParamByName will return the named parameter if it can be found. The error
// will be set if not
func (ps *PSet) GetParamByName(name string) (p *ByName, err error) {
	name = strings.TrimSpace(name)

	p, exists := ps.nameToParam[name]
	if !exists {
		return nil, fmt.Errorf("parameter %q does not exist", name)
	}

	return p, nil
}

// GetParamByPos will return the positional parameter if it exists. The error
// will be set if not.
func (ps *PSet) GetParamByPos(idx int) (p *ByPos, err error) {
	if idx < 0 || idx >= len(ps.byPos) {
		return nil, fmt.Errorf("parameter %d does not exist", idx)
	}

	return ps.byPos[idx], nil
}

// CountByPosParams will return the number of positional parameters
func (ps *PSet) CountByPosParams() int {
	return len(ps.byPos)
}

// AddFinalCheck will add a function to the list of functions to be called
// after all the parameters have been set. Note that multiple functions can
// be set and they will be called in the order that they are added. Each
// function should return an error (or nil) to be added to the list of errors
// detected. All the checks will be called even if one of them returns an
// error
func (ps *PSet) AddFinalCheck(fcf FinalCheckFunc) {
	ps.finalChecks = append(ps.finalChecks, fcf)
}

// SetTerminalParam sets the value of the parameter that is used to terminate
// the processing of parameters. This can be used to override the default
// value which is set to DfltTerminalParam
func (ps *PSet) SetTerminalParam(s string) { ps.terminalParam = s }

// TerminalParam will return the current value of the terminal
// parameter. This is the parameter which can be given to indicate that any
// following parameters should be handled by the RemHandler (see
// SetRemHandler and SetNamedRemHandler). Unless SetTerminalParam has been
// called this will return the default value: DfltTerminalParam
func (ps *PSet) TerminalParam() string { return ps.terminalParam }

// fixGroups checks that the parameter groups correctly reflect the
// parameters. For instance, the parameter can be in the wrong group if the
// parameter group name is changed after it has been added to the param
// set. Likewise the DontShowInStdUsage attribute can be set after the
// parameter has been added. I take it to be unlikely that the package wil be
// used in this way but better safe than sorry. This will also clear out any
// groups with no parameters, set the HiddenCount on each group and finally
// sort the parameters in each group in alphabetical order.
func (ps *PSet) fixGroups() {
	for name, g := range ps.groups {
		var badGroup bool
		for _, p := range g.Params {
			if p.groupName != name {
				badGroup = true
			}
		}

		if badGroup {
			params := g.Params
			g.Params = nil
			for _, p := range params {
				ps.addByNameToGroup(p)
			}
		}
	}

	for name, g := range ps.groups {
		if g.Params == nil {
			delete(ps.groups, name)
		}
	}

	for _, g := range ps.groups {
		g.SetHiddenCount()
	}

	for _, g := range ps.groups {
		sort.Slice(g.Params, func(i, j int) bool {
			return g.Params[i].name < g.Params[j].name
		})
	}
}

// GetGroups returns a slice of Groups sorted by group name. Each
// Group element has a slice of ByName parameters and these are sorted
// by the primary parameter name.
func (ps *PSet) GetGroups() []*Group {
	ps.fixGroups()
	grpParams := make([]*Group, 0, len(ps.groups))
	for _, g := range ps.groups {
		grpParams = append(grpParams, g)
	}
	sort.Slice(grpParams, func(i, j int) bool {
		return grpParams[i].Name < grpParams[j].Name
	})

	return grpParams
}

// GetGroupByName returns a pointer to the details for the named Group. If
// the name is not recognised then the pointer will be nil
func (ps PSet) GetGroupByName(name string) *Group {
	return ps.groups[name]
}

// HasAltSources returns true if there are any alternative sources
// (configuration files, either general or group-specific, or environment
// variable prefixes) false otherwise
func (ps PSet) HasAltSources() bool {
	if len(ps.configFiles) > 0 {
		return true
	}
	if len(ps.envPrefixes) > 0 {
		return true
	}
	for _, g := range ps.groups {
		if len(g.ConfigFiles) > 0 {
			return true
		}
	}
	return false
}

// HasEnvPrefixes returns true if there are any environment variable prefixes
// for this program, false otherwise
func (ps PSet) HasEnvPrefixes() bool {
	return len(ps.envPrefixes) > 0
}

// HasGlobalConfigFiles returns true if there are any non-group-specific
// config files for this program, false otherwise
func (ps PSet) HasGlobalConfigFiles() bool {
	return len(ps.configFiles) > 0
}

// FindMatchingNamedParams returns a, possibly empty, slice of parameter
// names which match the pattern and an error which will only be non-nil if
// the pattern is malformed.
func (ps PSet) FindMatchingNamedParams(pattern string) ([]string, error) {
	matches := []string{}

	for name := range ps.nameToParam {
		ok, err := path.Match(pattern, name)
		if err != nil {
			return []string{}, err
		}
		if ok {
			matches = append(matches, name)
		}
	}

	return matches, nil
}

// FindMatchingGroups returns a, possibly empty, slice of group names which
// match the pattern and an error which will only be non-nil if the pattern
// is malformed.
func (ps PSet) FindMatchingGroups(pattern string) ([]string, error) {
	matches := []string{}

	for name := range ps.groups {
		ok, err := path.Match(pattern, name)
		if err != nil {
			return []string{}, err
		}
		if ok {
			matches = append(matches, name)
		}
	}

	return matches, nil
}

// FindMatchingNotes returns a, possibly empty, slice of note names which
// match the pattern and an error which will only be non-nil if the pattern
// is malformed.
func (ps PSet) FindMatchingNotes(pattern string) ([]string, error) {
	matches := []string{}

	for name := range ps.notes {
		ok, err := path.Match(pattern, name)
		if err != nil {
			return []string{}, err
		}
		if ok {
			matches = append(matches, name)
		}
	}

	return matches, nil
}
