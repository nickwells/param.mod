package psetter_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/nickwells/param.mod/v2/param"
	"github.com/nickwells/param.mod/v2/param/psetter"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestCheck(t *testing.T) {
	var b bool
	var dur time.Duration
	var emptyStrList []string
	var goodStrList = []string{"aval"}
	var badStrList = []string{"bad val"}
	var strToBoolMap = make(map[string]bool)
	var strToBoolMapNil map[string]bool
	var goodStr = "aval"
	var badStr = "bad val"
	var anyStr = ""
	var f float64
	var i int64
	var intList []int64
	var re *regexp.Regexp
	var timeLoc *time.Location

	nilValueMsg := "Check failed: the Value to be set is nil"
	tooFewAValsMsg := []string{
		"Check failed: the allowed values map has ",
		"entries.",
		"It should have more than 1",
	}
	badInitialVal := []string{
		"Check failed: the initial value",
		"is not valid",
	}
	badInitialList := []string{
		"Check failed: element",
		"in the current list of entries is invalid",
	}
	mapNotCreatedMsg := "Check failed: the map has not been created"

	testCases := []struct {
		name          string
		s             param.Setter
		panicExpected bool
		expVals       []string
	}{
		{
			name: "Bool - ok",
			s:    &psetter.Bool{Value: &b},
		},
		{
			name:          "Bool - bad",
			s:             &psetter.Bool{},
			panicExpected: true,
			expVals:       []string{"test: psetter.Bool " + nilValueMsg},
		},
		{
			name: "Duration - ok",
			s:    &psetter.Duration{Value: &dur},
		},
		{
			name:          "Duration - bad",
			s:             &psetter.Duration{},
			panicExpected: true,
			expVals:       []string{"test: psetter.Duration " + nilValueMsg},
		},
		{
			name: "EnumList - ok - no strings",
			s: &psetter.EnumList{
				Value: &emptyStrList,
				AllowedVals: psetter.AValMap{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
		},
		{
			name: "EnumList - ok - good strings",
			s: &psetter.EnumList{
				Value: &goodStrList,
				AllowedVals: psetter.AValMap{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
		},
		{
			name: "EnumList - bad initial value",
			s: &psetter.EnumList{
				Value: &badStrList,
				AllowedVals: psetter.AValMap{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
			panicExpected: true,
			expVals: append([]string{"test: psetter.EnumList "},
				badInitialList...),
		},
		{
			name:          "EnumList - bad - no value",
			s:             &psetter.EnumList{},
			panicExpected: true,
			expVals:       []string{"test: psetter.EnumList " + nilValueMsg},
		},
		{
			name: "EnumList - bad - no allowedValues",
			s: &psetter.EnumList{
				Value: &emptyStrList,
			},
			panicExpected: true,
			expVals: append([]string{"test: psetter.EnumList "},
				tooFewAValsMsg...),
		},
		{
			name: "EnumList - bad - empty allowedValues",
			s: &psetter.EnumList{
				Value:       &emptyStrList,
				AllowedVals: psetter.AValMap{},
			},
			panicExpected: true,
			expVals: append([]string{"test: psetter.EnumList "},
				tooFewAValsMsg...),
		},
		{
			name: "EnumList - bad - allowedValues - only one entry",
			s: &psetter.EnumList{
				Value: &emptyStrList,
				AllowedVals: psetter.AValMap{
					"aval": "desc",
				},
			},
			panicExpected: true,
			expVals: append([]string{"test: psetter.EnumList "},
				tooFewAValsMsg...),
		},
		{
			name: "EnumMap - ok",
			s: &psetter.EnumMap{
				Value: &strToBoolMap,
				AllowedVals: psetter.AValMap{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
		},
		{
			name:          "EnumMap - bad - no value",
			s:             &psetter.EnumMap{},
			panicExpected: true,
			expVals:       []string{"test: psetter.EnumMap " + nilValueMsg},
		},
		{
			name: "EnumMap - bad - nil map",
			s: &psetter.EnumMap{
				Value: &strToBoolMapNil,
				AllowedVals: psetter.AValMap{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
			panicExpected: true,
			expVals:       []string{"test: psetter.EnumMap " + mapNotCreatedMsg},
		},
		{
			name: "EnumMap - bad - no allowedValues",
			s: &psetter.EnumMap{
				Value: &strToBoolMap,
			},
			panicExpected: true,
			expVals: append([]string{"test: psetter.EnumMap "},
				tooFewAValsMsg...),
		},
		{
			name: "EnumMap - bad - empty allowedValues",
			s: &psetter.EnumMap{
				Value:       &strToBoolMap,
				AllowedVals: psetter.AValMap{},
			},
			panicExpected: true,
			expVals: append([]string{"test: psetter.EnumMap "},
				tooFewAValsMsg...),
		},
		{
			name: "Enum - ok",
			s: &psetter.Enum{
				Value: &goodStr,
				AllowedVals: psetter.AValMap{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
		},
		{
			name: "Enum - bad initial value",
			s: &psetter.Enum{
				Value: &badStr,
				AllowedVals: psetter.AValMap{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
			panicExpected: true,
			expVals:       append([]string{"test: psetter.Enum "}, badInitialVal...),
		},
		{
			name:          "Enum - bad - no value",
			s:             &psetter.Enum{},
			panicExpected: true,
			expVals:       []string{"test: psetter.Enum " + nilValueMsg},
		},
		{
			name: "Enum - bad - no allowedValues",
			s: &psetter.Enum{
				Value: &anyStr,
			},
			panicExpected: true,
			expVals: append([]string{"test: psetter.Enum "},
				tooFewAValsMsg...),
		},
		{
			name: "Enum - bad - empty allowedValues",
			s: &psetter.Enum{
				Value:       &anyStr,
				AllowedVals: psetter.AValMap{},
			},
			panicExpected: true,
			expVals: append([]string{"test: psetter.Enum "},
				tooFewAValsMsg...),
		},
		{
			name: "Float64 - ok",
			s:    &psetter.Float64{Value: &f},
		},
		{
			name:          "Float64 - bad",
			s:             &psetter.Float64{},
			panicExpected: true,
			expVals:       []string{"test: psetter.Float64 " + nilValueMsg},
		},
		{
			name: "Int64 - ok",
			s:    &psetter.Int64{Value: &i},
		},
		{
			name:          "Int64 - bad",
			s:             &psetter.Int64{},
			panicExpected: true,
			expVals:       []string{"test: psetter.Int64 " + nilValueMsg},
		},
		{
			name: "Int64List - ok",
			s:    &psetter.Int64List{Value: &intList},
		},
		{
			name:          "Int64List - bad",
			s:             &psetter.Int64List{},
			panicExpected: true,
			expVals:       []string{"test: psetter.Int64List " + nilValueMsg},
		},
		{
			name: "Map - ok",
			s: &psetter.Map{
				Value: &strToBoolMap,
			},
		},
		{
			name:          "Map - bad - no value",
			s:             &psetter.Map{},
			panicExpected: true,
			expVals:       []string{"test: psetter.Map " + nilValueMsg},
		},
		{
			name: "Map - bad - nil map",
			s: &psetter.Map{
				Value: &strToBoolMapNil,
			},
			panicExpected: true,
			expVals:       []string{"test: psetter.Map " + mapNotCreatedMsg},
		},
		{
			name: "Nil - ok",
			s:    &psetter.Nil{},
		},
		{
			name: "Pathname - ok",
			s:    &psetter.Pathname{Value: &anyStr},
		},
		{
			name:          "Pathname - bad",
			s:             &psetter.Pathname{},
			panicExpected: true,
			expVals:       []string{"test: psetter.Pathname " + nilValueMsg},
		},
		{
			name: "Regexp - ok",
			s:    &psetter.Regexp{Value: &re},
		},
		{
			name:          "Regexp - bad",
			s:             &psetter.Regexp{},
			panicExpected: true,
			expVals:       []string{"test: psetter.Regexp " + nilValueMsg},
		},
		{
			name: "StrList - ok",
			s:    &psetter.StrList{Value: &emptyStrList},
		},
		{
			name:          "StrList - bad",
			s:             &psetter.StrList{},
			panicExpected: true,
			expVals:       []string{"test: psetter.StrList " + nilValueMsg},
		},
		{
			name: "String - ok",
			s:    &psetter.String{Value: &anyStr},
		},
		{
			name:          "String - bad",
			s:             &psetter.String{},
			panicExpected: true,
			expVals:       []string{"test: psetter.String " + nilValueMsg},
		},
		{
			name: "TimeLocation - ok",
			s:    &psetter.TimeLocation{Value: &timeLoc},
		},
		{
			name:          "TimeLocation - bad",
			s:             &psetter.TimeLocation{},
			panicExpected: true,
			expVals:       []string{"test: psetter.TimeLocation " + nilValueMsg},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		panicked, panicVal := panicSafeCheck(tc.s)

		testhelper.PanicCheckString(t, tcID,
			panicked, tc.panicExpected,
			panicVal, tc.expVals)
	}
}
