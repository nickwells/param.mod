package psetter_test

import (
	"testing"

	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/paramtest"
	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

const (
	updFlagNameEnumMap     = "upd-gf-EnumMap"
	keepBadFlagNameEnumMap = "keep-bad-EnumMap"
)

var commonEnumMapGFC = testhelper.GoldenFileCfg{
	DirNames:               []string{"testdata", "Setters", "EnumMap"},
	Pfx:                    "gf",
	Sfx:                    "txt",
	UpdFlagName:            updFlagNameEnumMap,
	KeepBadResultsFlagName: keepBadFlagNameEnumMap,
}

func init() {
	commonEnumMapGFC.AddUpdateFlag()
	commonEnumMapGFC.AddKeepBadResultsFlag()
}

func TestSetterEnumMap(t *testing.T) {
	var nilMap map[string]bool

	emptyMap := map[string]bool{}
	zvalMap := map[string]bool{
		"z": true,
	}

	dfltAllowedVals := psetter.AllowedVals[string]{
		"x": "x desc",
		"y": "y desc",
		"z": "z desc",
	}

	testCases := []paramtest.Setter{
		{
			ID: testhelper.MkID("no-val"),
			ExpPanic: testhelper.MkExpPanic(
				"psetter.EnumMap[string] Check failed: ",
				"the Value to be set is nil"),
			PSetter: psetter.EnumMap[string]{},
		},
		{
			ID: testhelper.MkID("no-allowed-vals"),
			PSetter: psetter.EnumMap[string]{
				Value: &nilMap,
			},
			ExpPanic: testhelper.MkExpPanic(
				"psetter.EnumMap[string] Check failed: ",
				"the map of allowed values has no entries.",
				" It should have at least 2"),
		},
		{
			ID: testhelper.MkID("only-one-allowed-val"),
			PSetter: psetter.EnumMap[string]{
				Value: &nilMap,
				AllowedVals: psetter.AllowedVals[string]{
					"x": "x-desc",
				},
			},
			ExpPanic: testhelper.MkExpPanic(
				"psetter.EnumMap[string] Check failed: ",
				"the map of allowed values has only 1 entry.",
				" It should have at least 2"),
		},
		{
			ID: testhelper.MkID("invalid-initial-entry"),
			PSetter: psetter.EnumMap[string]{
				Value: &zvalMap,
				AllowedVals: psetter.AllowedVals[string]{
					"x": "x-desc",
					"y": "y-desc",
				},
			},
			ExpPanic: testhelper.MkExpPanic(
				"psetter.EnumMap[string] Check failed: ",
				`the map entry with key "z" is invalid`,
				" - it is not in the allowed values map"),
		},
		{
			ID: testhelper.MkID("bad-alias-equals-aval"),
			PSetter: psetter.EnumMap[string]{
				Value: &zvalMap,
				AllowedVals: psetter.AllowedVals[string]{
					"x": "x-desc",
					"y": "y-desc",
				},
				Aliases: psetter.Aliases[string]{
					"x": []string{"y"},
				},
			},
			ExpPanic: testhelper.MkExpPanic(
				"psetter.EnumMap[string] Check failed: " +
					`bad alias: "x": []string{"y"} -` +
					` an allowed value has the same name as the alias`),
		},
		{
			ID: testhelper.MkID("bad-alias-empty"),
			PSetter: psetter.EnumMap[string]{
				Value: &zvalMap,
				AllowedVals: psetter.AllowedVals[string]{
					"x": "x-desc",
					"y": "y-desc",
				},
				Aliases: psetter.Aliases[string]{
					"z": []string{},
				},
			},
			ExpPanic: testhelper.MkExpPanic(
				"psetter.EnumMap[string] Check failed: " +
					`bad alias: "z": []string{} - the alias maps to no values`),
		},
		{
			ID: testhelper.MkID("bad-alias-duplicate"),
			PSetter: psetter.EnumMap[string]{
				Value: &zvalMap,
				AllowedVals: psetter.AllowedVals[string]{
					"x": "x-desc",
					"y": "y-desc",
				},
				Aliases: psetter.Aliases[string]{
					"z": []string{"x", "x"},
				},
			},
			ExpPanic: testhelper.MkExpPanic(
				"psetter.EnumMap[string] Check failed: " +
					`bad alias: "z": []string{"x", "x"} -` +
					` "x" appears more than once (at index 0 and 1)`),
		},
		{
			ID: testhelper.MkID("bad-alias-not-an-aval"),
			PSetter: psetter.EnumMap[string]{
				Value: &zvalMap,
				AllowedVals: psetter.AllowedVals[string]{
					"x": "x-desc",
					"y": "y-desc",
				},
				Aliases: psetter.Aliases[string]{
					"z": []string{"a"},
				},
			},
			ExpPanic: testhelper.MkExpPanic(
				"psetter.EnumMap[string] Check failed: " +
					`bad alias: "z": []string{"a"} -` +
					` "a" (at index 0) is unknown`),
		},
		{
			ID: testhelper.MkID("hidden-initial-entry"),
			PSetter: psetter.EnumMap[string]{
				Value: &zvalMap,
				AllowedVals: psetter.AllowedVals[string]{
					"x": "x-desc",
					"y": "y-desc",
				},
				AllowHiddenMapEntries: true,
			},
			ParamVal: "x=true",
		},
		{
			ID: testhelper.MkID("nil-map-z"),
			PSetter: psetter.EnumMap[string]{
				Value:       &nilMap,
				AllowedVals: dfltAllowedVals,
			},
			ParamVal: "z",
			ExtraTest: func(t *testing.T, s paramtest.Setter) {
				t.Helper()
				if nilMap == nil {
					t.Log(s.IDStr())
					t.Error("\t: The map should have been initialised\n")
				}
			},
		},
		{
			ID: testhelper.MkID("zval-map-z-false"),
			PSetter: psetter.EnumMap[string]{
				Value:       &zvalMap,
				AllowedVals: dfltAllowedVals,
			},
			ParamVal: "z=false",
		},
		{
			ID: testhelper.MkID("zval-map-blank-val"),
			PSetter: psetter.EnumMap[string]{
				Value:       &zvalMap,
				AllowedVals: dfltAllowedVals,
			},
			SetWithValErr: testhelper.MkExpErr(
				"empty value. Some value must be given"),
		},
		{
			ID: testhelper.MkID("empty-map-z-true"),
			PSetter: psetter.EnumMap[string]{
				Value:       &emptyMap,
				AllowedVals: dfltAllowedVals,
			},
			ParamVal: "z",
		},
		{
			ID: testhelper.MkID("empty-map-z-bad"),
			PSetter: psetter.EnumMap[string]{
				Value:       &emptyMap,
				AllowedVals: dfltAllowedVals,
			},
			ParamVal: "z=bad",
			SetWithValErr: testhelper.MkExpErr(
				`bad value: "z=bad": part: 1 ("z=bad") is invalid.`,
				`The value ("bad") cannot be interpreted as true or false:`,
				`strconv.ParseBool: parsing "bad": invalid syntax`),
		},
		{
			ID: testhelper.MkID("empty-map-aliases"),
			PSetter: psetter.EnumMap[string]{
				Value:       &emptyMap,
				AllowedVals: dfltAllowedVals,
				Aliases: psetter.Aliases[string]{
					"all": []string{"x", "y", "z"},
					"xy":  []string{"x", "y"},
					"xz":  []string{"x", "z"},
					"yz":  []string{"y", "z"},
				},
			},
			ParamVal: "all=true",
		},
		{
			ID: testhelper.MkID("empty-map-z-true-sep-dot"),
			PSetter: psetter.EnumMap[string]{
				Value:            &emptyMap,
				AllowedVals:      dfltAllowedVals,
				StrListSeparator: psetter.StrListSeparator{Sep: "."},
			},
			ParamVal: "z",
		},
	}

	for _, tc := range testCases {
		tc.GFC = commonEnumMapGFC
		tc.ValDescriber = true

		if tc.ParamName == "" {
			tc.ParamName = "set-enum-map"
		}

		tc.SetVR(param.Mandatory)

		nilMap = nil
		emptyMap = map[string]bool{}
		zvalMap = map[string]bool{
			"z": true,
		}

		tc.Test(t)
	}
}
