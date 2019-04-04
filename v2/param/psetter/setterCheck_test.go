package psetter_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/nickwells/param.mod/v2/param"
	"github.com/nickwells/param.mod/v2/param/psetter"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestCheck(t *testing.T) {
	var goodStr = "aval"
	var goodStrAlt = "aval-alt"
	var badStr = "bad val"
	var anyStr = ""
	var b bool
	var dur time.Duration
	var emptyStrList []string
	var goodStrList = []string{goodStr}
	var badStrList = []string{badStr}
	var strToBoolMap = make(map[string]bool)
	var strToBoolMapNil map[string]bool
	var strToBoolMapWithEntriesGood = map[string]bool{goodStr: true}
	var strToBoolMapWithEntriesBad = map[string]bool{
		goodStr: true,
		badStr:  true,
	}
	var f float64
	var i int64
	var intList []int64
	var re *regexp.Regexp
	var timeLoc *time.Location
	var avalMapEmpty = psetter.AValMap{}
	var avalMapOneEntry = psetter.AValMap{goodStr: "desc"}
	var avalMapGood = psetter.AValMap{
		goodStr:    "desc",
		goodStrAlt: "desc",
	}

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
		testhelper.ID
		testhelper.ExpPanic
		s param.Setter
	}{
		{
			ID: testhelper.MkID("Bool - ok"),
			s:  &psetter.Bool{Value: &b},
		},
		{
			ID: testhelper.MkID("Bool - bad"),
			s:  &psetter.Bool{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.Bool " +
				nilValueMsg),
		},
		{
			ID: testhelper.MkID("Duration - ok"),
			s:  &psetter.Duration{Value: &dur},
		},
		{
			ID: testhelper.MkID("Duration - bad"),
			s:  &psetter.Duration{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.Duration " +
				nilValueMsg),
		},
		{
			ID: testhelper.MkID("EnumList - ok - no strings"),
			s: &psetter.EnumList{
				Value:       &emptyStrList,
				AllowedVals: avalMapGood,
			},
		},
		{
			ID: testhelper.MkID("EnumList - ok - good strings"),
			s: &psetter.EnumList{
				Value:       &goodStrList,
				AllowedVals: avalMapGood,
			},
		},
		{
			ID: testhelper.MkID("EnumList - bad initial value"),
			s: &psetter.EnumList{
				Value:       &badStrList,
				AllowedVals: avalMapGood,
			},
			ExpPanic: testhelper.MkExpPanic(
				append([]string{"test: psetter.EnumList "},
					badInitialList...)...),
		},
		{
			ID: testhelper.MkID("EnumList - bad - no value"),
			s:  &psetter.EnumList{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.EnumList " +
				nilValueMsg),
		},
		{
			ID: testhelper.MkID("EnumList - bad - no allowedValues"),
			s: &psetter.EnumList{
				Value: &emptyStrList,
			},
			ExpPanic: testhelper.MkExpPanic(
				append([]string{"test: psetter.EnumList "},
					tooFewAValsMsg...)...),
		},
		{
			ID: testhelper.MkID("EnumList - bad - empty allowedValues"),
			s: &psetter.EnumList{
				Value:       &emptyStrList,
				AllowedVals: avalMapEmpty,
			},
			ExpPanic: testhelper.MkExpPanic(
				append([]string{"test: psetter.EnumList "},
					tooFewAValsMsg...)...),
		},
		{
			ID: testhelper.MkID("EnumList - bad - allowedValues" +
				" - only one entry"),
			s: &psetter.EnumList{
				Value:       &emptyStrList,
				AllowedVals: avalMapOneEntry,
			},
			ExpPanic: testhelper.MkExpPanic(
				append([]string{"test: psetter.EnumList "},
					tooFewAValsMsg...)...),
		},
		{
			ID: testhelper.MkID("EnumMap - ok"),
			s: &psetter.EnumMap{
				Value:       &strToBoolMap,
				AllowedVals: avalMapGood,
			},
		},
		{
			ID: testhelper.MkID("EnumMap - ok - Values has entries"),
			s: &psetter.EnumMap{
				Value:       &strToBoolMapWithEntriesGood,
				AllowedVals: avalMapGood,
			},
		},
		{
			ID: testhelper.MkID(
				"EnumMap - ok - Values has entries - missing allowed"),
			s: &psetter.EnumMap{
				Value:                 &strToBoolMapWithEntriesBad,
				AllowedVals:           avalMapGood,
				AllowHiddenMapEntries: true,
			},
		},
		{
			ID: testhelper.MkID(
				"EnumMap - bad - Values has entries - missing not allowed"),
			s: &psetter.EnumMap{
				Value:       &strToBoolMapWithEntriesBad,
				AllowedVals: avalMapGood,
			},
			ExpPanic: testhelper.MkExpPanic(
				append([]string{"test: psetter.EnumMap"},
					"the map entry with key",
					"is invalid - it is not in the allowed values map")...),
		},
		{
			ID: testhelper.MkID("EnumMap - bad - no value"),
			s:  &psetter.EnumMap{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.EnumMap " +
				nilValueMsg),
		},
		{
			ID: testhelper.MkID("EnumMap - bad - nil map"),
			s: &psetter.EnumMap{
				Value:       &strToBoolMapNil,
				AllowedVals: avalMapGood,
			},
			ExpPanic: testhelper.MkExpPanic("test: psetter.EnumMap " +
				mapNotCreatedMsg),
		},
		{
			ID: testhelper.MkID("EnumMap - bad - no allowedValues"),
			s: &psetter.EnumMap{
				Value: &strToBoolMap,
			},
			ExpPanic: testhelper.MkExpPanic(
				append([]string{"test: psetter.EnumMap "},
					tooFewAValsMsg...)...),
		},
		{
			ID: testhelper.MkID("EnumMap - bad - empty allowedValues"),
			s: &psetter.EnumMap{
				Value:       &strToBoolMap,
				AllowedVals: avalMapEmpty,
			},
			ExpPanic: testhelper.MkExpPanic(
				append([]string{"test: psetter.EnumMap "},
					tooFewAValsMsg...)...),
		},
		{
			ID: testhelper.MkID("Enum - ok"),
			s: &psetter.Enum{
				Value:       &goodStr,
				AllowedVals: avalMapGood,
			},
		},
		{
			ID: testhelper.MkID("Enum - bad initial value"),
			s: &psetter.Enum{
				Value:       &badStr,
				AllowedVals: avalMapGood,
			},
			ExpPanic: testhelper.MkExpPanic(
				append([]string{"test: psetter.Enum "}, badInitialVal...)...),
		},
		{
			ID: testhelper.MkID("Enum - bad - no value"),
			s:  &psetter.Enum{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.Enum " +
				nilValueMsg),
		},
		{
			ID: testhelper.MkID("Enum - bad - no allowedValues"),
			s: &psetter.Enum{
				Value: &anyStr,
			},
			ExpPanic: testhelper.MkExpPanic(
				append([]string{"test: psetter.Enum "}, tooFewAValsMsg...)...),
		},
		{
			ID: testhelper.MkID("Enum - bad - empty allowedValues"),
			s: &psetter.Enum{
				Value:       &anyStr,
				AllowedVals: avalMapEmpty,
			},
			ExpPanic: testhelper.MkExpPanic(
				append([]string{"test: psetter.Enum "}, tooFewAValsMsg...)...),
		},
		{
			ID: testhelper.MkID("Float64 - ok"),
			s:  &psetter.Float64{Value: &f},
		},
		{
			ID: testhelper.MkID("Float64 - bad"),
			s:  &psetter.Float64{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.Float64 " +
				nilValueMsg),
		},
		{
			ID: testhelper.MkID("Int64 - ok"),
			s:  &psetter.Int64{Value: &i},
		},
		{
			ID: testhelper.MkID("Int64 - bad"),
			s:  &psetter.Int64{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.Int64 " +
				nilValueMsg),
		},
		{
			ID: testhelper.MkID("Int64List - ok"),
			s:  &psetter.Int64List{Value: &intList},
		},
		{
			ID: testhelper.MkID("Int64List - bad"),
			s:  &psetter.Int64List{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.Int64List " +
				nilValueMsg),
		},
		{
			ID: testhelper.MkID("Map - ok"),
			s: &psetter.Map{
				Value: &strToBoolMap,
			},
		},
		{
			ID: testhelper.MkID("Map - bad - no value"),
			s:  &psetter.Map{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.Map " +
				nilValueMsg),
		},
		{
			ID: testhelper.MkID("Map - bad - nil map"),
			s: &psetter.Map{
				Value: &strToBoolMapNil,
			},
			ExpPanic: testhelper.MkExpPanic("test: psetter.Map " +
				mapNotCreatedMsg),
		},
		{
			ID: testhelper.MkID("Nil - ok"),
			s:  &psetter.Nil{},
		},
		{
			ID: testhelper.MkID("Pathname - ok"),
			s:  &psetter.Pathname{Value: &anyStr},
		},
		{
			ID: testhelper.MkID("Pathname - bad"),
			s:  &psetter.Pathname{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.Pathname " +
				nilValueMsg),
		},
		{
			ID: testhelper.MkID("Regexp - ok"),
			s:  &psetter.Regexp{Value: &re},
		},
		{
			ID: testhelper.MkID("Regexp - bad"),
			s:  &psetter.Regexp{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.Regexp " +
				nilValueMsg),
		},
		{
			ID: testhelper.MkID("StrList - ok"),
			s:  &psetter.StrList{Value: &emptyStrList},
		},
		{
			ID: testhelper.MkID("StrList - bad"),
			s:  &psetter.StrList{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.StrList " +
				nilValueMsg),
		},
		{
			ID: testhelper.MkID("String - ok"),
			s:  &psetter.String{Value: &anyStr},
		},
		{
			ID: testhelper.MkID("String - bad"),
			s:  &psetter.String{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.String " +
				nilValueMsg),
		},
		{
			ID: testhelper.MkID("TimeLocation - ok"),
			s:  &psetter.TimeLocation{Value: &timeLoc},
		},
		{
			ID: testhelper.MkID("TimeLocation - bad"),
			s:  &psetter.TimeLocation{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.TimeLocation " +
				nilValueMsg),
		},
	}

	for _, tc := range testCases {
		panicked, panicVal := panicSafeCheck(tc.s)
		testhelper.CheckExpPanic(t, panicked, panicVal, tc)
	}
}
