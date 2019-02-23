package psetter

import (
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
	if len(av) > 1 {
		return nil
	}
	return fmt.Errorf(
		"the allowed values map has %d entries. It should have more than 1",
		len(av))
}
