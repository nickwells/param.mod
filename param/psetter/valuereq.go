package psetter

import (
	"fmt"

	"github.com/nickwells/param.mod/v5/param"
)

// ValueReqMandatory is a mixin type that can be embedded in a Setter to
// provide suitable default values for a Setter where the parameter must have
// a value following it.
type ValueReqMandatory struct{}

// ValueReq returns the Mandatory value of the ValueReq type to indicate that
// a value must follow the parameter for this setter.
func (v ValueReqMandatory) ValueReq() param.ValueReq { return param.Mandatory }

// Set returns an error because if the value is Mandatory then a value must
// follow the parameter for this setter.
func (v ValueReqMandatory) Set(name string) error {
	return fmt.Errorf("a value must follow this parameter: %q,"+
		" either following an '=' or as the next parameter", name)
}

// ValueReqOptional is a mixin type that can be embedded in a Setter to
// provide suitable default values for a Setter where the following parameter
// is optional.
type ValueReqOptional struct{}

// ValueReq returns the Optional value of the ValueReq type to indicate that
// a value may follow the parameter for this setter but need not.
func (v ValueReqOptional) ValueReq() param.ValueReq { return param.Optional }

// ValueReqNone is a mixin type that can be embedded in a Setter to provide
// suitable default values for a Setter where the parameter must not have a
// following value.
type ValueReqNone struct{}

// ValueReq returns the None value of the ValueReq type to indicate that
// a value must not follow the parameter for this setter.
func (v ValueReqNone) ValueReq() param.ValueReq { return param.None }

// SetWithVal returns an error because if the value is None then a value must
// not follow the parameter for this setter.
func (v ValueReqNone) SetWithVal(name, _ string) error {
	return fmt.Errorf("a value must not follow this parameter: %q"+
		"Remove the '=' and any following text", name)
}
