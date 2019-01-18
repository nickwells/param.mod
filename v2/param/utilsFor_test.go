package param_test

import (
	"errors"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/nickwells/param.mod/v2/param"
)

// =======================================================

// failingSetter will return an error in all cases - it's intended for use in
// test cases where we want a setter to fail
type failingSetter struct {
	errMsg string
}

// ValueReq returns Mandatory indicating that some value must follow
// the parameter
func (s failingSetter) ValueReq() param.ValueReq {
	return param.Mandatory
}

// Set (called when there is no following value) returns an error
func (s failingSetter) Set(_ string) error {
	return errors.New("no value given")
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

type posParamInitialiser struct {
	setter param.Setter
	name   string
	desc   string
	opts   []param.PosOptFunc
}

// paramInitialisers holds a pointer to either a namedParamInitialiser or a
// posParamInitialiser pointer (either, neither or both could be nil)
type paramInitialisers struct {
	npi            *namedParamInitialiser
	npiShouldExist bool

	ppi            *posParamInitialiser
	ppiShouldExist bool
}

// panicSafeTestAddByPos adds a ByPos parameter to a parameter set and
// catches any panics. Then it returns the added parameter (if any), a
// boolean indicating if a panic occured and the value recovered from the
// panic
func panicSafeTestAddByPos(ps *param.ParamSet, ppi *posParamInitialiser) (pp *param.ByPos, panicked bool, panicVal interface{}) {
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
// boolean indicating if a panic occured and the value recovered from the
// panic
func panicSafeTestAddByName(ps *param.ParamSet, npi *namedParamInitialiser) (p *param.ByName, panicked bool, panicVal interface{}) {
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
// if a panic occured and the value recovered from the panic
func panicSafeTestParse(ps *param.ParamSet, params []string) (
	errMap param.ErrMap, panicked bool, panicVal interface{}, stackTrace []byte) {
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
// panics. Then it returns a boolean indicating if a panic occured and the
// value recovered from the panic
func panicSafeSetGroupDescription(ps *param.ParamSet, groupName, desc string) (
	panicked bool, panicVal interface{}, stackTrace []byte) {
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

// logErrMap reports the contents of the error map returned by param.Parse(...)
func logErrMap(t *testing.T, errMap param.ErrMap) {
	t.Helper()

	for k, v := range errMap {
		t.Log("\t\t", k, ":\n")
		for _, err := range v {
			t.Log("\t\t\t", err, "\n")
		}
	}
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
			t.Errorf("\t: there are unexpected errors for: '%s':", k)
			for _, err := range errs {
				t.Logf("\t\t: %s", err)
			}
		} else {
			for _, s := range expStrs {
				count := 0
				for _, err := range errs {
					if strings.Contains(err.Error(), s) {
						count++
					}
				}
				if count == 0 {
					nameLogged = logName(t, nameLogged, testID)
					t.Errorf("\t: errors for '%s' should contain '%s' but don't",
						k, s)
				}
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
		for k, v := range errMap {
			t.Logf("\t: Errors for: %s", k)
			for _, e := range v {
				t.Logf("\t:     %v", e)
			}
		}
	}
}
