package phelp

import (
	"bytes"
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
	"github.com/nickwells/twrap.mod/twrap"
)

func TestShowNotesRefsFmtStd(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		name      string
		refs      []string
		expOutput string
	}{
		{
			ID: testhelper.MkID("no refs - no output expected"),
		},
		{
			ID:        testhelper.MkID("one ref"),
			name:      "title",
			refs:      []string{"hello world"},
			expOutput: "            See title: hello world\n",
		},
		{
			ID:        testhelper.MkID("multi-refs"),
			name:      "title",
			refs:      []string{"hello", "world"},
			expOutput: "            See titles: hello, world\n",
		},
		{
			ID:   testhelper.MkID("multi-refs - one too long"),
			name: "title",
			refs: []string{
				"hello",                         // len: 5
				"world",                         // len: 5
				"a long string that would wrap", // len: 29
			},
			expOutput: "            See titles: hello, world,\n" +
				"                        a long string that would wrap\n",
		},
	}

	for _, tc := range testCases {
		buf := new(bytes.Buffer)
		twc := twrap.NewTWConfOrPanic(
			twrap.SetWriter(buf),
			twrap.SetTargetLineLen(60))
		showNotesRefsFmtStd(twc, tc.refs, tc.name)
		testhelper.DiffString[string](t,
			tc.IDStr(), "notes reference",
			buf.String(), tc.expOutput)
	}
}
