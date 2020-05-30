package psetter_test

import (
	"testing"

	"github.com/nickwells/param.mod/v5/param/psetter"
)

func TestStrListSeparator(t *testing.T) {
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
			t.Error("GetSeparator() returned: '" + val +
				"' but '" + tc.expectedSep + "' was expected")
		}
	}
}
