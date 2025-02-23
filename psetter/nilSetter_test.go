package psetter_test

import (
	"testing"

	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
)

func TestNilSetter(t *testing.T) {
	ns := psetter.Nil{}
	if ns.ValueReq() != param.None {
		t.Errorf("psetter.Nil should require no arguments\n")
	}

	if ns.Set("dummy") != nil {
		t.Errorf("psetter.Nil should return no error when being set\n")
	}

	if ns.AllowedValues() != "none" {
		t.Errorf("psetter.Nil should return 'none' from AllowedValues()\n")
	}

	if ns.CurrentValue() != "none" {
		t.Errorf("psetter.Nil should return 'none' from CurrentValue()\n")
	}
}
