package ptypes

// OptFunc represents an option function - something taking a pointer to a
// value with the intention of transforming it in some way and returning an
// error if the transformation failed. These functions can be passed as
// variadic parameters to other functions (typically, but not necessarily, a
// function that creates an initial value - a constructor)
type OptFunc[T any] func(*T) error
