package param_test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/nickwells/param.mod/param"
	"github.com/nickwells/param.mod/param/paramset"
	"github.com/nickwells/param.mod/param/psetter"
	"github.com/nickwells/testhelper.mod/testhelper"
	"testing"
)

// TestParamSet ...
func TestParamSet(t *testing.T) {
	var buff bytes.Buffer

	testCases := []struct {
		testName    string
		psOpts      []param.ParamSetOptFunc
		errExpected bool
		expEStr     string
	}{
		{
			testName: "nil",
		},
		{
			testName: "set writers",
			psOpts: []param.ParamSetOptFunc{
				param.SetStdWriter(&buff),
				param.SetErrWriter(&buff),
			},
		},
		{
			testName: "bad error writer",
			psOpts: []param.ParamSetOptFunc{
				param.SetErrWriter(&buff),
				param.SetErrWriter(nil),
			},
			errExpected: true,
			expEStr:     "param.SetErrWriter cannot take a nil value",
		},
		{
			testName: "bad std writer",
			psOpts: []param.ParamSetOptFunc{
				param.SetErrWriter(&buff),
				param.SetStdWriter(nil),
			},
			errExpected: true,
			expEStr:     "param.SetStdWriter cannot take a nil value",
		},
		{
			testName: "setopt error",
			psOpts: []param.ParamSetOptFunc{
				param.SetErrWriter(&buff),
				func(ps *param.ParamSet) error { return errors.New("whoops") },
			},
			errExpected: true,
			expEStr:     "whoops",
		},
	}

	for i, tc := range testCases {
		opts := make([]param.ParamSetOptFunc, 1, 1+len(tc.psOpts))
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

// TestParamSet_SetGroupDescription sets group descriptions and tests the
// resulting ParamSet matches expectations
func TestParamSet_SetGroupDescription(t *testing.T) {
	testCases := []struct {
		name             string
		sgdParams        []groupNameAndDesc
		panicExpected    bool
		panicMsgContains []string
		expectedDescs    []groupNameAndDesc
		groupsExpected   map[string]bool
	}{
		{
			name: "all good",
			sgdParams: []groupNameAndDesc{
				groupNameAndDesc{name: "a", desc: "group A desc"},
				groupNameAndDesc{name: "b", desc: "group B desc"},
				groupNameAndDesc{name: "c", desc: "group C desc"},
			},
			expectedDescs: []groupNameAndDesc{
				groupNameAndDesc{name: "a", desc: "group A desc"},
				groupNameAndDesc{name: "b", desc: "group B desc"},
				groupNameAndDesc{name: "c", desc: "group C desc"},
				groupNameAndDesc{name: "d", desc: ""},
			},
			groupsExpected: map[string]bool{
				"a": true,
				"b": true,
				"c": true,
				"d": false,
			},
		},
		{
			name: "reset description",
			sgdParams: []groupNameAndDesc{
				groupNameAndDesc{name: "a", desc: "group A desc"},
				groupNameAndDesc{name: "b", desc: "group B desc"},
				groupNameAndDesc{name: "b", desc: "other group B desc"},
			},
			panicExpected: true,
			panicMsgContains: []string{
				"description for param group b is already set to:",
				"group B desc",
			},
			expectedDescs: []groupNameAndDesc{
				groupNameAndDesc{name: "a", desc: "group A desc"},
				groupNameAndDesc{name: "b", desc: "group B desc"},
				groupNameAndDesc{name: "d", desc: ""},
			},
			groupsExpected: map[string]bool{
				"a": true,
				"b": true,
				"d": false,
			},
		},
		{
			name: "bad group name",
			sgdParams: []groupNameAndDesc{
				groupNameAndDesc{name: "99", desc: "group 99 desc"},
			},
			panicExpected: true,
			panicMsgContains: []string{
				"Invalid group name:",
				"the group name '99' is invalid. It must match",
			},
			expectedDescs: []groupNameAndDesc{
				groupNameAndDesc{name: "a", desc: ""},
				groupNameAndDesc{name: "99", desc: ""},
			},
			groupsExpected: map[string]bool{
				"a":  false,
				"99": false,
			},
		},
	}

	for i, tc := range testCases {
		testName := fmt.Sprintf("%d: %s", i, tc.name)

		ps, err := paramset.NewNoHelpNoExitNoErrRpt()
		if err != nil {
			t.Fatal(testName, " : couldn't construct the ParamSet: ", err)
		}

		var panicked bool
		var panicVal interface{}

		for _, sgdp := range tc.sgdParams {
			panicked, panicVal = panicSafeSetGroupDescription(ps,
				sgdp.name, sgdp.desc)
			if panicked {
				break
			}
		}
		testhelper.PanicCheckString(t, testName,
			panicked, tc.panicExpected,
			panicVal, tc.panicMsgContains)

		for _, gd := range tc.expectedDescs {
			desc := ps.GetGroupDesc(gd.name)
			if desc != gd.desc {
				t.Errorf("%s : the group description for '%s'"+
					" was expected to be: '%s' but was '%s'",
					testName, gd.name, gd.desc, desc)
			}
		}

		for gName, expected := range tc.groupsExpected {
			hasName := ps.HasGroupName(gName)
			if hasName != expected {
				t.Logf("%s\n", testName)
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

// TestParamSet_SetTerminalParam sets override values for the terminal parameter
func TestParamSet_SetTerminalParam(t *testing.T) {
	testCases := []struct {
		name  string
		tpVal string
		setTP bool
	}{
		{
			name:  "don't set",
			tpVal: param.DfltTerminalParam,
		},
		{
			name:  "new val",
			tpVal: "xxx",
			setTP: true,
		},
	}

	for i, tc := range testCases {
		testName := fmt.Sprintf("%d: %s", i, tc.name)

		ps, err := paramset.NewNoHelpNoExitNoErrRpt()
		if err != nil {
			t.Fatal(testName, " : couldn't construct the ParamSet: ", err)
		}

		if tc.setTP {
			ps.SetTerminalParam(tc.tpVal)
		}
		if ps.TerminalParam() != tc.tpVal {
			t.Errorf("%s : TerminalParam was expected to be: '%s' but was '%s'",
				testName, tc.tpVal, ps.TerminalParam())
		}
	}
}

// ExampleParamSet_Add shows the usage of the Add method of the
// ParamSet. This is used to add new parameters into the set.
func ExampleParamSet_Add() {
	ps, _ := paramset.New()

	// we declare f here for the purposes of the example but typically it
	// would be declared in package scope somewhere or in the main() func
	var f float64

	p := ps.Add(
		"param-name",
		psetter.Float64Setter{Value: &f},
		"a parameter description",
		param.GroupName("test.group"),
		param.Attrs(param.DontShowInStdUsage))

	fmt.Printf("%3.1f\n", f)
	fmt.Printf("group name: %s\n", p.GroupName())
	fmt.Printf("param name: %s\n", p.Name())
	fmt.Printf("CommandLineOnly: %t\n", p.AttrIsSet(param.CommandLineOnly))
	fmt.Printf("MustBeSet: %t\n", p.AttrIsSet(param.MustBeSet))
	fmt.Printf("SetOnlyOnce: %t\n", p.AttrIsSet(param.SetOnlyOnce))
	fmt.Printf("DontShowInStdUsage: %t\n", p.AttrIsSet(param.DontShowInStdUsage))

	// Output: 0.0
	// group name: test.group
	// param name: param-name
	// CommandLineOnly: false
	// MustBeSet: false
	// SetOnlyOnce: false
	// DontShowInStdUsage: true
}

// GroupAndParams holds the expected group name and associated parameter names
type GroupAndParams struct {
	groupName  string
	paramNames []string
}

type paramGroupTC struct {
	name            string
	npi             []*namedParamInitialiser
	expectedResults []GroupAndParams
}

// reportParamGroup prints the param group details
func reportParamGroup(t *testing.T, paramGroups []*param.ParamGroup) {
	t.Helper()
	for _, pg := range paramGroups {
		t.Logf("\t: Group: %s\n", pg.GroupName)
		for _, p := range pg.Params {
			t.Logf("\t\t%s\n", p.Name())
		}
	}
}

// checkParamGroup confirms that the param groups are as expected
func checkParamGroup(t *testing.T, i int, tc paramGroupTC, ps *param.ParamSet) {
	t.Helper()

	tcID := fmt.Sprintf("%d: %s", i, tc.name)
	paramGroups := ps.GetParamGroups()
	if len(paramGroups) != len(tc.expectedResults) {
		t.Logf("%s: the number of ParamGroups returned is unexpected\n", tcID)
		t.Logf("\t: expected: %d\n", len(tc.expectedResults))
		t.Logf("\t:      got: %d\n", len(paramGroups))
		reportParamGroup(t, paramGroups)
		t.Error("\t: Failed")
		return
	}
	for idx, pg := range paramGroups {
		if pg.GroupName != tc.expectedResults[idx].groupName {
			t.Logf("%s: the group name for group %d is unexpected\n", tcID, idx)
			t.Logf("\t: expected: %s\n", tc.expectedResults[idx].groupName)
			t.Logf("\t:      got: %s\n", pg.GroupName)
			reportParamGroup(t, paramGroups)
			t.Error("\t: Failed")
		}
		if len(pg.Params) != len(tc.expectedResults[idx].paramNames) {
			t.Logf("%s: the number of parameters for group %d is unexpected\n",
				tcID, idx)
			t.Logf("\t: expected: %d\n",
				len(tc.expectedResults[idx].paramNames))
			t.Logf("\t:      got: %d\n", len(pg.Params))
			reportParamGroup(t, paramGroups)
			t.Error("\t: Failed")
			return
		}
	}
}

func TestGetParamGroups(t *testing.T) {
	var boolVar bool
	testCases := []paramGroupTC{
		{
			name: "no params",
		},
		{
			name: "one param, default group",
			npi: []*namedParamInitialiser{
				&namedParamInitialiser{
					name:   "param",
					setter: psetter.BoolSetter{Value: &boolVar},
					desc:   "desc",
				},
			},
			expectedResults: []GroupAndParams{
				GroupAndParams{
					groupName: param.DfltGroupName,
					paramNames: []string{
						"param",
					},
				},
			},
		},
		{
			name: "two params, default group",
			npi: []*namedParamInitialiser{
				&namedParamInitialiser{
					name:   "param",
					setter: psetter.BoolSetter{Value: &boolVar},
					desc:   "desc",
				},
				&namedParamInitialiser{
					name:   "param2",
					setter: psetter.BoolSetter{Value: &boolVar},
					desc:   "desc",
				},
			},
			expectedResults: []GroupAndParams{
				GroupAndParams{
					groupName: param.DfltGroupName,
					paramNames: []string{
						"param",
						"param2",
					},
				},
			},
		},
		{
			name: "two params, two groups",
			npi: []*namedParamInitialiser{
				&namedParamInitialiser{
					name:   "param",
					setter: psetter.BoolSetter{Value: &boolVar},
					desc:   "desc",
					opts:   []param.OptFunc{param.GroupName("abc")},
				},
				&namedParamInitialiser{
					name:   "param2",
					setter: psetter.BoolSetter{Value: &boolVar},
					desc:   "desc",
					opts:   []param.OptFunc{param.GroupName("xyz")},
				},
			},
			expectedResults: []GroupAndParams{
				GroupAndParams{
					groupName: "abc",
					paramNames: []string{
						"param",
					},
				},
				GroupAndParams{
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
			t.Fatal("Cannot construct the ParamSet:", err.Error())
		}
		for _, npi := range tc.npi {
			ps.Add(npi.name, npi.setter, npi.desc, npi.opts...)
		}
		checkParamGroup(t, i, tc, ps)
	}

}
