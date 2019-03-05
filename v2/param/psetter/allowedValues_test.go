package psetter_test

import (
	"flag"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/nickwells/param.mod/v2/param"
	"github.com/nickwells/param.mod/v2/param/psetter"
	"github.com/nickwells/testhelper.mod/testhelper"
)

const (
	allowedValsDir    = "testdata"
	allowedValsSubDir = "allowedVals"
)

var updateAVals = flag.Bool("upd-avals", false,
	"update the files holding the allowed values messages")

func TestAllowedValues(t *testing.T) {
	gfc := testhelper.GoldenFileCfg{
		DirNames: []string{allowedValsDir, allowedValsSubDir},
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
		name string
		s    param.Setter
	}{
		{
			name: "Bool",
			s:    &psetter.Bool{Value: &b},
		},
		{
			name: "Duration",
			s:    &psetter.Duration{Value: &dur},
		},
		{
			name: "EnumList",
			s: &psetter.EnumList{
				Value: &emptyStrList,
				AllowedVals: psetter.AValMap{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
		},
		{
			name: "EnumMap",
			s: &psetter.EnumMap{
				Value: &strToBoolMap,
				AllowedVals: psetter.AValMap{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
		},
		{
			name: "Enum",
			s: &psetter.Enum{
				Value: &goodStr,
				AllowedVals: psetter.AValMap{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
		},
		{
			name: "Float64",
			s:    &psetter.Float64{Value: &f},
		},
		{
			name: "Int64",
			s:    &psetter.Int64{Value: &i},
		},
		{
			name: "Int64List",
			s:    &psetter.Int64List{Value: &intList},
		},
		{
			name: "Map",
			s: &psetter.Map{
				Value: &strToBoolMap,
			},
		},
		{
			name: "Nil",
			s:    &psetter.Nil{},
		},
		{
			name: "Pathname",
			s:    &psetter.Pathname{Value: &anyStr},
		},
		{
			name: "Regexp",
			s:    &psetter.Regexp{Value: &re},
		},
		{
			name: "StrList",
			s:    &psetter.StrList{Value: &emptyStrList},
		},
		{
			name: "String",
			s:    &psetter.String{Value: &anyStr},
		},
		{
			name: "TimeLocation",
			s:    &psetter.TimeLocation{Value: &timeLoc},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s.AllowedValues", i, tc.name)

		val := []byte(tc.s.AllowedValues())
		testhelper.CheckAgainstGoldenFile(t, tcID, val,
			gfc.PathName(tc.name), *updateAVals)
	}
}
