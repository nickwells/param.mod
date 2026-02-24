package psetter

import "github.com/nickwells/param.mod/v7/ptypes"

// AllowedVals is a type alias for [ptypes.AllowedVals]
type AllowedVals[T ~string] = ptypes.AllowedVals[T]

// Aliases is a type alias for [ptypes.Aliases]
type Aliases[T ~string] = ptypes.Aliases[T]

// ValDescriber is a type alias for [ptypes.ValDescriber]
type ValDescriber = ptypes.ValDescriber
