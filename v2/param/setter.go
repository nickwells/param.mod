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

// =============================================

// Setter is the interface that wraps the necessary methods for a
// paramSetter
//
// Each paramSetter must implement these methods
//
// Set takes the name of the parameter and no value. If the value
// requirements are such that an argument is not needed  then this
// can set the value or do whatever the paramSetter is supposed to do.
// Otherwise it should return an error
//
// SetWithVal takes the name of the parameter and the associated value.
// If the value requirements are such that an argument is needed or is
// optional then this can set the value or do whatever the paramSetter
// is supposed to do. Otherwise it should return an error
//
// ValueReq returns the ValueReq for the parameter: one of
// Mandatory, Optional or None.
//
// AllowedValues returns a string documenting the parameters that a parameter
// value can take
//
// CurrentValue returns a string showing the current value of the
// parameter. It is called before any arguments have been parsed in order to
// get the initial value for use in help messages
//
// CheckSetter is called to ensure that the setter has been correctly
// created, for instance that the pointer to the value is not nil or, if the
// pointer is to a map, that the map being pointed at has been
// created. Correct behaviour of this func would be to panic if the setter
// has not been properly set up.
type Setter interface {
	Set(string) error
	SetWithVal(string, string) error
	ValueReq() ValueReq
	AllowedValues() string
	CurrentValue() string
	CheckSetter(name string)
}

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
