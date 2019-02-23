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

func TestSetterCheck(t *testing.T) {
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
			name:          "BoolSetter - ok",
			s:             &psetter.BoolSetter{Value: &b},
			panicExpected: false,
		},
		{
			name:          "BoolSetter - bad",
			s:             &psetter.BoolSetter{},
			panicExpected: true,
			expVals:       []string{"test: BoolSetter " + nilValueMsg},
		},
		{
			name:          "BoolSetterNot - ok",
			s:             &psetter.BoolSetterNot{Value: &b},
			panicExpected: false,
		},
		{
			name:          "BoolSetterNot - bad",
			s:             &psetter.BoolSetterNot{},
			panicExpected: true,
			expVals:       []string{"test: BoolSetterNot " + nilValueMsg},
		},
		{
			name:          "DurationSetter - ok",
			s:             &psetter.DurationSetter{Value: &dur},
			panicExpected: false,
		},
		{
			name:          "DurationSetter - bad",
			s:             &psetter.DurationSetter{},
			panicExpected: true,
			expVals:       []string{"test: DurationSetter " + nilValueMsg},
		},
		{
			name: "EnumListSetter - ok - no strings",
			s: &psetter.EnumListSetter{
				Value: &emptyStrList,
				AllowedVals: psetter.AValMap{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
			panicExpected: false,
		},
		{
			name: "EnumListSetter - ok - good strings",
			s: &psetter.EnumListSetter{
				Value: &goodStrList,
				AllowedVals: psetter.AValMap{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
			panicExpected: false,
		},
		{
			name: "EnumListSetter - bad initial value",
			s: &psetter.EnumListSetter{
				Value: &badStrList,
				AllowedVals: psetter.AValMap{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
			panicExpected: true,
			expVals: append([]string{"test: EnumListSetter "},
				badInitialList...),
		},
		{
			name:          "EnumListSetter - bad - no value",
			s:             &psetter.EnumListSetter{},
			panicExpected: true,
			expVals:       []string{"test: EnumListSetter " + nilValueMsg},
		},
		{
			name: "EnumListSetter - bad - no allowedValues",
			s: &psetter.EnumListSetter{
				Value: &emptyStrList,
			},
			panicExpected: true,
			expVals: append([]string{"test: EnumListSetter "},
				tooFewAValsMsg...),
		},
		{
			name: "EnumListSetter - bad - empty allowedValues",
			s: &psetter.EnumListSetter{
				Value:       &emptyStrList,
				AllowedVals: psetter.AValMap{},
			},
			panicExpected: true,
			expVals: append([]string{"test: EnumListSetter "},
				tooFewAValsMsg...),
		},
		{
			name: "EnumListSetter - bad - allowedValues - only one entry",
			s: &psetter.EnumListSetter{
				Value: &emptyStrList,
				AllowedVals: psetter.AValMap{
					"aval": "desc",
				},
			},
			panicExpected: true,
			expVals: append([]string{"test: EnumListSetter "},
				tooFewAValsMsg...),
		},
		{
			name: "EnumMapSetter - ok",
			s: &psetter.EnumMapSetter{
				Value: &strToBoolMap,
				AllowedVals: psetter.AValMap{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
			panicExpected: false,
		},
		{
			name:          "EnumMapSetter - bad - no value",
			s:             &psetter.EnumMapSetter{},
			panicExpected: true,
			expVals:       []string{"test: EnumMapSetter " + nilValueMsg},
		},
		{
			name: "EnumMapSetter - bad - nil map",
			s: &psetter.EnumMapSetter{
				Value: &strToBoolMapNil,
				AllowedVals: psetter.AValMap{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
			panicExpected: true,
			expVals:       []string{"test: EnumMapSetter " + mapNotCreatedMsg},
		},
		{
			name: "EnumMapSetter - bad - no allowedValues",
			s: &psetter.EnumMapSetter{
				Value: &strToBoolMap,
			},
			panicExpected: true,
			expVals: append([]string{"test: EnumMapSetter "},
				tooFewAValsMsg...),
		},
		{
			name: "EnumMapSetter - bad - empty allowedValues",
			s: &psetter.EnumMapSetter{
				Value:       &strToBoolMap,
				AllowedVals: psetter.AValMap{},
			},
			panicExpected: true,
			expVals: append([]string{"test: EnumMapSetter "},
				tooFewAValsMsg...),
		},
		{
			name: "EnumSetter - ok",
			s: &psetter.EnumSetter{
				Value: &goodStr,
				AllowedVals: psetter.AValMap{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
			panicExpected: false,
		},
		{
			name: "EnumSetter - bad initial value",
			s: &psetter.EnumSetter{
				Value: &badStr,
				AllowedVals: psetter.AValMap{
					"aval":     "desc",
					"aval-alt": "desc",
				},
			},
			panicExpected: true,
			expVals:       append([]string{"test: EnumSetter "}, badInitialVal...),
		},
		{
			name:          "EnumSetter - bad - no value",
			s:             &psetter.EnumSetter{},
			panicExpected: true,
			expVals:       []string{"test: EnumSetter " + nilValueMsg},
		},
		{
			name: "EnumSetter - bad - no allowedValues",
			s: &psetter.EnumSetter{
				Value: &anyStr,
			},
			panicExpected: true,
			expVals: append([]string{"test: EnumSetter "},
				tooFewAValsMsg...),
		},
		{
			name: "EnumSetter - bad - empty allowedValues",
			s: &psetter.EnumSetter{
				Value:       &anyStr,
				AllowedVals: psetter.AValMap{},
			},
			panicExpected: true,
			expVals: append([]string{"test: EnumSetter "},
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
			name: "MapSetter - ok",
			s: &psetter.MapSetter{
				Value: &strToBoolMap,
			},
			panicExpected: false,
		},
		{
			name:          "MapSetter - bad - no value",
			s:             &psetter.MapSetter{},
			panicExpected: true,
			expVals:       []string{"test: MapSetter " + nilValueMsg},
		},
		{
			name: "MapSetter - bad - nil map",
			s: &psetter.MapSetter{
				Value: &strToBoolMapNil,
			},
			panicExpected: true,
			expVals:       []string{"test: MapSetter " + mapNotCreatedMsg},
		},
		{
			name:          "NilSetter - ok",
			s:             &psetter.NilSetter{},
			panicExpected: false,
		},
		{
			name:          "PathnameSetter - ok",
			s:             &psetter.PathnameSetter{Value: &anyStr},
			panicExpected: false,
		},
		{
			name:          "PathnameSetter - bad",
			s:             &psetter.PathnameSetter{},
			panicExpected: true,
			expVals:       []string{"test: PathnameSetter " + nilValueMsg},
		},
		{
			name:          "RegexpSetter - ok",
			s:             &psetter.RegexpSetter{Value: &re},
			panicExpected: false,
		},
		{
			name:          "RegexpSetter - bad",
			s:             &psetter.RegexpSetter{},
			panicExpected: true,
			expVals:       []string{"test: RegexpSetter " + nilValueMsg},
		},
		{
			name:          "StrListSetter - ok",
			s:             &psetter.StrListSetter{Value: &emptyStrList},
			panicExpected: false,
		},
		{
			name:          "StrListSetter - bad",
			s:             &psetter.StrListSetter{},
			panicExpected: true,
			expVals:       []string{"test: StrListSetter " + nilValueMsg},
		},
		{
			name:          "StringSetter - ok",
			s:             &psetter.StringSetter{Value: &anyStr},
			panicExpected: false,
		},
		{
			name:          "StringSetter - bad",
			s:             &psetter.StringSetter{},
			panicExpected: true,
			expVals:       []string{"test: StringSetter " + nilValueMsg},
		},
		{
			name:          "TimeLocationSetter - ok",
			s:             &psetter.TimeLocationSetter{Value: &timeLoc},
			panicExpected: false,
		},
		{
			name:          "TimeLocationSetter - bad",
			s:             &psetter.TimeLocationSetter{},
			panicExpected: true,
			expVals:       []string{"test: TimeLocationSetter " + nilValueMsg},
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
