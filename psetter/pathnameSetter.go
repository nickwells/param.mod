package psetter

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/english.mod/english"
	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/fileparse.mod/fileparse"
	"github.com/nickwells/strdist.mod/strdist"
)

// Pathname allows you to give a parameter that can be used to set a pathname
// value.
type Pathname struct {
	ValueReqMandatory

	// You must set a Value, the program will panic if not. This is the
	// pathname that the setter is setting.
	Value *string
	// Expectation allows you to set some file-specific checks.
	Expectation filecheck.Provisos
	// The Checks, if any, are applied to the supplied parameter value and
	// the new parameter will be applied only if they all return a nil error.
	Checks []check.String
	// ForceAbsolute, if set, causes any pathname value to be passed
	// to filepath.Abs before setting the value.
	ForceAbsolute bool
}

// CountChecks returns the number of check functions this setter has
func (s Pathname) CountChecks() int {
	return len(s.Checks)
}

// findAlternatives searches the directory, base, for entries similar to
// badName and then checks that the full path with the bad entry replaced
// satisfies the Expectation. If so it will add it to the list of
// alternatives.
func (s Pathname) findAlternatives(base, badName, tail string) string {
	f, err := os.Open(base)
	if err != nil {
		return fmt.Sprintf(", cannot open the directory for reading: %s", err)
	}
	defer f.Close()

	names, err := f.Readdirnames(1000)
	if err != nil && errors.Is(err, io.EOF) {
		return fmt.Sprintf(", cannot read the directory: %s", err)
	}
	alts := strdist.CaseBlindCosineFinder.FindNStrLike(3, badName, names...)
	var altStrs []string
	for _, alt := range alts {
		altStr := filepath.Join(base, alt, tail)
		if s.Expectation.StatusCheck(altStr) == nil {
			altStrs = append(altStrs, altStr)
		}
	}
	if len(altStrs) > 0 {
		return fmt.Sprintf(`, did you mean "%s"`,
			english.Join(altStrs, `", "`, `" or "`))
	}
	return ""
}

// expandError expands on the error, it tries to find the component of the
// path that does not exist.
func (s Pathname) expandError(pathname string, err error) error {
	suggestions := ""

	tailName := filepath.Base(pathname)
	pathname = filepath.Dir(pathname)
	tailPath := []string{}

	for tailName != pathname {
		info, err := s.Expectation.GetFileInfo(pathname)
		if err == nil {
			if !info.IsDir() {
				suggestions = fmt.Sprintf(
					"; %q exists but is not a directory", pathname)
				break
			}
			suggestions = fmt.Sprintf("; %q exists but %q does not",
				pathname, tailName)
			slices.Reverse[[]string, string](tailPath)
			tail := filepath.Join(tailPath...)
			suggestions += s.findAlternatives(pathname, tailName, tail)
			break
		}

		tailPath = append(tailPath, tailName)
		tailName = filepath.Base(pathname)
		pathname = filepath.Dir(pathname)
	}

	return fmt.Errorf("%w%s", err, suggestions)
}

// SetWithVal (called when a value follows the parameter) checks first that
// the value can be converted into a pathname (a tilda at the start of the
// path is converted to the appropriate home directory). Then it confirms
// that the file conforms to the supplied provisos. The Checks, if any, are
// run and if any check returns a non-nil error the Value is not updated and
// the error is returned. Only if the value is converted successfully, the
// Expectations are all met and no checks fail is the Value set and a nil
// error is returned.
func (s Pathname) SetWithVal(_ string, paramVal string) error {
	pathname, err := fileparse.FixFileName(paramVal)
	if err != nil {
		return err
	}

	if s.ForceAbsolute {
		pathname, err = filepath.Abs(pathname)
		if err != nil {
			return err
		}
	}

	err = s.Expectation.StatusCheck(pathname)
	if err != nil {
		if errors.Is(err, filecheck.ErrShouldExistButDoesNot) {
			err = s.expandError(pathname, err)
		}
		return err
	}

	for _, check := range s.Checks {
		err := check(pathname)
		if err != nil {
			return err
		}
	}

	*s.Value = pathname
	return nil
}

// AllowedValues returns a string describing the allowed values
func (s Pathname) AllowedValues() string {
	rval := "a pathname" + HasChecks(s)

	extras := s.Expectation.String()
	if extras != "" {
		rval += ". " + extras
	}

	return rval
}

// CurrentValue returns the current setting of the parameter value
func (s Pathname) CurrentValue() string {
	return fmt.Sprintf("%v", *s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or if it has nil Checks.
func (s Pathname) CheckSetter(name string) {
	// Check the value is not nil
	if s.Value == nil {
		panic(NilValueMessage(name, fmt.Sprintf("%T", s)))
	}

	// Check there are no nil Check funcs
	for i, check := range s.Checks {
		if check == nil {
			panic(NilCheckMessage(name, fmt.Sprintf("%T", s), i))
		}
	}
}
