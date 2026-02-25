package param_test

import (
	"testing"

	"github.com/nickwells/param.mod/v7/param"
	"github.com/nickwells/param.mod/v7/paramset"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

type pfxFunc struct {
	pfx string
	f   func(*param.PSet, string)
}

func TestSetEnv(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		seq              []pfxFunc
		panicExpected    bool
		panicMustContain []string
	}{
		{
			ID: testhelper.MkID("empty set - SetEnvPrefix - empty prefix"),
			seq: []pfxFunc{
				{
					pfx: "",
					f:   (*param.PSet).SetEnvPrefix,
				},
			},
			panicExpected: true,
			panicMustContain: []string{
				`environment variable prefix must not be empty`,
			},
		},
		{
			ID: testhelper.MkID("empty set - AddEnvPrefix - empty prefix"),
			seq: []pfxFunc{
				{
					pfx: "",
					f:   (*param.PSet).AddEnvPrefix,
				},
			},
			panicExpected: true,
			panicMustContain: []string{
				`environment variable prefix must not be empty`,
			},
		},
		{
			ID: testhelper.MkID("one prefix - AddEnvPrefix - empty prefix"),
			seq: []pfxFunc{
				{
					pfx: "somePfx_",
					f:   (*param.PSet).SetEnvPrefix,
				},
				{
					pfx: "",
					f:   (*param.PSet).AddEnvPrefix,
				},
			},
			panicExpected: true,
			panicMustContain: []string{
				`environment variable prefix must not be empty`,
			},
		},
		{
			ID: testhelper.MkID("empty set - SetEnvPrefix - good prefix"),
			seq: []pfxFunc{
				{
					pfx: "somePfx_",
					f:   (*param.PSet).SetEnvPrefix,
				},
			},
			panicExpected: false,
		},
		{
			ID: testhelper.MkID("one prefix - AddEnvPrefix - good prefix"),
			seq: []pfxFunc{
				{
					pfx: "somePfx_",
					f:   (*param.PSet).SetEnvPrefix,
				},
				{
					pfx: "another_good_prefix_",
					f:   (*param.PSet).AddEnvPrefix,
				},
			},
			panicExpected: false,
		},
		{
			ID: testhelper.MkID("one prefix - AddEnvPrefix - bad prefix"),
			seq: []pfxFunc{
				{
					pfx: "some_Pfx_",
					f:   (*param.PSet).SetEnvPrefix,
				},
				{
					pfx: "some_",
					f:   (*param.PSet).AddEnvPrefix,
				},
			},
			panicExpected: true,
			panicMustContain: []string{
				`invalid environment variable prefix: "some_": ` +
					`it's a prefix of the already added: "some_Pfx_"`,
			},
		},
		{
			ID: testhelper.MkID("one prefix - AddEnvPrefix - bad prefix"),
			seq: []pfxFunc{
				{
					pfx: "some_",
					f:   (*param.PSet).SetEnvPrefix,
				},
				{
					pfx: "some_Pfx_",
					f:   (*param.PSet).AddEnvPrefix,
				},
			},
			panicExpected: true,
			panicMustContain: []string{
				`invalid environment variable prefix: "some_Pfx_": ` +
					`the already added: "some_" is a prefix of it`,
			},
		},
	}
	for _, tc := range testCases {
		ps := paramset.NewNoHelpNoExitNoErrRpt()

		panicked, panicVal := panicEnvPrefix(t, tc.seq, ps)

		testhelper.PanicCheckError(t, tc.IDStr(),
			panicked, tc.panicExpected,
			panicVal, tc.panicMustContain)
	}
}

func panicEnvPrefix(t *testing.T, seq []pfxFunc, ps *param.PSet,
) (panicked bool, panicVal any) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()

	for _, pf := range seq {
		pf.f(ps, pf.pfx)
	}

	return panicked, panicVal
}
