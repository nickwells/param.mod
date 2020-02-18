package param

import (
	"errors"
	"fmt"
	"regexp"
)

// =============================================

var nameCheckRE *regexp.Regexp

// =============================================

func init() {
	nameCheckRE = regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]*$")
}

// nameCheck returns an error if the name is invalid or if it has already
// been used.
func (ps *PSet) nameCheck(name, whereAdded string) error {
	if !nameCheckRE.MatchString(name) {
		return fmt.Errorf(
			"the parameter name %q is invalid."+
				" It must match: %q.\nFix this at: %s",
			name, nameCheckRE.String(), whereAdded)
	}

	altP, exists := ps.nameToParam[name]
	if exists {
		errDesc := fmt.Sprintf("parameter name %q has already been used", name)
		if altP.altNames[0] != name {
			errDesc += fmt.Sprintf(" as an alternative to %q", altP.altNames[0])
		}

		errDesc += fmt.Sprintf(" (a member of parameter group %q)",
			altP.groupName)
		errDesc += "\n  this param added at: " + whereAdded
		errDesc += "\n  originally added at: " + altP.whereAdded

		return errors.New(errDesc)
	}

	return nil
}
