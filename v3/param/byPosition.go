package param

import (
	"fmt"
	"io"

	"github.com/nickwells/location.mod/location"
)

// ByPos represents a positional parameter. There are numerous strict
// rules about how this can be used. A positional parameter will only be
// checked on the command line (or a slice of strings); it cannot be set by
// an environment variable or in a config file. There must be enough command
// line arguments for all the positional parameters to be set. Only the last
// parameter can be a terminal parameter (a terminal parameter ends the
// parsing and returns). Having a parameter as a terminal parameter will
// allow different parameter sets to be used depending on the value of the
// positional parameter.
type ByPos struct {
	ps           *PSet
	setter       Setter
	name         string
	description  string
	initialValue string
	isTerminal   bool
}

// Name returns the parameter name
func (bp ByPos) Name() string { return bp.name }

// Description returns the parameter description
func (bp ByPos) Description() string { return bp.description }

// InitialValue returns the initialValue of the ByPos parameter
func (p ByPos) InitialValue() string { return p.initialValue }

// AllowedValues returns a description of the values that the ByPos
// parameter can accept
func (p ByPos) AllowedValues() string {
	return p.setter.AllowedValues()
}

// AllowedValuesMap returns the map (which may be nil) of values to
// descriptions for the values that the ByPos parameter can accept
func (p ByPos) AllowedValuesMap() AValMap {
	return p.setter.AllowedValuesMap()
}

// =============================================

// PosOptFunc is the type of a option func used to set various flags on a
// positional parameter
type PosOptFunc func(bp *ByPos) error

// AddByPos will add a new positional parameter to the set of parameters. The
// setter defines the function that should be performed when the parameter is
// processed and will typically be a parameter setter from the paramSetter
// package that will set the value of an associated variable
//
// Various other features of the parameter can be set by the OptFuncs which
// may be passed after the description.
//
// Unlike with the ByName parameter the name given here is purely for
// documentation purposes and should be a very short value just used as a
// hint at the intended purpose. The name should be expanded and explained by
// the description.
func (ps *PSet) AddByPos(name string,
	setter Setter,
	desc string,
	opts ...PosOptFunc) *ByPos {
	if ps.parsed {
		panic("Parameters have already been parsed." +
			" A new positional parameter (" + name + ") cannot be added.")
	}

	if setter.ValueReq() == None {
		panic(fmt.Sprintf(
			"Couldn't add positional parameter %d (%q):"+
				" a positional parameter must take a value",
			len(ps.byPos)+1, name))
	}

	setter.CheckSetter(name)

	checkTerminalFlags(ps)

	bp := &ByPos{
		ps:           ps,
		setter:       setter,
		name:         name,
		description:  desc,
		initialValue: setter.CurrentValue(),
	}

	for _, optFunc := range opts {
		if err := optFunc(bp); err != nil {
			panic(fmt.Sprintf(
				"Couldn't set the options for positional parameter %d (%q): %s",
				len(ps.byPos)+1, name, err))
		}
	}

	ps.byPos = append(ps.byPos, bp)
	return bp
}

// IsTerminal returns true if the ByPos parameter is marked as terminal
func (bp ByPos) IsTerminal() bool { return bp.isTerminal }

// SetAsTerminal is a function which can be passed as a PosOptFunc. It sets
// the flag on the positional parameter indicating that it is terminal. Only
// the last positional parameter can be terminal; this is checked separately
// later.
func SetAsTerminal(bp *ByPos) error {
	if len(bp.ps.byName) > 0 {
		return fmt.Errorf(
			"The param set has %d non-positional parameters."+
				" It cannot also have a terminal positional parameter as"+
				" the non-positional parameters will never be used."+
				" The addition of the standard parameters should be"+
				" turned off when the PSet is created if"+
				" positional parameters are wanted.",
			len(bp.ps.byName))
	}
	bp.isTerminal = true
	return nil
}

// checkTerminalFlags checks the positional parameters in the param set and
// if one of them has the flag indicating that the parameter is terminal then
// it panics. it should be called before adding any extra positional
// parameters
func checkTerminalFlags(ps *PSet) {
	for i, bp := range ps.byPos {
		if bp.isTerminal {
			panic(fmt.Sprintf(
				"Positional parameter %d is marked as terminal"+
					" but is not the last positional parameter", i))
		}
	}
}

// processParam will call the parameter's setter processor and then record
// any errors
func (bp *ByPos) processParam(loc *location.L, val string) {
	err := bp.setter.SetWithVal(bp.name, val)

	if err != nil {
		name := fmt.Sprintf("Positional parameter: %d (%s)",
			loc.Idx(), bp.name)
		bp.ps.errors[name] = append(bp.ps.errors[name],
			loc.Errorf("%s", err.Error()))
	}
}

// StdWriter returns the standard writer of the PSet that this parameter
// belongs to
func (bp ByPos) StdWriter() io.Writer {
	return bp.ps.StdWriter()
}

// ErrWriter returns the error writer of the PSet that this parameter
// belongs to
func (bp ByPos) ErrWriter() io.Writer {
	return bp.ps.ErrWriter()
}
