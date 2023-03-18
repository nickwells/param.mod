package psetter

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Aliases - this maps strings to lists of strings. It is expected that the
// keys are not in the set of allowed values and the entries in the
// associated value are allowed.
//
// It can be used as a mixin type that can be embedded in a Setter to provide
// alternative names for allowed values or to provide several names in one.
//
// It is recommended that you should use string constants for setting the
// aliases and the entries in the slice of values they correspond to. This
// will avoid possible errors.
//
// The advantages of const values are:
//
// - typos become compilation errors rather than silently failing.
//
// - the name of the constant value can distinguish between the string value
// and it's meaning as a semantic element representing a flag used to choose
// program behaviour.
//
// - the name that you give the const value can distinguish between identical
// strings and show which of various flags with the same string value you
// actually mean.
type Aliases map[string][]string

// AllowedValuesAliasMapper is the interface to be satisfied by a type having
// aliases
type AllowedValuesAliasMapper interface {
	AllowedValuesAliasMap() Aliases
}

// Keys returns an unsorted list of keys to the Aliases map and the
// length of the longest key.
func (a Aliases) Keys() ([]string, int) {
	keys := make([]string, 0, len(a))
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
		avals += sep + fmt.Sprintf("   %-*s: ", maxKeyLen, k) +
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
		pfx := fmt.Sprintf("Bad alias: %q: %q - ", ak, v)
		if len(v) == 0 {
			return errors.New(pfx + "it has an empty value")
		}

		if _, ok := av[ak]; ok {
			return errors.New(pfx + "an allowed value has the same name")
		}

		if ak == "" {
			return errors.New(pfx + "the alias name may not be blank")
		}
		if strings.ContainsRune(ak, '=') {
			return errors.New(pfx + "the alias name may not contain '=': ")
		}

		seenBefore := map[string]bool{}
		for _, avk := range v {
			if seenBefore[avk] {
				return fmt.Errorf("%s%q appears more than once", pfx, avk)
			}
			seenBefore[avk] = true
			if _, ok := av[avk]; !ok {
				return fmt.Errorf("%s%q is not an allowed value",
					pfx, avk)
			}
		}
	}
	return nil
}

// AllowedValuesAliasMap returns a copy of the map of aliases. This will be
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
