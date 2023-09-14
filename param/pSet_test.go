package param_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/paramset"
	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

// TestPSet ...
func TestPSet(t *testing.T) {
	var buff bytes.Buffer

	testCases := []struct {
		testName    string
		psOpts      []param.PSetOptFunc
		errExpected bool
		expEStr     string
	}{
		{
			testName: "nil",
		},
		{
			testName: "set writers",
			psOpts: []param.PSetOptFunc{
				param.SetStdWriter(&buff),
				param.SetErrWriter(&buff),
			},
		},
		{
			testName: "bad error writer",
			psOpts: []param.PSetOptFunc{
				param.SetErrWriter(&buff),
				param.SetErrWriter(nil),
			},
			errExpected: true,
			expEStr:     "param.SetErrWriter cannot take a nil value",
		},
		{
			testName: "bad std writer",
			psOpts: []param.PSetOptFunc{
				param.SetErrWriter(&buff),
				param.SetStdWriter(nil),
			},
			errExpected: true,
			expEStr:     "param.SetStdWriter cannot take a nil value",
		},
		{
			testName: "setopt error",
			psOpts: []param.PSetOptFunc{
				param.SetErrWriter(&buff),
				func(ps *param.PSet) error { return errors.New("whoops") },
			},
			errExpected: true,
			expEStr:     "whoops",
		},
	}

	for i, tc := range testCases {
		opts := make([]param.PSetOptFunc, 1, 1+len(tc.psOpts))
		opts[0] = param.DontExitOnParamSetupErr
		opts = append(opts, tc.psOpts...)

		ps, err := paramset.NewNoHelpNoExitNoErrRpt(opts...)
		if err != nil {
			if !tc.errExpected {
				t.Errorf("test %d: %s : returned an unexpected error: %s",
					i, tc.testName, err)
			} else if err.Error() != tc.expEStr {
				t.Errorf(
					"test %d: %s : err was expected to be: %s\n\t: but was: %s",
					i, tc.testName, tc.expEStr, err)
			}
		} else {
			if tc.errExpected {
				t.Errorf("test %d: %s : didn't return an expected error",
					i, tc.testName)
			}

			if ps.AreSet() {
				t.Errorf("test %d: %s : the parsed flag is unexpectedly set",
					i, tc.testName)
			}
		}
	}
}

// groupNameAndDesc holds a pair of name and description for passing to
// SetGroupDescription
type groupNameAndDesc struct {
	name string
	desc string
}

// TestPSet_SetGroupDescription sets group descriptions and tests the
// resulting PSet matches expectations
func TestPSet_SetGroupDescription(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		sgdParams        []groupNameAndDesc
		panicExpected    bool
		panicMsgContains []string
		expectedDescs    []groupNameAndDesc
		groupsExpected   map[string]bool
	}{
		{
			ID: testhelper.MkID("all good"),
			sgdParams: []groupNameAndDesc{
				{name: "a", desc: "group A desc"},
				{name: "b", desc: "group B desc"},
				{name: "c", desc: "group C desc"},
			},
			expectedDescs: []groupNameAndDesc{
				{name: "a", desc: "group A desc"},
				{name: "b", desc: "group B desc"},
				{name: "c", desc: "group C desc"},
				{name: "d", desc: ""},
			},
			groupsExpected: map[string]bool{
				"a": true,
				"b": true,
				"c": true,
				"d": false,
			},
		},
		{
			ID: testhelper.MkID("reset description"),
			sgdParams: []groupNameAndDesc{
				{name: "a", desc: "group A desc"},
				{name: "b", desc: "group B desc"},
				{name: "b", desc: "other group B desc"},
			},
			panicExpected: true,
			panicMsgContains: []string{
				"description for param group b is already set to:",
				"group B desc",
			},
			expectedDescs: []groupNameAndDesc{
				{name: "a", desc: "group A desc"},
				{name: "b", desc: "group B desc"},
				{name: "d", desc: ""},
			},
			groupsExpected: map[string]bool{
				"a": true,
				"b": true,
				"d": false,
			},
		},
		{
			ID: testhelper.MkID("bad group name"),
			sgdParams: []groupNameAndDesc{
				{name: "99", desc: "group 99 desc"},
			},
			panicExpected: true,
			panicMsgContains: []string{
				"Invalid group name:",
				"the group name '99' is invalid. It must match",
			},
			expectedDescs: []groupNameAndDesc{
				{name: "a", desc: ""},
				{name: "99", desc: ""},
			},
			groupsExpected: map[string]bool{
				"a":  false,
				"99": false,
			},
		},
	}

	for _, tc := range testCases {
		ps, err := paramset.NewNoHelpNoExitNoErrRpt()
		if err != nil {
			t.Fatal(tc.IDStr(), " : couldn't construct the PSet: ", err)
		}

		var panicked bool
		var panicVal any
		var stackTrace []byte

		for _, sgdp := range tc.sgdParams {
			panicked, panicVal, stackTrace = panicSafeAddGroup(ps,
				sgdp.name, sgdp.desc)
			if panicked {
				break
			}
		}
		testhelper.PanicCheckStringWithStack(t, tc.IDStr(),
			panicked, tc.panicExpected,
			panicVal, tc.panicMsgContains, stackTrace)

		for _, gd := range tc.expectedDescs {
			g, ok := ps.GetGroup(gd.name)
			if !ok {
				if gd.desc != "" {
					t.Log(tc.IDStr())
					t.Errorf("\t : group %q was not found", gd.name)
				}
			} else if g.Desc() != gd.desc {
				t.Log(tc.IDStr())
				t.Logf("\t: expected: %s", gd.desc)
				t.Logf("\t:  but was: %s", g.Desc())
				t.Errorf("\t : bad group description for '%s'", gd.name)
			}
		}

		for gName, expected := range tc.groupsExpected {
			_, ok := ps.GetGroup(gName)
			if ok != expected {
				t.Log(tc.IDStr())
				if expected {
					t.Errorf("\t: the group description for '%s'"+
						" was not found when expected",
						gName)
				} else {
					t.Errorf("\t: the group description for '%s'"+
						" was found when not expected",
						gName)
				}
			}
		}
	}
}

// TestPSet_SetTerminalParam sets override values for the terminal parameter
func TestPSet_SetTerminalParam(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		tpVal string
		setTP bool
	}{
		{
			ID:    testhelper.MkID("don't set"),
			tpVal: param.DfltTerminalParam,
		},
		{
			ID:    testhelper.MkID("new val"),
			tpVal: "xxx",
			setTP: true,
		},
	}

	for _, tc := range testCases {
		ps, err := paramset.NewNoHelpNoExitNoErrRpt()
		if err != nil {
			t.Fatal(tc.IDStr(), " : couldn't construct the PSet: ", err)
		}

		if tc.setTP {
			ps.SetTerminalParam(tc.tpVal)
		}
		if ps.TerminalParam() != tc.tpVal {
			t.Log(tc.IDStr())
			t.Logf("\t: expected: %s", tc.tpVal)
			t.Logf("\t:      got: %s", ps.TerminalParam())
			t.Errorf("\t: Bad TerminalParam")
		}
	}
}

// GroupAndParams holds the expected group name and associated parameter names
type GroupAndParams struct {
	groupName   string
	paramNames  []string
	hiddenCount int
	allHidden   bool
}

type paramGroupTC struct {
	testhelper.ID
	npi             []*namedParamInitialiser
	expectedResults []GroupAndParams
}

// reportParamGroup prints the param group details
func reportParamGroup(t *testing.T, groups []*param.Group) {
	t.Helper()
	for _, g := range groups {
		t.Logf("\t: Group: %s\n", g.Name())
		for _, p := range g.Params() {
			t.Logf("\t\t%s\n", p.Name())
		}
	}
}

// checkParamGroup confirms that the param groups are as expected
func checkParamGroup(t *testing.T, i int, tc paramGroupTC, ps *param.PSet) {
	t.Helper()

	groups := ps.GetGroups()
	if len(groups) != len(tc.expectedResults) {
		t.Log(tc.IDStr())
		t.Logf("\t: expected: %d", len(tc.expectedResults))
		t.Logf("\t:      got: %d", len(groups))
		reportParamGroup(t, groups)
		t.Error("\t: the number of Groups returned is unexpected")
		return
	}
	for idx, g := range groups {
		tcIDGrp := tc.IDStr() + fmt.Sprintf(" - group %d", idx)
		if g.Name() != tc.expectedResults[idx].groupName {
			t.Log(tcIDGrp)
			t.Logf("\t: expected: %s", tc.expectedResults[idx].groupName)
			t.Logf("\t:      got: %s", g.Name())
			reportParamGroup(t, groups)
			t.Error("\t: the group name is unexpected")
		}
		if len(g.Params()) != len(tc.expectedResults[idx].paramNames) {
			t.Log(tcIDGrp)
			t.Logf("\t: expected: %d", len(tc.expectedResults[idx].paramNames))
			t.Logf("\t:      got: %d", len(g.Params()))
			reportParamGroup(t, groups)
			t.Error("\t: the number of parameters is unexpected")
		}
		if g.HiddenCount() != tc.expectedResults[idx].hiddenCount {
			t.Log(tcIDGrp)
			t.Logf("\t: expected: %d", tc.expectedResults[idx].hiddenCount)
			t.Logf("\t:      got: %d", g.HiddenCount())
			reportParamGroup(t, groups)
			t.Error("\t: the number of hidden parameters is unexpected")
		}
		if g.AllParamsHidden() != tc.expectedResults[idx].allHidden {
			t.Log(tcIDGrp)
			t.Logf("\t: expected: %v", tc.expectedResults[idx].allHidden)
			t.Logf("\t:      got: %v", g.AllParamsHidden())
			reportParamGroup(t, groups)
			t.Error("\t: the number of hidden parameters is unexpected")
		}
	}
}

func TestGetParamGroups(t *testing.T) {
	var boolVar bool
	testCases := []paramGroupTC{
		{
			ID: testhelper.MkID("no params"),
		},
		{
			ID: testhelper.MkID("one param, default group"),
			npi: []*namedParamInitialiser{
				{
					name:   "param",
					setter: psetter.Bool{Value: &boolVar},
					desc:   "desc",
				},
			},
			expectedResults: []GroupAndParams{
				{
					groupName: param.DfltGroupName,
					paramNames: []string{
						"param",
					},
				},
			},
		},
		{
			ID: testhelper.MkID("two params, default group"),
			npi: []*namedParamInitialiser{
				{
					name:   "param",
					setter: psetter.Bool{Value: &boolVar},
					desc:   "desc",
				},
				{
					name:   "param2",
					setter: psetter.Bool{Value: &boolVar},
					desc:   "desc",
				},
			},
			expectedResults: []GroupAndParams{
				{
					groupName: param.DfltGroupName,
					paramNames: []string{
						"param",
						"param2",
					},
				},
			},
		},
		{
			ID: testhelper.MkID("two params, two groups"),
			npi: []*namedParamInitialiser{
				{
					name:   "param",
					setter: psetter.Bool{Value: &boolVar},
					desc:   "desc",
					opts:   []param.OptFunc{param.GroupName("abc")},
				},
				{
					name:   "param2",
					setter: psetter.Bool{Value: &boolVar},
					desc:   "desc",
					opts:   []param.OptFunc{param.GroupName("xyz")},
				},
			},
			expectedResults: []GroupAndParams{
				{
					groupName: "abc",
					paramNames: []string{
						"param",
					},
				},
				{
					groupName: "xyz",
					paramNames: []string{
						"param2",
					},
				},
			},
		},
		{
			ID: testhelper.MkID("three params, two hidden, two groups"),
			npi: []*namedParamInitialiser{
				{
					name:   "aaa",
					setter: psetter.Bool{Value: &boolVar},
					desc:   "desc",
					opts: []param.OptFunc{
						param.GroupName("abc"),
						param.Attrs(param.DontShowInStdUsage),
					},
				},
				{
					name:   "aab",
					setter: psetter.Bool{Value: &boolVar},
					desc:   "desc",
					opts: []param.OptFunc{
						param.GroupName("abc"),
						param.Attrs(param.DontShowInStdUsage),
					},
				},
				{
					name:   "param2",
					setter: psetter.Bool{Value: &boolVar},
					desc:   "desc",
					opts:   []param.OptFunc{param.GroupName("xyz")},
				},
			},
			expectedResults: []GroupAndParams{
				{
					groupName: "abc",
					paramNames: []string{
						"aaa",
						"aab",
					},
					hiddenCount: 2,
					allHidden:   true,
				},
				{
					groupName: "xyz",
					paramNames: []string{
						"param2",
					},
				},
			},
		},
	}

	for i, tc := range testCases {
		ps, err := paramset.NewNoHelpNoExit()
		if err != nil {
			t.Fatal("Cannot construct the PSet:", err.Error())
		}
		for _, npi := range tc.npi {
			ps.Add(npi.name, npi.setter, npi.desc, npi.opts...)
		}
		checkParamGroup(t, i, tc, ps)
	}
}
