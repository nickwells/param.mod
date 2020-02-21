package psetter_test

import (
	"fmt"
	"testing"

	"github.com/nickwells/param.mod/v3/param"
	"github.com/nickwells/param.mod/v3/param/psetter"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestBoolSet(t *testing.T) {
	var b bool
	bs := psetter.Bool{Value: &b}
	bsInv := psetter.Bool{
		Value:  &b,
		Invert: true,
	}

	paramTestCases := [...]struct {
		name        string
		ps          param.Setter
		expectedVal bool
	}{
		{"Bool", bs, true},
		{"Bool (inverted)", bsInv, false},
	}
	for _, tc := range paramTestCases {
		b = !tc.expectedVal // force the value to != the expected value

		err := tc.ps.Set("dummy")
		if err != nil {
			t.Error(tc.name,
				"Set(...) returned an unexpected error:",
				err)
		}
		if b != tc.expectedVal {
			t.Error(tc.name,
				"Set(...) did not set the value to ",
				tc.expectedVal)
		}
	}
}

func TestBoolSetVal(t *testing.T) {
	var b bool
	bs := psetter.Bool{Value: &b}
	bsInv := psetter.Bool{
		Value:  &b,
		Invert: true,
	}

	testCases := []struct {
		name             string
		ps               param.Setter
		paramVal         string
		expectedVal      bool
		errExpected      bool
		errShouldContain []string
	}{
		{
			name:        "good - set to true",
			ps:          bs,
			paramVal:    "true",
			expectedVal: true,
		},
		{
			name:        "good - set to false",
			ps:          bs,
			paramVal:    "false",
			expectedVal: false,
		},
		{
			name:        "bad - invalid param",
			ps:          bs,
			paramVal:    "blah",
			errExpected: true,
			errShouldContain: []string{
				"cannot interpret 'blah' as either true or false",
			},
		},
		{
			name:        "good (inverted) - set to true",
			ps:          bsInv,
			paramVal:    "true",
			expectedVal: false,
		},
		{
			name:        "good (inverted) - set to false",
			ps:          bsInv,
			paramVal:    "false",
			expectedVal: true,
		},
		{
			name:        "bad (inverted) - invalid param",
			ps:          bsInv,
			paramVal:    "blah",
			errExpected: true,
			errShouldContain: []string{
				"cannot interpret 'blah' as either true or false",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		b = !tc.expectedVal
		err := tc.ps.SetWithVal("dummy", tc.paramVal)
		ok := testhelper.CheckError(t, tcID,
			err, tc.errExpected, tc.errShouldContain)
		if ok {
			if err != nil {
				if b != !tc.expectedVal {
					t.Log(tcID)
					t.Errorf("\t: value was changed despite error\n")
				}
			} else {
				if b != tc.expectedVal {
					t.Log(tcID)
					t.Errorf("\t: value was expected to be %v\n", b)
				}
			}
		}
	}
}
