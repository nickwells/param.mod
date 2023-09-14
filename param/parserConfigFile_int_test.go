package param

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

type noHelpNoExitNoErrRpt struct{}

func (nh noHelpNoExitNoErrRpt) ProcessArgs(_ *PSet)            {}
func (nh noHelpNoExitNoErrRpt) Help(_ *PSet, _ ...string)      {}
func (nh noHelpNoExitNoErrRpt) AddParams(_ *PSet)              {}
func (nh noHelpNoExitNoErrRpt) ErrorHandler(_ *PSet, _ ErrMap) {}

var nhnenr noHelpNoExitNoErrRpt

type i64 struct {
	Value *int64
}

// SetWithVal (called when a value follows the parameter) checks that the
// value can be parsed to an integer, if it cannot be parsed successfully it
// returns an error. If there are checks and any check is violated it returns
// an error. Only if the value is parsed successfully and no checks are
// violated is the Value set.
func (s i64) SetWithVal(_ string, paramVal string) error {
	v, err := strconv.ParseInt(paramVal, 0, 0)
	if err != nil {
		return fmt.Errorf("could not parse %q as an integer value: %s",
			paramVal, err)
	}

	*s.Value = v
	return nil
}

// Set (called when a value follows the parameter) checks that the
// value can be parsed to an integer, if it cannot be parsed successfully it
// returns an error. If there are checks and any check is violated it returns
// an error. Only if the value is parsed successfully and no checks are
// violated is the Value set.
func (s i64) Set(name string) error {
	return fmt.Errorf("a value must follow this parameter: %q,"+
		" either following an '=' or as a next parameter", name)
}

// ValueReq returns Mandatory
func (s i64) ValueReq() ValueReq {
	return Mandatory
}

// AllowedValues returns a string describing the allowed values
func (s i64) AllowedValues() string {
	return "any value that can be read as a whole number"
}

// CurrentValue returns the current setting of the parameter value
func (s i64) CurrentValue() string {
	return fmt.Sprintf("%v", *s.Value)
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s i64) CheckSetter(name string) {
	if s.Value == nil {
		panic("No value pointer for: " + name + ": i64")
	}
}

func TestSplitParamName(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		pName        string
		expProgNames []string
		expParamName string
	}{
		{
			ID: testhelper.MkID("nil"),
		},
		{
			ID:           testhelper.MkID("param only"),
			pName:        "param",
			expParamName: "param",
		},
		{
			ID:           testhelper.MkID("param only (with whitespace)"),
			pName:        "   param   ",
			expParamName: "param",
		},
		{
			ID:           testhelper.MkID("param only (with embedded whitespace)"),
			pName:        "   param missing-equals   ",
			expParamName: "param missing-equals",
		},
		{
			ID:           testhelper.MkID("progname only"),
			pName:        "progname/",
			expProgNames: []string{"progname"},
		},
		{
			ID:           testhelper.MkID("progname only (with whitespace)"),
			pName:        "  progname  /  ",
			expProgNames: []string{"progname"},
		},
		{
			ID:           testhelper.MkID("progname and param"),
			pName:        "progname/param",
			expProgNames: []string{"progname"},
			expParamName: "param",
		},
		{
			ID: testhelper.MkID(
				"progname and param (with whitespace)"),
			pName:        "  progname  /  param  ",
			expProgNames: []string{"progname"},
			expParamName: "param",
		},
		{
			ID:           testhelper.MkID("progname only"),
			pName:        "progname1,progname2/",
			expProgNames: []string{"progname1", "progname2"},
		},
		{
			ID:           testhelper.MkID("progname only (with whitespace)"),
			pName:        "  progname1 , progname2  /  ",
			expProgNames: []string{"progname1", "progname2"},
		},
		{
			ID:           testhelper.MkID("progname and param"),
			pName:        "progname1,progname2/param",
			expProgNames: []string{"progname1", "progname2"},
			expParamName: "param",
		},
		{
			ID: testhelper.MkID(
				"progname and param (with whitespace)"),
			pName:        "  progname1 , progname2  /  param  ",
			expProgNames: []string{"progname1", "progname2"},
			expParamName: "param",
		},
	}

	for _, tc := range testCases {
		progNames, paramName := splitParamName(tc.pName)
		if testhelper.StringSliceDiff(progNames, tc.expProgNames) {
			t.Log(tc.IDStr())
			t.Logf("\t: expected: %v\n", tc.expProgNames)
			t.Logf("\t:      got: %v\n", progNames)
			t.Errorf("\t: Unexpected program names")
		}
		if paramName != tc.expParamName {
			t.Log(tc.IDStr())
			t.Logf("\t: expected: %s\n", tc.expParamName)
			t.Logf("\t:      got: %s\n", paramName)
			t.Errorf("\t: Unexpected param name")
		}
	}
}

func TestParamLineParser(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		line  string
		pname string
	}{
		{
			ID: testhelper.MkID("string has no equals"),
			ExpErr: testhelper.MkExpErr(
				"this is not a parameter of this program.",
				"Did you mean: ival ?"),
			line:  "ival 5",
			pname: "ival 5",
		},
	}
	var iVal int64
	for _, tc := range testCases {
		iVal = 0
		ps, err := NewSet(SetHelper(nhnenr))
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("\t: could not construct the paramset\n")
		}
		ps.Add("ival", i64{Value: &iVal}, "help...", Attrs(MustBeSet))

		pflp := paramLineParser{
			ps:    ps,
			eRule: paramMustExist,
		}
		loc := location.New("test case")
		_ = pflp.ParseLine(tc.line, loc)
		errs := ps.errors[tc.pname]
		if len(errs) == 0 && tc.ErrExpected() {
			t.Log(tc.IDStr())
			t.Errorf("\t: an error was expected but not found\n")
		} else if len(errs) != 1 {
			t.Log(tc.IDStr())
			for _, err := range errs {
				t.Log("\t: ", err)
			}
			t.Errorf("\t: too many errors were seen\n")
		} else {
			testhelper.CheckExpErr(t, errs[0], tc)
		}
	}
}
