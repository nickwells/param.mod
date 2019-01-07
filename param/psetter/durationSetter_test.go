package psetter_test

import (
	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/param/psetter"
	"testing"
	"time"
)

func TestDurationSetter(t *testing.T) {
	testCases := []struct {
		testName      string
		val           string
		expDuration   time.Duration
		check         func(d time.Duration) error
		errorExpected bool
	}{
		{
			testName:      "bad duration",
			val:           "blah",
			errorExpected: true,
		},
		{
			testName:      "good duration - 1 hour",
			val:           "1h",
			expDuration:   time.Duration(1) * time.Hour,
			errorExpected: false,
		},
		{
			testName:      "bad duration - 1 hour LT check",
			val:           "1h",
			check:         check.DurationLT(time.Duration(1) * time.Hour),
			errorExpected: true,
		},
		{
			testName:      "good duration - 1 hour LT check",
			val:           "1m",
			expDuration:   time.Duration(1) * time.Minute,
			check:         check.DurationLT(time.Duration(1) * time.Hour),
			errorExpected: false,
		},
		{
			testName:      "good duration - 2 hour GT check",
			val:           "2h",
			expDuration:   time.Duration(2) * time.Hour,
			check:         check.DurationGT(time.Duration(1) * time.Hour),
			errorExpected: false,
		},
		{
			testName:      "bad duration - 1 hour GT check",
			val:           "1h",
			check:         check.DurationGT(time.Duration(1) * time.Hour),
			errorExpected: true,
		},
	}

	for i, tc := range testCases {
		var d time.Duration
		ds := psetter.DurationSetter{
			Value:  &d,
			Checks: []check.Duration{tc.check},
		}

		err := ds.SetWithVal("", tc.val)
		if err != nil {
			if !tc.errorExpected {
				t.Errorf("test %d: %s : an unexpected error was returned when processing '%s': %s",
					i, tc.testName, tc.val, err)
			}
		} else if tc.errorExpected {
			t.Errorf("test %d: %s : an error was expected when processing '%s' but none was returned",
				i, tc.testName, tc.val)
		} else if d != tc.expDuration {
			t.Errorf("test %d: %s : the duration was not as expected, got %v, expected %v",
				i, tc.testName, d, tc.expDuration)
		}
	}

}
