package param

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/nickwells/location.mod/location"
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
type ErrMap map[string][]error

// ParamSet represents a collection of parameters to be parsed together. A
// program will typically only have one ParamSet but having this makes it
// easier to retain control. You would create the ParamSet in main using the
// paramset.New func which will automatically set the standard help
// member. This lets you know precisely which parameters have been enabled
// before calling Parse
type ParamSet struct {
	parsed          bool
	parseCalledFrom string

	progName     string
	progBaseName string
	progDesc     string

	byPos           []*ByPos
	byName          []*ByName
	nameToParam     map[string]*ByName
	paramGroups     map[string]GroupDesc
	unusedParams    map[string][]string
	errors          ErrMap
	finalChecks     []FinalCheckFunc
	envPrefixes     []string
	configFiles     []ConfigFileDetails
	groupCfgFiles   map[string][]ConfigFileDetails
	remainingParams []string
	terminalParam   string
	maxParamNameLen int

	errWriter io.Writer
	stdWriter io.Writer

	helper Helper

	exitOnParamSetupErr bool
}

// ParamSetOptFunc is the type of a function that can be passed to NewSet
type ParamSetOptFunc func(ps *ParamSet) error

// SetHelper returns a ParamSetOptFunc which can be passed to NewSet. This
// sets the value of the helper to be used by the parameter set and adds the
// parameters that the helper needs. Note that the helper must be set and so
// you must pass some such function. To avoid lots of duplicate code there
// are sensible defaults which can be used in the param/paramset package.
//
// This can only be set once and this will return an error if the helper has
// already been set
func SetHelper(h Helper) ParamSetOptFunc {
	return func(ps *ParamSet) error {
		if ps.helper != nil {
			return errors.New("The helper has already been set")
		}
		ps.helper = h
		h.AddParams(ps)
		return nil
	}
}

// DontExitOnParamSetupErr turns off the standard behaviour of exiting if an
// error is detected while initialising the param set. The error is reported
// in either case
func DontExitOnParamSetupErr(ps *ParamSet) error {
	ps.exitOnParamSetupErr = false
	return nil
}

// SetErrWriter returns a ParamSetOptFunc which can be passed to NewSet. It
// sets the Writer to which error messages are written
func SetErrWriter(w io.Writer) ParamSetOptFunc {
	return func(ps *ParamSet) error {
		if w == nil {
			return fmt.Errorf("param.SetErrWriter cannot take a nil value")
		}
		ps.errWriter = w
		return nil
	}
}

// SetStdWriter returns a ParamSetOptFunc which can be passed to NewSet. It
// sets the Writer to which standard messages are written
func SetStdWriter(w io.Writer) ParamSetOptFunc {
	return func(ps *ParamSet) error {
		if w == nil {
			return fmt.Errorf("param.SetStdWriter cannot take a nil value")
		}
		ps.stdWriter = w
		return nil
	}
}

// StdWriter returns the current value of the StdWriter to which standard
// messages should be written
func (ps *ParamSet) StdWriter() io.Writer {
	return ps.stdWriter
}

// ErrWriter returns the current value of the ErrWriter to which error
// messages should be written
func (ps *ParamSet) ErrWriter() io.Writer {
	return ps.errWriter
}

// SetProgramDescription returns a ParamSetOptFunc which can be passed to
// NewSet. It will set the program description
func SetProgramDescription(desc string) ParamSetOptFunc {
	return func(ps *ParamSet) error {
		ps.progDesc = desc
		return nil
	}
}

// SetProgramDescription sets the program description
func (ps *ParamSet) SetProgramDescription(desc string) {
	ps.progDesc = desc
}

// ProgDesc returns the program description
func (ps *ParamSet) ProgDesc() string {
	return ps.progDesc
}

// NewSet creates a new ParamSet with the various maps and slices
// initialised
func NewSet(psof ...ParamSetOptFunc) (*ParamSet, error) {
	ps := &ParamSet{
		parseCalledFrom: "Parse() not yet called",
		progName:        DfltProgName,
		nameToParam:     make(map[string]*ByName),
		paramGroups:     make(map[string]GroupDesc),
		unusedParams:    make(map[string][]string),
		errors:          make(ErrMap),
		finalChecks:     make([]FinalCheckFunc, 0),

		envPrefixes:   make([]string, 0, 1),
		configFiles:   make([]ConfigFileDetails, 0, 1),
		groupCfgFiles: make(map[string][]ConfigFileDetails),

		terminalParam: DfltTerminalParam,

		errWriter: os.Stderr,
		stdWriter: os.Stdout,

		exitOnParamSetupErr: true,
	}

	for _, f := range psof {
		err := f(ps)
		if err != nil {
			fmt.Fprintf(
				ps.errWriter,
				"An error was detected while creating the ParamSet: %s\n",
				err)
			if ps.exitOnParamSetupErr {
				os.Exit(1)
			}

			return nil, err
		}
	}

	if ps.helper == nil {
		var err = errors.New("A helper must be passed when creating a ParamSet")
		fmt.Fprintln(ps.errWriter, err)
		if ps.exitOnParamSetupErr {
			os.Exit(1)
		}

		return nil, err
	}

	return ps, nil
}

// Errors returns the map of errors for the param set
func (ps ParamSet) Errors() ErrMap { return ps.errors }

// Help will call the helper's Help function
func (ps *ParamSet) Help(message ...string) {
	ps.helper.Help(ps, message...)
}

// ProgName returns the name of the program - the value of the
// zeroth argument
func (ps *ParamSet) ProgName() string { return ps.progName }

// AreSet will return true if Parse has been called or false otherwise
func (ps *ParamSet) AreSet() bool { return ps.parsed }

// UnusedParams returns a copy of the map of unused parameter names. The map
// associates a parameter name with a slice of strings which records where
// the parameter has been set. Unused parameters are always from config files
// or environment variables; unrecognised parameters given on the command
// line are reported as errors.
func (ps *ParamSet) UnusedParams() map[string][]string {
	up := make(map[string][]string, len(ps.unusedParams))
	for pName := range ps.unusedParams {
		copy(up[pName], ps.unusedParams[pName])
	}
	return up
}

// markAsUnused will add the named parameter to the list of unused parameters
func (ps *ParamSet) markAsUnused(name string, loc *location.L) {
	ps.unusedParams[name] = append(ps.unusedParams[name], loc.String())
}

// cmdLineOnly returns true if the parameter should only be set from the
// command line.
func (ps *ParamSet) cmdLineOnly(p *ByName) bool {
	if p.attributes&CommandLineOnly == CommandLineOnly {
		return true
	}
	return false
}

// recordCmdLineOnlyErr records as an error the attempt to set a command-line
// only parameter from a non-command line source
func (ps *ParamSet) recordCmdLineOnlyErr(paramName string, loc *location.L) {
	ps.errors[paramName] = append(ps.errors[paramName],
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
func (ps *ParamSet) recordUnexpectedParam(paramName string, loc *location.L) {
	msg := "this is not a parameter of this program."

	bestSuggestion := ps.findClosestMatch(paramName)
	if bestSuggestion != "" {
		msg += "\n\nDid you mean: " + bestSuggestion + " ?"
	}

	ps.errors[paramName] = append(ps.errors[paramName], loc.Error(msg))
}

func (ps *ParamSet) setNonCommandLineValue(paramParts []string, source string, loc *location.L) bool {
	paramName := paramParts[0]
	p, exists := ps.nameToParam[paramName]

	if !exists {
		ps.markAsUnused(paramName, loc)
		return false
	}

	if ps.cmdLineOnly(p) {
		ps.recordCmdLineOnlyErr(paramName, loc)
		return false
	}

	paramParts = cleanParamParts(p, paramParts)

	p.processParam(source, loc, paramParts)
	return true
}

type existanceRule int

const (
	paramMustExist existanceRule = iota
	paramNeedNotExist
)

func (ps *ParamSet) setValueFromGroupFile(paramParts []string, loc *location.L, gName string) {
	//XXX - needs to be changed
	paramName := paramParts[0]
	p, exists := ps.nameToParam[paramName]

	if !exists {
		ps.recordUnexpectedParam(paramName, loc)
		return
	}
	if p.groupName != gName {
		ps.errors[paramName] = append(ps.errors[paramName],
			loc.Error("this parameter is not a member of group: "+gName))
		return
	}

	if ps.cmdLineOnly(p) {
		ps.recordCmdLineOnlyErr(paramName, loc)
		return
	}

	paramParts = cleanParamParts(p, paramParts)

	p.processParam("group-specific parameter configuration file",
		loc, paramParts)
	return
}

func (ps *ParamSet) setValueFromFile(paramParts []string, loc *location.L, eRule existanceRule) {
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

	if ps.cmdLineOnly(p) {
		ps.recordCmdLineOnlyErr(paramName, loc)
		return
	}

	paramParts = cleanParamParts(p, paramParts)

	p.processParam("parameter configuration file", loc, paramParts)
	return
}

// GetParamByName will return the named parameter if it can be found. The error
// will be set if not
func (ps *ParamSet) GetParamByName(name string) (p *ByName, err error) {
	name = strings.TrimSpace(name)

	p, exists := ps.nameToParam[name]
	if !exists {
		return nil, fmt.Errorf("parameter %s does not exist", name)
	}

	return p, nil
}

// GetParamByPos will return the positional parameter if it exists. The error
// will be set if not.
func (ps *ParamSet) GetParamByPos(idx int) (p *ByPos, err error) {
	if idx < 0 || idx >= len(ps.byPos) {
		return nil, fmt.Errorf("parameter %d does not exist", idx)
	}

	return ps.byPos[idx], nil
}

// AddFinalCheck will add a function to the list of functions to be called
// after all the parameters have been set. Note that multiple functions can
// be set and they will be called in the order that they are added. Each
// function should return an error (or nil) to be added to the list of errors
// detected. All the checks will be called even if one of them returns an
// error
func (ps *ParamSet) AddFinalCheck(fcf FinalCheckFunc) {
	ps.finalChecks = append(ps.finalChecks, fcf)
}

// SetTerminalParam sets the value of the parameter that is used to terminate
// the processing of parameters. This can be used to override the default
// value which is set to DfltTerminalParam
func (ps *ParamSet) SetTerminalParam(s string) { ps.terminalParam = s }

// TerminalParam will return the current value of the terminal parameter.
func (ps *ParamSet) TerminalParam() string { return ps.terminalParam }

// ParamGroup holds details about a group of parameters
type ParamGroup struct {
	GroupName   string
	Desc        string
	Params      []*ByName
	HiddenCount int
	ConfigFiles []ConfigFileDetails
}

// AllParamsHidden returns true if all the parameters are marked as not to be
// shown in the standard usage message, false otherwise
func (pg ParamGroup) AllParamsHidden() bool {
	return len(pg.Params) == pg.HiddenCount
}

// GetParamGroups returns a slice of ParamGroups sorted by group name. Each
// ParamGroup element has a slice of ByName parameters and these are sorted
// by the primary parameter name.
func (ps *ParamSet) GetParamGroups() []*ParamGroup {
	gpMap := make(map[string]ParamGroup)
	for _, p := range ps.byName {
		gp := gpMap[p.groupName]
		gp.GroupName = p.groupName
		gp.Desc = ps.GetGroupDesc(p.groupName)
		gp.Params = append(gp.Params, p)
		if p.attributes&DontShowInStdUsage == DontShowInStdUsage {
			gp.HiddenCount++
		}
		gpMap[p.groupName] = gp
	}

	grpParams := make([]*ParamGroup, 0, len(gpMap))
	for gName := range gpMap {
		gp := gpMap[gName]

		if cfd, ok := ps.groupCfgFiles[gName]; ok {
			gp.ConfigFiles = make([]ConfigFileDetails, len(cfd))
			copy(gp.ConfigFiles, cfd)
		}

		sort.Slice(gp.Params, func(i, j int) bool {
			return gp.Params[i].name < gp.Params[j].name
		})

		grpParams = append(grpParams, &gp)
	}
	sort.Slice(grpParams, func(i, j int) bool {
		return grpParams[i].GroupName < grpParams[j].GroupName
	})

	return grpParams
}
