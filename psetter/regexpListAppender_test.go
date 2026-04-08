package psetter_test

import (
	"regexp"
	"testing"

	"github.com/nickwells/param.mod/v7/param"
	"github.com/nickwells/param.mod/v7/paramtest"
	"github.com/nickwells/param.mod/v7/psetter"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

const (
	updFlagNameRegexpListAppender     = "upd-gf-RegexpListAppender"
	keepBadFlagNameRegexpListAppender = "keep-bad-RegexpListAppender"
)

var commonGFCRegexpListAppender = testhelper.GoldenFileCfg{
	DirNames:               []string{"testdata", "RegexpListAppender"},
	Pfx:                    "gf",
	Sfx:                    "txt",
	UpdFlagName:            updFlagNameRegexpListAppender,
	KeepBadResultsFlagName: keepBadFlagNameRegexpListAppender,
}

func init() {
	commonGFCRegexpListAppender.AddUpdateFlag()
	commonGFCRegexpListAppender.AddKeepBadResultsFlag()
}

func TestSetterRegexpListAppender(t *testing.T) {
	const dfltParamName = "param-name-regexp-list-appender"

	reList1 := []*regexp.Regexp{
		regexp.MustCompile(`.*`),
	}
	reList2 := []*regexp.Regexp{
		regexp.MustCompile(`.*`),
	}
	reList3 := []*regexp.Regexp{
		regexp.MustCompile(`.*`),
	}

	testCases := []paramtest.Setter{
		{
			ID: testhelper.MkID("bad-setter-no-value"),
			ExpPanic: testhelper.MkExpPanic(dfltParamName +
				": psetter.RegexpListAppender" +
				" Check failed: the Value to be set is nil"),
			PSetter: psetter.RegexpListAppender{},
		},
		{
			ID: testhelper.MkID("good-setter-empty-param-value"),
			PSetter: psetter.RegexpListAppender{
				Value: &reList1,
			},
		},
		{
			ID: testhelper.MkID("good-setter-valid-param-value"),
			PSetter: psetter.RegexpListAppender{
				Value: &reList2,
			},
			ParamVal: `^hello.*world$`,
		},
		{
			ID: testhelper.MkID("good-setter-bad-param-value"),
			PSetter: psetter.RegexpListAppender{
				Value: &reList3,
			},
			ParamVal: `*`,
			SetWithValErr: testhelper.MkExpErr(
				`could not parse "*" into a regular expression:` +
					` error parsing regexp:` +
					` missing argument to repetition operator:` + " `*`"),
		},
	}

	for _, tc := range testCases {
		f := func(t *testing.T) {
			tc.GFC = commonGFCRegexpListAppender

			if tc.ParamName == "" {
				tc.ParamName = dfltParamName
			}

			tc.SetVR(param.Mandatory)

			tc.Test(t)
		}

		t.Run(tc.IDStr(), f)
	}
}
