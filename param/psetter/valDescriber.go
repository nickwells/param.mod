package psetter

// ValDescriber is an interface for a function which can be added as a method
// on an existing Setter and if it is present its return value will be used
// to describe the parameter value to be supplied. It only makes sense to
// provide such a method if the parameter can take a value.
type ValDescriber interface {
	ValDescribe() string
}
