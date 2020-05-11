package param

// =============================================

// Setter is the interface that wraps the necessary methods for a
// parameter.
//
// Set takes just the name of the parameter. If the value requirements are
// such that an argument is not needed then this can set the value or do
// whatever the Setter is supposed to do.  Otherwise it should return an
// error.
//
// SetWithVal takes the name of the parameter and the associated value.  If
// the value requirements are such that an argument is needed or is optional
// then this can set the value or do whatever the Setter is supposed to
// do. Otherwise it should return an error.
//
// ValueReq returns the ValueReq for the parameter: one of Mandatory,
// Optional or None. It is used when generating a help message.
//
// AllowedValues returns a string documenting the parameters that a parameter
// value can take. It is used when generating a help message.
//
// CurrentValue returns a string showing the current value of the
// parameter. It is called when the parameter is added, before any arguments
// have been parsed in order to get the initial value for use in help
// messages.
//
// CheckSetter is called when the parameter is added, to ensure that the
// Setter has been correctly created, for instance that the pointer to the
// value is not nil or, if the pointer is to a map, that the map being
// pointed at has been created. Correct behaviour of this func is to panic if
// the Setter has not been properly set up.
//
// When creating your own Setter implementation you may find it useful to use
// one of the ValueReq types as an embedded type. So, for instance, if your
// Setter must have a following value then you can embed the
// param.ValueReqMandatory struct in your struct. This will provide an
// appropriate ValueReq method and a Set method that will return an error;
// you only need to write the SetWithVal method. Similarly if your Setter
// must not have a following value embed the param.ValueReqNone struct and it
// will provide an appropriate ValueReq method and a SetWithVal method that
// will return an error; you only need to write the Set method. For examples
// of how this is done see the Setter instances in the psetter package.
type Setter interface {
	Set(string) error
	SetWithVal(string, string) error
	ValueReq() ValueReq
	AllowedValues() string
	CurrentValue() string
	CheckSetter(name string)
}
