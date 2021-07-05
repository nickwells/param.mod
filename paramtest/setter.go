package paramtest

import (
	"testing"

	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/psetter"
	"github.com/nickwells/testhelper.mod/testhelper"
)

// Setter holds the values needed to test a param.Setter and provides methods
// to perform the tests. It is used to test the collection of Setters
// provided in param.psetter and can also be used to test custom setters.
type Setter struct {
	testhelper.ID
	testhelper.ExpPanic // for CheckSetter

	GFC testhelper.GoldenFileCfg

	PSetter param.Setter

	VRExp param.ValueReq

	ParamName     string
	SetErr        testhelper.ExpErr
	ParamVal      string
	SetWithValErr testhelper.ExpErr

	ValDescriber bool

	ExtraTest func(*testing.T, Setter)
}

// Test performs all the tests
func (s Setter) Test(t *testing.T) {
	t.Helper()

	panicked, panicVal := testhelper.PanicSafe(
		func() {
			s.PSetter.CheckSetter(s.Name)
		})
	testhelper.CheckExpPanic(t, panicked, panicVal, s)
	if panicked || s.PanicExpected() {
		return
	}

	if actVR := s.PSetter.ValueReq(); actVR != s.VRExp {
		t.Log(s.IDStr())
		t.Logf("\t: expected ValueReq setting: %s\n", s.VRExp)
		t.Logf("\t:   actual ValueReq setting: %s\n", actVR)
		t.Error("\t: ValueReq is incorrect\n")
	}

	av := s.PSetter.AllowedValues()
	s.GFC.Check(t,
		s.IDStr()+" - Allowed Values",
		s.Name+".aval", []byte(av))

	cv := s.PSetter.CurrentValue()
	s.GFC.Check(t,
		s.IDStr()+" - Initial Value",
		s.Name+".init-val", []byte(cv))

	err := s.PSetter.Set(s.ParamName)
	testhelper.CheckExpErrWithID(t,
		s.IDStr()+" - Set", err, s.SetErr)
	cv = s.PSetter.CurrentValue()
	s.GFC.Check(t,
		s.IDStr()+" - Value after Set",
		s.Name+".val-postSet", []byte(cv))

	err = s.PSetter.SetWithVal(s.ParamName, s.ParamVal)
	testhelper.CheckExpErrWithID(t,
		s.IDStr()+" - SetWithVal", err, s.SetWithValErr)
	cv = s.PSetter.CurrentValue()
	s.GFC.Check(t,
		s.IDStr()+" - Value after SetWithVal",
		s.Name+".val-postSetWithVal", []byte(cv))

	if s.ValDescriber {
		if vd, ok := s.PSetter.(psetter.ValDescriber); !ok {
			t.Log(s.IDStr())
			t.Error("\t: should have a ValDescribe method but doesn't\n")
		} else {
			desc := vd.ValDescribe()
			s.GFC.Check(t,
				s.IDStr()+" - Value Description",
				s.Name+".val-description", []byte(desc))
		}
	}
	if s.ExtraTest != nil {
		s.ExtraTest(t, s)
	}
}
