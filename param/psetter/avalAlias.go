package psetter

import (
	"fmt"
	"sort"
	"strings"
)

// Aliases - this maps strings to lists of strings. It is expected that the
// keys are not in the set of allowed values and the entries in the
// associated value are allowed.
//
// It can be used as a mixin type that can be embedded in a Setter to provide
// alternative names for allowed values or to provide several names in one
type Aliases map[string][]string

// AllowedValuesAliasMapper is the interface to be satisfied by a type having
// aliases
type AllowedValuesAliasMapper interface {
	AllowedValuesAliasMap() Aliases
}

// Keys returns an unsorted list of keys to the Aliases map and the
// length of the longest key.
func (a Aliases) Keys() ([]string, int) {
	var keys = make([]string, 0, len(a))
	var maxKeyLen int

	for k := range a {
		keys = append(keys, k)
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}

	return keys, maxKeyLen
}

// String returns a string documenting the entries in the map - each entry is
// on a separate line
func (a Aliases) String() string {
	if a == nil {
		return ""
	}
	var avals string
	keys, maxKeyLen := a.Keys()
	sort.Strings(keys)

	sep := ""
	for _, k := range keys {
		avals += sep + fmt.Sprintf("%-*s: ", maxKeyLen, k) +
			strings.Join(a[k], ", ")
		sep = "\n"
	}
	return avals
}

// Check returns a nil error if the map is "good" or an error with an
// explanation of the problem otherwise.
//
// A map is "good" if each key does not exist in the AllowedVals but each
// entry in the associated list is in the AllowedVals. Also, an empty alias
// is not allowed.
func (a Aliases) Check(av AllowedVals) error {
	for ak, v := range a {
		pfx := fmt.Sprintf("Alias %q is invalid", ak)
		if len(v) == 0 {
			return fmt.Errorf("%s - it has an empty value", pfx)
		}

		if _, ok := av[ak]; ok {
			return fmt.Errorf("%s - there is an allowed value of the same name",
				pfx)
		}

		seenBefore := map[string]bool{}
		for _, avk := range v {
			if seenBefore[avk] {
				return fmt.Errorf("%s - %q is in the list twice", pfx, avk)
			}
			seenBefore[avk] = true
			if _, ok := av[avk]; !ok {
				return fmt.Errorf("%s - %q is not in the allowed values",
					pfx, avk)
			}
		}
	}
	return nil
}

// AllowedValuesMap returns a copy of the map of allowed values. This will be
// used by the standard help package to generate a list of allowed values.
func (a Aliases) AllowedValuesAliasMap() Aliases {
	rval := make(map[string][]string)
	for k, v := range a {
		rval[k] = v
	}
	return rval
}

// IsAnAlias returns true if the passed value is a key in the aliases map
func (a Aliases) IsAnAlias(val string) bool {
	_, ok := a[val]
	return ok
}

// AliasVal returns a copy of the value of the alias
func (a Aliases) AliasVal(name string) []string {
	rval := make([]string, len(a[name]))
	copy(rval, a[name])
	return rval
}
