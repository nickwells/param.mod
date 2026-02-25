package param

import (
	"fmt"
	"path"
	"slices"
	"sort"
	"strings"

	"github.com/nickwells/errutil.mod/errutil"
	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v7/ptypes"
)

// dfltTerminalParam is the default value of the parameter that will stop
// command-line parameters from being processed. Any parameters found after
// this value will be available through the TrailingParams() func. This default
// value can be overridden through the SetTerminalParam func
const dfltTerminalParam = "--"

// dfltProgName is the program name that will be returned if Parse has not
// yet been called
const dfltProgName = "PROGRAM NAME UNKNOWN"

// FinalCheckFunc is the type of a function to be called after all the
// parameters have been set
type FinalCheckFunc func() error

// PSet holds the collection of parameters to be parsed together. Typical
// usage would involve creating a PSet (usually by calling paramset.New) and
// then adding the parameters to the PSet before calling Parse to set the
// parameter values.
type PSet struct {
	parseCalledFrom string

	progName     string
	progBaseName string
	progDesc     string

	byPos          []*ByPos
	byName         []*ByName
	nameToParam    map[string]*ByName
	nameToPosParam map[string]*ByPos
	groups         map[string]*Group
	unusedParams   map[string][]string
	errMap         errutil.ErrMap
	errorCount     int
	finalChecks    []FinalCheckFunc
	envPrefixes    []string
	configFiles    []ConfigFileDetails
	examples       []Example
	references     []Reference
	notes          map[string]*Note

	paramPrefixes  []string
	shortestPrefix string

	trailingParams         []string
	terminalParam          string
	terminalParamSeen      bool
	trailingParamsExpected bool
	trailingParamsName     string

	helper Helper

	helpRequired bool
	shouldExit   bool
	exitStatus   int

	parsed bool
}

// PSetOptFunc is the type of a function that can be passed to
// [NewSet]. These functions can be used to set optional behaviour on the
// parameter set.
type PSetOptFunc = ptypes.OptFunc[PSet]

// TrailingParamsExpected returns true if a remainder handler has been set
// successfully and false otherwise
func (ps *PSet) TrailingParamsExpected() bool {
	return ps.trailingParamsExpected
}

// SetTrailingParamsExpected sets the flag notifying the PSet that trailing
// parameters are allowed and should not be flagged as an error. See also
// [SetTrailingParamsExpected] (an option function that can be passed to
// [NewSet]).
//
// This must be called before the parameters are parsed; this will
// panic otherwise.
func (ps *PSet) SetTrailingParamsExpected() {
	ps.panicIfAlreadyParsed(
		"the trailing parameters expected flag may not be set")

	ps.trailingParamsExpected = true
}

// SetTrailingParamsExpected is a PSetOptFunc that sets the flag notifying
// the PSet that trailing parameters are allowed and should not be flagged as
// an error. See also the [PSet.SetTrailingParamsExpected] method.
func SetTrailingParamsExpected(ps *PSet) error {
	ps.trailingParamsExpected = true

	return nil
}

// TrailingParamsName returns the name that has been given to the trailing
// parameters (if any)
func (ps *PSet) TrailingParamsName() string { return ps.trailingParamsName }

// SetTrailingParamsName sets the name to be given to the trailing parameters
// (if any) in help messages. It also sets the flag notifying the PSet that
// trailing parameters are allowed and should not be flagged as an error. See
// also [SetTrailingParamsName] (which returns an option function that can
// be passed to [NewSet]).
//
// This must be called before the parameters are parsed; this will
// panic otherwise.
func (ps *PSet) SetTrailingParamsName(name string) {
	ps.panicIfAlreadyParsed("the trailing parameters name may not be set")

	ps.trailingParamsName = name
	ps.trailingParamsExpected = true
}

// SetTrailingParamsName returns a PSetOptFunc which can be passed to
// [NewSet]. It will set the name to be given to any trailing parameters. It
// also sets the flag notifying the PSet that trailing parameters are allowed
// and should not be flagged as an error. See also the
// [PSet.SetTrailingParamsName] method.
func SetTrailingParamsName(name string) PSetOptFunc {
	return func(ps *PSet) error {
		ps.SetTrailingParamsName(name)

		return nil
	}
}

// SetProgramDescription returns a PSetOptFunc which can be passed to
// [NewSet]. It will set the program description. See also the
// [PSet.SetProgramDescription] method.
func SetProgramDescription(desc string) PSetOptFunc {
	return func(ps *PSet) error {
		ps.progDesc = desc

		return nil
	}
}

// SetProgramDescription sets the program description. See also
// [SetProgramDescription] (an option function that can be passed to
// [NewSet]).
//
// This must be called before the parameters are parsed; this will
// panic otherwise.
func (ps *PSet) SetProgramDescription(desc string) {
	ps.progDesc = desc
}

// ProgDesc returns the program description
func (ps *PSet) ProgDesc() string {
	ps.panicIfAlreadyParsed("the program description may not be set")

	return ps.progDesc
}

// SetParamPrefixes returns a PSetOptFunc which can be passed to [NewSet]. It
// will set the list of allowed parameter prefixes that are to be removed
// before parameter processing. See also the [PSet.SetParamPrefixes] method.
func SetParamPrefixes(pp ...string) PSetOptFunc {
	return func(ps *PSet) error {
		ps.SetParamPrefixes(pp)

		return nil
	}
}

// SetParamPrefixes sets the list of allowed parameter prefixes that are to
// be removed before parameter processing. The prefixes are sorted into
// longest-first order (preserving the order of equal length prefixes). Note
// that the prefixes used do not participate in the parameter processing but
// are discarded unseen. See also [SetParamPrefixes] (an option function that
// can be passed to [NewSet]).
//
// The prefixes must be set before the parameters are parsed; this will
// panic otherwise.
func (ps *PSet) SetParamPrefixes(pp []string) {
	ps.panicIfAlreadyParsed("parameter prefixes may not be set")

	slices.SortStableFunc(pp, func(a, b string) int { return len(b) - len(a) })

	ps.paramPrefixes = pp
	ps.shortestPrefix = ""

	if len(pp) > 0 {
		ps.shortestPrefix = pp[len(pp)-1]
	}
}

// ParamPrefixes returns the allowed parameter prefixes for this parameter
// set.
func (ps *PSet) ParamPrefixes() []string {
	return ps.paramPrefixes
}

// ShortestPrefix returns the shortest parameter prefix for this parameter
// set. This will be used as the parameter name prefix in the standard help
// message.
func (ps *PSet) ShortestPrefix() string {
	return ps.shortestPrefix
}

// NewSet creates a new PSet with the various internal maps and slices
// initialised. Generally you would be better off creating a PSet through the
// paramset.New function which will automatically set the default helper. If
// there are any problems constructing the PSet then the function will panic
// with the error.
func NewSet(h Helper, psof ...PSetOptFunc) *PSet {
	ps := &PSet{
		parseCalledFrom: "Parse() not yet called",
		progName:        dfltProgName,
		progBaseName:    dfltProgName,
		nameToParam:     make(map[string]*ByName),
		nameToPosParam:  make(map[string]*ByPos),
		groups:          make(map[string]*Group),
		notes:           make(map[string]*Note),
		unusedParams:    make(map[string][]string),
		errMap:          *(errutil.NewErrMap()),
		finalChecks:     make([]FinalCheckFunc, 0),

		envPrefixes: make([]string, 0, 1),
		configFiles: make([]ConfigFileDetails, 0, 1),

		terminalParam:  dfltTerminalParam,
		paramPrefixes:  []string{"--", "-"},
		shortestPrefix: "-",

		helper: h,
	}

	h.AddParams(ps)

	for _, f := range psof {
		err := f(ps)
		if err != nil {
			panic(fmt.Errorf("while creating the PSet: %w", err))
		}
	}

	return ps
}

// TrailingParams returns any arguments that come after the terminal
// parameter. See also the [PSet.SetTrailingParamsName] and
// [PSet.SetTrailingParamsExpected] methods.
func (ps *PSet) TrailingParams() []string { return ps.trailingParams }

// Errors returns the map of errors for the param set
func (ps PSet) Errors() errutil.ErrMap { return ps.errMap }

// AddErr adds the errors to the named entry in the Error Map. Any nil errors
// are filtered out of the slice and if the slice is empty no change is made
func (ps *PSet) AddErr(name string, errs ...error) {
	errs = slices.DeleteFunc(
		errs, func(e error) bool { return e == nil })

	if len(errs) == 0 {
		return
	}

	ps.errorCount += len(errs)
	ps.errMap[name] = append(ps.errMap[name], errs...)
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

// AlreadyParsed will return a non-nil error if Parse has been called.
func (ps *PSet) AlreadyParsed() error {
	if ps.parsed {
		return fmt.Errorf("param.PSet.Parse has already been called, from: %s",
			ps.parseCalledFrom)
	}

	return nil
}

// panicIfAlreadyParsed will panic if the parameters have already been
// parsed.
func (ps *PSet) panicIfAlreadyParsed(msg string) {
	if err := ps.AlreadyParsed(); err != nil {
		if msg == "" {
			panic(err)
		}

		panic(fmt.Errorf("%s: %w", msg, err))
	}
}

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

func (ps *PSet) setValue(
	paramParts []string, loc *location.L, eRule existenceRule, gName string,
) {
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
		ps.errMap[paramName] = append(ps.errMap[paramName],
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
func (ps *PSet) GetParamByName(name string) (*ByName, error) {
	name = strings.TrimSpace(name)

	if p, exists := ps.nameToParam[name]; exists {
		return p, nil
	}

	return nil, fmt.Errorf("named parameter %q does not exist", name)
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
// error.
//
// This will panic if called after the parameters have been parsed.
func (ps *PSet) AddFinalCheck(fcf FinalCheckFunc) {
	ps.panicIfAlreadyParsed("cannot add a final check function")

	ps.finalChecks = append(ps.finalChecks, fcf)
}

// SetTerminalParam sets the value of the parameter that is used to terminate
// the processing of parameters. This can be used to override the default
// value which is set to dfltTerminalParam. See also [SetTerminalParam]
// (which returns an option function that can be passed to [NewSet]).
//
// This will panic if called after the parameters have been parsed.
func (ps *PSet) SetTerminalParam(s string) {
	ps.panicIfAlreadyParsed("cannot set the terminal parameter")

	ps.terminalParam = s
}

// SetTerminalParam returns a PSetOptFunc which can be passed to NewSet. It
// will set the value of the parameter that is used to terminate the
// processing of parameters. This can be used to override the default value
// which is set to dfltTerminalParam. See also the [PSet.SetTerminalParam]
// method.
func SetTerminalParam(s string) PSetOptFunc {
	return func(ps *PSet) error {
		ps.SetTerminalParam(s)

		return nil
	}
}

// TerminalParam will return the current value of the terminal
// parameter. This is the parameter which can be given to indicate that
// argument parsing should stop and any subsequent arguments will be
// available through the PSet.TrailingParams method. Unless SetTerminalParam
// has been called this will return the default value: DfltTerminalParam
func (ps *PSet) TerminalParam() string { return ps.terminalParam }

// fixGroups checks that the parameter groups correctly reflect the
// parameters. For instance, the parameter can be in the wrong group if the
// parameter group name is changed after it has been added to the param
// set. Likewise the DontShowInStdUsage attribute can be set after the
// parameter has been added. I take it to be unlikely that the package will
// be used in this way but better safe than sorry. This will also clear out
// any groups with no parameters and sort the parameters in each group in
// alphabetical order.
func (ps *PSet) fixGroups() {
	for name, g := range ps.groups {
		var badGroup bool

		for _, p := range g.params {
			if p.groupName != name {
				badGroup = true
			}
		}

		if badGroup {
			params := g.params
			g.params = nil

			for _, p := range params {
				ps.addByNameToGroup(p)
			}
		}
	}

	for name, g := range ps.groups {
		if g.params == nil {
			delete(ps.groups, name)
		}
	}

	for _, g := range ps.groups {
		sort.Slice(g.params, func(i, j int) bool {
			return g.params[i].name < g.params[j].name
		})
	}
}

// GetGroups returns a slice of Groups sorted by group name. Each
// Group element has a slice of ByName parameters and these are sorted
// by the primary parameter name.
func (ps *PSet) GetGroups() []*Group {
	ps.fixGroups()

	groups := make([]*Group, 0, len(ps.groups))
	for _, g := range ps.groups {
		groups = append(groups, g)
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].name < groups[j].name
	})

	return groups
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
	if ps.HasGlobalConfigFiles() {
		return true
	}

	if ps.HasEnvPrefixes() {
		return true
	}

	for _, g := range ps.groups {
		if len(g.configFiles) > 0 {
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

// HelpRequired sets the help required flag which will cause PSet.Parse to
// call the Helper's Help method.
func (ps *PSet) HelpRequired() {
	ps.helpRequired = true
}

// ShouldExit sets the shouldExit flag which will cause PSet.Parse to exit.
func (ps *PSet) ShouldExit() {
	ps.shouldExit = true
}

// SetExitStatus sets the exitStatus to the supplied value (if it is greater
// than the prior status) and sets the shouldExit flag which will cause
// PSet.Parse to exit.
func (ps *PSet) SetExitStatus(s int) {
	ps.ShouldExit()

	if s > ps.exitStatus {
		ps.exitStatus = s
	}
}
