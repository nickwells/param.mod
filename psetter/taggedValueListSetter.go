package psetter

import (
	"fmt"
	"strings"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v7/phelputils"
	"github.com/nickwells/param.mod/v7/ptypes"
	"github.com/nickwells/twrap.mod/twrap"
)

const taggedValueVTSep = "="

// TaggedValue represents a tagged, enumerated value. This is what is
// populated by the [TaggedValueList] setter. In your code you can define a
// type alias for the specific type (with the E and T type parameters filled
// in). This will make your code easier to read and should give an indication
// as to the purpose of the resulting values.
type TaggedValue[E, T ~string] struct {
	Value E
	Tags  []T
}

// String returns a formatted representation of the TaggedValue value
func (tv TaggedValue[E, T]) String(tagSep string) string {
	sb := strings.Builder{}
	sb.WriteString(string(tv.Value))

	sep := taggedValueVTSep

	for _, t := range tv.Tags {
		sb.WriteString(sep + string(t))
		sep = tagSep
	}

	return sb.String()
}

// TaggedValueList sets the values in a slice of TaggedValues. The values
// must be in the allowed values map. This is similar to the [EnumList] type
// except that each entry can have associated with it, a list of tags.
//
// It is recommended that you should use string constants for setting the
// list entries and for initialising the allowed values map to avoid possible
// errors.
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
type TaggedValueList[E, T ~string] struct {
	ValueReqMandatory
	// The AllowedVals must be set, the program will panic if not. These are
	// the only values that will be allowed in the slice of strings.
	ptypes.AllowedVals[E]

	// The TagAllowedVals must be set, the program will panic if not. These are
	// the only values that will be allowed in the slice of tags.
	TagAllowedVals ptypes.AllowedVals[T]

	// The Aliases need not be given but if they are then each alias must not
	// be in AllowedVals and all of the resulting values must be in
	// AllowedVals.
	ptypes.Aliases[E]

	// The TagAliases need not be given but if they are then each alias must not
	// be in TagAllowedVals and all of the resulting values must be in
	// TagAllowedVals.
	TagAliases ptypes.Aliases[T]

	// Value must be set, the program will panic if not. This is the slice of
	// values that this setter is setting.
	Value *[]TaggedValue[E, T]

	// The StrListSeparator allows you to override the default separator
	// between list elements.
	StrListSeparator

	// The TagListSeparator allows you to override the default separator
	// between tag list elements. This must be set to a different value to
	// the StrListSeparator, the program wil panic if not.
	TagListSeparator StrListSeparator

	// The Checks, if any, are applied to the list of new values and the
	// Value will only be updated if they all return a nil error.
	Checks []check.ValCk[[]TaggedValue[E, T]]

	// The TagChecks, if any, are applied to the list tags of each value and
	// the Value will only be updated if they all return a nil error.
	TagChecks []check.ValCk[[]T]
}

// CountChecks returns the number of check functions this setter has
func (s TaggedValueList[E, T]) CountChecks() int {
	return len(s.Checks)
}

// mkTagList makes a slice of tags from the tag-string, ts. If any of the
// parsed tags is invalid or if the resulting slice is invalid, an error is
// returned.
func (s TaggedValueList[E, T]) mkTagList(ts string) ([]T, error) {
	tags := []T{}
	tagListSep := s.TagListSeparator.GetSeparator()

	for t := range strings.SplitSeq(ts, tagListSep) {
		if t == "" {
			continue
		}

		if s.TagAllowedVals.ValueAllowed(t) {
			tags = append(tags, T(t))
			continue
		}

		if !s.TagAliases.IsAnAlias(t) {
			return []T{}, fmt.Errorf("tag value is not allowed: %q", t)
		}

		for _, av := range s.TagAliases.AliasVal(T(t)) {
			tags = append(tags, T(av))
		}
	}

	for _, tCheck := range s.TagChecks {
		err := tCheck(tags)
		if err != nil {
			return []T{}, err
		}
	}

	return tags, nil
}

// SetWithVal (called when a value follows the parameter) splits the value
// using the list separator. It then checks all the values for validity and
// only if all the values are in the allowed values list does it add them
// to the slice of strings pointed to by the Value. It returns a error for
// the first invalid value or if a check is breached.
func (s TaggedValueList[E, T]) SetWithVal(_ string, paramVal string) error {
	teVals := []TaggedValue[E, T]{}

	sep := s.GetSeparator()

	vals := strings.SplitSeq(paramVal, sep)
	for v := range vals {
		ev, tagVals, _ := strings.Cut(v, taggedValueVTSep)

		tags, err := s.mkTagList(tagVals)
		if err != nil {
			return err
		}

		if s.ValueAllowed(ev) {
			tv := TaggedValue[E, T]{Value: E(ev), Tags: tags}
			teVals = append(teVals, tv)

			continue
		}

		if !s.IsAnAlias(ev) {
			return fmt.Errorf("bad value: %q", ev)
		}

		for _, av := range s.AliasVal(E(ev)) {
			tv := TaggedValue[E, T]{Value: E(av), Tags: tags}
			teVals = append(teVals, tv)
		}

		tv := TaggedValue[E, T]{Value: E(ev), Tags: tags}
		teVals = append(teVals, tv)
	}

	for _, check := range s.Checks {
		err := check(teVals)
		if err != nil {
			return err
		}
	}

	*s.Value = teVals

	return nil
}

// AllowedValues returns a string listing the allowed values
func (s TaggedValueList[E, T]) AllowedValues() string {
	return s.ListValDesc("name"+taggedValueVTSep+"tags") +
		" where 'tags' is " + s.TagListSeparator.ListValDesc("tag-values") +
		" and the name is one of the allowed values" +
		HasChecks(s) + "."
}

// CurrentValue returns the current setting of the parameter value
func (s TaggedValueList[E, T]) CurrentValue() string {
	var str strings.Builder

	sep := ""
	tagSep := s.TagListSeparator.GetSeparator()

	for _, v := range *s.Value {
		str.WriteString(sep)
		str.WriteString(v.String(tagSep))

		sep = s.GetSeparator()
	}

	return str.String()
}

// checkChecks checks the Checks and the TagChecks functions for nil
// functions.
func (s TaggedValueList[E, T]) checkChecks(name string) {
	// Check there are no nil Check funcs
	for i, check := range s.Checks {
		if check == nil {
			panic(NilCheckMessage(name, fmt.Sprintf("%T.Checks", s), i))
		}
	}

	// Check there are no nil TagCheck funcs
	for i, check := range s.TagChecks {
		if check == nil {
			panic(NilCheckMessage(name, fmt.Sprintf("%T.TagChecks", s), i))
		}
	}
}

// checkAllowedVals checks the AllowedVals and Aliases for both the embedded
// values and the Tag-specific ones.
func (s TaggedValueList[E, T]) checkAllowedVals(name string) {
	// Check that the AllowedVals map is well formed
	if err := s.AllowedVals.Check(); err != nil {
		panic(BadSetterMessage(name,
			fmt.Sprintf("%T.AllowedVals", s),
			err.Error()))
	}

	// Check the alias values
	if err := s.Aliases.Check(s.AllowedVals); err != nil {
		panic(BadSetterMessage(name,
			fmt.Sprintf("%T.Aliases", s),
			err.Error()))
	}

	// Check that the TagAllowedVals map is well formed
	if err := s.TagAllowedVals.Check(); err != nil {
		panic(BadSetterMessage(name,
			fmt.Sprintf("%T.TagAllowedVals", s),
			"tag: "+err.Error()))
	}

	// Check the tag alias values
	if err := s.TagAliases.Check(s.TagAllowedVals); err != nil {
		panic(BadSetterMessage(name,
			fmt.Sprintf("%T.TagAliases", s),
			"tag: "+err.Error()))
	}
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or there are no allowed values or the initial value is not
// allowed.
func (s TaggedValueList[E, T]) CheckSetter(name string) {
	// Check the value is not nil
	if s.Value == nil {
		panic(NilValueMessage(name, fmt.Sprintf("%T", s)))
	}

	// Check that the AllowedValsand Aliases are well formed
	s.checkAllowedVals(name)

	// Check that the current values are all allowed
	for i, tev := range *s.Value {
		if _, ok := s.AllowedVals[tev.Value]; !ok {
			panic(
				BadValueMessage(
					name,
					fmt.Sprintf("%T", s),
					fmt.Sprintf(
						"element %d"+
							" in the list of initial values is invalid:"+
							" bad value: %q",
						i, tev.Value)))
		}

		for j, t := range tev.Tags {
			if _, ok := s.TagAllowedVals[t]; !ok {
				panic(
					BadValueMessage(
						name,
						fmt.Sprintf("%T", s),
						fmt.Sprintf(
							"element %d"+
								" in the list of initial values is invalid:"+
								" bad tag: %d: %q",
							i, j, t)))
			}
		}
	}

	s.checkChecks(name)

	sep := s.GetSeparator()
	tagSep := s.TagListSeparator.GetSeparator()

	if sep == tagSep {
		listSep := fmt.Sprintf("the list separator (%q)", sep)
		tagListSep := fmt.Sprintf("the tag list separator (%q)", tagSep)
		panic(fmt.Errorf("%s must differ from %s", listSep, tagListSep))
	}
}

// ValDescribe returns a string describing the allowed values
func (s TaggedValueList[E, T]) ValDescribe() string {
	return fmt.Sprintf("name%stag%stag...%sname%stag...",
		taggedValueVTSep, s.TagListSeparator.GetSeparator(),
		s.GetSeparator(), taggedValueVTSep)
}

// ExtraHelp provides additional help text showing the allowed tag values.
func (s TaggedValueList[E, T]) ExtraHelp(
	twc *twrap.TWConf,
	indent, extraIndent int,
) {
	prefix := "Allowed tag values: "

	parts := []string{
		phelputils.MakeAllowedValueDesc("tag",
			s.TagAllowedVals.AllowedValuesMap()),
	}

	aliases := phelputils.MakeAliasDesc("tag",
		s.TagAliases.AllowedValuesAliasMap())
	if aliases != "" {
		parts = append(parts, aliases)
	}

	twc.WrapPrefixed(prefix, parts[0], indent)

	indent2 := indent + len(prefix)
	indent3 := indent2 + extraIndent

	for _, p := range parts[1:] {
		twc.Wrap2Indent(p, indent2, indent3)
	}
}
