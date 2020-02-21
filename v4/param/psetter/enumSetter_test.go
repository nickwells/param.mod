package psetter_test

import (
	"testing"

	"github.com/nickwells/param.mod/v3/param"
	"github.com/nickwells/param.mod/v3/param/psetter"
)

func TestEnum(t *testing.T) {
	var value string
	es := psetter.Enum{
		Value: &value,
		AllowedVals: param.AllowedVals{
			"e1": "E1 explained",
			"e2": "E2 explained",
			"e3": "E3 explained",
		},
	}

	if pvr := es.ValueReq(); pvr != param.Mandatory {
		t.Error(
			"Enum should need a value. ValueReq() returned ",
			pvr.String())
	}

	if err := es.Set(""); err == nil {
		t.Error("Enum should have returned an error" +
			" when Set(...) was called")
	}

	testCases := [...]struct {
		val         string
		expectedVal string
		errExpected bool
	}{
		{"e1", "e1", false},
		{"e2", "e2", false},
		{"e3", "e3", false},
		{"e4", "", true},
		{"", "", true},
	}

	for _, tc := range testCases {
		value = ""
		err := es.SetWithVal("", tc.val)
		if err == nil && tc.errExpected {
			t.Error("processing value: '" + tc.val + "'" +
				" should have raised an error but didn't")
		} else if err != nil && !tc.errExpected {
			t.Error("processing value: '"+tc.val+"'"+
				" should not have raised an error but did: ", err)
		}
		if value != tc.expectedVal {
			t.Error("processing value: '" + tc.val +
				"' should have set value to '" + tc.expectedVal +
				"' but instead set it to : '" + value + "'")
		}
	}
}

func TestEnumList(t *testing.T) {
	var value []string
	els := psetter.EnumList{
		Value: &value,
		AllowedVals: param.AllowedVals{
			"e1": "E1 explained",
			"e2": "E2 explained",
			"e3": "E3 explained",
		},
	}

	if pvr := els.ValueReq(); pvr != param.Mandatory {
		t.Error(
			"EnumList should need a value. ValueReq() returned ", pvr.String())
	}

	if err := els.Set(""); err == nil {
		t.Error("EnumList should have returned an error when Set(...) was called")
	}

	testCases := [...]struct {
		val         string
		expectedVal []string
		errExpected bool
	}{
		{"e1",
			[]string{"e1"},
			false},
		{"e2",
			[]string{"e2"},
			false},
		{"e3",
			[]string{"e3"},
			false},
		{"e3,e1",
			[]string{"e3", "e1"},
			false},
		{"e4",
			[]string{""},
			true},
		{"e2,e4",
			[]string{""},
			true},
		{"",
			[]string{""},
			true},
	}

	for _, tc := range testCases {
		value = value[:0]
		err := els.SetWithVal("", tc.val)
		if err == nil && tc.errExpected {
			t.Error("processing value: '" + tc.val + "'" +
				" should have raised an error but didn't")
		} else if err != nil && !tc.errExpected {
			t.Error("processing value: '"+tc.val+"'"+
				" should not have raised an error but did: ", err)
		}
		for i, v := range value {
			if v != tc.expectedVal[i] {
				t.Error("processing value: '"+tc.val+
					"' should have set value at", i, " to '"+tc.expectedVal[i]+
					"' but instead set it to : '"+v+"'")
			}
		}
	}
}

func TestEnumMap(t *testing.T) { // nolint: gocyclo
	var value map[string]bool
	ems := psetter.EnumMap{
		Value: &value,
		AllowedVals: param.AllowedVals{
			"e1": "E1 explained",
			"e2": "E2 explained",
			"e3": "E3 explained",
		},
	}

	if pvr := ems.ValueReq(); pvr != param.Mandatory {
		t.Error(
			"EnumMap should need a value. ValueReq() returned ",
			pvr.String())
	}

	if err := ems.Set(""); err == nil {
		t.Error("EnumMap should have returned an error" +
			" when Set(...) was called")
	}

	testCases := [...]struct {
		val         string
		expectedVal map[string]bool
		errExpected bool
	}{
		{"e1",
			map[string]bool{"e1": true},
			false},
		{"e2",
			map[string]bool{"e2": true},
			false},
		{"e3",
			map[string]bool{"e3": true},
			false},
		{"e3,e1",
			map[string]bool{"e3": true, "e1": true},
			false},
		{"e4",
			map[string]bool{},
			true},
		{"e2,e4",
			map[string]bool{},
			true},
		{"",
			map[string]bool{},
			true},
	}

	for _, tc := range testCases {
		value = map[string]bool{}
		err := ems.SetWithVal("", tc.val)
		if err == nil && tc.errExpected {
			t.Error("processing value: '" + tc.val + "'" +
				" should have raised an error but didn't")
		} else if err != nil && !tc.errExpected {
			t.Error("processing value: '"+tc.val+"'"+
				" should not have raised an error but did: ", err)
		}
		for k := range value {
			if !tc.expectedVal[k] {
				t.Error("processing value: '" + tc.val +
					"' the map entry for '" + k + "'" +
					" was set and should not have been")
			}
		}
		for k := range tc.expectedVal {
			if !value[k] {
				t.Error("processing value: '" + tc.val +
					"' the map entry for '" + k + "' was not set and should have been")
			}
		}
	}
}
