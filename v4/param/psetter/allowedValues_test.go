package psetter_test

import (
	"flag"
	"regexp"
	"testing"
	"time"

	"github.com/nickwells/param.mod/v4/param"
	"github.com/nickwells/param.mod/v4/param/psetter"
	"github.com/nickwells/testhelper.mod/testhelper"
)

const (
	testDataDir       = "testdata"
	allowedValsSubDir = "allowedVals"
)

var updateAVals = flag.Bool("upd-avals", false,
	"update the files holding the allowed values messages")

func TestAllowedValues(t *testing.T) {
	gfc := testhelper.GoldenFileCfg{
		DirNames: []string{testDataDir, allowedValsSubDir},
		Sfx:      "txt",
	}
	var b bool
	var dur time.Duration
	var emptyStrList []string
	var strToBoolMap = make(map[string]bool)
	var goodStr = "aval"
	var anyStr = ""
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
				AllowedVals: param.AllowedVals{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
		},
		{
			ID: testhelper.MkID("EnumMap"),
			s: &psetter.EnumMap{
				Value: &strToBoolMap,
				AllowedVals: param.AllowedVals{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
		},
		{
			ID: testhelper.MkID("Enum"),
			s: &psetter.Enum{
				Value: &goodStr,
				AllowedVals: param.AllowedVals{
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
		testhelper.CheckAgainstGoldenFile(t, tc.IDStr(), val,
			gfc.PathName(tc.ID.Name), *updateAVals)
	}
}
