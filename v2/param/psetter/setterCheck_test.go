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
			name:          "Bool - ok",
			s:             &psetter.Bool{Value: &b},
			panicExpected: false,
		},
		{
			name:          "Bool - bad",
			s:             &psetter.Bool{},
			panicExpected: true,
			expVals:       []string{"test: Bool " + nilValueMsg},
		},
		{
			name:          "BoolNot - ok",
			s:             &psetter.BoolNot{Value: &b},
			panicExpected: false,
		},
		{
			name:          "BoolNot - bad",
			s:             &psetter.BoolNot{},
			panicExpected: true,
			expVals:       []string{"test: BoolNot " + nilValueMsg},
		},
		{
			name:          "Duration - ok",
			s:             &psetter.Duration{Value: &dur},
			panicExpected: false,
		},
		{
			name:          "Duration - bad",
			s:             &psetter.Duration{},
			panicExpected: true,
			expVals:       []string{"test: Duration " + nilValueMsg},
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
			panicExpected: false,
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
			panicExpected: false,
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
			expVals: append([]string{"test: EnumList "},
				badInitialList...),
		},
		{
			name:          "EnumList - bad - no value",
			s:             &psetter.EnumList{},
			panicExpected: true,
			expVals:       []string{"test: EnumList " + nilValueMsg},
		},
		{
			name: "EnumList - bad - no allowedValues",
			s: &psetter.EnumList{
				Value: &emptyStrList,
			},
			panicExpected: true,
			expVals: append([]string{"test: EnumList "},
				tooFewAValsMsg...),
		},
		{
			name: "EnumList - bad - empty allowedValues",
			s: &psetter.EnumList{
				Value:       &emptyStrList,
				AllowedVals: psetter.AValMap{},
			},
			panicExpected: true,
			expVals: append([]string{"test: EnumList "},
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
			expVals: append([]string{"test: EnumList "},
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
			panicExpected: false,
		},
		{
			name:          "EnumMap - bad - no value",
			s:             &psetter.EnumMap{},
			panicExpected: true,
			expVals:       []string{"test: EnumMap " + nilValueMsg},
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
			expVals:       []string{"test: EnumMap " + mapNotCreatedMsg},
		},
		{
			name: "EnumMap - bad - no allowedValues",
			s: &psetter.EnumMap{
				Value: &strToBoolMap,
			},
			panicExpected: true,
			expVals: append([]string{"test: EnumMap "},
				tooFewAValsMsg...),
		},
		{
			name: "EnumMap - bad - empty allowedValues",
			s: &psetter.EnumMap{
				Value:       &strToBoolMap,
				AllowedVals: psetter.AValMap{},
			},
			panicExpected: true,
			expVals: append([]string{"test: EnumMap "},
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
			panicExpected: false,
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
			expVals:       append([]string{"test: Enum "}, badInitialVal...),
		},
		{
			name:          "Enum - bad - no value",
			s:             &psetter.Enum{},
			panicExpected: true,
			expVals:       []string{"test: Enum " + nilValueMsg},
		},
		{
			name: "Enum - bad - no allowedValues",
			s: &psetter.Enum{
				Value: &anyStr,
			},
			panicExpected: true,
			expVals: append([]string{"test: Enum "},
				tooFewAValsMsg...),
		},
		{
			name: "Enum - bad - empty allowedValues",
			s: &psetter.Enum{
				Value:       &anyStr,
				AllowedVals: psetter.AValMap{},
			},
			panicExpected: true,
			expVals: append([]string{"test: Enum "},
				tooFewAValsMsg...),
		},
		{
			name:          "Float64Setter - ok",
			s:             &psetter.Float64Setter{Value: &f},
			panicExpected: false,
		},
		{
			name:          "Float64Setter - bad",
			s:             &psetter.Float64Setter{},
			panicExpected: true,
			expVals:       []string{"test: Float64Setter " + nilValueMsg},
		},
		{
			name:          "Int64Setter - ok",
			s:             &psetter.Int64Setter{Value: &i},
			panicExpected: false,
		},
		{
			name:          "Int64Setter - bad",
			s:             &psetter.Int64Setter{},
			panicExpected: true,
			expVals:       []string{"test: Int64Setter " + nilValueMsg},
		},
		{
			name:          "Int64ListSetter - ok",
			s:             &psetter.Int64ListSetter{Value: &intList},
			panicExpected: false,
		},
		{
			name:          "Int64ListSetter - bad",
			s:             &psetter.Int64ListSetter{},
			panicExpected: true,
			expVals:       []string{"test: Int64ListSetter " + nilValueMsg},
		},
		{
			name: "Map - ok",
			s: &psetter.Map{
				Value: &strToBoolMap,
			},
			panicExpected: false,
		},
		{
			name:          "Map - bad - no value",
			s:             &psetter.Map{},
			panicExpected: true,
			expVals:       []string{"test: Map " + nilValueMsg},
		},
		{
			name: "Map - bad - nil map",
			s: &psetter.Map{
				Value: &strToBoolMapNil,
			},
			panicExpected: true,
			expVals:       []string{"test: Map " + mapNotCreatedMsg},
		},
		{
			name:          "Nil - ok",
			s:             &psetter.Nil{},
			panicExpected: false,
		},
		{
			name:          "Pathname - ok",
			s:             &psetter.Pathname{Value: &anyStr},
			panicExpected: false,
		},
		{
			name:          "Pathname - bad",
			s:             &psetter.Pathname{},
			panicExpected: true,
			expVals:       []string{"test: Pathname " + nilValueMsg},
		},
		{
			name:          "Regexp - ok",
			s:             &psetter.Regexp{Value: &re},
			panicExpected: false,
		},
		{
			name:          "Regexp - bad",
			s:             &psetter.Regexp{},
			panicExpected: true,
			expVals:       []string{"test: Regexp " + nilValueMsg},
		},
		{
			name:          "StrList - ok",
			s:             &psetter.StrList{Value: &emptyStrList},
			panicExpected: false,
		},
		{
			name:          "StrList - bad",
			s:             &psetter.StrList{},
			panicExpected: true,
			expVals:       []string{"test: StrList " + nilValueMsg},
		},
		{
			name:          "String - ok",
			s:             &psetter.String{Value: &anyStr},
			panicExpected: false,
		},
		{
			name:          "String - bad",
			s:             &psetter.String{},
			panicExpected: true,
			expVals:       []string{"test: String " + nilValueMsg},
		},
		{
			name:          "TimeLocation - ok",
			s:             &psetter.TimeLocation{Value: &timeLoc},
			panicExpected: false,
		},
		{
			name:          "TimeLocation - bad",
			s:             &psetter.TimeLocation{},
			panicExpected: true,
			expVals:       []string{"test: TimeLocation " + nilValueMsg},
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
