package param

import (
	"testing"

	"github.com/nickwells/testhelper.mod/testhelper"
)

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
