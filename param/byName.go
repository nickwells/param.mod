package param

import (
	"fmt"
	"strings"

	"github.com/nickwells/location.mod/location"
)

// =============================================

// ByName represents a parameter which is set by matching a name - it
// includes details such as:
//
// - the parameter name
//
// - any alternate names by which it can be set
//
// - the name of the group to which it belongs
//
// - the action to take when it is observed (typically setting a value)
//
// - the description (used in the help message)
//
// - the attributes
//
// - the place(s) where it has been set - the last one takes precedence
//
// This should not be created directly. The Add method on the PSet will
// generate a ByName and add it to the set of program parameters, A pointer
// to the ByName is returned.
//
// Most of the values and methods on this are for the purposes of generating
// the help message and enforcing usage restrictions when parsing the program
// parameters.
//
// For anyone not writing a bespoke help class the only useful methods on
// this class are the HasBeenSet and WhereSet methods. You can record the
// ByName pointer returned by the PSet.Add method and then in a FinalCheck
// function you can test whether or not this and other parameters were set
// and confirm that the combination of parameters is allowed. See the
// PSet.AddFinalCheck method.
type ByName struct {
	BaseParam
	altNames        []string
	groupName       string
	whereIsParamSet []string
	attributes      Attributes
}

// AltNames returns a copy of the alternative names of the ByName parameter
func (p ByName) AltNames() []string {
	an := make([]string, len(p.altNames))

	copy(an, p.altNames)

	return an
}

// HasBeenSet will return true if the parameter has been set.
func (p ByName) HasBeenSet() bool {
	return len(p.whereIsParamSet) > 0
}

// WhereSet returns a copy of the list of places where the ByName parameter
// has been set
func (p ByName) WhereSet() []string {
	ws := make([]string, len(p.whereIsParamSet))

	copy(ws, p.whereIsParamSet)

	return ws
}

// GroupName returns the groupName of the ByName parameter
func (p ByName) GroupName() string { return p.groupName }

// Attributes records various flags that can be set on a ByName parameter
type Attributes int32

const (
	// CommandLineOnly means that the parameter can only be set on the
	// command line. Note that this also includes being set through values
	// passed to the Parse func as a slice of strings. You might want to set
	// this attribute on parameters which would always be different between
	// command invocations or where setting it would make the program
	// terminate. For instance, it is set on the standard help attributes as
	// setting these in a configuration file would never allow the program to
	// execute.
	CommandLineOnly Attributes = 1 << iota
	// MustBeSet means that the parameter must be given - it cannot be
	// omitted
	MustBeSet
	// SetOnlyOnce means that only the first time it is set will have any
	// effect and any subsequent attempts to set it will be ignored. You can
	// control the behaviour when multiple attempts are made through a
	// SetterFunc (see the SetOnce type in the paction package). You might
	// want to set this on a parameter that you want to set for all users in
	// a global configuration file that only the system administrator can
	// edit. This would allow you to set a system-wide policy.
	SetOnlyOnce
	// DontShowInStdUsage means that the parameter name will be suppressed
	// when the usage message is printed unless the expanded usage message
	// has been requested
	DontShowInStdUsage
	// IsTerminalParam means that when this parameter is encountered
	// command-line parameter processing will stop. Any further parameters
	// will be ignored and added to the slice of remaining params for
	// subsequent processing by the application. Setting it will also set the
	// CommandLineOnly attribute.
	IsTerminalParam
)

// AttrIsSet will return true if the supplied attribute is set on the
// param. Multiple attributes may be given in which case they must all be set
func (p ByName) AttrIsSet(attr Attributes) bool {
	return p.attributes&attr == attr
}

// =============================================

// ByNameOptFunc is the type of a option func used to set various flags etc on a
// parameter.
type ByNameOptFunc func(p *ByName) error

// Add will add a new named parameter to the set that will be recognised. The
// setter defines the function that should be performed when the parameter is
// processed and will typically be a parameter setter from the psetter
// package that will set the value of an associated variable.
//
// Any leading or trailing spaces are silently removed. Add will panic if the
// parameter has already been used. Add will also panic if the name doesn't
// start with a letter or if it contains any other character than a letter, a
// digit or a dash.
//
// Various other features of the parameter can be set by the ByNameOptFunc's
// which may be passed after the description.
func (ps *PSet) Add(
	name string, setter Setter, desc string, opts ...ByNameOptFunc,
) *ByName {
	if err := ps.AlreadyParsed(); err != nil {
		panic(fmt.Errorf("named parameter (%q) can't be added: %w", name, err))
	}

	setter.CheckSetter(name)

	ppCount := len(ps.byPos)
	if ppCount > 0 &&
		ps.byPos[ppCount-1].isTerminal {
		panic(
			fmt.Errorf("named parameter (%q) can't be added:"+
				" it can never be used as"+
				" the param set has a terminal positional parameter",
				name))
	}

	name = strings.TrimSpace(name)

	whereAdded := caller()
	if err := ps.nameCheck(name, whereAdded); err != nil {
		panic(fmt.Errorf("named parameter (%q) can't be added: %w", name, err))
	}

	p := &ByName{
		BaseParam: mkBaseParam(ps, name, setter, desc, whereAdded),
		groupName: DfltGroupName,
	}
	ps.nameToParam[name] = p
	ps.byName = append(ps.byName, p)
	p.altNames = append(p.altNames, name)

	for _, optFunc := range opts {
		if err := optFunc(p); err != nil {
			panic(fmt.Errorf("named parameter (%q) can't be added: %w",
				name, err))
		}
	}

	ps.addByNameToGroup(p)

	return p
}

// addByNameToGroup adds the ByName parameter to the appropriate parameter
// group in the PSet
func (ps *PSet) addByNameToGroup(p *ByName) {
	g, ok := ps.groups[p.groupName]
	if !ok {
		g = &Group{name: p.groupName}
		ps.groups[p.groupName] = g
	}

	g.params = append(g.params, p)
}

// Attrs returns a ByNameOptFunc which will set the attributes of the
// parameter to the passed value. Note that if the IsTerminal attribute is
// set then the CommandLineOnly attribute is forced on as well.
func Attrs(attrs Attributes) ByNameOptFunc {
	return func(p *ByName) error {
		if attrs&IsTerminalParam == IsTerminalParam {
			attrs |= CommandLineOnly

			p.ps.TrailingParamsExpected()
		}

		p.attributes = attrs

		return nil
	}
}

// AltNames returns a ByNameOptFunc which will attach multiple alternative
// names to the parameter.  It will return an error if any alternative name
// has already been used
func AltNames(altNames ...string) ByNameOptFunc {
	return func(p *ByName) error {
		for _, altName := range altNames {
			altName = strings.TrimSpace(altName)

			if err := p.ps.nameCheck(altName, p.whereAdded); err != nil {
				return err
			}

			p.ps.nameToParam[altName] = p
			p.altNames = append(p.altNames, altName)
		}

		return nil
	}
}

// ValueName returns a ByNameOptFunc which will set the short value name used
// in the parameter summary (it follows the "=" after the parameter
// name). See BaseParam.SetValueName.
func ValueName(vName string) ByNameOptFunc {
	return func(p *ByName) error {
		return (&p.BaseParam).SetValueName(vName)
	}
}

// SeeAlso returns a ByNameOptFunc which will add the names of parameters to
// the list of parameters to be referenced when showing the help message. See
// BaseParam.SetSeeAlsoRefs.
func SeeAlso(refs ...string) ByNameOptFunc {
	source := caller()

	return func(p *ByName) error {
		return (&p.BaseParam).SetSeeAlsoRefs(source, refs...)
	}
}

// SeeNote returns a ByNameOptFunc which will add the names of notes to the
// list of notes to be referenced when showing the help message. See
// BaseParam.SetSeeNotes.
func SeeNote(notes ...string) ByNameOptFunc {
	source := caller()

	return func(p *ByName) error {
		return (&p.BaseParam).SetSeeNotes(source, notes...)
	}
}

// GroupName returns a ByNameOptFunc which will set the parameter group name
// for the parameter. The group name is stripped of any leading or trailing
// white space and it is checked for validity; an error is returned if it is
// not valid.  A parameter group can be used to collect related parameters
// together, this grouping will be reflected when the usage message is
// displayed
func GroupName(name string) ByNameOptFunc {
	return func(p *ByName) error {
		name = strings.TrimSpace(name)

		if err := GroupNameCheck(name); err != nil {
			return err
		}

		p.groupName = name

		return nil
	}
}

// processParam will call the parameter's setter processor and then record
// any errors, record where it was set and call any associated post actions
func (p *ByName) processParam(loc *location.L, paramParts []string) {
	if p.AttrIsSet(SetOnlyOnce) && p.HasBeenSet() {
		p.ps.AddErr(p.name,
			loc.Error(fmt.Sprintf(
				"This may only be set once but has already been set at %s",
				p.whereIsParamSet[0])))

		return
	}

	if p.AttrIsSet(IsTerminalParam) {
		p.ps.terminalParamSeen = true
	}

	var err error

	const (
		nameOnly = 1
		hasValue = 2
	)

	switch len(paramParts) {
	case nameOnly:
		err = p.setter.Set(paramParts[0])
	case hasValue:
		err = p.setter.SetWithVal(paramParts[0], paramParts[1])
	default:
		err = fmt.Errorf("bad parameter: %q", paramParts)
	}

	if err != nil {
		p.ps.AddErr(p.name, loc.Error(err.Error()))
		return
	}

	p.whereIsParamSet = append(p.whereIsParamSet, loc.String())

	for _, action := range p.postAction {
		err = action(*loc, (&p.BaseParam), paramParts)
		if err != nil {
			p.ps.AddErr(p.name, loc.Error(err.Error()))
		}
	}
}
