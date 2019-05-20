package param

import "errors"

// ValueReq encodes whether or not a value is required after a
// parameter
type ValueReq int

// Mandatory means that a value must follow the parameter.
//
// Optional means that a value may follow the parameter but need not in which
// case the default value will be used.
//
// None means that a value must not follow the parameter.
const (
	Mandatory ValueReq = iota
	Optional
	None
)

// ValueReqMandatory is a mixin type that can be embedded in a Setter to
// provide suitable default values for a Setter where the parameter must have
// a value following it.
type ValueReqMandatory struct{}

// ValueReq returns the Mandatory value of the ValueReq type to indicate that
// a value must follow the parameter for this setter
func (v ValueReqMandatory) ValueReq() ValueReq { return Mandatory }

// Set returns an error because if the value is Mandatory then a value must
// follow the parameter for this setter
func (v ValueReqMandatory) Set(_ string) error {
	return errors.New("a value must follow this parameter," +
		" either following an '=' or as a next parameter")
}

// ValueReqOptional is a mixin type that can be embedded in a Setter to
// provide suitable default values for a Setter where the following parameter
// is optional.
type ValueReqOptional struct{}

// ValueReq returns the Optional value of the ValueReq type to indicate that
// a value must follow the parameter for this setter
func (v ValueReqOptional) ValueReq() ValueReq { return Optional }

// ValueReqNone is a mixin type that can be embedded in a Setter to provide
// suitable default values for a Setter where the parameter must not have a
// following value.
type ValueReqNone struct{}

// ValueReq returns the None value of the ValueReq type to indicate that
// a value must not follow the parameter for this setter
func (v ValueReqNone) ValueReq() ValueReq { return None }

// SetWithVal returns an error because if the value is None then a value must
// not follow the parameter for this setter
func (v ValueReqNone) SetWithVal(_, _ string) error {
	return errors.New("a value must not follow this parameter. " +
		"Remove the '=' and any following text")
}
