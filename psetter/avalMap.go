package psetter

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// AllowedVals - this maps allowed values for an enumerated parameter to
// explanatory text. It forms part of the usage documentation of the
// program. It is also used to validate a supplied parameter.
//
// It can be used as a mixin type that can be embedded in a Setter to provide
// a restricted set of allowed values.
//
// It is recommended that you should use string constants for setting the
// map entries.
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
type AllowedVals[T ~string] map[T]string

// AllowedValuesMapper is the interface to be satisfied by a type having a
// map of allowed values.
type AllowedValuesMapper interface {
	AllowedValuesMap() AllowedVals[string]
}

// Keys returns an unsorted list of keys to the AllowedVals map and the
// length of the longest key.
func (av AllowedVals[T]) Keys() ([]string, int) {
	keys := make([]string, 0, len(av))
	var maxKeyLen int

	for k := range av {
		keys = append(keys, string(k))
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}

	return keys, maxKeyLen
}

// String returns a string documenting the entries in the map - each entry is
// on a separate line
func (av AllowedVals[T]) String() string {
	if av == nil {
		return ""
	}
	var avals string
	keys, maxKeyLen := av.Keys()
	sort.Strings(keys)

	sep := ""
	for _, k := range keys {
		avals += sep + fmt.Sprintf("   %-*s: ", maxKeyLen, k) + av[T(k)]
		sep = "\n"
	}
	return avals
}

// Check returns a nil error if the map is "good" or an error with an
// explanation of the problem otherwise.
//
// A map is "good" if it has more than one entry. A set of allowed values
// with one or fewer entries is obviously a mistake: if no entries are valid
// then the parameter can never be set correctly and if it only has a single
// entry then the current (initial) value is the only allowed value and so
// there is no need for a parameter as no alternative can ever be allowed.
func (av AllowedVals[T]) Check() error {
	minEntries := "It should have at least 2"
	switch len(av) {
	case 0:
		return errors.New("the map of allowed values has no entries. " +
			minEntries)
	case 1:
		return errors.New("the map of allowed values has only 1 entry. " +
			minEntries)
	}

	for k, v := range av {
		pfx := fmt.Sprintf("Bad allowed value: %q: %q - ", k, v)
		if k == "" {
			return errors.New(pfx + "the allowed value may not be blank")
		}
		if strings.ContainsRune(string(k), '=') {
			return errors.New(pfx + "the allowed value may not contain '=': ")
		}
	}
	return nil
}

// AllowedValuesMap returns a copy of the map of allowed values. This will be
// used by the standard help package to generate a list of allowed values.
func (av AllowedVals[T]) AllowedValuesMap() AllowedVals[string] {
	rval := make(map[string]string)
	for k, v := range av {
		rval[string(k)] = v
	}
	return rval
}

// ValueAllowed returns true if the passed value is a key in the allowed
// values map
func (av AllowedVals[T]) ValueAllowed(val string) bool {
	_, ok := av[T(val)]
	return ok
}
