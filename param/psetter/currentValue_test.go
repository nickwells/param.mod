package psetter_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/psetter"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestCurrentValue(t *testing.T) {
	vBool := true
	setterBool := psetter.Bool{Value: &vBool}
	setterBoolInverted := psetter.Bool{Value: &vBool, Invert: true}

	vDuration := 1 * time.Millisecond
	setterDuration := psetter.Duration{Value: &vDuration}

	vEnumList := []string{"hello", "world"}
	setterEnumList := psetter.EnumList[string]{Value: &vEnumList}

	vEnumMap := map[string]bool{
		"hello": true,
	}
	setterEnumMap := psetter.EnumMap[string]{Value: &vEnumMap}

	vEnum := "Hello, World!"
	setterEnum := psetter.Enum[string]{Value: &vEnum}

	vFloat64 := float64(3.14159)
	setterFloat64 := psetter.Float[float64]{Value: &vFloat64}

	vInt64List := []int64{1, 2}
	setterInt64List := psetter.IntList[int64]{Value: &vInt64List}

	vInt64 := int64(42)
	setterInt64 := psetter.Int[int64]{Value: &vInt64}

	vUint64 := uint64(42)
	setterUint64 := psetter.Uint[uint64]{Value: &vUint64}

	vMap := map[string]bool{
		"hello": true,
	}
	setterMap := psetter.Map{Value: &vMap}

	setterNil := psetter.Nil{}

	vPathname := "/a/b/c"
	setterPathname := psetter.Pathname{Value: &vPathname}

	vRegexp := regexp.MustCompile("[a-z]*")
	setterRegexp := psetter.Regexp{Value: &vRegexp}
	setterRegexpNil := psetter.Regexp{}

	vString := "Hello, World!"
	setterString := psetter.String[string]{Value: &vString}

	vStrList := []string{"hello", "world"}
	setterStrList := psetter.StrList{Value: &vStrList}

	vTimeLocation, _ := time.LoadLocation("UTC")
	setterTimeLocation := psetter.TimeLocation{Value: &vTimeLocation}

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
			ID:            testhelper.MkID("Uint64 - 42"),
			s:             setterUint64,
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
