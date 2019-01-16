package psetter_test

import (
	"github.com/nickwells/param.mod/v2/param/psetter"
	"testing"
)

func TestStringList(t *testing.T) {
	var sls psetter.StrListSeparator

	slsTestCases := [...]struct {
		sep         string
		expectedSep string
	}{
		{"", psetter.StrListDefaultSep},
		{":", ":"},
	}
	for _, tc := range slsTestCases {
		sls.Sep = tc.sep
		val := sls.GetSeparator()
		if val != tc.expectedSep {
			t.Error("GetSeparator() returned: '" + val + "' but '" + tc.expectedSep + "' was expected")
		}
	}
}
