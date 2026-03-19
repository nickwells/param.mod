package psetter_test

import (
	"testing"

	"github.com/nickwells/param.mod/v7/param"
	"github.com/nickwells/param.mod/v7/paramtest"
	"github.com/nickwells/param.mod/v7/psetter"
	"github.com/nickwells/param.mod/v7/ptypes"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

const (
	updFlagNameTaggedValueList     = "upd-gf-TaggedValueList"
	keepBadFlagNameTaggedValueList = "keep-bad-TaggedValueList"
)

var commonGFCTaggedValueList = testhelper.GoldenFileCfg{
	DirNames:               []string{"testdata", "TaggedValueList"},
	Pfx:                    "gf",
	Sfx:                    "txt",
	UpdFlagName:            updFlagNameTaggedValueList,
	KeepBadResultsFlagName: keepBadFlagNameTaggedValueList,
}

func init() {
	commonGFCTaggedValueList.AddUpdateFlag()
	commonGFCTaggedValueList.AddKeepBadResultsFlag()
}

func TestSetterTaggedValueList(t *testing.T) {
	const dfltParamName = "param-name"

	const (
		allowedVal1 = "av1"
		allowedVal2 = "av2"
		allowedVal3 = "av3"

		allowedTagVal1 = "tag1"
		allowedTagVal2 = "tag2"
		allowedTagVal3 = "tag3"

		badVal = "bad-value"

		goodAlias = "good-alias"
		badAlias  = "noCorrespondingValue"
	)

	allowedVals := ptypes.AllowedVals[string]{
		allowedVal1: "notes for av1",
		allowedVal2: "notes for av2",
		allowedVal3: "notes for av3",
	}

	tagAllowedVals := ptypes.AllowedVals[string]{
		allowedTagVal1: "notes for tag1",
		allowedTagVal2: "notes for tag2",
		allowedTagVal3: "notes for tag3",
	}

	badAdValsOneEntry := ptypes.AllowedVals[string]{
		allowedVal1: "notes for av1",
	}

	badAliases := ptypes.Aliases[string]{
		badAlias: []string{badVal},
	}

	goodAliases := ptypes.Aliases[string]{
		goodAlias: []string{allowedVal1},
	}

	var (
		l1 = []psetter.TaggedValue[string, string]{
			{Value: allowedVal1},
			{Value: allowedVal2},
		}
		l2 = []psetter.TaggedValue[string, string]{
			{Value: allowedVal1},
			{Value: allowedVal2},
		}
		l3 = []psetter.TaggedValue[string, string]{
			{Value: allowedVal1},
			{Value: allowedVal2},
		}
		l4 = []psetter.TaggedValue[string, string]{
			{Value: allowedVal1},
			{Value: allowedVal2},
		}
		l5 = []psetter.TaggedValue[string, string]{
			{Value: allowedVal1},
			{Value: allowedVal2},
		}
	)

	testCases := []paramtest.Setter{
		{
			ID: testhelper.MkID("good-setter"),
			PSetter: psetter.TaggedValueList[string, string]{
				Value:            &l1,
				AllowedVals:      allowedVals,
				TagAllowedVals:   tagAllowedVals,
				TagListSeparator: psetter.StrListSeparator{Sep: "|"},
			},
			ParamVal: allowedVal3,
		},
		{
			ID: testhelper.MkID("good-setter-bad-val"),
			PSetter: psetter.TaggedValueList[string, string]{
				Value:            &l2,
				AllowedVals:      allowedVals,
				TagAllowedVals:   tagAllowedVals,
				TagListSeparator: psetter.StrListSeparator{Sep: "|"},
			},
			ParamVal: badVal,
			SetWithValErr: testhelper.MkExpErr(
				`bad value: "` + badVal + `"`),
		},
		{
			ID: testhelper.MkID("good-setter-empty-val"),
			PSetter: psetter.TaggedValueList[string, string]{
				Value:            &l2,
				AllowedVals:      allowedVals,
				TagAllowedVals:   tagAllowedVals,
				TagListSeparator: psetter.StrListSeparator{Sep: "|"},
			},
			ParamVal: "",
			SetWithValErr: testhelper.MkExpErr(
				`bad value: ""`),
		},
		{
			ID: testhelper.MkID("good-setter-multi-vals"),
			PSetter: psetter.TaggedValueList[string, string]{
				Value:            &l3,
				AllowedVals:      allowedVals,
				TagAllowedVals:   tagAllowedVals,
				TagListSeparator: psetter.StrListSeparator{Sep: "|"},
			},
			ParamVal: allowedVal1 + "," + allowedVal3,
		},
		{
			ID: testhelper.MkID("good-setter-alias-val"),
			PSetter: psetter.TaggedValueList[string, string]{
				Value:            &l4,
				AllowedVals:      allowedVals,
				Aliases:          goodAliases,
				TagAllowedVals:   tagAllowedVals,
				TagListSeparator: psetter.StrListSeparator{Sep: "|"},
			},
			ParamVal: goodAlias,
		},
		{
			ID: testhelper.MkID("bad-setter-nil-value"),
			ExpPanic: testhelper.MkExpPanic(dfltParamName +
				": psetter.TaggedValueList[string,string]" +
				" Check failed: the Value to be set is nil"),
			PSetter: psetter.TaggedValueList[string, string]{
				AllowedVals:      allowedVals,
				TagAllowedVals:   tagAllowedVals,
				TagListSeparator: psetter.StrListSeparator{Sep: "|"},
			},
		},
		{
			ID: testhelper.MkID("bad-setter-no-avals"),
			ExpPanic: testhelper.MkExpPanic(dfltParamName +
				": psetter.TaggedValueList[string,string].AllowedVals" +
				" Check failed:" +
				" the Setter is improperly constructed:" +
				" the map of allowed values has no entries." +
				" It should have at least 2"),
			PSetter: psetter.TaggedValueList[string, string]{
				Value: &l5,
			},
		},
		{
			ID: testhelper.MkID("bad-setter-one-aval"),
			ExpPanic: testhelper.MkExpPanic(dfltParamName +
				": psetter.TaggedValueList[string,string].AllowedVals" +
				" Check failed:" +
				" the Setter is improperly constructed:" +
				" the map of allowed values has only 1 entry." +
				" It should have at least 2"),
			PSetter: psetter.TaggedValueList[string, string]{
				AllowedVals: badAdValsOneEntry,
				Value:       &l5,
			},
		},
		{
			ID: testhelper.MkID("bad-setter-bad-aliases"),
			ExpPanic: testhelper.MkExpPanic(dfltParamName +
				": psetter.TaggedValueList[string,string].Aliases" +
				" Check failed:" +
				" the Setter is improperly constructed:" +
				` bad alias: "` + badAlias + `":` +
				` []string{"bad-value"} - ` +
				`"bad-value" (at index 0) is unknown`),
			PSetter: psetter.TaggedValueList[string, string]{
				AllowedVals: allowedVals,
				Value:       &l5,
				Aliases:     badAliases,
			},
		},
	}

	for _, tc := range testCases {
		f := func(t *testing.T) {
			tc.GFC = commonGFCTaggedValueList

			if tc.ParamName == "" {
				tc.ParamName = dfltParamName
			}

			tc.SetVR(param.Mandatory)
			tc.ValDescriber = true

			tc.Test(t)
		}

		t.Run(tc.IDStr(), f)
	}
}
