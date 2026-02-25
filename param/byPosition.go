package param

import (
	"fmt"

	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v7/ptypes"
)

// ByPos represents a positional parameter. There are numerous strict
// rules about how this can be used. A positional parameter will only be
// checked on the command line (or a slice of strings); it cannot be set by
// an environment variable or in a config file. There must be enough command
// line arguments for all the positional parameters to be set. Only the last
// parameter can be a terminal parameter (a terminal parameter ends the
// parsing and returns).
//
// Having a parameter as a terminal parameter will allow any following
// parameters to be parsed with different parameter sets. So, for instance,
// the programmer can switch on the value of the positional parameters and
// choose a different PSet for parsing the remaining parameters. This allows
// support for tools with an interface like 'git' or the 'go' command itself.
type ByPos struct {
	BaseParam
	isTerminal bool
}

// =============================================

// ByPosOptFunc is the type of a option func used to set various flags on a
// positional parameter
type ByPosOptFunc = ptypes.OptFunc[ByPos]

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
func (ps *PSet) AddByPos(name string, setter Setter,
	desc string, opts ...ByPosOptFunc,
) *ByPos {
	panicPrefix := fmt.Sprintf("can't add positional parameter %d: %q",
		len(ps.byPos)+1, name)

	ps.panicIfAlreadyParsed(panicPrefix)

	if setter.ValueReq() == None {
		panic(fmt.Errorf("%s: the Setter must take a value", panicPrefix))
	}

	setter.CheckSetter(name)

	checkTerminalFlags(ps)

	bp := &ByPos{
		BaseParam: mkBaseParam(ps, name, setter, desc, caller()),
	}

	for _, optFunc := range opts {
		if err := optFunc(bp); err != nil {
			panic(fmt.Errorf("%s: %w", panicPrefix, err))
		}
	}

	ps.nameToPosParam[name] = bp
	ps.byPos = append(ps.byPos, bp)

	return bp
}

// ByPosValueName returns a ByPosOptFunc which will set the short value name
// used in the parameter summary (it follows the "=" after the parameter
// name). See BaseParam.SetValueName.
func ByPosValueName(vName string) ByPosOptFunc {
	return func(p *ByPos) error {
		return (&p.BaseParam).SetValueName(vName)
	}
}

// ByPosSeeAlso returns a ByPosOptFunc which will add the names of parameters
// to the list of parameters to be referenced when showing the help
// message. See BaseParam.SetSeeAlsoRefs.
func ByPosSeeAlso(refs ...string) ByPosOptFunc {
	source := caller()

	return func(p *ByPos) error {
		return (&p.BaseParam).SetSeeAlsoRefs(source, refs...)
	}
}

// ByPosSeeNote returns a ByPosOptFunc which will add the names of notes to
// the list of notes to be referenced when showing the help message. See
// BaseParam.SetSeeNotes.
func ByPosSeeNote(notes ...string) ByPosOptFunc {
	source := caller()

	return func(p *ByPos) error {
		return (&p.BaseParam).SetSeeNotes(source, notes...)
	}
}

// IsTerminal returns true if the ByPos parameter is marked as terminal
func (bp ByPos) IsTerminal() bool { return bp.isTerminal }

// SetAsTerminal is a function which can be passed as a ByPosOptFunc. It sets
// the flag on the positional parameter indicating that it is terminal. Only
// the last positional parameter can be terminal; this is checked separately
// later.
func SetAsTerminal(bp *ByPos) error {
	if len(bp.ps.byName) > 0 {
		return fmt.Errorf(
			"the param set has %d non-positional parameters."+
				" It cannot also have a terminal positional parameter as"+
				" the non-positional parameters will never be used."+
				" The addition of the standard parameters should be"+
				" turned off when the PSet is created if"+
				" positional parameters are wanted",
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
			panic(fmt.Errorf(
				"positional parameter %d (%q) is marked as terminal"+
					" but is not the last positional parameter", i, bp.name))
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
		bp.ps.AddErr(name,
			loc.Errorf("%s", err.Error()))
	}

	for _, action := range bp.postAction {
		err = action(*loc, (&bp.BaseParam), []string{val})
		if err != nil {
			bp.ps.AddErr(bp.name, loc.Error(err.Error()))
		}
	}
}
