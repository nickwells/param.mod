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
func (ps *ParamSet) nameCheck(name string) error {
	if !nameCheckRE.MatchString(name) {
		return fmt.Errorf(
			"the parameter name '%s' is invalid. It must match: '%s'",
			name, nameCheckRE.String())
	}

	altP, exists := ps.nameToParam[name]
	if exists {
		errDesc := fmt.Sprintf("parameter name '%s' has already been used",
			name)
		if altP.altNames[0] != name {
			errDesc += fmt.Sprintf(" as an alternative to '%s'",
				altP.altNames[0])
		}
		if altP.groupName != "" {
			errDesc += fmt.Sprintf(" (a member of parameter group %s)",
				altP.groupName)
		}
		return errors.New(errDesc)
	}
	return nil
}
