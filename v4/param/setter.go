package param

// =============================================

// Setter is the interface that wraps the necessary methods for a
// paramSetter
//
// Each paramSetter must implement these methods
//
// Set takes the name of the parameter and no value. If the value
// requirements are such that an argument is not needed  then this
// can set the value or do whatever the paramSetter is supposed to do.
// Otherwise it should return an error.
//
// SetWithVal takes the name of the parameter and the associated value.
// If the value requirements are such that an argument is needed or is
// optional then this can set the value or do whatever the paramSetter
// is supposed to do. Otherwise it should return an error.
//
// ValueReq returns the ValueReq for the parameter: one of
// Mandatory, Optional or None.
//
// AllowedValues returns a string documenting the parameters that a parameter
// value can take. A function generating a help message should call both this
// and the AllowedValuesMap method and ignore any nil map returned.
//
// AllowedValuesMap returns a map of allowed values to descriptions. This
// will return nil if the parameter does not have an associated map of
// values. A function generating a help message should call both this and the
// AllowedValues method and ignore any nil map returned.
//
// CurrentValue returns a string showing the current value of the
// parameter. It is called before any arguments have been parsed in order to
// get the initial value for use in help messages.
//
// CheckSetter is called to ensure that the setter has been correctly
// created, for instance that the pointer to the value is not nil or, if the
// pointer is to a map, that the map being pointed at has been
// created. Correct behaviour of this func would be to panic if the setter
// has not been properly set up.
//
// When creating your own Setter implementation you may find it useful to use
// one of the ValueReq types as an embedded type.So, for instance, if your
// Setter must have a following value then you can embed the
// param.ValueReqMandatory struct in your struct and this will provide
// correct ValueReq and Set methods. Similarly if your Setter must not have a
// following value embed the param.ValueReqNone struct and it will provide
// correct ValueReq and SetWithVal methods. For examples of how this is done
// see any of the Setter instances in the psetter package.
//
// Likewise the NilAVM type can be embedded in a Setter implementation and it
// will provide a default AllowedValuesMap method which will return a nil
// map.
type Setter interface {
	Set(string) error
	SetWithVal(string, string) error
	ValueReq() ValueReq
	AllowedValues() string
	CurrentValue() string
	CheckSetter(name string)
}

// AllowedValuesMapper
type AllowedValuesMapper interface {
	AllowedValuesMap() AllowedVals
}
