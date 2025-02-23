package psetter_test

import (
	"testing"

	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestAllowedVals(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		av      psetter.AllowedVals[string]
		expAVal string
	}{
		{
			ID: testhelper.MkID("nil"),
		},
		{
			ID: testhelper.MkID("empty"),
			av: map[string]string{},
		},
		{
			ID: testhelper.MkID("one entry"),
			av: map[string]string{
				"name": "desc",
			},
			expAVal: `   name: desc`,
		},
		{
			ID: testhelper.MkID("two entries"),
			av: map[string]string{
				"name":      "desc",
				"long name": "long name desc",
			},
			expAVal: `   long name: long name desc
   name     : desc`,
		},
	}

	for _, tc := range testCases {
		s := tc.av.String()
		testhelper.DiffString(t, tc.IDStr(), "allowed values", s, tc.expAVal)
	}
}

func TestAllowedValuesMap(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		av     psetter.AllowedVals[string]
		expVal psetter.AllowedVals[string]
	}{
		{
			ID:     testhelper.MkID("empty"),
			av:     psetter.AllowedVals[string]{},
			expVal: psetter.AllowedVals[string]{},
		},
		{
			ID: testhelper.MkID("not empty"),
			av: psetter.AllowedVals[string]{
				"hello": "world",
			},
			expVal: psetter.AllowedVals[string]{
				"hello": "world",
			},
		},
	}

	for _, tc := range testCases {
		actVal := tc.av.AllowedValuesMap()

		err := testhelper.DiffVals(actVal, tc.expVal)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("\t: bad copy: %s\n", err)
		}
	}
}

func TestAllowedValuesCheck(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		av psetter.AllowedVals[string]
	}{
		{
			ID: testhelper.MkID("good"),
			av: psetter.AllowedVals[string]{
				"av1": "allowed value 1",
				"av2": "allowed value 2",
			},
		},
		{
			ID: testhelper.MkID("bad - empty"),
			ExpErr: testhelper.MkExpErr(
				"the map of allowed values has no entries." +
					" It should have at least 2"),
			av: psetter.AllowedVals[string]{},
		},
		{
			ID: testhelper.MkID("bad - one entry"),
			ExpErr: testhelper.MkExpErr(
				"the map of allowed values has only 1 entry." +
					" It should have at least 2"),
			av: psetter.AllowedVals[string]{
				"hello": "world",
			},
		},
		{
			ID: testhelper.MkID("bad - blank value"),
			ExpErr: testhelper.MkExpErr(
				`Bad allowed value: "": "blank" -` +
					` the allowed value must not be blank`),
			av: psetter.AllowedVals[string]{
				"":    "blank",
				"av1": "allowed value 1",
				"av2": "allowed value 2",
			},
		},
		{
			ID: testhelper.MkID("bad - value containing '='"),
			ExpErr: testhelper.MkExpErr(
				`Bad allowed value: "av0=x": "contains '='" -` +
					` the allowed value must not contain '='`),
			av: psetter.AllowedVals[string]{
				"av0=x": "contains '='",
				"av1":   "allowed value 1",
				"av2":   "allowed value 2",
			},
		},
	}

	for _, tc := range testCases {
		err := tc.av.Check()
		testhelper.CheckExpErr(t, err, tc)
	}
}
