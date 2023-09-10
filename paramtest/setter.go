package paramtest

import (
	"testing"

	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
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

// SetVR sets the VRExp value to the supplied vr and generates appropriate
// SetErr or SetWithValErr entries accordingly. The errors are generated
// assuming that you have used the psetter.ValueReq... types when building
// your Setter; if you have not used these types you may need to generate the
// errors yourself.
//
// Note that this uses the Setter's ParamName field so that must have its
// final value before this method is called.
func (s *Setter) SetVR(vr param.ValueReq) {
	s.VRExp = vr

	switch vr {
	case param.Mandatory:
		var vrm psetter.ValueReqMandatory
		err := vrm.Set(s.ParamName)
		s.SetErr = testhelper.MkExpErr(err.Error())
	case param.None:
		var vrm psetter.ValueReqNone
		err := vrm.SetWithVal(s.ParamName, "")
		s.SetWithValErr = testhelper.MkExpErr(err.Error())
	}
}

// Test performs all the tests
func (s Setter) Test(t *testing.T) {
	t.Helper()

	panicked, panicVal := testhelper.PanicSafe(
		func() {
			s.PSetter.CheckSetter(s.ParamName)
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

	var (
		setStr        = `Set("` + s.ParamName + `")`
		setWithValStr = `SetWithVal("` + s.ParamName + `", "` + s.ParamVal + `")`
	)
	err := s.PSetter.Set(s.ParamName)
	testhelper.CheckExpErrWithID(t,
		s.IDStr()+" - "+setStr, err, s.SetErr)
	cv = s.PSetter.CurrentValue()
	s.GFC.Check(t,
		s.IDStr()+` - Value after `+setStr,
		s.Name+".val-postSet", []byte(cv))

	err = s.PSetter.SetWithVal(s.ParamName, s.ParamVal)
	testhelper.CheckExpErrWithID(t,
		s.IDStr()+" - "+setWithValStr, err, s.SetWithValErr)
	cv = s.PSetter.CurrentValue()
	s.GFC.Check(t,
		s.IDStr()+` - Value after `+setWithValStr,
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
