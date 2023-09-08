package psetter_test

import (
	"testing"

	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestAllowedVals(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		avMap   psetter.AllowedVals[string]
		expAVal string
	}{
		{
			ID:    testhelper.MkID("empty"),
			avMap: map[string]string{},
		},
		{
			ID: testhelper.MkID("one entry"),
			avMap: map[string]string{
				"name": "desc",
			},
			expAVal: `   name: desc`,
		},
		{
			ID: testhelper.MkID("two entries"),
			avMap: map[string]string{
				"name":      "desc",
				"long name": "long name desc",
			},
			expAVal: `   long name: long name desc
   name     : desc`,
		},
	}

	for _, tc := range testCases {
		s := tc.avMap.String()
		testhelper.DiffString(t, tc.IDStr(), "allowed values", s, tc.expAVal)
	}
}
