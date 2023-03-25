package param_test

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/paramset"
)

// =======================================================

// failingSetter will return an error in all cases - it's intended for use in
// test cases where we want a setter to fail
type failingSetter struct {
	errMsg string
}

// ValueReq returns the Mandatory value of the ValueReq type to indicate that
// a value must follow the parameter for this setter
func (failingSetter) ValueReq() param.ValueReq { return param.Mandatory }

// Set returns an error because if the value is Mandatory then a value must
// follow the parameter for this setter
func (failingSetter) Set(name string) error {
	return fmt.Errorf("a value must follow this parameter: %q,"+
		" either following an '=' or as a next parameter", name)
}

// SetWithVal (called when a value follows the parameter) always returns an
// error.
func (s failingSetter) SetWithVal(_ string, _ string) error {
	return errors.New("failingSetter: " + s.errMsg)
}

// AllowedValues returns a string describing the allowed values
func (s failingSetter) AllowedValues() string {
	return "none - all values cause an error"
}

// CurrentValue returns the current setting of the parameter value
func (s failingSetter) CurrentValue() string {
	return "none (failingSetter)"
}

// CheckSetter checks that the setter has been properly created
func (s failingSetter) CheckSetter(_ string) {
}

// =======================================================

type namedParamInitialiser struct {
	name   string
	setter param.Setter
	desc   string
	opts   []param.OptFunc
}

// compare compares the parameter with the initialiser values and reports an
// error if it differs
func (npi namedParamInitialiser) compare(
	t *testing.T, tID string, p *param.ByName, sbs ShouldBeSetType,
) {
	t.Helper()

	if p == nil {
		if sbs == ShouldBeSet {
			t.Log(tID)
			t.Errorf("\t: the param should be set but it is nil\n")
		}
		return
	}

	if p.Name() != npi.name {
		t.Log(tID)
		t.Errorf("\t: the name did not match: '%s' != '%s'\n",
			p.Name(), npi.name)
		return
	}
	if p.Description() != npi.desc {
		t.Log(tID)
		t.Errorf("\t: the description did not match: '%s' != '%s'\n",
			p.Description(), npi.desc)
		return
	}
	if p.HasBeenSet() {
		if sbs == ShouldNotBeSet {
			t.Log(tID)
			t.Errorf("\t: param has been set but should not be\n")
		}
	} else if sbs == ShouldBeSet {
		t.Log(tID)
		t.Errorf("\t: param has not been set but should be\n")
	}
}

type posParamInitialiser struct {
	setter param.Setter
	name   string
	desc   string
	opts   []param.PosOptFunc
}

// paramInitialisers holds a pointer to either a namedParamInitialiser or a
// posParamInitialiser pointer (either, neither or both could be nil)
type paramInitialisers struct {
	npi    *namedParamInitialiser
	npiSBS ShouldBeSetType

	ppi    *posParamInitialiser
	ppiSBS ShouldBeSetType
}

// panicSafeTestAddByPos adds a ByPos parameter to a parameter set and
// catches any panics. Then it returns the added parameter (if any), a
// boolean indicating if a panic occurred and the value recovered from the
// panic
func panicSafeTestAddByPos(ps *param.PSet, ppi *posParamInitialiser,
) (pp *param.ByPos, panicked bool, panicVal any,
) {
	if ppi == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()
	pp = ps.AddByPos(ppi.name, ppi.setter, ppi.desc, ppi.opts...)
	return pp, panicked, panicVal
}

// panicSafeTestAddByName adds a ByName parameter to a parameter set and
// catches any panics. Then it returns the added parameter (if any), a
// boolean indicating if a panic occurred and the value recovered from the
// panic
func panicSafeTestAddByName(ps *param.PSet, npi *namedParamInitialiser,
) (p *param.ByName, panicked bool, panicVal any,
) {
	if npi == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()
	p = ps.Add(npi.name, npi.setter, npi.desc, npi.opts...)
	return p, panicked, panicVal
}

// panicSafeTestParse parses the supplied parameters and catches any
// panics. Then it returns a map of the errors (if any), a boolean indicating
// if a panic occurred and the value recovered from the panic
func panicSafeTestParse(ps *param.PSet, params []string,
) (errMap param.ErrMap, panicked bool, panicVal any, stackTrace []byte,
) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
			stackTrace = debug.Stack()
		}
	}()
	stackTrace = []byte{}
	errMap = ps.Parse(params)
	return errMap, panicked, panicVal, stackTrace
}

// panicSafeSetGroupDescription sets the group description and catches any
// panics. Then it returns a boolean indicating if a panic occurred and the
// value recovered from the panic
func panicSafeSetGroupDescription(ps *param.PSet, groupName, desc string,
) (panicked bool, panicVal any, stackTrace []byte,
) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
			stackTrace = debug.Stack()
		}
	}()
	ps.SetGroupDescription(groupName, desc)
	return panicked, panicVal, stackTrace
}

// =======================================================

// logErrs prints the errors out one per line
func logErrs(t *testing.T, msg string, errs []error) {
	t.Helper()

	t.Log(msg)
	for _, err := range errs {
		t.Log("\t:\t", err, "\n")
	}
}

// logErrMap reports the contents of the error map returned by param.Parse(...)
func logErrMap(t *testing.T, errMap param.ErrMap) bool {
	if len(errMap) == 0 {
		return false
	}

	t.Helper()

	for k, errs := range errMap {
		logErrs(t, "\t: Errors for: "+k+":\n", errs)
	}
	return true
}

// logName logs the name if it hasn't already been logged and returns true to
// set the nameLogged flag
func logName(t *testing.T, nameLogged bool, name string) bool {
	t.Helper()
	if !nameLogged {
		t.Log(name)
	}
	return true
}

// errsContainStr returns true if any of the errors contains the string,
// false otherwise
func errsContainStr(errs []error, s string) bool {
	for _, err := range errs {
		if strings.Contains(err.Error(), s) {
			return true
		}
	}
	return false
}

// errMapCheck checks the error map and reports any discrepancies with the
// expected values
func errMapCheck(t *testing.T, testID string, errMap param.ErrMap, expected map[string][]string) {
	t.Helper()

	var nameLogged bool

	if len(errMap) != len(expected) {
		nameLogged = logName(t, nameLogged, testID)
		t.Errorf(
			"\t: the error map had %d entries, it was expected to have %d\n",
			len(errMap), len(expected))
	}
	for k, errs := range errMap {
		expStrs, ok := expected[k]
		if !ok {
			nameLogged = logName(t, nameLogged, testID)
			logErrs(t, "\t: there are unexpected errors for: "+k+":\n", errs)
			continue
		}
		for _, s := range expStrs {
			if !errsContainStr(errs, s) {
				nameLogged = logName(t, nameLogged, testID)
				t.Errorf("\t: errors for '%s' should contain '%s' but don't",
					k, s)
			}
		}
	}

	for k := range expected {
		if _, ok := errMap[k]; !ok {
			nameLogged = logName(t, nameLogged, testID)
			t.Errorf("\t: error map should contain '%s' but doesn't", k)
		}
	}

	if nameLogged {
		logErrMap(t, errMap)
	}
}

// makePSetOrFatal creates the PSet and if any error is detected it reports a
// fatal error
func makePSetOrFatal(t *testing.T, testID string) *param.PSet {
	ps, err := paramset.NewNoHelpNoExitNoErrRpt()
	if err != nil {
		t.Fatal(testID, " : couldn't construct the PSet: ", err)
	}

	return ps
}
