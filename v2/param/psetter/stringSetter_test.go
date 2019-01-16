package psetter_test

import (
	"github.com/nickwells/param.mod/v2/param"
	"github.com/nickwells/param.mod/v2/param/psetter"
	"testing"
)

func TestStringSetter(t *testing.T) {
	var value string = "initialValue"
	ss := psetter.StringSetter{
		Value: &value,
	}

	if pvr := ss.ValueReq(); pvr != param.Mandatory {
		t.Error(
			"StringSetter should need a value. ValueReq() returned ",
			pvr.String())
	}

	if err := ss.Set(""); err == nil {
		t.Error("StringSetter should have returned an error" +
			" when Set(...) was called")
	}
	if value != "initialValue" {
		t.Error("StringSetter should not have changed the value" +
			" when Set(...) was called. New value: '" + value + "'")
	}

	if err := ss.SetWithVal("", "test"); err != nil {
		t.Error("StringSetter should not have returned an error"+
			" when SetWithVal(...) was called, err: ", err)
	}
	if value != "test" {
		t.Error("StringSetter should have set the value to 'test'" +
			" when SetWithVal(...) was called." +
			" Actual value: '" + value + "'")
	}
}

func TestStringListSetter(t *testing.T) {
	var value []string = make([]string, 0)
	ss := psetter.StrListSetter{Value: &value}

	if pvr := ss.ValueReq(); pvr != param.Mandatory {
		t.Error("StrListSetter should need a value."+
			" ValueReq() returned ", pvr.String())
	}

	if err := ss.Set(""); err == nil {
		t.Error("StrListSetter should have returned an error" +
			" when Set(...) was called")
	}
	if len(value) != 0 {
		t.Errorf("StrListSetter should not have changed the value"+
			" when Set(...) was called. New value: %v'", value)
	}

	testCases := [...]struct {
		val             string
		shouldReturnErr bool
		expectedVal     []string
	}{
		{"test", false, []string{"test"}},
		{"test,test2", false, []string{"test", "test2"}},
		{"", false, []string{""}},
	}

	for i, tc := range testCases {
		value = value[:0]
		err := ss.SetWithVal("", tc.val)
		if tc.shouldReturnErr && err == nil {
			t.Errorf("case %d: StrListSetter should have returned an error"+
				" when SetWithVal(..., '%s') was called but didn't",
				i, tc.val)

		} else if !tc.shouldReturnErr {
			if err != nil {
				t.Errorf("case %d: StrListSetter should not return an error"+
					" when SetWithVal(..., '%s') was called but did. Err: %s",
					i, tc.val, err)
			} else {
				if len(value) != len(tc.expectedVal) {
					t.Errorf("case %d: StrListSetter should have set"+
						" %d entries in the value list"+
						" when SetWithVal(..., '%s') was called"+
						" but %d values were set",
						i, len(tc.expectedVal), tc.val, len(value))
				}
			}
		}
	}
}
