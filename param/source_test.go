package param_test

import (
	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/param"
	"github.com/nickwells/param.mod/param/paramset"
	"github.com/nickwells/param.mod/param/psetter"
	"testing"
)

func TestSource(t *testing.T) {
	var b bool
	ps, _ := paramset.NewNoHelpNoExitNoErrRpt()
	p := ps.Add("p", psetter.BoolSetter{Value: &b}, "desc")
	loc := location.New("loc")
	loc.Incr()
	testCases := []struct {
		name    string
		from    string
		loc     *location.L
		pVals   []string
		p       *param.ByName
		expStr  string
		expDesc string
	}{
		{
			name:    "ok",
			from:    "source",
			loc:     loc,
			pVals:   []string{"p"},
			p:       p,
			expStr:  "Param: p (at loc:1)",
			expDesc: "source (at loc:1) [p]",
		},
		{
			name:    "ok",
			from:    "source",
			loc:     loc,
			pVals:   []string{"p", "true"},
			p:       p,
			expStr:  "Param: p (at loc:1)",
			expDesc: "source (at loc:1) [p=true]",
		},
	}

	for i, tc := range testCases {
		s := param.Source{
			From:      tc.from,
			Loc:       *tc.loc,
			ParamVals: tc.pVals,
			Param:     tc.p,
		}
		if s.String() != tc.expStr {
			t.Logf("test %d: %s :\n", i, tc.name)
			t.Logf("\t:   string: %s\n", s.String())
			t.Logf("\t: expected: %s\n", tc.expStr)
			t.Errorf("\t: bad string\n")
		}
		if s.Desc() != tc.expDesc {
			t.Logf("test %d: %s :\n", i, tc.name)
			t.Logf("\t:     desc: %s\n", s.Desc())
			t.Logf("\t: expected: %s\n", tc.expDesc)
			t.Errorf("\t: bad desc\n")
		}
	}
}

func TestSources(t *testing.T) {
	var b bool
	ps, _ := paramset.NewNoHelpNoExitNoErrRpt()
	p := ps.Add("p", psetter.BoolSetter{Value: &b}, "desc")
	loc := location.New("loc")
	loc.Incr()
	expStr := "Param: p (at loc:1), Param: p (at loc:1)"

	s := param.Source{
		From:      "source",
		Loc:       *loc,
		ParamVals: []string{"p", "true"},
		Param:     p,
	}
	sources := param.Sources{s, s}
	if sources.String() != expStr {
		t.Logf("\t:   string: %s\n", sources.String())
		t.Logf("\t: expected: %s\n", expStr)
		t.Errorf("\t: bad string\n")
	}
}
