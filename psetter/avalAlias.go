package psetter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nickwells/english.mod/english"
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
type Aliases[T ~string] map[T][]T

// AllowedValuesAliasMapper is the interface to be satisfied by a type having
// aliases
type AllowedValuesAliasMapper interface {
	AllowedValuesAliasMap() Aliases[string]
}

// Keys returns an unsorted list of keys to the Aliases map and the
// length of the longest key.
func (a Aliases[T]) Keys() ([]string, int) {
	keys := make([]string, 0, len(a))

	var maxKeyLen int

	for k := range a {
		keys = append(keys, string(k))

		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}

	return keys, maxKeyLen
}

// String returns a string documenting the entries in the map - each entry is
// on a separate line
func (a Aliases[T]) String() string {
	if a == nil {
		return ""
	}

	var avals string

	keys, maxKeyLen := a.Keys()

	sort.Strings(keys)

	sep := ""

	for _, k := range keys {
		kav := convertToStringSlice(a[T(k)])
		avals += sep + fmt.Sprintf("   %-*s: ", maxKeyLen, k) +
			strings.Join(kav, ", ")
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
func (a Aliases[T]) Check(av AllowedVals[T]) error {
	aliasKeys := []T{}
	for ak := range a {
		aliasKeys = append(aliasKeys, ak)
	}

	sort.Slice(aliasKeys,
		func(i, j int) bool { return aliasKeys[i] < aliasKeys[j] })

	allProblems := []string{}

	for _, name := range aliasKeys {
		aliasProblems := a.aliasNameProblems(name, av)
		aliasProblems = append(aliasProblems, a.aliasValueProblems(name, av)...)

		if len(aliasProblems) > 0 {
			sep := " - "

			if len(aliasProblems) > 1 {
				sep = "\n    - "
			}

			sort.Strings(aliasProblems)
			allProblems = append(allProblems,
				fmt.Sprintf("%q: %#v%s%s",
					name, a[name],
					sep,
					strings.Join(aliasProblems, sep)))
		}
	}

	if len(allProblems) > 0 {
		sep := " "
		if len(allProblems) > 1 {
			sep = fmt.Sprintf(" (%d)\n", len(allProblems))
		}

		return fmt.Errorf("bad %s:%s%s",
			english.Plural("alias", len(allProblems)),
			sep,
			strings.Join(allProblems, "\n"))
	}

	return nil
}

// aliasValueProblems checks the alias value for validity. It must
//   - not be empty
//   - not contain  duplicates
//   - not contain invalid values
//
// It returns all the problems found as a slice of strings
func (a Aliases[T]) aliasValueProblems(name T, av AllowedVals[T]) []string {
	if len(a[name]) == 0 {
		return []string{"the alias maps to no values"}
	}

	indexes := map[T][]int{}
	badValues := map[T][]int{}

	for i, avk := range a[name] {
		indexes[avk] = append(indexes[avk], i)

		if _, ok := av[avk]; !ok {
			badValues[avk] = append(badValues[avk], i)
		}
	}

	valueProblems := a.reportBadAliases(badValues)
	valueProblems = append(valueProblems, a.reportDuplicateVals(indexes)...)

	return valueProblems
}

// aliasNameProblems checks the alias name for validity. It must
//   - not be in the set of allowed values
//   - not be blank
//   - not contain  '='
//
// It returns all the problems found as a slice of strings
func (a Aliases[T]) aliasNameProblems(name T, av AllowedVals[T]) []string {
	if name == "" {
		return []string{"the alias name must not be blank"}
	}

	if strings.ContainsRune(string(name), '=') {
		return []string{"the alias name must not contain '='"}
	}

	if _, ok := av[name]; ok {
		return []string{"an allowed value has the same name as the alias"}
	}

	return nil
}

// reportBadAliases generates a string listing all the invalid alias values.
func (a Aliases[T]) reportBadAliases(badVals map[T][]int) []string {
	bvKeys := []T{}
	for k := range badVals {
		bvKeys = append(bvKeys, k)
	}

	sort.Slice(bvKeys,
		func(i, j int) bool { return bvKeys[i] < bvKeys[j] })

	problems := []string{}

	for _, k := range bvKeys {
		iVals := []string{}
		for _, i := range badVals[k] {
			iVals = append(iVals, fmt.Sprintf("%d", i))
		}

		problems = append(problems,
			fmt.Sprintf("%q (at index %s) is unknown",
				k, english.Join(iVals, ", ", " and ")))
	}

	return problems
}

// reportDuplicateVals generates a string listing all the duplicate alias
// values.
func (a Aliases[T]) reportDuplicateVals(indexes map[T][]int) []string {
	iKeys := []T{}
	for k := range indexes {
		iKeys = append(iKeys, k)
	}

	sort.Slice(iKeys,
		func(i, j int) bool { return iKeys[i] < iKeys[j] })

	problems := []string{}

	for _, k := range iKeys {
		if len(indexes[k]) > 1 {
			iVals := []string{}

			for _, i := range indexes[k] {
				iVals = append(iVals, fmt.Sprintf("%d", i))
			}

			problems = append(problems,
				fmt.Sprintf("%q appears more than once (at index %s)",
					k, english.Join(iVals, ", ", " and ")))
		}
	}

	return problems
}

// AllowedValuesAliasMap returns a copy of the map of aliases. This will be
// used by the standard help package to generate a list of allowed values.
func (a Aliases[T]) AllowedValuesAliasMap() Aliases[string] {
	rval := make(map[string][]string)

	for k, v := range a {
		strVals := make([]string, 0, len(v))
		for _, tv := range v {
			strVals = append(strVals, string(tv))
		}

		rval[string(k)] = strVals
	}

	return rval
}

// IsAnAlias returns true if the passed value is a key in the aliases map
func (a Aliases[T]) IsAnAlias(val string) bool {
	_, ok := a[T(val)]
	return ok
}

// AliasVal returns a copy of the value of the alias
func (a Aliases[T]) AliasVal(name T) []T {
	rval := make([]T, len(a[name]))
	copy(rval, a[name])

	return rval
}

// convertToStringSlice returns a copy of the passed slice with the values
// converted to string
func convertToStringSlice[T ~string](ts []T) []string {
	ss := make([]string, 0, len(ts))
	for _, v := range ts {
		ss = append(ss, string(v))
	}

	return ss
}
