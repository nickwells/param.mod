package param_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nickwells/param.mod/v2/param"
	"github.com/nickwells/param.mod/v2/param/paramset"
	"github.com/nickwells/param.mod/v2/param/psetter"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestParamAdd(t *testing.T) {
	var p1 int64
	testCases := []struct {
		name             string
		npi              *namedParamInitialiser
		panicExpected    bool
		panicMsgContains []string
		paramShouldExist bool
	}{
		{
			name: "bad name - empty",
			npi: &namedParamInitialiser{
				setter: &psetter.Int64Setter{Value: &p1},
			},
			panicExpected: true,
			panicMsgContains: []string{
				"the parameter name",
				"is invalid. It must match",
			},
		},
		{
			name: "bad name - bad char",
			npi: &namedParamInitialiser{
				name:   "?",
				setter: &psetter.Int64Setter{Value: &p1},
			},
			panicExpected: true,
			panicMsgContains: []string{
				"the parameter name",
				"is invalid. It must match",
			},
		},
		{
			name: "bad name - bad first char",
			npi: &namedParamInitialiser{
				name:   "-hello",
				setter: &psetter.Int64Setter{Value: &p1},
			},
			panicExpected: true,
			panicMsgContains: []string{
				"the parameter name",
				"is invalid. It must match",
			},
		},
		{
			name: "good name",
			npi: &namedParamInitialiser{
				name:   "param-1",
				setter: &psetter.Int64Setter{Value: &p1},
				opts: []param.OptFunc{
					param.AltName("param-1-alt"),
					param.GroupName("test"),
				},
			},
			paramShouldExist: true,
		},
		{
			name: "bad name - duplicate",
			npi: &namedParamInitialiser{
				name:   "param-1",
				setter: &psetter.Int64Setter{Value: &p1},
			},
			panicExpected: true,
			panicMsgContains: []string{
				"parameter name",
				"has already been used",
			},
			paramShouldExist: true,
		},
		{
			name: "bad alt name - already used",
			npi: &namedParamInitialiser{
				name:   "param-2",
				setter: &psetter.Int64Setter{Value: &p1},
				opts:   []param.OptFunc{param.AltName("param-1")},
			},
			panicExpected: true,
			panicMsgContains: []string{
				"parameter name",
				"has already been used",
			},
			paramShouldExist: true,
		},
		{
			name: "bad alt name - invalid",
			npi: &namedParamInitialiser{
				name:   "param-3",
				setter: &psetter.Int64Setter{Value: &p1},
				opts:   []param.OptFunc{param.AltName("?")},
			},
			panicExpected: true,
			panicMsgContains: []string{
				"parameter name",
				"is invalid. It must match",
			},
			paramShouldExist: true,
		},
		{
			name: "bad alt name - already used as alt",
			npi: &namedParamInitialiser{
				name:   "param-4",
				setter: &psetter.Int64Setter{Value: &p1},
				opts:   []param.OptFunc{param.AltName("param-1-alt")},
			},
			panicExpected: true,
			panicMsgContains: []string{
				"parameter name",
				"has already been used as an alternative to 'param-1'",
				"a member of parameter group test",
			},
			paramShouldExist: true,
		},
	}

	ps, err := paramset.NewNoHelpNoExitNoErrRpt()
	if err != nil {
		t.Fatal("couldn't construct the ParamSet: ", err)
	}
	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		p, panicked, panicVal := panicSafeTestAddByName(ps, tc.npi)
		testhelper.PanicCheckString(t, tcID,
			panicked, tc.panicExpected,
			panicVal, tc.panicMsgContains)

		p2, err := ps.GetParamByName(tc.npi.name)
		if tc.paramShouldExist {
			if p2 == nil {
				t.Log(tcID)
				t.Errorf("\t: param: '%s' should exist\n",
					tc.npi.name)
			}
			if err != nil {
				t.Log(tcID)
				t.Errorf("\t: GetParamByName(...) returned an error: %s\n",
					err)
			}
		} else {
			if p2 != nil {
				t.Log(tcID)
				t.Errorf(
					"\t: param: '%s' should not exist if Add(...) panicked\n",
					tc.npi.name)
			}
			if err == nil {
				t.Log(tcID)
				t.Errorf(
					"\t: GetParamByName(...) should have returned an error\n")
			}
		}

		if p != nil {
			if p.Name() != tc.npi.name {
				t.Log(tcID)
				t.Errorf("\t: the name did not match: '%s' != '%s'\n",
					p.Name(), tc.npi.name)
			} else if p.Description() != tc.npi.desc {
				t.Log(tcID)
				t.Errorf("\t: the description did not match: '%s' != '%s'\n",
					p.Description(), tc.npi.desc)
			} else if p.HasBeenSet() {
				t.Log(tcID)
				t.Errorf(
					"\t: param has been set but params haven't been parsed\n")
			}
		}
	}
}

func TestParamAddPos(t *testing.T) {
	var p1 int64
	testCases := []struct {
		name             string
		pi               []paramInitialisers
		panicExpected    bool
		panicMsgContains []string

		paramsToParse     []string
		errsExpected      map[string][]string
		remainderExpected []string
	}{
		{
			name: "good params",
			pi: []paramInitialisers{
				paramInitialisers{
					npi: &namedParamInitialiser{
						name:   "param-1",
						setter: &psetter.Int64Setter{Value: &p1},
						opts: []param.OptFunc{
							param.AltName("param-1-alt"),
							param.GroupName("test"),
						},
					},
					npiShouldExist: true,
				},
				paramInitialisers{
					npi: &namedParamInitialiser{
						name:   "param-2",
						setter: &psetter.Int64Setter{Value: &p1},
						opts: []param.OptFunc{
							param.AltName("param-2-alt"),
							param.GroupName("test"),
						},
					},
					npiShouldExist: true,
				},
				paramInitialisers{
					ppi: &posParamInitialiser{
						name:   "ppi1",
						setter: &psetter.Int64Setter{Value: &p1},
					},
					ppiShouldExist: true,
				},
			},
			paramsToParse:     []string{"1", "--", "another param"},
			remainderExpected: []string{"another param"},
		},
		{
			name: "bad params - name empty",
			pi: []paramInitialisers{
				paramInitialisers{
					npi: &namedParamInitialiser{
						name:   "",
						setter: &psetter.Int64Setter{Value: &p1},
						opts: []param.OptFunc{
							param.GroupName("test"),
						},
					},
				},
			},
			panicExpected: true,
			panicMsgContains: []string{
				"the parameter name",
				"is invalid. It must match",
			},
		},
		{
			name: "bad params - terminal ByPos and pre-existing ByName",
			pi: []paramInitialisers{
				paramInitialisers{
					npi: &namedParamInitialiser{
						name:   "param-1",
						setter: &psetter.Int64Setter{Value: &p1},
						opts: []param.OptFunc{
							param.AltName("param-1-alt"),
							param.GroupName("test"),
						},
					},
					npiShouldExist: true,
				},
				paramInitialisers{
					ppi: &posParamInitialiser{
						name:   "ppi1",
						setter: &psetter.Int64Setter{Value: &p1},
						opts: []param.PosOptFunc{
							param.SetAsTerminal,
						},
					},
				},
			},
			panicExpected: true,
			panicMsgContains: []string{
				"error setting the options for positional parameter 0:",
				"The param set has 1 non-positional parameters.",
				" It cannot also have a terminal positional parameter as" +
					" the non-positional parameters will never be used.",
			},
		},
		{
			name: "bad params - terminal ByPos and then add ByName",
			pi: []paramInitialisers{
				paramInitialisers{
					ppi: &posParamInitialiser{
						name:   "ppi1",
						setter: &psetter.Int64Setter{Value: &p1},
						opts: []param.PosOptFunc{
							param.SetAsTerminal,
						},
					},
					ppiShouldExist: true,
				},
				paramInitialisers{
					npi: &namedParamInitialiser{
						name:   "param-1",
						setter: &psetter.Int64Setter{Value: &p1},
						opts: []param.OptFunc{
							param.AltName("param-1-alt"),
							param.GroupName("test"),
						},
					},
				},
			},
			panicExpected: true,
			panicMsgContains: []string{
				"The param set has a terminal positional parameter.",
				"The non-positional parameter param-1 cannot be added" +
					" as it will never be used",
			},
		},
		{
			name: "bad params - terminal ByPos not the last",
			pi: []paramInitialisers{
				paramInitialisers{
					ppi: &posParamInitialiser{
						name:   "ppi1",
						setter: &psetter.Int64Setter{Value: &p1},
						opts: []param.PosOptFunc{
							param.SetAsTerminal,
						},
					},
					ppiShouldExist: true,
				},
				paramInitialisers{
					ppi: &posParamInitialiser{
						name:   "ppi2",
						setter: &psetter.Int64Setter{Value: &p1},
						opts:   []param.PosOptFunc{},
					},
				},
			},
			panicExpected: true,
			panicMsgContains: []string{
				"Positional parameter 0 is marked as terminal but" +
					" is not the last positional parameter",
			},
		},
		{
			name: "Parse(...) errors - ByPos: setter.Process(...) failure",
			pi: []paramInitialisers{
				paramInitialisers{
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
				"Positional parameter: 1 (ppi1)": []string{
					"failingSetter: ByPos param",
				},
			},
		},
		{
			name: "Parse(...) errors - ByName: setter.Process(...) failure",
			pi: []paramInitialisers{
				paramInitialisers{
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
				"test99": []string{
					"error with parameter:",
					"failingSetter: ByName param",
				},
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		ps, err := paramset.NewNoHelpNoExitNoErrRpt()
		if err != nil {
			t.Fatal(tcID, " : couldn't construct the ParamSet: ", err)
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
						t.Log(tcID)
						t.Errorf("\t: named parameter: '%s'"+
							" should not exist but does\n",
							pi.npi.name)
					} else {
						t.Log(tcID)
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
						t.Log(tcID)
						t.Errorf("\t: positional parameter: %d"+
							" should not exist but does\n",
							posIdx)
					} else {
						t.Log(tcID)
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
		testhelper.PanicCheckString(t, tcID,
			panicked, tc.panicExpected,
			panicVal, tc.panicMsgContains)

		if !panicked {
			errMap, panicked, panicVal, stackTrace :=
				panicSafeTestParse(ps, tc.paramsToParse)
			if !testhelper.PanicCheckStringWithStack(t, tcID,
				panicked, false,
				panicVal, []string{}, stackTrace) {
				errMapCheck(t, tcID, errMap, tc.errsExpected)

				if tc.remainderExpected != nil &&
					len(tc.remainderExpected) == 0 {
					tc.remainderExpected = nil
				}

				if !reflect.DeepEqual(ps.Remainder(), tc.remainderExpected) {
					t.Log(tcID)
					t.Logf("\t: remainder received: %v\n", ps.Remainder())
					t.Logf("\t: remainder expected: %v\n", tc.remainderExpected)
					t.Errorf("\t: unexpected remainder\n")
				}
			}
		}
	}
}

func TestParamParse1(t *testing.T) {
	var p1 int64
	var p2 int64

	ps, err := paramset.NewNoHelpNoExitNoErrRpt()

	if err != nil {
		t.Fatal("couldn't construct the ParamSet: ", err)
	}

	testName := "param.Parse - first Parse"
	ps.Add("param1", &psetter.Int64Setter{Value: &p1}, "a param")
	ps.Add("param2", &psetter.Int64Setter{Value: &p2}, "a param")

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
			setter: &psetter.Int64Setter{Value: &p1},
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
			setter: &psetter.Int64Setter{Value: &p1},
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
			"": []string{
				"param.Parse has already been called, previously from:",
			},
		})
	}
}

func TestParamParse(t *testing.T) {
	var p1 int64
	testCases := []struct {
		testName     string
		expectedEMap map[string][]string
		paramsPassed []string
		params       []*namedParamInitialiser
	}{
		{
			testName: "one param, no error - separate param and value",
			params: []*namedParamInitialiser{
				&namedParamInitialiser{
					name:   "test1",
					setter: &psetter.Int64Setter{Value: &p1},
				},
			},
			paramsPassed: []string{
				"-test1", "99",
			},
		},
		{
			testName: "one param, no error - param=value",
			params: []*namedParamInitialiser{
				&namedParamInitialiser{
					name:   "test1",
					setter: &psetter.Int64Setter{Value: &p1},
				},
			},
			paramsPassed: []string{
				"-test1=99",
			},
		},
		{
			testName: "bad param, doesn't match",
			params: []*namedParamInitialiser{
				&namedParamInitialiser{
					name:   "test1",
					setter: &psetter.Int64Setter{Value: &p1},
				},
			},
			paramsPassed: []string{
				"-test2",
			},
			expectedEMap: map[string][]string{
				"test2": []string{"is not a parameter"},
			},
		},
		{
			testName: "bad param, no second part",
			params: []*namedParamInitialiser{
				&namedParamInitialiser{
					name:   "test1",
					setter: &psetter.Int64Setter{Value: &p1},
				},
			},
			paramsPassed: []string{
				"-test1",
			},
			expectedEMap: map[string][]string{
				"test1": []string{"error with parameter"},
			},
		},
		{
			testName: "bad param, not a number",
			params: []*namedParamInitialiser{
				&namedParamInitialiser{
					name:   "test1",
					setter: &psetter.Int64Setter{Value: &p1},
				},
			},
			paramsPassed: []string{
				"-test1", "this is not a number",
			},
			expectedEMap: map[string][]string{
				"test1": []string{"error with parameter"},
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.testName)
		ps, err := paramset.NewNoHelpNoExitNoErrRpt()
		if err != nil {
			t.Fatal("An error was detected while constructing the ParamSet:",
				err)
		}
		for _, p := range tc.params {
			panicSafeTestAddByName(ps, p)
		}
		errMap, panicked, panicVal, stackTrace := panicSafeTestParse(ps, tc.paramsPassed)
		if !testhelper.ReportUnexpectedPanic(t, tcID,
			panicked, panicVal, stackTrace) {
			errMapCheck(t, tcID, errMap, tc.expectedEMap)
		}
	}

}

func TestParamByName(t *testing.T) {
	var val1 int64 = 123
	val1InitialVal := fmt.Sprint(val1)
	ps, err := paramset.NewNoHelpNoExitNoErrRpt()
	if err != nil {
		t.Fatal("Couldn't create the ParamSet: ", err)
	}

	const (
		param1Name    = "param1"
		param1AltName = "p1"
	)

	p := ps.Add(param1Name,
		psetter.Int64Setter{
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
			t.Error("the parameter should have had as a pricipal name: ",
				param1Name, " was ", altNames[0])
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
		t.Error("a parameter value should be required for an Int64Setter"+
			" ValueReq() returned: ", pvr.String())
	}

	if p.AllowedValues() == "" {
		t.Error(
			"a non-empty allowed values string is expected for an Int64Setter")
	}
}
