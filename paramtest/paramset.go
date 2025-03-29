package paramtest

import (
	"testing"

	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/paramset"
)

// MakeParamSetOrFatal constructs and returns a param set with a minimal
// Helper, no error reporting and no exit on failure. If the PSet cannot be
// constructed it reports a Fatal testing error
func MakeParamSetOrFatal(t *testing.T, id string) *param.PSet {
	t.Helper()

	ps, err := paramset.NewNoHelpNoExitNoErrRpt()
	if err != nil {
		t.Fatal(id, " : couldn't construct the PSet: ", err)
	}

	return ps
}
