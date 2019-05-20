package param

import (
	"testing"
)

func TestSplitParamName(t *testing.T) {
	testCases := []struct {
		testName     string
		pName        string
		expProgName  string
		expParamName string
	}{
		{
			testName:     "nil",
			pName:        "",
			expProgName:  "",
			expParamName: "",
		},
		{
			testName:     "param only",
			pName:        "param",
			expProgName:  "",
			expParamName: "param",
		},
		{
			testName:     "param only (with whitespace)",
			pName:        "   param   ",
			expProgName:  "",
			expParamName: "param",
		},
		{
			testName:     "progname only",
			pName:        "progname/",
			expProgName:  "progname",
			expParamName: "",
		},
		{
			testName:     "progname only (with whitespace)",
			pName:        "  progname  /  ",
			expProgName:  "progname",
			expParamName: "",
		},
		{
			testName:     "progname and param",
			pName:        "progname/param",
			expProgName:  "progname",
			expParamName: "param",
		},
		{
			testName:     "progname and param (with whitespace)",
			pName:        "  progname  /  param  ",
			expProgName:  "progname",
			expParamName: "param",
		},
	}

	for i, tc := range testCases {
		progName, paramName := splitParamName(tc.pName)
		if progName != tc.expProgName {
			t.Logf("test %d: %s : Unexpected program name.\n",
				i, tc.testName)
			t.Errorf("\t: Got %s, expected: %s\n",
				progName, tc.expProgName)
		}
		if paramName != tc.expParamName {
			t.Logf("test %d: %s : Unexpected param name.\n",
				i, tc.testName)
			t.Errorf("\t: Got %s, expected: %s\n",
				paramName, tc.expParamName)
		}
	}

}
