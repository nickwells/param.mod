package param

import "github.com/nickwells/location.mod/location"

// ActionFunc is the type of a function to be called when the ByName
// parameter is encountered.
//
// loc gives details on where the param was seen, this includes a description
// of the source (for instance "command line")
//
// p is the parameter which was matched.
//
// paramValues will have one or possibly two entries: the name used to match
// the param and (possibly) the value string.
type ActionFunc func(loc location.L, p *BaseParam, paramValues []string) error

// PostAction will return an option function which will add an action
// function to the list of functions to be called after the value has been
// set.
func PostAction(action ActionFunc) ByNameOptFunc {
	return func(p *ByName) error {
		(&p.BaseParam).postAction = append(p.postAction, action)
		return nil
	}
}

// ByPosPostAction will return an option function which will add an action
// function to the list of functions to be called after the value has been
// set.
func ByPosPostAction(action ActionFunc) ByPosOptFunc {
	return func(p *ByPos) error {
		(&p.BaseParam).postAction = append(p.postAction, action)
		return nil
	}
}
