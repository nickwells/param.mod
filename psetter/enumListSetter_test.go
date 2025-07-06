package psetter_test

import (
	"testing"

	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/paramtest"
	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

const (
	updFlagNameEnumList     = "upd-gf-EnumList"
	keepBadFlagNameEnumList = "keep-bad-EnumList"
)

var commonGFCEnumList = testhelper.GoldenFileCfg{
	DirNames:               []string{"testdata", "EnumList"},
	Pfx:                    "gf",
	Sfx:                    "txt",
	UpdFlagName:            updFlagNameEnumList,
	KeepBadResultsFlagName: keepBadFlagNameEnumList,
}

func init() {
	commonGFCEnumList.AddUpdateFlag()
	commonGFCEnumList.AddKeepBadResultsFlag()
}

func TestSetterEnumList(t *testing.T) {
	const dfltParamName = "param-name"

	const (
		allowedVal1 = "av1"
		allowedVal2 = "av2"
		allowedVal3 = "av3"

		badVal = "bad-value"

		goodAlias = "good-alias"
		badAlias  = "noCorrespondingValue"
	)

	allowedVals := psetter.AllowedVals[string]{
		allowedVal1: "notes for av1",
		allowedVal2: "notes for av2",
		allowedVal3: "notes for av3",
	}
	badAdValsOneEntry := psetter.AllowedVals[string]{
		allowedVal1: "notes for av1",
	}

	badAliases := psetter.Aliases[string]{
		badAlias: []string{badVal},
	}

	goodAliases := psetter.Aliases[string]{
		goodAlias: []string{allowedVal1},
	}

	var (
		l1 = []string{allowedVal1, allowedVal2}
		l2 = []string{allowedVal1, allowedVal2}
		l3 = []string{allowedVal1, allowedVal2}
		l4 = []string{allowedVal1, allowedVal2}
		l5 = []string{allowedVal1, allowedVal2}
	)

	testCases := []paramtest.Setter{
		{
			ID: testhelper.MkID("good-setter"),
			PSetter: psetter.EnumList[string]{
				Value:       &l1,
				AllowedVals: allowedVals,
			},
			ParamVal: allowedVal3,
		},
		{
			ID: testhelper.MkID("good-setter-bad-val"),
			PSetter: psetter.EnumList[string]{
				Value:       &l2,
				AllowedVals: allowedVals,
			},
			ParamVal: badVal,
			SetWithValErr: testhelper.MkExpErr(
				`value is not allowed: "` + badVal + `"`),
		},
		{
			ID: testhelper.MkID("good-setter-empty-val"),
			PSetter: psetter.EnumList[string]{
				Value:       &l2,
				AllowedVals: allowedVals,
			},
			ParamVal: "",
			SetWithValErr: testhelper.MkExpErr(
				`value is not allowed: ""`),
		},
		{
			ID: testhelper.MkID("good-setter-multi-vals"),
			PSetter: psetter.EnumList[string]{
				Value:       &l3,
				AllowedVals: allowedVals,
			},
			ParamVal: allowedVal1 + "," + allowedVal3,
		},
		{
			ID: testhelper.MkID("good-setter-alias-val"),
			PSetter: psetter.EnumList[string]{
				Value:       &l4,
				AllowedVals: allowedVals,
				Aliases:     goodAliases,
			},
			ParamVal: goodAlias,
		},
		{
			ID: testhelper.MkID("bad-setter-nil-value"),
			ExpPanic: testhelper.MkExpPanic(dfltParamName +
				": psetter.EnumList[string]" +
				" Check failed: the Value to be set is nil"),
			PSetter: psetter.EnumList[string]{
				AllowedVals: allowedVals,
			},
		},
		{
			ID: testhelper.MkID("bad-setter-no-avals"),
			ExpPanic: testhelper.MkExpPanic(dfltParamName +
				": psetter.EnumList[string]" +
				" Check failed: the map of allowed values has no entries." +
				" It should have at least 2"),
			PSetter: psetter.EnumList[string]{
				Value: &l5,
			},
		},
		{
			ID: testhelper.MkID("bad-setter-one-aval"),
			ExpPanic: testhelper.MkExpPanic(dfltParamName +
				": psetter.EnumList[string]" +
				" Check failed: the map of allowed values has only 1 entry." +
				" It should have at least 2"),
			PSetter: psetter.EnumList[string]{
				AllowedVals: badAdValsOneEntry,
				Value:       &l5,
			},
		},
		{
			ID: testhelper.MkID("bad-setter-bad-aliases"),
			ExpPanic: testhelper.MkExpPanic(dfltParamName +
				": psetter.EnumList[string]" +
				" Check failed: " +
				`bad alias: "` + badAlias + `":` +
				` []string{"bad-value"} - ` +
				`"bad-value" (at index 0) is unknown`),
			PSetter: psetter.EnumList[string]{
				AllowedVals: allowedVals,
				Value:       &l5,
				Aliases:     badAliases,
			},
		},
	}

	for _, tc := range testCases {
		f := func(t *testing.T) {
			tc.GFC = commonGFCEnumList

			if tc.ParamName == "" {
				tc.ParamName = dfltParamName
			}

			tc.SetVR(param.Mandatory)

			tc.Test(t)
		}

		t.Run(tc.IDStr(), f)
	}
}
