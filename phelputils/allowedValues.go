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
		var part string

		avm := sAVM.AllowedValuesMap()
		if len(avm) == 1 {
			part = "The value must be:\n" + avm.String()
		} else if len(avm) > 1 {
			part = "The value must be one of the following:\n" +
				avm.String()
		}

		parts = append(parts, part)
	}

	if sAVAM, ok := s.(ptypes.AllowedValuesAliasMapper); ok {
		var part string

		avam := sAVAM.AllowedValuesAliasMap()
		if len(avam) == 1 {
			part = "The following alias is available:\n" + avam.String()
		} else if len(avam) > 1 {
			part = "The following aliases are available:\n" + avam.String()
		}

		parts = append(parts, part)
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
