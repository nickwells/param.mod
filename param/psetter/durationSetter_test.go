package psetter_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/v4/param/psetter"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestDuration(t *testing.T) {
	testCases := []struct {
		name          string
		val           string
		expDuration   time.Duration
		check         func(d time.Duration) error
		errorExpected bool
	}{
		{
			name:          "bad duration",
			val:           "blah",
			errorExpected: true,
		},
		{
			name:          "good duration - 1 hour",
			val:           "1h",
			expDuration:   time.Duration(1) * time.Hour,
			errorExpected: false,
		},
		{
			name:          "bad duration - 1 hour LT check",
			val:           "1h",
			check:         check.DurationLT(time.Duration(1) * time.Hour),
			errorExpected: true,
		},
		{
			name:          "good duration - 1 hour LT check",
			val:           "1m",
			expDuration:   time.Duration(1) * time.Minute,
			check:         check.DurationLT(time.Duration(1) * time.Hour),
			errorExpected: false,
		},
		{
			name:          "good duration - 2 hour GT check",
			val:           "2h",
			expDuration:   time.Duration(2) * time.Hour,
			check:         check.DurationGT(time.Duration(1) * time.Hour),
			errorExpected: false,
		},
		{
			name:          "bad duration - 1 hour GT check",
			val:           "1h",
			check:         check.DurationGT(time.Duration(1) * time.Hour),
			errorExpected: true,
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		var d time.Duration
		ds := psetter.Duration{
			Value:  &d,
			Checks: []check.Duration{tc.check},
		}

		err := ds.SetWithVal("", tc.val)
		ok := testhelper.CheckError(t, tcID, err, tc.errorExpected, []string{})
		if ok && err == nil {
			if d != tc.expDuration {
				t.Log(tcID)
				t.Logf("\t: duration expected: %v\n", tc.expDuration)
				t.Logf("\t:               was: %v\n", d)
				t.Errorf("\t: unexpected duration\n")
			}
		}
	}
}

func TestDurationCountChecks(t *testing.T) {
	testCases := []struct {
		name     string
		checks   []check.Duration
		expCount int
	}{
		{
			name:     "no checks",
			checks:   []check.Duration{},
			expCount: 0,
		},
		{
			name:     "one check",
			checks:   []check.Duration{check.DurationLT(99)},
			expCount: 1,
		},
		{
			name: "two checks",
			checks: []check.Duration{
				check.DurationLT(99),
				check.DurationGT(9),
			},
			expCount: 2,
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		var d time.Duration
		ds := psetter.Duration{
			Value:  &d,
			Checks: tc.checks,
		}
		if ds.CountChecks() != tc.expCount {
			t.Log(tcID)
			t.Logf("\t: check count expected: %d\n", tc.expCount)
			t.Logf("\t:                  was: %d\n", ds.CountChecks())
			t.Errorf("\t: unexpected check count\n")
		}
	}
}
