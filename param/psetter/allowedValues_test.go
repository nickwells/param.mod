package psetter_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/psetter"
	"github.com/nickwells/testhelper.mod/testhelper"
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
			s: &psetter.EnumList{
				Value: &emptyStrList,
				AllowedVals: psetter.AllowedVals{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
		},
		{
			ID: testhelper.MkID("EnumMap"),
			s: &psetter.EnumMap{
				Value: &strToBoolMap,
				AllowedVals: psetter.AllowedVals{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
		},
		{
			ID: testhelper.MkID("Enum"),
			s: &psetter.Enum{
				Value: &goodStr,
				AllowedVals: psetter.AllowedVals{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
		},
		{
			ID: testhelper.MkID("Float64"),
			s:  &psetter.Float64{Value: &f},
		},
		{
			ID: testhelper.MkID("Int64"),
			s:  &psetter.Int64{Value: &i},
		},
		{
			ID: testhelper.MkID("Int64List"),
			s:  &psetter.Int64List{Value: &intList},
		},
		{
			ID: testhelper.MkID("Map"),
			s: &psetter.Map{
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
			s:  &psetter.StrList{Value: &emptyStrList},
		},
		{
			ID: testhelper.MkID("StrListAppender"),
			s:  &psetter.StrListAppender{Value: &emptyStrList},
		},
		{
			ID: testhelper.MkID("StrListAppender-Prepend"),
			s:  &psetter.StrListAppender{Value: &emptyStrList, Prepend: true},
		},
		{
			ID: testhelper.MkID("String"),
			s:  &psetter.String{Value: &anyStr},
		},
		{
			ID: testhelper.MkID("TimeLocation"),
			s:  &psetter.TimeLocation{Value: &timeLoc},
		},
	}

	for _, tc := range testCases {
		val := []byte(tc.s.AllowedValues())
		gfc.Check(t, tc.IDStr(), tc.ID.Name, val)
	}
}
