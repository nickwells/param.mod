package psetter_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/param.mod/v4/param"
	"github.com/nickwells/param.mod/v4/param/psetter"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestSetWithVal(t *testing.T) {
	paramName := "setval-test"
	var vBool bool = true
	var setterBool = psetter.Bool{Value: &vBool}

	var vDuration time.Duration = 1 * time.Millisecond
	var setterDuration = psetter.Duration{Value: &vDuration}

	var vEnumList []string = []string{"hello", "world"}
	var setterEnumList = psetter.EnumList{
		Value: &vEnumList,
		AllowedVals: psetter.AllowedVals{
			"hello": "hello description",
			"world": "world description",
			"bye":   "bye description",
		},
	}
	var setterEnumListWithChecks = psetter.EnumList{
		Value: &vEnumList,
		AllowedVals: psetter.AllowedVals{
			"hello": "hello description",
			"world": "world description",
			"bye":   "bye description",
		},
		Checks: []check.StringSlice{
			check.StringSliceLenEQ(2),
			nil,
		},
	}

	var vEnumMap = map[string]bool{
		"hello": true,
	}
	var setterEnumMap = psetter.EnumMap{
		Value: &vEnumMap,
		AllowedVals: psetter.AllowedVals{
			"hello": "hello description",
			"world": "world description",
			"bye":   "bye description",
		},
	}

	var vEnum string = "hello"
	var setterEnum = psetter.Enum{
		Value: &vEnum,
		AllowedVals: psetter.AllowedVals{
			"hello": "hello description",
			"world": "world description",
			"bye":   "bye description",
		},
	}

	var vFloat64 float64 = 3.14159
	var setterFloat64 = psetter.Float64{Value: &vFloat64}
	var setterFloat64WithChecks = psetter.Float64{
		Value: &vFloat64,
		Checks: []check.Float64{
			nil,
			check.Float64GT(5),
		},
	}

	var vInt64List []int64 = []int64{1, 2}
	var setterInt64List = psetter.Int64List{Value: &vInt64List}
	var setterInt64ListWithChecks = psetter.Int64List{
		Value: &vInt64List,
		Checks: []check.Int64Slice{
			check.Int64SliceNoDups,
		},
	}

	var vInt64 int64 = 42
	var setterInt64 = psetter.Int64{Value: &vInt64}
	var setterInt64WithChecks = psetter.Int64{
		Value: &vInt64,
		Checks: []check.Int64{
			check.Int64GT(6),
		},
	}

	var vMap = map[string]bool{
		"hello": true,
	}
	var setterMap = psetter.Map{Value: &vMap}

	var vMapEmpty = map[string]bool{}
	var setterMapEmpty = psetter.Map{Value: &vMapEmpty}

	var setterNil = psetter.Nil{}

	var vPathname string
	var setterPathname = psetter.Pathname{Value: &vPathname}
	var vPathname2 string
	var setterPathnameWithExpectation = psetter.Pathname{
		Value:       &vPathname2,
		Expectation: filecheck.Provisos{Existence: filecheck.MustExist},
	}
	var vPathname3 string
	var setterPathnameWithChecks = psetter.Pathname{
		Value: &vPathname3,
		Checks: []check.String{
			nil,
			check.StringHasPrefix("testdata"),
		},
	}

	var vRegexp *regexp.Regexp
	var setterRegexp = psetter.Regexp{Value: &vRegexp}

	var vString string
	var setterStringWithChecks = psetter.String{
		Value: &vString,
		Checks: []check.String{
			nil,
			check.StringHasPrefix("hello"),
		},
	}

	var vStrList []string
	var setterStrListWithChecks = psetter.StrList{
		Value: &vStrList,
		Checks: []check.StringSlice{
			nil,
			check.StringSliceLenEQ(2),
		},
	}

	var vTimeLocation *time.Location
	var setterTimeLocationWithChecks = psetter.TimeLocation{
		Value:  &vTimeLocation,
		Checks: []check.TimeLocation{nil},
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
				"could not parse 'whoops - a bad duration' as a duration: "),
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
			ExpErr: testhelper.MkExpErr(
				"the length of the list (1) must equal 2"),
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
				`list entries: 0 and 1 are duplicates, both are: 5`),
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
				"should exist but doesn't"),
			s:     setterPathnameWithExpectation,
			value: "testdata//pathname/nonesuch",
		},
		{
			ID: testhelper.MkID("pathname with checks - bad"),
			ExpErr: testhelper.MkExpErr(
				"'a/b/c' should have 'testdata' as a prefix"),
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
				"'goodbye,cruel,world' should have 'hello' as a prefix"),
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
			ExpErr: testhelper.MkExpErr(
				"the length of the list (3) must equal 2"),
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
				"could not find \"nonesuch\" as a time location"),
			s:     setterTimeLocationWithChecks,
			value: "nonesuch",
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
