package psetter_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestSetWithVal(t *testing.T) {
	paramName := "setval-test"
	vBool := true
	setterBool := psetter.Bool{Value: &vBool}

	vDuration := 1 * time.Millisecond
	setterDuration := psetter.Duration{Value: &vDuration}

	vEnumList := []string{"hello", "world"}
	setterEnumList := psetter.EnumList[string]{
		Value: &vEnumList,
		AllowedVals: psetter.AllowedVals[string]{
			"hello": "hello description",
			"world": "world description",
			"bye":   "bye description",
		},
	}
	setterEnumListWithChecks := psetter.EnumList[string]{
		Value: &vEnumList,
		AllowedVals: psetter.AllowedVals[string]{
			"hello": "hello description",
			"world": "world description",
			"bye":   "bye description",
		},
		Checks: []check.StringSlice{
			check.SliceLength[[]string](check.ValEQ(2)),
		},
	}

	vEnumMap := map[string]bool{
		"hello": true,
	}
	setterEnumMap := psetter.EnumMap[string]{
		Value: &vEnumMap,
		AllowedVals: psetter.AllowedVals[string]{
			"hello": "hello description",
			"world": "world description",
			"bye":   "bye description",
		},
	}

	vEnum := "hello"
	setterEnum := psetter.Enum[string]{
		Value: &vEnum,
		AllowedVals: psetter.AllowedVals[string]{
			"hello": "hello description",
			"world": "world description",
			"bye":   "bye description",
		},
	}

	vFloat64 := float64(3.14159)
	setterFloat64 := psetter.Float[float64]{Value: &vFloat64}
	setterFloat64WithChecks := psetter.Float[float64]{
		Value: &vFloat64,
		Checks: []check.Float64{
			check.ValGT[float64](5),
		},
	}

	vInt64List := []int64{1, 2}
	setterInt64List := psetter.IntList[int64]{Value: &vInt64List}
	setterInt64ListWithChecks := psetter.IntList[int64]{
		Value: &vInt64List,
		Checks: []check.Int64Slice{
			check.SliceHasNoDups[[]int64, int64],
		},
	}

	var vInt64 int64 = 42

	setterInt64 := psetter.Int[int64]{Value: &vInt64}
	setterInt64WithChecks := psetter.Int[int64]{
		Value: &vInt64,
		Checks: []check.Int64{
			check.ValGT[int64](6),
		},
	}

	vMap := map[string]bool{
		"hello": true,
	}
	setterMap := psetter.Map[string]{Value: &vMap}

	vMapEmpty := map[string]bool{}
	setterMapEmpty := psetter.Map[string]{Value: &vMapEmpty}

	setterNil := psetter.Nil{}

	var (
		vPathname  string
		vPathname2 string
		vPathname3 string
	)

	setterPathname := psetter.Pathname{Value: &vPathname}
	setterPathnameWithExpectation := psetter.Pathname{
		Value:       &vPathname2,
		Expectation: filecheck.DirExists(),
	}
	setterPathnameWithChecks := psetter.Pathname{
		Value: &vPathname3,
		Checks: []check.String{
			check.StringHasPrefix[string]("testdata"),
		},
	}

	var vRegexp *regexp.Regexp

	setterRegexp := psetter.Regexp{Value: &vRegexp}

	var vString string

	setterStringWithChecks := psetter.String[string]{
		Value: &vString,
		Checks: []check.String{
			check.StringHasPrefix[string]("hello"),
		},
	}

	var vStrList []string

	setterStrListWithChecks := psetter.StrList[string]{
		Value: &vStrList,
		Checks: []check.StringSlice{
			check.SliceLength[[]string](check.ValEQ(2)),
		},
	}

	var vTimeLocation *time.Location

	setterTimeLocationWithChecks := psetter.TimeLocation{
		Value:  &vTimeLocation,
		Checks: []check.TimeLocation{},
	}

	var vCSInt int

	calcSetterWithChecks := psetter.Calculated[int]{
		Value: &vCSInt,
		CalcMap: map[string]psetter.NamedCalc[int]{
			"std": {
				Name: "standard",
				Calc: func(_, _ string) (int, error) { return 42, nil },
			},
			"std2": {
				Name: "bad standard",
				Calc: func(_, _ string) (int, error) { return 24, nil },
			},
		},
		NoDefault: true,
		Checks:    []check.ValCk[int]{check.ValEQ[int](42)},
	}
	calcSetterWithBadTC := psetter.Calculated[int]{
		Value: &vCSInt,
		CalcMap: map[string]psetter.NamedCalc[int]{
			"std": {
				Name: "standard",
				Calc: func(_, _ string) (int, error) { return 42, nil },
			},
			"std2": {
				Name: "bad standard",
				Calc: func(_, _ string) (int, error) {
					return 24, errors.New("bad value")
				},
			},
		},
		NoDefault: true,
	}

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr

		s      param.Setter
		value  string
		expVal string
	}{
		{
			ID:     testhelper.MkID("bool - false"),
			s:      setterBool,
			value:  "false",
			expVal: "false",
		},
		{
			ID:     testhelper.MkID("bool - true"),
			s:      setterBool,
			value:  "true",
			expVal: "true",
		},
		{
			ID:     testhelper.MkID("duration - 1h"),
			s:      setterDuration,
			value:  "1.5s",
			expVal: "1.5s",
		},
		{
			ID: testhelper.MkID("duration - bad value"),
			ExpErr: testhelper.MkExpErr(
				`could not parse "whoops - a bad duration" as a duration: `),
			s:     setterDuration,
			value: "whoops - a bad duration",
		},
		{
			ID:     testhelper.MkID("enum list - good value"),
			s:      setterEnumList,
			value:  "hello",
			expVal: "hello",
		},
		{
			ID:     testhelper.MkID("enum list - bad value"),
			ExpErr: testhelper.MkExpErr("value is not allowed: \"whoops\""),
			s:      setterEnumList,
			value:  "whoops",
		},
		{
			ID:     testhelper.MkID("enum list with checks - good value"),
			s:      setterEnumListWithChecks,
			value:  "hello,world",
			expVal: "hello,world",
		},
		{
			ID: testhelper.MkID("enum list with checks - bad value"),
			ExpErr: testhelper.MkExpErr("the length of the list (1)" +
				" is incorrect: the value (1) must equal 2"),
			s:     setterEnumListWithChecks,
			value: "hello",
		},
		{
			ID:     testhelper.MkID("enum map - good value"),
			s:      setterEnumMap,
			value:  "hello=false",
			expVal: "hello=false",
		},
		{
			ID: testhelper.MkID("enum map - bad name"),
			ExpErr: testhelper.MkExpErr(
				`bad value: "whoops=false":` +
					` part: 1 ("whoops=false") is invalid.` +
					` The name ("whoops") is not allowed`),
			s:     setterEnumMap,
			value: "whoops=false",
		},
		{
			ID: testhelper.MkID("enum map - bad bool"),
			ExpErr: testhelper.MkExpErr(
				`bad value: "hello=whoops":` +
					` part: 1 ("hello=whoops") is invalid.` +
					` The value ("whoops") cannot be` +
					` interpreted as true or false:`),
			s:     setterEnumMap,
			value: "hello=whoops",
		},
		{
			ID:     testhelper.MkID("enum - good"),
			s:      setterEnum,
			value:  "hello",
			expVal: "hello",
		},
		{
			ID: testhelper.MkID("enum - bad"),
			ExpErr: testhelper.MkExpErr(
				`value not allowed: "whoops"`),
			s:     setterEnum,
			value: "whoops",
		},
		{
			ID:     testhelper.MkID("float64 - good"),
			s:      setterFloat64,
			value:  "1.234",
			expVal: "1.234",
		},
		{
			ID: testhelper.MkID("float64 - bad"),
			ExpErr: testhelper.MkExpErr(
				`could not interpret "whoops" as a number`),
			s:     setterFloat64,
			value: "whoops",
		},
		{
			ID: testhelper.MkID("float64 with checks - bad"),
			ExpErr: testhelper.MkExpErr(
				"the value ",
				" must be greater than"),
			s:     setterFloat64WithChecks,
			value: "4",
		},
		{
			ID:     testhelper.MkID("float64 with checks - good"),
			s:      setterFloat64WithChecks,
			value:  "5.1",
			expVal: "5.1",
		},
		{
			ID:     testhelper.MkID("int64 list - good"),
			s:      setterInt64List,
			value:  "5,6",
			expVal: "5,6",
		},
		{
			ID: testhelper.MkID("int64 list - bad"),
			ExpErr: testhelper.MkExpErr(
				`bad value: "5,whoops": part: 2 (whoops) cannot` +
					` be interpreted as a whole number`),
			s:     setterInt64List,
			value: "5,whoops",
		},
		{
			ID: testhelper.MkID("int64 list with checks - bad"),
			ExpErr: testhelper.MkExpErr(
				`duplicate list entries: 0 and 1 are both: 5`),
			s:     setterInt64ListWithChecks,
			value: "5,5",
		},
		{
			ID:     testhelper.MkID("int64 list with checks - good"),
			s:      setterInt64ListWithChecks,
			value:  "5,6",
			expVal: "5,6",
		},
		{
			ID:     testhelper.MkID("int64 - good"),
			s:      setterInt64,
			value:  "5",
			expVal: "5",
		},
		{
			ID: testhelper.MkID("int64 - bad"),
			ExpErr: testhelper.MkExpErr(
				`could not interpret "whoops" as a whole number: `),
			s:     setterInt64,
			value: "whoops",
		},
		{
			ID:     testhelper.MkID("int64 with checks - good"),
			s:      setterInt64WithChecks,
			value:  "9",
			expVal: "9",
		},
		{
			ID: testhelper.MkID("int64 with checks - bad"),
			ExpErr: testhelper.MkExpErr(
				`the value (5) must be greater than 6`),
			s:     setterInt64WithChecks,
			value: "5",
		},
		{
			ID:     testhelper.MkID("map - with val - good"),
			s:      setterMap,
			value:  "hello",
			expVal: "hello=true",
		},
		{
			ID:     testhelper.MkID("map - with val - good"),
			s:      setterMap,
			value:  "hello=false",
			expVal: "hello=false",
		},
		{
			ID:     testhelper.MkID("map (initially empty) - with val - good"),
			s:      setterMapEmpty,
			value:  "hello=false",
			expVal: "hello=false",
		},
		{
			ID: testhelper.MkID("map - with val - bad"),
			ExpErr: testhelper.MkExpErr(
				`bad value: "hello=whoops": part: 1 ("hello=whoops")` +
					` is invalid. The value ("whoops") cannot be` +
					` interpreted as true or false:`),
			s:     setterMap,
			value: "hello=whoops",
		},
		{
			ID: testhelper.MkID("map - with val - bad"),
			ExpErr: testhelper.MkExpErr(
				`bad value: "hello=whoops": part: 1 ("hello=whoops")` +
					` is invalid. The value ("whoops") cannot be` +
					` interpreted as true or false:`),
			s:     setterMap,
			value: "hello=whoops",
		},
		{
			ID: testhelper.MkID("nil - bad"),
			ExpErr: testhelper.MkExpErr(
				`a value must not follow this parameter: "`+paramName+`"`,
				"Remove the '=' and any following text"),
			s:     setterNil,
			value: "anything",
		},
		{
			ID:     testhelper.MkID("pathname - good"),
			s:      setterPathname,
			value:  "a/b/c",
			expVal: "a/b/c",
		},
		{
			ID:     testhelper.MkID("pathname with expectation - good"),
			s:      setterPathnameWithExpectation,
			value:  "testdata//pathname",
			expVal: "testdata/pathname",
		},
		{
			ID: testhelper.MkID("pathname with expectation - bad"),
			ExpErr: testhelper.MkExpErr(
				"should exist but does not",
				`"testdata/pathname" exists but "nonesuch" does not`),
			s:     setterPathnameWithExpectation,
			value: "testdata//pathname/nonesuch",
		},
		{
			ID: testhelper.MkID("pathname with checks - bad"),
			ExpErr: testhelper.MkExpErr(
				`"a/b/c" should have "testdata" as a prefix`),
			s:     setterPathnameWithChecks,
			value: "a/b/c",
		},
		{
			ID:     testhelper.MkID("pathname with checks - good"),
			s:      setterPathnameWithChecks,
			value:  "testdata/b/c",
			expVal: "testdata/b/c",
		},
		{
			ID:     testhelper.MkID("regexp - good"),
			s:      setterRegexp,
			value:  "[a-z]*",
			expVal: "[a-z]*",
		},
		{
			ID: testhelper.MkID("regexp - bad"),
			ExpErr: testhelper.MkExpErr(
				"could not parse \"[a-z*\" into a regular expression"),
			s:     setterRegexp,
			value: "[a-z*",
		},
		{
			ID: testhelper.MkID("String - bad"),
			ExpErr: testhelper.MkExpErr(
				`"goodbye,cruel,world" should have "hello" as a prefix`),
			s:     setterStringWithChecks,
			value: "goodbye,cruel,world",
		},
		{
			ID:     testhelper.MkID("String - good"),
			s:      setterStringWithChecks,
			value:  "hello,world",
			expVal: "hello,world",
		},
		{
			ID: testhelper.MkID("StringList - bad"),
			ExpErr: testhelper.MkExpErr("the length of the list (3)" +
				" is incorrect: the value (3) must equal 2"),
			s:     setterStrListWithChecks,
			value: "hello,cruel,world",
		},
		{
			ID:     testhelper.MkID("StringList - good"),
			s:      setterStrListWithChecks,
			value:  "hello,world",
			expVal: "hello,world",
		},
		{
			ID:     testhelper.MkID("TimeLocation - good"),
			s:      setterTimeLocationWithChecks,
			value:  "America/New York",
			expVal: "America/New_York",
		},
		{
			ID: testhelper.MkID("TimeLocation - bad"),
			ExpErr: testhelper.MkExpErr(
				`bad timezone "nonesuch"`),
			s:     setterTimeLocationWithChecks,
			value: "nonesuch",
		},
		{
			ID:     testhelper.MkID("Calculated[int] with checks - good"),
			s:      calcSetterWithChecks,
			value:  "std",
			expVal: "42",
		},
		{
			ID:     testhelper.MkID("Calculated[int] with checks - bad"),
			ExpErr: testhelper.MkExpErr("the value (24) must equal 42"),
			s:      calcSetterWithChecks,
			value:  "std2",
		},
		{
			ID:     testhelper.MkID("Calculated[int] with TC - good"),
			s:      calcSetterWithBadTC,
			value:  "std",
			expVal: "42",
		},
		{
			ID:     testhelper.MkID("Calculated[int] with TC - bad"),
			ExpErr: testhelper.MkExpErr("bad value"),
			s:      calcSetterWithBadTC,
			value:  "std2",
		},
	}

	for _, tc := range testCases {
		err := tc.s.SetWithVal(paramName, tc.value)
		if testhelper.CheckExpErr(t, err, tc) && err == nil {
			newVal := tc.s.CurrentValue()
			if newVal != tc.expVal {
				t.Log(tc.IDStr())
				t.Log("\t: expected: ", tc.expVal)
				t.Log("\t:      got: ", newVal)
				t.Errorf("\t: unexpected value\n")
			}
		}
	}
}
