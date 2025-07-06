package psetter_test

import (
	"testing"

	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
)

func TestEnum(t *testing.T) {
	value := ""
	es := psetter.Enum[string]{
		Value: &value,
		AllowedVals: psetter.AllowedVals[string]{
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

	testCases := []struct {
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
