package phelp

import (
	"errors"
	"fmt"

	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

// configFileSetter sets a config file from which to read configuration
// parameters
type configFileSetter struct {
	psetter.ValueReqMandatory

	seenBefore map[string]bool
}

var configFileProvisos = filecheck.FileNonEmpty()

// SetWithVal (called when a value follows the parameter) does some minimal
// checking of the parameter - the processing of the file is all done in the
// post action function. It records whether the file has been seen before and
// reports it as an error if so.
func (s configFileSetter) SetWithVal(_ string, paramVal string) error {
	if paramVal == "" {
		return errors.New("no file name has been given")
	}

	if err := configFileProvisos.StatusCheck(paramVal); err != nil {
		return err
	}

	if s.seenBefore[paramVal] {
		return fmt.Errorf("the file name (%q) has been seen before", paramVal)
	}

	s.seenBefore[paramVal] = true
	return nil
}

// AllowedValues returns a string describing the allowed values
func (s configFileSetter) AllowedValues() string {
	return "a pathname to a file which must exist," +
		" containing configuration parameters"
}

// CurrentValue returns none
func (s configFileSetter) CurrentValue() string {
	return "none"
}

// CheckSetter checks that the seenBefore map has been initialised.
func (s *configFileSetter) CheckSetter(_ string) {
	if s.seenBefore == nil {
		s.seenBefore = make(map[string]bool)
	}
}

// ValDescribe returns the short name of the value expected
func (s configFileSetter) ValDescribe() string {
	return "filename"
}
