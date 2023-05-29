package paramtest

import (
	"errors"
	"testing"

	"github.com/nickwells/errutil.mod/errutil"
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

// Parser holds the values needed to test that the parameters given in a PSet
// correctly fill in the expected values in a interface when a given set of
// arguments is passed.
//
// When creating multiple Parsers you should give each instance a new PSet as
// the PSet remembers whether Parse has previously been called and will panic
// if called twice.
//
// If the CmpFunc returns a non-nil error the test is taken to have
// failed. The CmpFunc will be passed the Val and ExpVal values and should
// check that the contents of Val after parsing the Args matches the
// ExpVal. You might find the testhelper.DiffVals func to be of help with
// this.
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
//
// If you need to perform some setup before running the test or cleanup
// afterwards you should set the Pre and Post runctions accordingly. If
// either is nil then it is not run. Any errors returned by these functions
// conmstitute a test failure and will be reported as such.
type Parser struct {
	testhelper.ID

	ExpParseErrors errutil.ErrMap
	Ps             *param.PSet

	Val       any
	ExpVal    any
	CheckFunc func(val, expVal any) error

	Pre  func() error
	Post func() error

	Args []string
}

// Test perfoms the test - it will call Parse on the PSet, passing the
// Args. It will check that the mapped errors match the ErrMap returned by
// Parse and then will call CmpFunc reporting any error that returns as a
// test error. It will return a non-nil error if the test failed.
func (p Parser) Test(t *testing.T) (err error) {
	t.Helper()

	if p.Pre != nil {
		err = p.Pre()
		if err != nil {
			t.Log(p.IDStr())
			t.Logf("\t: Unexpected test-setup error: %s\n", err)
			t.Error("\t: error returned by the test 'Pre' func")
			return err
		}
	}
	if p.Post != nil {
		defer func() {
			postErr := p.Post()
			if postErr != nil {
				t.Log(p.IDStr())
				t.Logf("\t: Unexpected test-cleanup error: %s\n", postErr)
				t.Error("\t: error returned by the test 'Post' func")
				err = errors.Join(err, postErr)
			}
		}()
	}

	errMap := errutil.ErrMap(p.Ps.Parse(p.Args))

	if err = errMap.Matches(p.ExpParseErrors); err != nil {
		t.Log(p.IDStr())
		t.Logf("\t: Unexpected parsing errors: %s\n", err)
		t.Logf("\t: Actual:\n%s\n", errMap)
		t.Logf("\t: Expected:\n%s\n", p.ExpParseErrors)
		t.Error("\t: Parsing failed in an unexpected way\n")
		return err
	}

	if err = p.CheckFunc(p.Val, p.ExpVal); err != nil {
		t.Log(p.IDStr())
		t.Logf("\t: Unexpected error: %s\n", err)
		t.Error("\t: The resultant value is not as expected\n")
		return err
	}
	return nil
}
