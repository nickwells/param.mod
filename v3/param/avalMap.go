package param

import (
	"errors"
	"fmt"
	"sort"
)

// AValMap - this maps allowed values for an enumerated parameter to
// explanatory text. It forms part of the usage documentation of the
// program and will appear when the -help parameter is given by the
// user. It is also used to validate a supplied parameter.
type AValMap map[string]string

// String returns a string documenting the entries in the map - each entry is
// on a separate line
func (av AValMap) String() string {
	var avals string
	var valNames = make([]string, 0, len(av))
	var maxNameLen int

	for k := range av {
		valNames = append(valNames, k)
		if len(k) > maxNameLen {
			maxNameLen = len(k)
		}
	}
	sort.Strings(valNames)
	sep := ""
	for _, v := range valNames {
		avals += sep + fmt.Sprintf("%-*s: ", maxNameLen, v) + av[v]
		sep = "\n"
	}
	return avals
}

// OK returns a nil error if the map is "good" or an error with an
// explanation of the problem otherwise.
//
// A map is "good" if it has more than one entry. A set of allowed values
// with one or fewer entries is obviously a mistake: if no entries are valid
// then the parameter can never be set correctly and if it only has a single
// entry then the current (initial) value is the only allowed value and so
// there is no need for a parameter as no alternative can ever be allowed.
func (av AValMap) OK() error {
	minEntries := "It should have at least 2"
	switch len(av) {
	case 0:
		return errors.New("the map of allowed values has no entries. " +
			minEntries)
	case 1:
		return errors.New("the map of allowed values has only 1 entry. " +
			minEntries)
	default:
		return nil
	}
}

// NilAVM is a mixin type that can be embedded in a Setter to provide a
// default implementation for the AllowedValuesMap method of the interface
type NilAVM struct{}

// AllowedValuesMap returns a nil map of allowed values - unless your setter
// has a map of allowed values you should probably use this
func (_ NilAVM) AllowedValuesMap() AValMap {
	return nil
}

// AVM is a mixin type that can be embedded in a Setter to provide an allowed
// values map for the Setter. It provides an implementation for the
// AllowedValuesMap method of the interface

type AVM struct{ AllowedVals AValMap }

// AllowedValuesMap returns the map of allowed values. Use this if your setter
// has a map of allowed values.
func (a AVM) AllowedValuesMap() AValMap {
	return a.AllowedVals
}

// ValueAllowed returns true if the passed value is a key in the allowed
// values map
func (a AVM) ValueAllowed(val string) bool {
	_, ok := a.AllowedVals[val]
	return ok
}

// ValueMapOK returns any error that the OK method on the allowed values map
// returns
func (a AVM) ValueMapOK() error {
	return a.AllowedVals.OK()
}
