package param

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
