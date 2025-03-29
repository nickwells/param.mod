package psetter_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

const (
	testDataDir       = "testdata"
	allowedValsSubDir = "allowedVals"
)

var gfc = testhelper.GoldenFileCfg{
	DirNames:               []string{testDataDir, allowedValsSubDir},
	Sfx:                    "txt",
	UpdFlagName:            "upd-help-files",
	KeepBadResultsFlagName: "keep-bad-results",
}

func init() {
	gfc.AddUpdateFlag()
	gfc.AddKeepBadResultsFlag()
}

func TestAllowedValues(t *testing.T) {
	var b bool

	var dur time.Duration

	var emptyStrList []string

	strToBoolMap := make(map[string]bool)

	var (
		goodStr = "aval"
		anyStr  = ""
	)

	var f float64

	var i int64

	var ui uint64

	var intList []int64

	var re *regexp.Regexp

	var timeLoc *time.Location

	testCases := []struct {
		testhelper.ID
		s param.Setter
	}{
		{
			ID: testhelper.MkID("Bool"),
			s:  &psetter.Bool{Value: &b},
		},
		{
			ID: testhelper.MkID("Duration"),
			s:  &psetter.Duration{Value: &dur},
		},
		{
			ID: testhelper.MkID("EnumList"),
			s: &psetter.EnumList[string]{
				Value: &emptyStrList,
				AllowedVals: psetter.AllowedVals[string]{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
		},
		{
			ID: testhelper.MkID("EnumMap"),
			s: &psetter.EnumMap[string]{
				Value: &strToBoolMap,
				AllowedVals: psetter.AllowedVals[string]{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
		},
		{
			ID: testhelper.MkID("Enum"),
			s: &psetter.Enum[string]{
				Value: &goodStr,
				AllowedVals: psetter.AllowedVals[string]{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
		},
		{
			ID: testhelper.MkID("Float64"),
			s:  &psetter.Float[float64]{Value: &f},
		},
		{
			ID: testhelper.MkID("Int64"),
			s:  &psetter.Int[int64]{Value: &i},
		},
		{
			ID: testhelper.MkID("Uint64"),
			s:  &psetter.Uint[uint64]{Value: &ui},
		},
		{
			ID: testhelper.MkID("Int64List"),
			s:  &psetter.IntList[int64]{Value: &intList},
		},
		{
			ID: testhelper.MkID("Map"),
			s: &psetter.Map[string]{
				Value: &strToBoolMap,
			},
		},
		{
			ID: testhelper.MkID("Nil"),
			s:  &psetter.Nil{},
		},
		{
			ID: testhelper.MkID("Pathname"),
			s:  &psetter.Pathname{Value: &anyStr},
		},
		{
			ID: testhelper.MkID("PathnameListAppender"),
			s:  &psetter.PathnameListAppender{Value: &emptyStrList},
		},
		{
			ID: testhelper.MkID("PathnameListAppender-Prepend"),
			s: &psetter.PathnameListAppender{
				Value:   &emptyStrList,
				Prepend: true,
			},
		},
		{
			ID: testhelper.MkID("Regexp"),
			s:  &psetter.Regexp{Value: &re},
		},
		{
			ID: testhelper.MkID("StrList"),
			s:  &psetter.StrList[string]{Value: &emptyStrList},
		},
		{
			ID: testhelper.MkID("StrListAppender"),
			s:  &psetter.StrListAppender[string]{Value: &emptyStrList},
		},
		{
			ID: testhelper.MkID("StrListAppender-Prepend"),
			s: &psetter.StrListAppender[string]{
				Value:   &emptyStrList,
				Prepend: true,
			},
		},
		{
			ID: testhelper.MkID("String"),
			s:  &psetter.String[string]{Value: &anyStr},
		},
		{
			ID: testhelper.MkID("TimeLocation"),
			s:  &psetter.TimeLocation{Value: &timeLoc},
		},
	}

	for _, tc := range testCases {
		val := []byte(tc.s.AllowedValues())
		gfc.Check(t, tc.IDStr(), tc.Name, val)
	}
}
