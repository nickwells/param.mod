package psetter_test

import (
	"testing"

	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/paramtest"
	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

const (
	updFlagNameEnum     = "upd-gf-Enum"
	keepBadFlagNameEnum = "keep-bad-Enum"
)

var commonGFCEnum = testhelper.GoldenFileCfg{
	DirNames:               []string{"testdata", "Enum"},
	Pfx:                    "gf",
	Sfx:                    "txt",
	UpdFlagName:            updFlagNameEnum,
	KeepBadResultsFlagName: keepBadFlagNameEnum,
}

func init() {
	commonGFCEnum.AddUpdateFlag()
	commonGFCEnum.AddKeepBadResultsFlag()
}

func TestSetterEnum(t *testing.T) {
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
		goodAlias: []string{allowedVal2},
	}

	var (
		v1 = allowedVal1
		v2 = allowedVal1
		v3 = allowedVal1
		v4 = allowedVal1
	)

	testCases := []paramtest.Setter{
		{
			ID: testhelper.MkID("good-setter"),
			PSetter: psetter.Enum[string]{
				Value:       &v1,
				AllowedVals: allowedVals,
			},
			ParamVal: allowedVal3,
		},
		{
			ID: testhelper.MkID("good-setter-bad-val"),
			PSetter: psetter.Enum[string]{
				Value:       &v2,
				AllowedVals: allowedVals,
			},
			ParamVal: badVal,
			SetWithValErr: testhelper.MkExpErr(
				`value is not allowed: "` + badVal + `"`),
		},
		{
			ID: testhelper.MkID("good-setter-empty-val"),
			PSetter: psetter.Enum[string]{
				Value:       &v2,
				AllowedVals: allowedVals,
			},
			ParamVal: "",
			SetWithValErr: testhelper.MkExpErr(
				`value is not allowed: ""`),
		},
		{
			ID: testhelper.MkID("good-setter-alias-val"),
			PSetter: psetter.Enum[string]{
				Value:       &v3,
				AllowedVals: allowedVals,
				Aliases:     goodAliases,
			},
			ParamVal: goodAlias,
		},
		{
			ID: testhelper.MkID("bad-setter-nil-value"),
			ExpPanic: testhelper.MkExpPanic(dfltParamName +
				": psetter.Enum[string]" +
				" Check failed: the Value to be set is nil"),
			PSetter: psetter.Enum[string]{
				AllowedVals: allowedVals,
			},
		},
		{
			ID: testhelper.MkID("bad-setter-no-avals"),
			ExpPanic: testhelper.MkExpPanic(dfltParamName +
				": psetter.Enum[string]" +
				" Check failed: the map of allowed values has no entries." +
				" It should have at least 2"),
			PSetter: psetter.Enum[string]{
				Value: &v4,
			},
		},
		{
			ID: testhelper.MkID("bad-setter-one-aval"),
			ExpPanic: testhelper.MkExpPanic(dfltParamName +
				": psetter.Enum[string]" +
				" Check failed: the map of allowed values has only 1 entry." +
				" It should have at least 2"),
			PSetter: psetter.Enum[string]{
				AllowedVals: badAdValsOneEntry,
				Value:       &v4,
			},
		},
		{
			ID: testhelper.MkID("bad-setter-bad-aliases"),
			ExpPanic: testhelper.MkExpPanic(dfltParamName +
				": psetter.Enum[string]" +
				" Check failed: " +
				`bad alias: "` + badAlias + `":` +
				` []string{"bad-value"} - ` +
				`"bad-value" (at index 0) is unknown`),
			PSetter: psetter.Enum[string]{
				AllowedVals: allowedVals,
				Value:       &v4,
				Aliases:     badAliases,
			},
		},
	}

	for _, tc := range testCases {
		f := func(t *testing.T) {
			tc.GFC = commonGFCEnum

			if tc.ParamName == "" {
				tc.ParamName = dfltParamName
			}

			tc.SetVR(param.Mandatory)

			tc.Test(t)
		}

		t.Run(tc.IDStr(), f)
	}
}
