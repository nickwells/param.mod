package psetter_test

import (
	"testing"

	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestTaggedCalc_Check(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		tag  string
		calc func(n, v string) (int, error)
	}{
		{
			ID:   testhelper.MkID("good"),
			tag:  "good",
			calc: func(n string, v string) (int, error) { return 42, nil },
		},
		{
			ID:     testhelper.MkID("bad name"),
			ExpErr: testhelper.MkExpErr("the Name must not be empty"),
			calc:   func(n string, v string) (int, error) { return 42, nil },
		},
		{
			ID:     testhelper.MkID("bad calc"),
			ExpErr: testhelper.MkExpErr("the Calc must not be nil"),
			tag:    "good",
		},
	}

	for _, tc := range testCases {
		v := psetter.NamedCalc[int]{Name: tc.tag, Calc: tc.calc}
		err := v.Check()
		testhelper.CheckExpErr(t, err, tc)
	}
}
