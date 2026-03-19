package phelputils

import (
	"strings"

	"github.com/nickwells/param.mod/v7/param"
	"github.com/nickwells/param.mod/v7/ptypes"
)

// AllowedValueParts returns separate blocks of text as are needed to
// describe the allowed value. If the value has already been shown for
// another parameter (as recorded in the cache) then a reference to that is
// shown instead if the resultant text is longer than a threshold value and
// the resultant reference text is shorter than the full value would be.
func AllowedValueParts(
	cache ptypes.AValCache,
	pName string,
	s param.Setter,
) []string {
	const allowedValueThresholdLength = 50

	parts := []string{s.AllowedValues()}

	if sAVM, ok := s.(ptypes.AllowedValuesMapper); ok {
		desc := MakeAllowedValueDesc("value", sAVM.AllowedValuesMap())
		if desc != "" {
			parts = append(parts, desc)
		}
	}

	if sAVAM, ok := s.(ptypes.AllowedValuesAliasMapper); ok {
		desc := MakeAliasDesc("", sAVAM.AllowedValuesAliasMap())
		if desc != "" {
			parts = append(parts, desc)
		}
	}

	keyStr := strings.Join(parts, "")
	if len(keyStr) > allowedValueThresholdLength {
		name, alreadyShown := cache[keyStr]
		if alreadyShown {
			ref := "(see parameter: " + name + ")"

			if len(ref) < len(keyStr) {
				return []string{ref}
			}
		}

		cache[keyStr] = pName
	}

	return parts
}

// MakeAllowedValueDesc returns a string describing the allowed values
func MakeAllowedValueDesc(name string, avm ptypes.AllowedVals[string]) string {
	if len(avm) == 0 {
		return ""
	}

	desc := "The " + name + " must be"

	if len(avm) > 1 {
		desc += " one of the following"
	}

	desc += ":\n" + avm.String()

	return desc
}

// MakeAliasDesc returns a string describing the aliases
func MakeAliasDesc(name string, am ptypes.Aliases[string]) string {
	if len(am) == 0 {
		return ""
	}

	desc := "The following"
	if name != "" {
		desc += " " + name
	}

	if len(am) == 1 {
		desc += " alias is"
	} else {
		desc += " aliases are"
	}

	desc += " available:\n" + am.String()

	return desc
}
