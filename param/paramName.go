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

// ParameterNameCheck checks that the given parameter name is valid and returns
// an error if not. A parameter name must start with a letter and be followed
// by zero or more letters, digits or dashes
func ParameterNameCheck(name string) error {
	if !nameCheckRE.MatchString(name) {
		return fmt.Errorf(
			"the parameter name %q is invalid."+
				" It must start with a letter and be followed by"+
				" zero or more letters, digits or dashes",
			name)
	}

	return nil
}

// nameCheck returns an error if the name is invalid or if it has already
// been used.
func (ps *PSet) nameCheck(name, whereAdded string) error {
	if err := ParameterNameCheck(name); err != nil {
		return fmt.Errorf("bad parameter name added at: %s. %w",
			whereAdded, err)
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
