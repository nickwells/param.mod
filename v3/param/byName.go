package param

import (
	"fmt"
	"io"
	"strings"

	"github.com/nickwells/location.mod/location"
)

// =============================================

// ByName represents a parameter which is set by matching a name - it
// includes details such as: the primary name; any alternate names by which it
// can be set; the name of the group to which it belongs; the action to take
// when it is observed (typically setting a value); the description (used in
// the help message); the place(s) where it has been set - the last one takes
// precedence and the attributes
type ByName struct {
	ps              *PSet
	name            string
	altNames        []string
	groupName       string
	setter          Setter
	description     string
	initialValue    string
	whereIsParamSet []string
	attributes      Attributes
	postAction      []ActionFunc
}

// Name returns the name of the ByName parameter
func (p ByName) Name() string { return p.name }

// AltNames returns a copy of the alternative names of the ByName parameter
func (p ByName) AltNames() []string {
	an := make([]string, len(p.altNames))
	copy(an, p.altNames)
	return an
}

// WhereSet returns a copy of the list of places where the ByName parameter
// has been set
func (p ByName) WhereSet() []string {
	ws := make([]string, len(p.whereIsParamSet))
	copy(ws, p.whereIsParamSet)
	return ws
}

// Description returns the description of the ByName parameter
func (p ByName) Description() string { return p.description }

// InitialValue returns the initialValue of the ByName parameter
func (p ByName) InitialValue() string { return p.initialValue }

// GroupName returns the groupName of the ByName parameter
func (p ByName) GroupName() string { return p.groupName }

// ValueReq returns the value requirements for the ByName parameter.
func (p ByName) ValueReq() ValueReq {
	return p.setter.ValueReq()
}

// AllowedValues returns a description of the values that the ByName
// parameter can accept
func (p ByName) AllowedValues() string {
	return p.setter.AllowedValues()
}

// AllowedValuesMap returns the map (which may be nil) of values to
// descriptions for the values that the ByName parameter can accept
func (p ByName) AllowedValuesMap() AValMap {
	return p.setter.AllowedValuesMap()
}

// Attributes holds the attributes of the ByName parameter
type Attributes int32

// Attributes values
const (
	// CommandLineOnly means that the parameter can only be set on the
	// command line. Note that this also includes being set through values
	// passed to the Parse func as a slice of strings
	CommandLineOnly Attributes = 1 << iota
	// MustBeSet means that the parameter must be given - it cannot be
	// omitted
	MustBeSet
	// SetOnlyOnce means that only the first time it is set will have any
	// effect and any subsequent attempts to set it will be ignored. You can
	// control the behaviour when multiple attempts are made through a
	// SetterFunc (see the SetOnce type in the paction package)
	SetOnlyOnce
	// DontShowInStdUsage means that the parameter name will be suppressed
	// when the usage message is printed unless the expanded usage message
	// has been requested
	DontShowInStdUsage
)

// AttrIsSet will return true if the supplied attribute is set on the
// param. Multiple attributes may be given in which case they must all be set
func (p ByName) AttrIsSet(attr Attributes) bool {
	return p.attributes&attr == attr
}

// =============================================

// FinalCheckFunc is the type of a function to be called after all the
// parameters have been set
type FinalCheckFunc func() error

// =============================================

// OptFunc is the type of a option func used to set various flags etc on a
// parameter.
type OptFunc func(p *ByName) error

// Add will add a new named parameter to the set that will be recognised. The
// setter defines the function that should be performed when the parameter is
// processed and will typically be a parameter setter from the psetter
// package that will set the value of an associated variable.
//
// Any leading or trailing spaces are silently removed. Add will panic if the
// parameter has already been used. Add will also panic if the name doesn't
// start with a letter or if it contains any other character than a letter,
// a digit or a dash.
//
// Various other features of the parameter can be set by the OptFuncs which
// may be passed after the description.
func (ps *PSet) Add(name string,
	setter Setter,
	desc string,
	opts ...OptFunc) *ByName {
	if ps.parsed {
		panic("Parameters have already been parsed." +
			" A new named parameter (" + name + ") cannot be added.")
	}

	setter.CheckSetter(name)

	ppCount := len(ps.byPos)
	if ppCount > 0 &&
		ps.byPos[ppCount-1].isTerminal {
		panic("The param set has a terminal positional parameter." +
			" The non-positional parameter " + name + " cannot be added as" +
			" it will never be used.")
	}

	name = strings.TrimSpace(name)

	if err := ps.nameCheck(name); err != nil {
		panic(err.Error())
	}

	p := &ByName{
		ps:           ps,
		name:         name,
		groupName:    DfltGroupName,
		setter:       setter,
		description:  desc,
		initialValue: setter.CurrentValue(),
	}
	ps.nameToParam[name] = p
	ps.byName = append(ps.byName, p)
	p.altNames = append(p.altNames, name)

	for _, optFunc := range opts {
		if err := optFunc(p); err != nil {
			panic(fmt.Sprintf(
				"Error setting the options for param %s: %s.",
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
		g = &Group{Name: p.groupName}
		ps.groups[p.groupName] = g
	}
	g.Params = append(g.Params, p)
}

// HasBeenSet will return true if the parameter has been set.
func (p *ByName) HasBeenSet() bool {
	return len(p.whereIsParamSet) > 0
}

// Attrs returns an OptFunc which will set the attributes of the parameter to
// the passed value.
func Attrs(attrs Attributes) OptFunc {
	return func(p *ByName) error {
		p.attributes = attrs
		return nil
	}
}

// AltName will attach an alternative name to the parameter.
// It will return an error if the alternative name has already been used
func AltName(altName string) OptFunc {
	return func(p *ByName) error {
		altName = strings.TrimSpace(altName)

		if err := p.ps.nameCheck(altName); err != nil {
			return err
		}

		p.ps.nameToParam[altName] = p
		p.altNames = append(p.altNames, altName)
		return nil
	}
}

// GroupName will set the parameter group name for the parameter. The group
// name is stripped of any leading or trailing white space and it is checked
// for validity; an error is returned if it is not valid.  A parameter group
// can be used to collect related parameters together, this grouping will be
// reflected when the usage message is displayed
func GroupName(name string) OptFunc {
	return func(p *ByName) error {
		name = strings.TrimSpace(name)
		err := groupNameCheck(name)
		if err != nil {
			return err
		}
		p.groupName = name
		return nil
	}
}

// processParam will call the parameter's setter processor and then record
// any errors, record where it was set and call any associated post actions
func (p *ByName) processParam(loc *location.L, paramParts []string) {
	var err error

	if (p.attributes&SetOnlyOnce) == SetOnlyOnce &&
		len(p.whereIsParamSet) > 0 {
		// it's already been set so don't process the value
	} else if len(paramParts) == 1 {
		err = p.setter.Set(paramParts[0])
	} else {
		err = p.setter.SetWithVal(paramParts[0], paramParts[1])
	}

	if err != nil {
		p.ps.errors[p.name] = append(p.ps.errors[p.name],
			loc.Error("error with parameter: "+err.Error()))
		return
	}

	p.whereIsParamSet = append(p.whereIsParamSet, loc.String())

	for _, action := range p.postAction {
		err = action(*loc, p, paramParts)

		if err != nil {
			p.ps.errors[p.name] = append(p.ps.errors[p.name],
				loc.Error("error with parameter: "+err.Error()))
		}
	}
}

// StdWriter returns the standard writer of the PSet that this parameter
// belongs to
func (p ByName) StdWriter() io.Writer {
	return p.ps.StdWriter()
}

// ErrWriter returns the error writer of the PSet that this parameter
// belongs to
func (p ByName) ErrWriter() io.Writer {
	return p.ps.ErrWriter()
}
