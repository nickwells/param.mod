package paction

import (
	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v6/param"
)

// Counter is used to record how many of a collection of params has
// been set
//
// It can be used to check, among other things, that at least one of a
// collection has been set or if more than one has been set
type Counter struct {
	ParamCount map[string]int

	ParamsSetAt param.Sources
}

// Count returns the count - the number of distinct parameters that have been
// set
func (c Counter) Count() int { return len(c.ParamCount) }

// Total returns the total number of times the ActionFunction was called
func (c Counter) Total() int {
	if c.ParamCount == nil {
		return 0
	}
	var tot int
	for _, v := range c.ParamCount {
		tot += v
	}
	return tot
}

// SetBy returns a string describing all the places that the parameter(s)
// were set
func (c Counter) SetBy() string { return c.ParamsSetAt.String() }

// MakeActionFunc returns a function suitable for passing to the PostAction
// method of a param.ByName object
func (c *Counter) MakeActionFunc() param.ActionFunc {
	return func(loc location.L, p *param.ByName, paramValues []string) error {
		if c.ParamCount == nil {
			c.ParamCount = make(map[string]int)
		}
		c.ParamCount[p.Name()]++

		c.ParamsSetAt = append(c.ParamsSetAt,
			param.Source{
				From:      loc.Source(),
				Loc:       loc,
				ParamVals: paramValues,
				Param:     p,
			})
		return nil
	}
}
