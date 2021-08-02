package paramtest

import (
	"testing"

	"github.com/nickwells/errutil.mod/errutil"
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/testhelper.mod/testhelper"
)

// Parser holds the values needed to test that the parameters given in a PSet
// correctly fill in the expected values in a interface when a given set of
// arguments is passed.
//
// When creating multiple Parsers you should give each instance a new PSet as
// the PSet remembers whether Parse has previously been called and will panic
// if called twice.
//
// If the CmpFunc returns
// a non-nil error the test is taken to have failed. The CmpFunc will be
// passed the Val and ExpVal values and should check that the contents of Val
// after parsing the Args matches the ExpVal. Strictly all you need do is check that the .
//
// You should construct the PSet, adding params that set values in Val so
// that when the supplied Args are Parsed they will update that value ready
// for it to be compared with ExpVal.
//
// When constructing the PSet you should use the NewNoHelpNoExitNoErrRpt()
// func from the paramset package so as not to have any interference in the
// test from the standard help parameters. Note that this will mean that any
// errors due to parameters whose names clash with standard parameters will
// not be caught but this test can be performed separately.
type Parser struct {
	testhelper.ID

	ExpParseErrors errutil.ErrMap
	Ps             *param.PSet

	Val       interface{}
	ExpVal    interface{}
	CheckFunc func(val, expVal interface{}) error

	Args []string
}

// Test perfoms the test - it will call Parse on the PSet, passing the
// Args. It will check that the mapped errors match the ErrMap returned by
// Parse and then will call CmpFunc reporting any error that returns as a
// test error.
func (p Parser) Test(t *testing.T) {
	t.Helper()

	errMap := errutil.ErrMap(p.Ps.Parse(p.Args))

	if err := errMap.Matches(p.ExpParseErrors); err != nil {
		t.Log(p.IDStr())
		t.Logf("\t: Unexpected parsing errors: %s\n", err)
		t.Error("\t: Parsing failed in an unexpected way\n")
		return
	}

	if err := p.CheckFunc(p.Val, p.ExpVal); err != nil {
		t.Log(p.IDStr())
		t.Logf("\t: Unexpected error: %s\n", err)
		t.Error("\t: The resultant value is not as expected\n")
	}
}
