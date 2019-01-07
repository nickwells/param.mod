package psetter

import "testing"

func TestAValMap(t *testing.T) {
	testCases := []struct {
		name        string
		avMap       AValMap
		allowedVals string
	}{
		{
			name:  "empty",
			avMap: map[string]string{},
		},
		{
			name: "one entry",
			avMap: map[string]string{
				"name": "desc",
			},
			allowedVals: `name: desc`,
		},
		{
			name: "two entries",
			avMap: map[string]string{
				"name":      "desc",
				"long name": "long name desc",
			},
			allowedVals: `long name: long name desc
name     : desc`,
		},
	}

	for i, tc := range testCases {
		s := allowedValues(tc.avMap)
		if s != tc.allowedVals {
			t.Logf("test %d: %s :\n", i, tc.name)
			t.Logf("\t: allowed values: %s\n", s)
			t.Logf("\t:       expected: %s\n", tc.allowedVals)
			t.Errorf("\t: bad allowed values string\n")
		}
	}

}
