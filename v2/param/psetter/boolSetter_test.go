package psetter_test

import (
	"github.com/nickwells/param.mod/v2/param"
	"github.com/nickwells/param.mod/v2/param/psetter"
	"testing"
)

func TestBoolSetter(t *testing.T) {
	var b bool
	bs := psetter.BoolSetter{Value: &b}
	bsn := psetter.BoolSetterNot{Value: &b}

	paramTestCases := [...]struct {
		name          string
		ps            param.Setter
		expectedValue bool
	}{
		{"BoolSetter", bs, true},
		{"BoolSetterNot", bsn, false},
	}
	for _, tc := range paramTestCases {
		b = !tc.expectedValue // force the value to != the expected value

		err := tc.ps.Set("")
		if err != nil {
			t.Error(tc.name,
				"Set(...) returned an unexpected error:",
				err)
		}
		if b != tc.expectedValue {
			t.Error(tc.name,
				"Set(...) did not set the value to ",
				tc.expectedValue)
		}

		err = tc.ps.SetWithVal("", "any")
		if err == nil {
			t.Error(tc.name,
				"SetWithVal(...) should have returned an error but didn't")
		}
	}
}
