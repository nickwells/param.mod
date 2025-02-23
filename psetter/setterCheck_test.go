package psetter_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestCheck(t *testing.T) {
	goodStr := "aval"
	goodStrAlt := "aval-alt"
	badStr := "bad val"
	anyStr := ""

	var b bool

	var dur time.Duration

	var emptyStrList []string

	goodStrList := []string{goodStr}
	badStrList := []string{badStr}
	strToBoolMap := make(map[string]bool)

	var strToBoolMapNil1 map[string]bool

	var strToBoolMapNil2 map[string]bool

	strToBoolMapWithEntriesGood := map[string]bool{goodStr: true}
	strToBoolMapWithEntriesBad := map[string]bool{
		goodStr: true,
		badStr:  true,
	}

	var f float64

	var i int64

	var intList []int64

	var re *regexp.Regexp

	var timeLoc *time.Location

	avalMapEmpty := psetter.AllowedVals[string]{}
	avalMapOneEntry := psetter.AllowedVals[string]{goodStr: "desc"}
	avalMapGood := psetter.AllowedVals[string]{
		goodStr:    "desc",
		goodStrAlt: "desc",
	}

	nilValueMsg := "Check failed: the Value to be set is nil"
	tooFewAValsMsg := []string{
		"Check failed: the map of allowed values has ",
		"It should have at least 2",
	}
	badInitialVal := []string{
		"Check failed: the initial value",
		"is not valid",
	}
	badInitialList := []string{
		"Check failed: element",
		"in the current list of entries is invalid",
	}

	nc := psetter.NamedCalc[int64]{
		Name: "name",
		Calc: func(_, _ string) (int64, error) {
			return 42, nil
		},
	}

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
			s: &psetter.EnumList[string]{
				Value:       &emptyStrList,
				AllowedVals: avalMapGood,
			},
		},
		{
			ID: testhelper.MkID("EnumList - ok - good strings"),
			s: &psetter.EnumList[string]{
				Value:       &goodStrList,
				AllowedVals: avalMapGood,
			},
		},
		{
			ID: testhelper.MkID("EnumList - bad initial value"),
			s: &psetter.EnumList[string]{
				Value:       &badStrList,
				AllowedVals: avalMapGood,
			},
			ExpPanic: testhelper.MkExpPanic(
				append([]string{"test: psetter.EnumList[string] "},
					badInitialList...)...),
		},
		{
			ID: testhelper.MkID("EnumList - bad - no value"),
			s:  &psetter.EnumList[string]{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.EnumList[string] " +
				nilValueMsg),
		},
		{
			ID: testhelper.MkID("EnumList - bad - no allowedValues"),
			s: &psetter.EnumList[string]{
				Value: &emptyStrList,
			},
			ExpPanic: testhelper.MkExpPanic(
				append([]string{"test: psetter.EnumList[string] "},
					tooFewAValsMsg...)...),
		},
		{
			ID: testhelper.MkID("EnumList - bad - empty allowedValues"),
			s: &psetter.EnumList[string]{
				Value:       &emptyStrList,
				AllowedVals: avalMapEmpty,
			},
			ExpPanic: testhelper.MkExpPanic(
				append([]string{"test: psetter.EnumList[string] "},
					tooFewAValsMsg...)...),
		},
		{
			ID: testhelper.MkID("EnumList - bad - allowedValues" +
				" - only one entry"),
			s: &psetter.EnumList[string]{
				Value:       &emptyStrList,
				AllowedVals: avalMapOneEntry,
			},
			ExpPanic: testhelper.MkExpPanic(
				append([]string{"test: psetter.EnumList[string] "},
					tooFewAValsMsg...)...),
		},
		{
			ID: testhelper.MkID("EnumMap - ok"),
			s: &psetter.EnumMap[string]{
				Value:       &strToBoolMap,
				AllowedVals: avalMapGood,
			},
		},
		{
			ID: testhelper.MkID("EnumMap - ok - Values has entries"),
			s: &psetter.EnumMap[string]{
				Value:       &strToBoolMapWithEntriesGood,
				AllowedVals: avalMapGood,
			},
		},
		{
			ID: testhelper.MkID(
				"EnumMap - ok - Values has entries - missing allowed"),
			s: &psetter.EnumMap[string]{
				Value:                 &strToBoolMapWithEntriesBad,
				AllowedVals:           avalMapGood,
				AllowHiddenMapEntries: true,
			},
		},
		{
			ID: testhelper.MkID(
				"EnumMap - bad - Values has entries - missing not allowed"),
			s: &psetter.EnumMap[string]{
				Value:       &strToBoolMapWithEntriesBad,
				AllowedVals: avalMapGood,
			},
			ExpPanic: testhelper.MkExpPanic(
				append([]string{"test: psetter.EnumMap[string]"},
					"the map entry with key",
					"is invalid - it is not in the allowed values map")...),
		},
		{
			ID: testhelper.MkID("EnumMap - bad - no value"),
			s:  &psetter.EnumMap[string]{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.EnumMap[string] " +
				nilValueMsg),
		},
		{
			ID: testhelper.MkID("EnumMap - good - nil map (created)"),
			s: &psetter.EnumMap[string]{
				Value:       &strToBoolMapNil1,
				AllowedVals: avalMapGood,
			},
		},
		{
			ID: testhelper.MkID("EnumMap - bad - no allowedValues"),
			s: &psetter.EnumMap[string]{
				Value: &strToBoolMap,
			},
			ExpPanic: testhelper.MkExpPanic(
				append([]string{"test: psetter.EnumMap[string] "},
					tooFewAValsMsg...)...),
		},
		{
			ID: testhelper.MkID("EnumMap - bad - empty allowedValues"),
			s: &psetter.EnumMap[string]{
				Value:       &strToBoolMap,
				AllowedVals: avalMapEmpty,
			},
			ExpPanic: testhelper.MkExpPanic(
				append([]string{"test: psetter.EnumMap[string] "},
					tooFewAValsMsg...)...),
		},
		{
			ID: testhelper.MkID("Enum - ok"),
			s: &psetter.Enum[string]{
				Value:       &goodStr,
				AllowedVals: avalMapGood,
			},
		},
		{
			ID: testhelper.MkID("Enum - bad initial value"),
			s: &psetter.Enum[string]{
				Value:       &badStr,
				AllowedVals: avalMapGood,
			},
			ExpPanic: testhelper.MkExpPanic(
				append([]string{"test: psetter.Enum[string] "},
					badInitialVal...)...),
		},
		{
			ID: testhelper.MkID("Enum - bad - no value"),
			s:  &psetter.Enum[string]{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.Enum[string] " +
				nilValueMsg),
		},
		{
			ID: testhelper.MkID("Enum - bad - no allowedValues"),
			s: &psetter.Enum[string]{
				Value: &anyStr,
			},
			ExpPanic: testhelper.MkExpPanic(
				append([]string{"test: psetter.Enum[string] "},
					tooFewAValsMsg...)...),
		},
		{
			ID: testhelper.MkID("Enum - bad - empty allowedValues"),
			s: &psetter.Enum[string]{
				Value:       &anyStr,
				AllowedVals: avalMapEmpty,
			},
			ExpPanic: testhelper.MkExpPanic(
				append([]string{"test: psetter.Enum[string] "},
					tooFewAValsMsg...)...),
		},
		{
			ID: testhelper.MkID("Float64 - ok"),
			s:  &psetter.Float[float64]{Value: &f},
		},
		{
			ID: testhelper.MkID("Float64 - bad"),
			s:  &psetter.Float[float64]{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.Float[float64] " +
				nilValueMsg),
		},
		{
			ID: testhelper.MkID("Int64 - ok"),
			s:  &psetter.Int[int64]{Value: &i},
		},
		{
			ID: testhelper.MkID("Int64 - bad"),
			s:  &psetter.Int[int64]{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.Int[int64] " +
				nilValueMsg),
		},
		{
			ID: testhelper.MkID("Int64List - ok"),
			s:  &psetter.IntList[int64]{Value: &intList},
		},
		{
			ID: testhelper.MkID("Int64List - bad"),
			s:  &psetter.IntList[int64]{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.IntList[int64] " +
				nilValueMsg),
		},
		{
			ID: testhelper.MkID("Map - ok"),
			s: &psetter.Map[string]{
				Value: &strToBoolMap,
			},
		},
		{
			ID: testhelper.MkID("Map - bad - no value"),
			s:  &psetter.Map[string]{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.Map[string] " +
				nilValueMsg),
		},
		{
			ID: testhelper.MkID("Map - good - nil map"),
			s: &psetter.Map[string]{
				Value: &strToBoolMapNil2,
			},
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
			s:  &psetter.StrList[string]{Value: &emptyStrList},
		},
		{
			ID: testhelper.MkID("StrList - bad"),
			s:  &psetter.StrList[string]{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.StrList[string] " +
				nilValueMsg),
		},
		{
			ID: testhelper.MkID("String - ok"),
			s:  &psetter.String[string]{Value: &anyStr},
		},
		{
			ID: testhelper.MkID("String - bad"),
			s:  &psetter.String[string]{},
			ExpPanic: testhelper.MkExpPanic("test: psetter.String[string] " +
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
		{
			ID: testhelper.MkID("Calculated - ok, with default"),
			s: &psetter.Calculated[int64]{
				Value: &i,
				CalcMap: map[string]psetter.NamedCalc[int64]{
					"s1": nc,
				},
				Default: nc,
			},
		},
		{
			ID: testhelper.MkID("Calculated - ok, no default"),
			s: &psetter.Calculated[int64]{
				Value: &i,
				CalcMap: map[string]psetter.NamedCalc[int64]{
					"s1": nc,
					"s2": nc,
				},
				NoDefault: true,
			},
		},
		{
			ID:       testhelper.MkID("Calculated - bad, empty CalcMap"),
			ExpPanic: testhelper.MkExpPanic("the CalcMap cannot be empty"),
			s: &psetter.Calculated[int64]{
				Value: &i,
			},
		},
		{
			ID: testhelper.MkID(
				"Calculated - bad, no default, len CalcMap < 2"),
			ExpPanic: testhelper.MkExpPanic(
				"with no default value the CalcMap must have" +
					" more than one entry"),
			s: &psetter.Calculated[int64]{
				Value: &i,
				CalcMap: map[string]psetter.NamedCalc[int64]{
					"s1": nc,
				},
				NoDefault: true,
			},
		},
		{
			ID: testhelper.MkID("Calculated - bad - bad default Name"),
			ExpPanic: testhelper.MkExpPanic(
				"the default NamedCalc is invalid:" +
					" the Name must not be empty"),
			s: &psetter.Calculated[int64]{
				Value: &i,
				CalcMap: map[string]psetter.NamedCalc[int64]{
					"s1": nc,
				},
			},
		},
		{
			ID: testhelper.MkID("Calculated - bad - bad default Calc"),
			ExpPanic: testhelper.MkExpPanic(
				"the default NamedCalc is invalid:" +
					" the Calc must not be nil"),
			s: &psetter.Calculated[int64]{
				Value: &i,
				CalcMap: map[string]psetter.NamedCalc[int64]{
					"s1": nc,
				},
				Default: psetter.NamedCalc[int64]{Name: "name"},
			},
		},
	}

	for _, tc := range testCases {
		panicked, panicVal := panicSafeCheck(tc.s)

		testhelper.CheckExpPanic(t, panicked, panicVal, tc)
	}
}
