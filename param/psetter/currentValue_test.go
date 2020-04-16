package psetter_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/nickwells/param.mod/v4/param"
	"github.com/nickwells/param.mod/v4/param/psetter"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestCurrentValue(t *testing.T) {
	var vBool bool = true
	var setterBool = psetter.Bool{Value: &vBool}
	var setterBoolInverted = psetter.Bool{Value: &vBool, Invert: true}

	var vDuration time.Duration = 1 * time.Millisecond
	var setterDuration = psetter.Duration{Value: &vDuration}

	var vEnumList []string = []string{"hello", "world"}
	var setterEnumList = psetter.EnumList{Value: &vEnumList}

	var vEnumMap = map[string]bool{
		"hello": true,
	}
	var setterEnumMap = psetter.EnumMap{Value: &vEnumMap}

	var vEnum string = "Hello, World!"
	var setterEnum = psetter.Enum{Value: &vEnum}

	var vFloat64 float64 = 3.14159
	var setterFloat64 = psetter.Float64{Value: &vFloat64}

	var vInt64List []int64 = []int64{1, 2}
	var setterInt64List = psetter.Int64List{Value: &vInt64List}

	var vInt64 int64 = 42
	var setterInt64 = psetter.Int64{Value: &vInt64}

	var vMap = map[string]bool{
		"hello": true,
	}
	var setterMap = psetter.Map{Value: &vMap}

	var setterNil = psetter.Nil{}

	var vPathname string = "/a/b/c"
	var setterPathname = psetter.Pathname{Value: &vPathname}

	var vRegexp *regexp.Regexp = regexp.MustCompile("[a-z]*")
	var setterRegexp = psetter.Regexp{Value: &vRegexp}
	var setterRegexpNil = psetter.Regexp{}

	var vString string = "Hello, World!"
	var setterString = psetter.String{Value: &vString}

	var vStrList = []string{"hello", "world"}
	var setterStrList = psetter.StrList{Value: &vStrList}

	vTimeLocation, _ := time.LoadLocation("UTC")
	var setterTimeLocation = psetter.TimeLocation{Value: &vTimeLocation}

	testCases := []struct {
		testhelper.ID
		s             param.Setter
		expectedValue string
	}{
		{
			ID:            testhelper.MkID("bool - true"),
			s:             setterBool,
			expectedValue: "true",
		},
		{
			ID:            testhelper.MkID("bool - true"),
			s:             setterBoolInverted,
			expectedValue: "false",
		},
		{
			ID:            testhelper.MkID("duration - 1ms"),
			s:             setterDuration,
			expectedValue: "1ms",
		},
		{
			ID:            testhelper.MkID("enumlist - hello,world"),
			s:             setterEnumList,
			expectedValue: "hello,world",
		},
		{
			ID:            testhelper.MkID("enummap - val=true"),
			s:             setterEnumMap,
			expectedValue: "hello=true",
		},
		{
			ID:            testhelper.MkID("enum - Hello, World!"),
			s:             setterEnum,
			expectedValue: "Hello, World!",
		},
		{
			ID:            testhelper.MkID("float64 - 3.14159"),
			s:             setterFloat64,
			expectedValue: "3.14159",
		},
		{
			ID:            testhelper.MkID("Int64List - 1,2"),
			s:             setterInt64List,
			expectedValue: "1,2",
		},
		{
			ID:            testhelper.MkID("Int64 - 42"),
			s:             setterInt64,
			expectedValue: "42",
		},
		{
			ID:            testhelper.MkID("Map - hello=true"),
			s:             setterMap,
			expectedValue: "hello=true",
		},
		{
			ID:            testhelper.MkID("Nil"),
			s:             setterNil,
			expectedValue: "none",
		},
		{
			ID:            testhelper.MkID("Pathname - /a/b/c"),
			s:             setterPathname,
			expectedValue: "/a/b/c",
		},
		{
			ID:            testhelper.MkID("Pathname - /a/b/c"),
			s:             setterPathname,
			expectedValue: "/a/b/c",
		},
		{
			ID:            testhelper.MkID("Regexp - [a-z]*"),
			s:             setterRegexp,
			expectedValue: "[a-z]*",
		},
		{
			ID:            testhelper.MkID("RegexpNil - Illegal value"),
			s:             setterRegexpNil,
			expectedValue: "Illegal value",
		},
		{
			ID:            testhelper.MkID("String - Hello, World!"),
			s:             setterString,
			expectedValue: "Hello, World!",
		},
		{
			ID:            testhelper.MkID("StrList - hello,world"),
			s:             setterStrList,
			expectedValue: "hello,world",
		},
		{
			ID:            testhelper.MkID("TimeLocation - UTC"),
			s:             setterTimeLocation,
			expectedValue: "UTC",
		},
	}

	for _, tc := range testCases {
		cv := tc.s.CurrentValue()
		if cv != tc.expectedValue {
			t.Log(tc.IDStr())
			t.Log("\t: expected: ", tc.expectedValue)
			t.Log("\t:      got: ", cv)
			t.Errorf("\t: unexpected return from CurrentValue\n")
		}
	}
}
