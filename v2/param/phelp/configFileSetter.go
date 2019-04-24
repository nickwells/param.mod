package phelp

import (
	"errors"

	"github.com/nickwells/param.mod/v2/param"
)

// configFileSetter sets a config file from which to read configuration
// parameters
type configFileSetter struct {
	param.ValueReqMandatory

	seenBefore map[string]bool
}

// SetWithVal (called when a value follows the parameter) does some minimal
// checking of the parameter - the processing of the file is all done in the
// post action function. It records whether the file has been seen before and
// reports it as an error if so.
func (s configFileSetter) SetWithVal(_ string, paramVal string) error {
	if paramVal == "" {
		return errors.New("no config file name has been given")
	}

	if s.seenBefore[paramVal] {
		return errors.New("the config file name has been seen before")
	}

	s.seenBefore[paramVal] = true
	return nil
}

// AllowedValues returns a string describing the allowed values
func (s configFileSetter) AllowedValues() string {
	return "a pathname to a file which must exist," +
		" containing config parameter settings"
}

// CurrentValue returns none
func (s configFileSetter) CurrentValue() string {
	return "none"
}

// CheckSetter checks that the seenBefore map has been initialised.
func (s configFileSetter) CheckSetter(name string) {
	if s.seenBefore == nil {
		panic(name + ": phelp.configFileSetter Check failed:" +
			" the map of previously seen config files has not been set")
	}
}
