package paction

import (
	"fmt"
	"os"

	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v5/param"
)

// SetOnce is used to record if a parameter has been previously set and take
// appropriate action if it has.
//
// This can be used if you have a parameter that you want to set in a system-
// wide config file (owned, for instance, by root) and which the users cannot
// override. It doesn't make a lot of sense to have this if you don't have a
// system-wide config file since the user can just change it themselves -
// it'll be more irritating than anything else.
//
// Note that this is an alternative mechanism to the SetOnlyOnce parameter
// attribute. It does NOT guarantee that the value won't be set multiple times,
// only that an error will be returned (if the action value is set to
// ErrorOnMultipleTries). If the SetOnlyOnce attribute is set on the
// parameter then this function will only be called the first time the
// parameter is set and will therefore have no effect.
type SetOnce struct {
	paramsSetAt param.Sources
}

// SetOnceErrAction records the behaviour to exhibit when the parameter has
// been set more than once
type SetOnceErrAction int32

// SetOnceErrAction values
const (
	// IgnoreMultipleTries means that no further action beyond recording that
	// the parameter has been set will be taken
	IgnoreMultipleTries SetOnceErrAction = iota
	// ErrorOnMultipleTries means that more than one attempt to set the value
	// will cause an error to be returned by the action function
	ErrorOnMultipleTries
	// ExitOnMultipleTries means that more than one attempt to set the value
	// will cause the program to exit
	ExitOnMultipleTries
)

// MakeActionFunc returns a function suitable for passing to the PostAction
// method of a param.ByName object. This uses a pre-existing SetOnce value
// and so it could be set on multiple different parameters but the expected
// use is to have a separate SetOnce value for each parameter that you want
// to protect
func (so *SetOnce) MakeActionFunc(action SetOnceErrAction) param.ActionFunc {
	return func(loc location.L, p *param.ByName, paramVals []string) error {
		if len(so.paramsSetAt) == 0 {
			so.paramsSetAt = append(so.paramsSetAt,
				param.Source{
					From:      loc.Source(),
					Loc:       loc,
					ParamVals: paramVals,
					Param:     p,
				})
			return nil
		}

		if action == ErrorOnMultipleTries {
			return fmt.Errorf(
				"parameter %s has been set already, firstly at: %s",
				p.Name(), so.paramsSetAt[0].Desc())
		}

		if action == ExitOnMultipleTries {
			fmt.Fprintf(os.Stderr,
				"parameter %s has been set already, at %s. Aborting",
				p.Name(), so.paramsSetAt[0].Desc())
			os.Exit(1)
		}

		// if action == IgnoreMultipleTries
		return nil
	}
}

// SetOnceActionFunc constructs a local SetOnce value and then makes an
// ActionFunc to operate on it. The advantage is that the SetOnce value is
// guaranteed to be unique to the parameter, the disadvantage is that you
// don't have access to the value if you want to handle it further
func SetOnceActionFunc(action SetOnceErrAction) param.ActionFunc {
	so := &SetOnce{}
	return so.MakeActionFunc(action)
}
