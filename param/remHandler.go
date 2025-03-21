package param

import (
	"math"

	"github.com/nickwells/location.mod/location"
)

// RemHandler describes how the remaining parameters should be handled. These
// are the parameters that remain after either a terminal positional
// parameter or else a terminal parameter (by default this is "--")
//
// The HandleRemainder func is for specifying how to handle any remaining
// arguments after the standard parsing has been completed. It takes as an
// argument the PSet on which Parse has been called. The remaining
// parameters can be retrieved for further processing through the Remainder
// method on that PSet. The second argument to HandleRemainder is the
// location that the previous parser reached in the parameters it had been
// passed.
//
// A typical way of using this mechanism might be to select a new pre-created
// PSet (or create a new one on the fly) based on the settings of any
// previous parameters by the previous Parse call. Then Parse can be called
// on the new PSet passing the remaining arguments given by the Remainder
// method on the original PSet. This would be a typical use case where a
// terminal positional parameter has been used; several common commands such
// as git and go itself use a similar style of command invocation.
//
// In the case where a terminal parameter is given and there is a following
// list of parameters it is more likely that the extra parameters are not
// intended as flags to control the operation of the program. In this case
// the remaining parameters might be taken as a list of values to be
// processed and a different HandleRemainder function would be appropriate.
type RemHandler interface {
	HandleRemainder(ps *PSet, loc *location.L)
}

type dfltRemHandler struct{}

// HandleRemainder calls the helper ErrorHandler to report that unexpected
// additional arguments were passed.
func (rh dfltRemHandler) HandleRemainder(ps *PSet, loc *location.L) {
	remCount := len(ps.Remainder())

	if remCount == 0 {
		return
	}

	args := ""
	sep := "'"
	end := "'"

	const maxLen = 20

	for i, r := range ps.Remainder() {
		charsToTake := int(math.Min(
			float64(len(r)),
			float64(maxLen-len(args)-len(sep))))
		if charsToTake <= 0 {
			args += "' ..."
			end = ""

			break
		}

		args += sep + r[:charsToTake]
		sep = "' '"

		if charsToTake < len(r) {
			args += "...'"
			if i < remCount-1 {
				args += " ..."
			}

			end = ""

			break
		}
	}

	args += end

	var err error
	if remCount == 1 {
		err = loc.Error("there was an unexpected extra parameter: " + args)
	} else {
		err = loc.Errorf("there were %d unexpected extra parameters: %s",
			remCount, args)
	}

	ps.helper.ErrorHandler(ps,
		ErrMap{
			"": []error{err},
		})
}

// NullRemHandler is a type which can be set as a remainder handler if you
// wish to perform the remainder handling yourself after the parsing is
// complete.
type NullRemHandler struct{}

// HandleRemainder does nothing, specifically it doesn't call the helper's
// ErrorHandler (which by default will terminate the program). If you set
// this as the RemHandler for the PSet then you will have to handle the
// remaining arguments after the program has been called by calling the
// Remainder method on the PSet.
func (rh NullRemHandler) HandleRemainder(_ *PSet, _ *location.L) {}
