package param_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nickwells/param.mod/v4/param"
	"github.com/nickwells/param.mod/v4/param/paramset"
	"github.com/nickwells/param.mod/v4/param/psetter"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestParamAdd(t *testing.T) { // nolint: gocyclo
	var p1 int64
	testCases := []struct {
		testhelper.ID
		testhelper.ExpPanic
		npi              *namedParamInitialiser
		paramShouldExist bool
	}{
		{
			ID: testhelper.MkID("bad name - empty"),
			npi: &namedParamInitialiser{
				setter: &psetter.Int64{Value: &p1},
			},
			ExpPanic: testhelper.MkExpPanic(
				"the parameter name",
				"is invalid. It must match"),
		},
		{
			ID: testhelper.MkID("bad name - bad char"),
			npi: &namedParamInitialiser{
				name:   "?",
				setter: &psetter.Int64{Value: &p1},
			},
			ExpPanic: testhelper.MkExpPanic(
				"the parameter name",
				"is invalid. It must match"),
		},
		{
			ID: testhelper.MkID("bad name - bad first char"),
			npi: &namedParamInitialiser{
				name:   "-hello",
				setter: &psetter.Int64{Value: &p1},
			},
			ExpPanic: testhelper.MkExpPanic(
				"the parameter name",
				"is invalid. It must match"),
		},
		{
			ID: testhelper.MkID("good name"),
			npi: &namedParamInitialiser{
				name:   "param-1",
				setter: &psetter.Int64{Value: &p1},
				opts: []param.OptFunc{
					param.AltName("param-1-alt"),
					param.GroupName("test"),
				},
			},
			paramShouldExist: true,
		},
		{
			ID: testhelper.MkID("bad name - duplicate"),
			npi: &namedParamInitialiser{
				name:   "param-1",
				setter: &psetter.Int64{Value: &p1},
			},
			ExpPanic: testhelper.MkExpPanic(
				"parameter name",
				"has already been used"),
			paramShouldExist: true,
		},
		{
			ID: testhelper.MkID("bad alt name - already used"),
			npi: &namedParamInitialiser{
				name:   "param-2",
				setter: &psetter.Int64{Value: &p1},
				opts:   []param.OptFunc{param.AltName("param-1")},
			},
			ExpPanic: testhelper.MkExpPanic(
				"parameter name",
				"has already been used"),
			paramShouldExist: true,
		},
		{
			ID: testhelper.MkID("bad alt name - invalid"),
			npi: &namedParamInitialiser{
				name:   "param-3",
				setter: &psetter.Int64{Value: &p1},
				opts:   []param.OptFunc{param.AltName("?")},
			},
			ExpPanic: testhelper.MkExpPanic(
				"parameter name",
				"is invalid. It must match"),
			paramShouldExist: true,
		},
		{
			ID: testhelper.MkID("bad alt name - already used as alt"),
			npi: &namedParamInitialiser{
				name:   "param-4",
				setter: &psetter.Int64{Value: &p1},
				opts:   []param.OptFunc{param.AltName("param-1-alt")},
			},
			ExpPanic: testhelper.MkExpPanic(
				"parameter name",
				`has already been used as an alternative to "param-1"`,
				`a member of parameter group "test"`),
			paramShouldExist: true,
		},
	}

	ps, err := paramset.NewNoHelpNoExitNoErrRpt()
	if err != nil {
		t.Fatal("couldn't construct the PSet: ", err)
	}
	for _, tc := range testCases {
		p, panicked, panicVal := panicSafeTestAddByName(ps, tc.npi)
		testhelper.CheckExpPanic(t, panicked, panicVal, tc)

		p2, err := ps.GetParamByName(tc.npi.name)
		if tc.paramShouldExist {
			if p2 == nil {
				t.Log(tc.IDStr())
				t.Errorf("\t: param: '%s' should exist\n",
					tc.npi.name)
			}
			if err != nil {
				t.Log(tc.IDStr())
				t.Errorf("\t: GetParamByName(...) returned an error: %s\n",
					err)
			}
		} else {
			if p2 != nil {
				t.Log(tc.IDStr())
				t.Errorf(
					"\t: param: '%s' should not exist if Add(...) panicked\n",
					tc.npi.name)
			}
			if err == nil {
				t.Log(tc.IDStr())
				t.Errorf(
					"\t: GetParamByName(...) should have returned an error\n")
			}
		}

		if p != nil {
			if p.Name() != tc.npi.name {
				t.Log(tc.IDStr())
				t.Errorf("\t: the name did not match: '%s' != '%s'\n",
					p.Name(), tc.npi.name)
			} else if p.Description() != tc.npi.desc {
				t.Log(tc.IDStr())
				t.Errorf("\t: the description did not match: '%s' != '%s'\n",
					p.Description(), tc.npi.desc)
			} else if p.HasBeenSet() {
				t.Log(tc.IDStr())
				t.Errorf(
					"\t: param has been set but params haven't been parsed\n")
			}
		}
	}
}

func TestParamAddPos(t *testing.T) { // nolint: gocyclo
	var p1 int64
	testCases := []struct {
		testhelper.ID
		testhelper.ExpPanic
		pi []paramInitialisers

		paramsToParse     []string
		errsExpected      map[string][]string
		remainderExpected []string
	}{
		{
			ID: testhelper.MkID("good params"),
			pi: []paramInitialisers{
				{
					npi: &namedParamInitialiser{
						name:   "param-1",
						setter: &psetter.Int64{Value: &p1},
						opts: []param.OptFunc{
							param.AltName("param-1-alt"),
							param.GroupName("test"),
						},
					},
					npiShouldExist: true,
				},
				{
					npi: &namedParamInitialiser{
						name:   "param-2",
						setter: &psetter.Int64{Value: &p1},
						opts: []param.OptFunc{
							param.AltName("param-2-alt"),
							param.GroupName("test"),
						},
					},
					npiShouldExist: true,
				},
				{
					ppi: &posParamInitialiser{
						name:   "ppi1",
						setter: &psetter.Int64{Value: &p1},
					},
					ppiShouldExist: true,
				},
			},
			paramsToParse:     []string{"1", "--", "another param"},
			remainderExpected: []string{"another param"},
		},
		{
			ID: testhelper.MkID("bad params - name empty"),
			pi: []paramInitialisers{
				{
					npi: &namedParamInitialiser{
						name:   "",
						setter: &psetter.Int64{Value: &p1},
						opts: []param.OptFunc{
							param.GroupName("test"),
						},
					},
				},
			},
			ExpPanic: testhelper.MkExpPanic(
				"the parameter name",
				"is invalid. It must match"),
		},
		{
			ID: testhelper.MkID(
				"bad params - terminal ByPos and pre-existing ByName"),
			pi: []paramInitialisers{
				{
					npi: &namedParamInitialiser{
						name:   "param-1",
						setter: &psetter.Int64{Value: &p1},
						opts: []param.OptFunc{
							param.AltName("param-1-alt"),
							param.GroupName("test"),
						},
					},
					npiShouldExist: true,
				},
				{
					ppi: &posParamInitialiser{
						name:   "ppi1",
						setter: &psetter.Int64{Value: &p1},
						opts: []param.PosOptFunc{
							param.SetAsTerminal,
						},
					},
				},
			},
			ExpPanic: testhelper.MkExpPanic(
				`Couldn't set the options for positional parameter 1 ("ppi1"):`,
				"The param set has 1 non-positional parameters.",
				" It cannot also have a terminal positional parameter as"+
					" the non-positional parameters will never be used."),
		},
		{
			ID: testhelper.MkID(
				"bad params - terminal ByPos and then add ByName"),
			pi: []paramInitialisers{
				{
					ppi: &posParamInitialiser{
						name:   "ppi1",
						setter: &psetter.Int64{Value: &p1},
						opts: []param.PosOptFunc{
							param.SetAsTerminal,
						},
					},
					ppiShouldExist: true,
				},
				{
					npi: &namedParamInitialiser{
						name:   "param-1",
						setter: &psetter.Int64{Value: &p1},
						opts: []param.OptFunc{
							param.AltName("param-1-alt"),
							param.GroupName("test"),
						},
					},
				},
			},
			ExpPanic: testhelper.MkExpPanic(
				"The param set has a terminal positional parameter.",
				"The non-positional parameter param-1 cannot be added"+
					" as it will never be used"),
		},
		{
			ID: testhelper.MkID("bad params - terminal ByPos not the last"),
			pi: []paramInitialisers{
				{
					ppi: &posParamInitialiser{
						name:   "ppi1",
						setter: &psetter.Int64{Value: &p1},
						opts: []param.PosOptFunc{
							param.SetAsTerminal,
						},
					},
					ppiShouldExist: true,
				},
				{
					ppi: &posParamInitialiser{
						name:   "ppi2",
						setter: &psetter.Int64{Value: &p1},
						opts:   []param.PosOptFunc{},
					},
				},
			},
			ExpPanic: testhelper.MkExpPanic(
				"Positional parameter 0 is marked as terminal but" +
					" is not the last positional parameter"),
		},
		{
			ID: testhelper.MkID("Parse(...) errors - ByPos: setter.Process(...) failure"),
			pi: []paramInitialisers{
				{
					ppi: &posParamInitialiser{
						name:   "ppi1",
						setter: &failingSetter{errMsg: "ByPos param"},
					},
					ppiShouldExist: true,
				},
			},
			paramsToParse:     []string{"1"},
			remainderExpected: []string{},
			errsExpected: map[string][]string{
				"Positional parameter: 1 (ppi1)": {
					"failingSetter: ByPos param",
				},
			},
		},
		{
			ID: testhelper.MkID("Parse(...) errors - ByName: setter.Process(...) failure"),
			pi: []paramInitialisers{
				{
					npi: &namedParamInitialiser{
						name:   "test99",
						setter: &failingSetter{errMsg: "ByName param"},
					},
					npiShouldExist: true,
				},
			},
			paramsToParse:     []string{"-test99", "val"},
			remainderExpected: []string{},
			errsExpected: map[string][]string{
				"test99": {
					"error with parameter:",
					"failingSetter: ByName param",
				},
			},
		},
	}

	for _, tc := range testCases {
		ps, err := paramset.NewNoHelpNoExitNoErrRpt()
		if err != nil {
			t.Fatal(tc.IDStr(), " : couldn't construct the PSet: ", err)
		}
		var panicked bool
		var panicVal interface{}

		var posIdx int
		for _, pi := range tc.pi {
			if pi.npi != nil {
				_, panicked, panicVal = panicSafeTestAddByName(ps, pi.npi)

				np, err := ps.GetParamByName(pi.npi.name)
				var exists bool
				if np != nil {
					exists = true
				}
				if pi.npiShouldExist != exists {
					if exists {
						t.Log(tc.IDStr())
						t.Errorf("\t: named parameter: '%s'"+
							" should not exist but does\n",
							pi.npi.name)
					} else {
						t.Log(tc.IDStr())
						t.Errorf("\t: named parameter: '%s'"+
							" should exist but doesn't. Err: %s\n",
							pi.npi.name, err)
					}
				}

				if panicked {
					break
				}
			}

			if pi.ppi != nil {
				_, panicked, panicVal = panicSafeTestAddByPos(ps, pi.ppi)

				pp, err := ps.GetParamByPos(posIdx)
				var exists bool
				if pp != nil {
					exists = true
				}
				if pi.ppiShouldExist != exists {
					if exists {
						t.Log(tc.IDStr())
						t.Errorf("\t: positional parameter: %d"+
							" should not exist but does\n",
							posIdx)
					} else {
						t.Log(tc.IDStr())
						t.Errorf("\t: positional parameter: %d"+
							" should exist but doesn't. Err: %s\n",
							posIdx, err)
					}
				}

				if panicked {
					break
				}
				posIdx++
			}
		}
		testhelper.CheckExpPanic(t, panicked, panicVal, tc)

		if !panicked {
			errMap, panicked, panicVal, stackTrace :=
				panicSafeTestParse(ps, tc.paramsToParse)
			if !testhelper.PanicCheckStringWithStack(t, tc.IDStr(),
				panicked, false,
				panicVal, []string{}, stackTrace) {
				errMapCheck(t, tc.IDStr(), errMap, tc.errsExpected)

				if tc.remainderExpected != nil &&
					len(tc.remainderExpected) == 0 {
					tc.remainderExpected = nil
				}

				if !reflect.DeepEqual(ps.Remainder(), tc.remainderExpected) {
					t.Log(tc.IDStr())
					t.Logf("\t: remainder received: %v\n", ps.Remainder())
					t.Logf("\t: remainder expected: %v\n", tc.remainderExpected)
					t.Errorf("\t: unexpected remainder\n")
				}
			}
		}
	}
}

func TestParamParse1(t *testing.T) { // nolint: gocyclo
	var p1 int64
	var p2 int64

	ps, err := paramset.NewNoHelpNoExitNoErrRpt()

	if err != nil {
		t.Fatal("couldn't construct the PSet: ", err)
	}

	testName := "param.Parse - first Parse"
	ps.Add("param1", &psetter.Int64{Value: &p1}, "a param")
	ps.Add("param2", &psetter.Int64{Value: &p2}, "a param")

	if ps.AreSet() {
		t.Errorf("%s: params haven't been set but AreSet() says they have",
			testName)
	}
	errMap, panicked, panicVal, stackTrace :=
		panicSafeTestParse(ps, []string{"-param1", "99"})

	if testhelper.ReportUnexpectedPanic(t, testName,
		panicked, panicVal, stackTrace) {
		return
	}

	if !ps.AreSet() {
		t.Log(testName)
		t.Errorf("\t: params have been set but AreSet() says they haven't")
	}

	if len(errMap) != 0 {
		t.Log(testName)
		for k, v := range errMap {
			t.Log("\t:", k, ":")
			for _, err := range v {
				t.Log("\t\t", err)
			}
		}
		t.Errorf("\t: unexpected errors were detected while parsing")
	} else if p1 != 99 {
		t.Log(testName)
		t.Errorf("\t: p1 (= %d) was expected to be == 99", p1)
	} else if p2 != 0 {
		t.Log(testName)
		t.Errorf("\t: p2 (= %d) was expected to be == 0", p2)
	} else if ps.ProgName() != param.DfltProgName {
		t.Log(testName)
		t.Errorf("\t: ps.ProgName() ('= %s') was expected to be == '%s'",
			ps.ProgName(), param.DfltProgName)
	}

	_, panicked, panicVal = panicSafeTestAddByName(ps,
		&namedParamInitialiser{
			name:   "test99",
			setter: &psetter.Int64{Value: &p1},
			desc:   "desc - this should not be added",
		})

	testhelper.PanicCheckString(t, "param.Add - adding a param after parsing",
		panicked, true,
		panicVal, []string{
			"Parameters have already been parsed." +
				" A new named parameter (test99) cannot be added",
		})

	_, panicked, panicVal = panicSafeTestAddByPos(ps,
		&posParamInitialiser{
			name:   "ppi1",
			setter: &psetter.Int64{Value: &p1},
		})
	testhelper.PanicCheckString(t,
		"Adding a positional param after parsing",
		panicked, true,
		panicVal, []string{
			"Parameters have already been parsed." +
				" A new positional parameter (ppi1) cannot be added.",
		})

	testName = "param.Parse - bad -  second parse attempt"
	errMap, panicked, panicVal, stackTrace = panicSafeTestParse(ps,
		[]string{"-param1", "99"})
	if !testhelper.ReportUnexpectedPanic(t, testName,
		panicked, panicVal, stackTrace) {
		errMapCheck(t, testName, errMap, map[string][]string{
			"": {
				"param.Parse has already been called, previously from:",
			},
		})
	}
}

func TestParamParse(t *testing.T) {
	var p1 int64
	testCases := []struct {
		testhelper.ID
		expectedEMap map[string][]string
		paramsPassed []string
		params       []*namedParamInitialiser
	}{
		{
			ID: testhelper.MkID("one param, no error - separate param and value"),
			params: []*namedParamInitialiser{
				{
					name:   "test1",
					setter: &psetter.Int64{Value: &p1},
				},
			},
			paramsPassed: []string{
				"-test1", "99",
			},
		},
		{
			ID: testhelper.MkID("one param, no error - param=value"),
			params: []*namedParamInitialiser{
				{
					name:   "test1",
					setter: &psetter.Int64{Value: &p1},
				},
			},
			paramsPassed: []string{
				"-test1=99",
			},
		},
		{
			ID: testhelper.MkID("bad param, doesn't match"),
			params: []*namedParamInitialiser{
				{
					name:   "test1",
					setter: &psetter.Int64{Value: &p1},
				},
			},
			paramsPassed: []string{
				"-test2",
			},
			expectedEMap: map[string][]string{
				"test2": {"is not a parameter"},
			},
		},
		{
			ID: testhelper.MkID("bad param, no second part"),
			params: []*namedParamInitialiser{
				{
					name:   "test1",
					setter: &psetter.Int64{Value: &p1},
				},
			},
			paramsPassed: []string{
				"-test1",
			},
			expectedEMap: map[string][]string{
				"test1": {"error with parameter"},
			},
		},
		{
			ID: testhelper.MkID("bad param, not a number"),
			params: []*namedParamInitialiser{
				{
					name:   "test1",
					setter: &psetter.Int64{Value: &p1},
				},
			},
			paramsPassed: []string{
				"-test1", "this is not a number",
			},
			expectedEMap: map[string][]string{
				"test1": {"error with parameter"},
			},
		},
	}

	for _, tc := range testCases {
		ps, err := paramset.NewNoHelpNoExitNoErrRpt()
		if err != nil {
			t.Fatal("An error was detected while constructing the PSet:",
				err)
		}
		for _, p := range tc.params {
			panicSafeTestAddByName(ps, p)
		}
		errMap, panicked, panicVal, stackTrace := panicSafeTestParse(ps, tc.paramsPassed)
		if !testhelper.ReportUnexpectedPanic(t, tc.IDStr(),
			panicked, panicVal, stackTrace) {
			errMapCheck(t, tc.IDStr(), errMap, tc.expectedEMap)
		}
	}

}

func TestParamByName(t *testing.T) {
	var val1 int64 = 123
	val1InitialVal := fmt.Sprint(val1)
	ps, err := paramset.NewNoHelpNoExitNoErrRpt()
	if err != nil {
		t.Fatal("Couldn't create the PSet: ", err)
	}

	const (
		param1Name    = "param1"
		param1AltName = "p1"
	)

	p := ps.Add(param1Name,
		psetter.Int64{
			Value: &val1,
		},
		"an int64 parameter",
		param.AltName(param1AltName))

	altNames := p.AltNames()
	if len(altNames) != 2 {
		t.Error("the parameter should have had two names, has: ",
			len(altNames))
	} else {
		if altNames[0] != param1Name {
			t.Error("the parameter should have had ",
				param1Name, " as a principal name: ",
				" was ", altNames[0])
		}
		if altNames[1] != param1AltName {
			t.Error("the parameter should have had as an alternative name: ",
				param1AltName, " was ", altNames[1])
		}
	}

	ws := p.WhereSet()
	if len(ws) != 0 {
		t.Error(
			"Parse has not been called but the parameter WhereSet list has: ",
			len(ws), " entries")
	}
	if p.InitialValue() != val1InitialVal {
		t.Error(
			"the initial value of the parameter was expected to be: ",
			val1InitialVal, "  but was: ", p.InitialValue())
	}

	if pvr := p.ValueReq(); pvr != param.Mandatory {
		t.Error("a parameter value should be required for an Int64"+
			" ValueReq() returned: ", pvr.String())
	}

	if p.AllowedValues() == "" {
		t.Error(
			"a non-empty allowed values string is expected for an Int64")
	}
}
