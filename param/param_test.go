package param_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nickwells/param.mod/v4/param"
	"github.com/nickwells/param.mod/v4/param/psetter"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestParamAdd(t *testing.T) {
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
				setter: psetter.Int64{Value: &p1},
			},
			ExpPanic: testhelper.MkExpPanic(
				"the parameter name",
				"is invalid. It must match"),
		},
		{
			ID: testhelper.MkID("bad name - bad char"),
			npi: &namedParamInitialiser{
				name:   "?",
				setter: psetter.Int64{Value: &p1},
			},
			ExpPanic: testhelper.MkExpPanic(
				"the parameter name",
				"is invalid. It must match"),
		},
		{
			ID: testhelper.MkID("bad name - bad first char"),
			npi: &namedParamInitialiser{
				name:   "-hello",
				setter: psetter.Int64{Value: &p1},
			},
			ExpPanic: testhelper.MkExpPanic(
				"the parameter name",
				"is invalid. It must match"),
		},
		{
			ID: testhelper.MkID("good name"),
			npi: &namedParamInitialiser{
				name:   "param-1",
				setter: psetter.Int64{Value: &p1},
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
				setter: psetter.Int64{Value: &p1},
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
				setter: psetter.Int64{Value: &p1},
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
				setter: psetter.Int64{Value: &p1},
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
				setter: psetter.Int64{Value: &p1},
				opts:   []param.OptFunc{param.AltName("param-1-alt")},
			},
			ExpPanic: testhelper.MkExpPanic(
				"parameter name",
				`has already been used as an alternative to "param-1"`,
				`a member of parameter group "test"`),
			paramShouldExist: true,
		},
	}

	ps := makePSetOrFatal(t, t.Name())
	for _, tc := range testCases {
		p, panicked, panicVal := panicSafeTestAddByName(ps, tc.npi)
		testhelper.CheckExpPanic(t, panicked, panicVal, tc)

		p2, err := ps.GetParamByName(tc.npi.name)
		if tc.paramShouldExist {
			checkByNameVal(t, tc.IDStr(), tc.npi.name, p2, err, ShouldBeSet)
		} else {
			checkByNameVal(t, tc.IDStr(), tc.npi.name, p2, err, ShouldNotBeSet)
		}

		tc.npi.compare(t, tc.IDStr(), p, ShouldNotBeSet)
	}
}

// checkByNameVal checks the ByName value and the associated error
func checkByNameVal(t *testing.T, tcID, paramName string,
	p *param.ByName, err error,
	sbs ShouldBeSetType) {
	t.Helper()

	if sbs == ShouldBeSet {
		if p == nil {
			t.Log(tcID)
			t.Errorf("\t: the ByName param: %q should exist\n",
				paramName)
		}
		if err != nil {
			t.Log(tcID)
			t.Errorf("\t: error seen when getting the ByName param: %q: %v\n",
				paramName, err)
		}
	} else {
		if p != nil {
			t.Log(tcID)
			t.Errorf("\t: the ByName param: %q should not exist\n",
				paramName)
		}
		if err == nil {
			t.Log(tcID)
			t.Errorf("\t: no error seen when getting the ByName param: %q\n",
				paramName)
		}
	}
}

// checkByPosVal checks the ByPos value and the associated error
func checkByPosVal(t *testing.T, tcID string, paramIdx int,
	p *param.ByPos, err error,
	sbs ShouldBeSetType) {
	t.Helper()

	if sbs == ShouldBeSet {
		if p == nil {
			t.Log(tcID)
			t.Errorf("\t: the ByPos param: %d should exist\n",
				paramIdx)
		}
		if err != nil {
			t.Log(tcID)
			t.Errorf("\t: error seen when getting the ByPos param: %d: %v\n",
				paramIdx, err)
		}
	} else {
		if p != nil {
			t.Log(tcID)
			t.Errorf("\t: the ByPos param: %d should not exist\n",
				paramIdx)
		}
		if err == nil {
			t.Log(tcID)
			t.Errorf("\t: missing error when getting the ByPos param: %d\n",
				paramIdx)
		}
	}
}

func TestParamAddPos(t *testing.T) {
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
						setter: psetter.Int64{Value: &p1},
						opts: []param.OptFunc{
							param.AltName("param-1-alt"),
							param.GroupName("test"),
						},
					},
					npiSBS: ShouldBeSet,
				},
				{
					npi: &namedParamInitialiser{
						name:   "param-2",
						setter: psetter.Int64{Value: &p1},
						opts: []param.OptFunc{
							param.AltName("param-2-alt"),
							param.GroupName("test"),
						},
					},
					npiSBS: ShouldBeSet,
				},
				{
					ppi: &posParamInitialiser{
						name:   "ppi1",
						setter: psetter.Int64{Value: &p1},
					},
					ppiSBS: ShouldBeSet,
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
						setter: psetter.Int64{Value: &p1},
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
						setter: psetter.Int64{Value: &p1},
						opts: []param.OptFunc{
							param.AltName("param-1-alt"),
							param.GroupName("test"),
						},
					},
					npiSBS: ShouldBeSet,
				},
				{
					ppi: &posParamInitialiser{
						name:   "ppi1",
						setter: psetter.Int64{Value: &p1},
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
						setter: psetter.Int64{Value: &p1},
						opts: []param.PosOptFunc{
							param.SetAsTerminal,
						},
					},
					ppiSBS: ShouldBeSet,
				},
				{
					npi: &namedParamInitialiser{
						name:   "param-1",
						setter: psetter.Int64{Value: &p1},
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
						setter: psetter.Int64{Value: &p1},
						opts: []param.PosOptFunc{
							param.SetAsTerminal,
						},
					},
					ppiSBS: ShouldBeSet,
				},
				{
					ppi: &posParamInitialiser{
						name:   "ppi2",
						setter: psetter.Int64{Value: &p1},
						opts:   []param.PosOptFunc{},
					},
				},
			},
			ExpPanic: testhelper.MkExpPanic(
				"Positional parameter 0 is marked as terminal but" +
					" is not the last positional parameter"),
		},
		{
			ID: testhelper.MkID(
				"Parse(...) errors - ByPos: setter.Process(...) failure"),
			pi: []paramInitialisers{
				{
					ppi: &posParamInitialiser{
						name:   "ppi1",
						setter: &failingSetter{errMsg: "ByPos param"},
					},
					ppiSBS: ShouldBeSet,
				},
			},
			paramsToParse: []string{"1"},
			errsExpected: map[string][]string{
				"Positional parameter: 1 (ppi1)": {
					"failingSetter: ByPos param",
				},
			},
		},
		{
			ID: testhelper.MkID(
				"Parse(...) errors - ByName: setter.Process(...) failure"),
			pi: []paramInitialisers{
				{
					npi: &namedParamInitialiser{
						name:   "test99",
						setter: &failingSetter{errMsg: "ByName param"},
					},
					npiSBS: ShouldBeSet,
				},
			},
			paramsToParse: []string{"-test99", "val"},
			errsExpected: map[string][]string{
				"test99": {
					"failingSetter: ByName param",
				},
			},
		},
	}

	for _, tc := range testCases {
		ps := makePSetOrFatal(t, tc.IDStr())
		var panicked bool
		var panicVal interface{}

		var posIdx int
		for _, pi := range tc.pi {
			if pi.npi != nil {
				_, panicked, panicVal = panicSafeTestAddByName(ps, pi.npi)

				np, err := ps.GetParamByName(pi.npi.name)
				checkByNameVal(t, tc.IDStr(), pi.npi.name, np, err, pi.npiSBS)

				if panicked {
					break
				}
			}

			if pi.ppi != nil {
				_, panicked, panicVal = panicSafeTestAddByPos(ps, pi.ppi)

				pp, err := ps.GetParamByPos(posIdx)
				checkByPosVal(t, tc.IDStr(), posIdx, pp, err, pi.ppiSBS)

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

// parsePSet will check that parsing hasn't happened yet and then parse the
// args with the pset and afterwards check that parsing has happened
func parsePSet(t *testing.T, ps *param.PSet, args []string) param.ErrMap {
	t.Helper()

	if ps.AreSet() {
		t.Log(t.Name())
		t.Errorf("\t: params haven't been set but AreSet() says they have")
	}

	errMap, panicked, panicVal, stackTrace := panicSafeTestParse(ps, args)

	if testhelper.ReportUnexpectedPanic(t, t.Name(),
		panicked, panicVal, stackTrace) {
		return errMap
	}

	if !ps.AreSet() {
		t.Log(t.Name())
		t.Errorf("\t: params have been set but AreSet() says they haven't")
	}

	return errMap
}

func TestParamParseTwice(t *testing.T) {
	ps := makePSetOrFatal(t, t.Name())
	errMap := parsePSet(t, ps, []string{})

	if logErrMap(t, errMap) {
		t.Errorf("\t: unexpected errors were detected while parsing")
	}

	_, panicked, panicVal, _ := panicSafeTestParse(ps, []string{})
	testhelper.PanicCheckString(t, t.Name(),
		panicked, true,
		panicVal, []string{
			"param.Parse has already been called, previously from:",
		})
}

func TestParamAddParamAfterParse(t *testing.T) {
	var p1 int64

	ps := makePSetOrFatal(t, t.Name())
	errMap := parsePSet(t, ps, []string{})

	if logErrMap(t, errMap) {
		t.Errorf("\t: unexpected errors were detected while parsing")
	}
	_, panicked, panicVal := panicSafeTestAddByName(ps,
		&namedParamInitialiser{
			name:   "test99",
			setter: psetter.Int64{Value: &p1},
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
			setter: psetter.Int64{Value: &p1},
		})
	testhelper.PanicCheckString(t,
		"Adding a positional param after parsing",
		panicked, true,
		panicVal, []string{
			"Parameters have already been parsed." +
				" A new positional parameter (ppi1) cannot be added.",
		})
}

func TestParamParse1(t *testing.T) {
	var p1 int64
	var p2 int64

	ps := makePSetOrFatal(t, t.Name())

	testCases := []struct {
		testhelper.ID
		name        string
		setter      psetter.Int64
		desc        string
		expectedVal int64
	}{
		{
			ID:          testhelper.MkID("param1"),
			name:        "param1",
			setter:      psetter.Int64{Value: &p1},
			desc:        "param 1",
			expectedVal: 99,
		},
		{
			ID:          testhelper.MkID("param2"),
			name:        "param2",
			setter:      psetter.Int64{Value: &p2},
			desc:        "param 2",
			expectedVal: 0,
		},
	}

	for _, tc := range testCases {
		ps.Add(tc.name, tc.setter, tc.desc)
	}

	errMap := parsePSet(t, ps, []string{"-param1", "99"})

	if logErrMap(t, errMap) {
		t.Errorf("\t: unexpected errors were detected while parsing")
	}
	for _, tc := range testCases {
		if *tc.setter.Value != tc.expectedVal {
			t.Log(tc.IDStr())
			t.Errorf("\t: the value (= %d) was expected to be %d",
				*tc.setter.Value, tc.expectedVal)
		}
	}
	if ps.ProgName() != param.DfltProgName {
		t.Log(t.Name())
		t.Errorf("\t: ps.ProgName() ('= %s') was expected to be == '%s'",
			ps.ProgName(), param.DfltProgName)
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
			ID: testhelper.MkID(
				"one param, no error - separate param and value"),
			params: []*namedParamInitialiser{
				{
					name:   "test1",
					setter: psetter.Int64{Value: &p1},
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
					setter: psetter.Int64{Value: &p1},
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
					setter: psetter.Int64{Value: &p1},
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
					setter: psetter.Int64{Value: &p1},
				},
			},
			paramsPassed: []string{
				"-test1",
			},
			expectedEMap: map[string][]string{
				"test1": {
					"a value must follow this parameter",
					"either following an '=' or as a next parameter",
				},
			},
		},
		{
			ID: testhelper.MkID("bad param, not a number"),
			params: []*namedParamInitialiser{
				{
					name:   "test1",
					setter: psetter.Int64{Value: &p1},
				},
			},
			paramsPassed: []string{
				"-test1", "this is not a number",
			},
			expectedEMap: map[string][]string{
				"test1": {`could not interpret "this is not a number"` +
					` as a whole number`},
			},
		},
	}

	for _, tc := range testCases {
		ps := makePSetOrFatal(t, tc.IDStr())

		for _, p := range tc.params {
			panicSafeTestAddByName(ps, p)
		}
		errMap, panicked, panicVal, stackTrace :=
			panicSafeTestParse(ps, tc.paramsPassed)
		if !testhelper.ReportUnexpectedPanic(t, tc.IDStr(),
			panicked, panicVal, stackTrace) {
			errMapCheck(t, tc.IDStr(), errMap, tc.expectedEMap)
		}
	}

}

func TestParamByName(t *testing.T) {
	var val1 int64 = 123
	val1InitialVal := fmt.Sprint(val1)
	ps := makePSetOrFatal(t, t.Name())

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

	if pvr := p.Setter().ValueReq(); pvr != param.Mandatory {
		t.Error("a parameter value should be required for an Int64"+
			" ValueReq() returned: ", pvr.String())
	}

	if p.Setter().AllowedValues() == "" {
		t.Error(
			"a non-empty allowed values string is expected for an Int64")
	}
}
