package psetter_test

import (
	"fmt"
	"testing"

	"github.com/nickwells/param.mod/v6/psetter"
)

type hasChecksTester struct {
	count int
}

// (hct hasChecksTester)CountChecks returns an int
func (hct hasChecksTester) CountChecks() int {
	return hct.count
}

func TestHasChecks(t *testing.T) {
	testCases := []struct {
		name   string
		count  int
		expVal string
	}{
		{
			name:   "no checks",
			count:  0,
			expVal: "",
		},
		{
			name:   "1 check",
			count:  1,
			expVal: " subject to checks",
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		hct := hasChecksTester{count: tc.count}

		if got := psetter.HasChecks(hct); got != tc.expVal {
			t.Log(tcID)
			t.Logf("\t:    count: %d\n", tc.count)
			t.Logf("\t: expected: %s\n", tc.expVal)
			t.Logf("\t:      got: %s\n", got)
			t.Errorf("\t: unexpected return from HasChecks(...)\n")
		}
	}
}
