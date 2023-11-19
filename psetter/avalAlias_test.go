package psetter_test

import (
	"testing"

	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestAliasesString(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		a      psetter.Aliases[string]
		expStr string
	}{
		{
			ID: testhelper.MkID("nil"),
		},
		{
			ID: testhelper.MkID("with entries"),
			a: psetter.Aliases[string]{
				"first":  {"one", "two", "three"},
				"second": {"blah", "blahTwo", "blahThree"},
				"third":  {},
			},
			expStr: "   first : one, two, three\n" +
				"   second: blah, blahTwo, blahThree\n" +
				"   third : ",
		},
	}

	for _, tc := range testCases {
		actStr := tc.a.String()
		testhelper.DiffString[string](t, tc.IDStr(), "String()",
			actStr, tc.expStr)
	}
}

func TestAliasesAllowedValuesAliasMap(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		a      psetter.Aliases[string]
		expVal psetter.Aliases[string]
	}{
		{
			ID:     testhelper.MkID("empty"),
			a:      psetter.Aliases[string]{},
			expVal: psetter.Aliases[string]{},
		},
		{
			ID: testhelper.MkID("not empty"),
			a: psetter.Aliases[string]{
				"hello": {"world"},
			},
			expVal: psetter.Aliases[string]{
				"hello": {"world"},
			},
		},
	}

	for _, tc := range testCases {
		actVal := tc.a.AllowedValuesAliasMap()
		err := testhelper.DiffVals(actVal, tc.expVal)
		if err != nil {
			t.Log(tc.IDStr())
			t.Errorf("\t: bad copy: %s\n", err)
		}
	}
}

func TestAliasesCheck(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		a  psetter.Aliases[string]
		av psetter.AllowedVals[string]
	}{
		{
			ID: testhelper.MkID("nil"),
		},
		{
			ID: testhelper.MkID("empty"),
			a:  psetter.Aliases[string]{},
		},
		{
			ID: testhelper.MkID("good"),
			a: psetter.Aliases[string]{
				"aval1": {"val1"},
			},
			av: psetter.AllowedVals[string]{
				"val1": "val1 description",
			},
		},
		{
			ID: testhelper.MkID("bad - empty alias value"),
			ExpErr: testhelper.MkExpErr(
				`bad alias: "aval1": []string{} - the alias maps to no values`),
			a: psetter.Aliases[string]{
				"aval1": {},
			},
			av: psetter.AllowedVals[string]{
				"val1": "val1 description",
			},
		},
		{
			ID: testhelper.MkID("bad - alias matches value"),
			ExpErr: testhelper.MkExpErr(
				`bad alias: "val1": []string{"aval1"} -` +
					` an allowed value has the same name`),
			a: psetter.Aliases[string]{
				"val1": {"aval1"},
			},
			av: psetter.AllowedVals[string]{
				"val1":  "val1 description",
				"aval1": "aval1 description",
			},
		},
		{
			ID: testhelper.MkID("bad - blank alias name"),
			ExpErr: testhelper.MkExpErr(
				`bad alias: "": []string{"val1"} -` +
					` the alias name must not be blank`),
			a: psetter.Aliases[string]{
				"": {"val1"},
			},
			av: psetter.AllowedVals[string]{
				"val1": "val1 description",
			},
		},
		{
			ID: testhelper.MkID("bad - bad alias name"),
			ExpErr: testhelper.MkExpErr(
				`bad alias: "aval1=x": []string{"val1"} -` +
					` the alias name must not contain '='`),
			a: psetter.Aliases[string]{
				"aval1=x": {"val1"},
			},
			av: psetter.AllowedVals[string]{
				"val1": "val1 description",
			},
		},
		{
			ID: testhelper.MkID("bad - duplicate alias value"),
			ExpErr: testhelper.MkExpErr(
				`bad alias: "aval1": []string{"val1", "val2", "val1"} -` +
					` "val1" appears more than once (at index 0 and 2)`),
			a: psetter.Aliases[string]{
				"aval1": {"val1", "val2", "val1"},
			},
			av: psetter.AllowedVals[string]{
				"val1": "val1 description",
				"val2": "val2 description",
			},
		},
		{
			ID: testhelper.MkID("bad - bad alias value"),
			ExpErr: testhelper.MkExpErr(
				`bad alias: "aval1": []string{"val1", "val2", "val3"} -` +
					` "val3" (at index 2) is unknown`),
			a: psetter.Aliases[string]{
				"aval1": {"val1", "val2", "val3"},
			},
			av: psetter.AllowedVals[string]{
				"val1": "val1 description",
				"val2": "val2 description",
			},
		},
		{
			ID: testhelper.MkID("bad - multiple problems"),
			ExpErr: testhelper.MkExpErr(
				`bad aliases: (6)
"": []string{"val1", "val2", "val99"}
    - "val99" (at index 2) is unknown
    - the alias name must not be blank
"aval1": []string{"val1", "val2", "val99"} - "val99" (at index 2) is unknown
"aval3": []string{} - the alias maps to no values
"aval4=": []string{"val1", "val2", "val99"}
    - "val99" (at index 2) is unknown
    - the alias name must not contain '='
"aval6": []string{"val1", "val1", "val1", "val99", "val99"}
    - "val1" appears more than once (at index 0, 1 and 2)
    - "val99" (at index 3 and 4) is unknown
    - "val99" appears more than once (at index 3 and 4)
"val5": []string{"val1"} - an allowed value has the same name as the alias`),
			a: psetter.Aliases[string]{
				"":       {"val1", "val2", "val99"},
				"aval1":  {"val1", "val2", "val99"},
				"aval3":  {},
				"aval4=": {"val1", "val2", "val99"},
				"aval6":  {"val1", "val1", "val1", "val99", "val99"},
				"val5":   {"val1"},
			},
			av: psetter.AllowedVals[string]{
				"val1": "val1 description",
				"val2": "val2 description",
				"val3": "val3 description",
				"val4": "val4 description",
				"val5": "val5 description",
				"val6": "val6 description",
				"val7": "val7 description",
			},
		},
	}

	for _, tc := range testCases {
		err := tc.a.Check(tc.av)
		testhelper.CheckExpErr(t, err, tc)
	}
}
