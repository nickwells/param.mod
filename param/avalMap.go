package param

import (
	"errors"
	"fmt"
	"sort"
)

// AllowedVals - this maps allowed values for an enumerated parameter to
// explanatory text. It forms part of the usage documentation of the program
// and will appear when the -help parameter is given by the user. It is also
// used to validate a supplied parameter.
type AllowedVals map[string]string

// String returns a string documenting the entries in the map - each entry is
// on a separate line
func (av AllowedVals) String() string {
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

// ValueMapOK returns a nil error if the map is "good" or an error with an
// explanation of the problem otherwise.
//
// A map is "good" if it has more than one entry. A set of allowed values
// with one or fewer entries is obviously a mistake: if no entries are valid
// then the parameter can never be set correctly and if it only has a single
// entry then the current (initial) value is the only allowed value and so
// there is no need for a parameter as no alternative can ever be allowed.
func (av AllowedVals) ValueMapOK() error {
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

// AllowedValuesMap returns a copy of the map of allowed values. This will be
// used by the standard help package to generate a list of allowed values.
func (av AllowedVals) AllowedValuesMap() AllowedVals {
	rval := make(map[string]string)
	for k, v := range av {
		rval[k] = v
	}
	return rval
}

// ValueAllowed returns true if the passed value is a key in the allowed
// values map
func (av AllowedVals) ValueAllowed(val string) bool {
	_, ok := av[val]
	return ok
}
