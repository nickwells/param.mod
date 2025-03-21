package phelp

import (
	"fmt"

	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/twrap.mod/twrap"
)

const (
	paramIndent       = 6
	paramLine2Indent  = 9
	descriptionIndent = 12
	textIndent        = 4
)

const (
	majorSectionSeparator = "\n===============\n\n"
	minorSectionSeparator = "---------------\n"
)

// printHelpMessages prints the messages
func (h StdHelp) printHelpMessages(twc *twrap.TWConf, messages ...string) {
	for _, message := range messages {
		twc.Wrap(message, 0)
	}
}

// Help prints any messages and then a standardised usage message based on
// the parameters supplied to the param set. If it is called directly (that
// is if the help style is set to noHelp) then the output will be written to
// the param.PSet's error writer (by default stderr) rather than to its
// standard writer (stdout) and os.Exit will be called with an exit status of
// 1 to indicate an error.
func (h StdHelp) Help(ps *param.PSet, messages ...string) {
	w := ps.StdW()

	if h.sectionsChosen.hasNothingChosen() {
		if err := h.setHelpSections(standardHelpSectionAlias); err != nil {
			panic(fmt.Sprint("Couldn't set the default help sections:", err))
		}

		w = ps.ErrW()
	}

	twc := twrap.NewTWConfOrPanic(
		twrap.SetWriter(w),
		twrap.SetTargetLineLen(h.helpLineLen))

	printSep := false
	if len(messages) > 0 {
		printSep = true

		h.printHelpMessages(twc, messages...)
	}

	for _, sec := range helpSectionsInOrder {
		if h.sectionsChosen[sec.name] {
			printSep = printSepIf(twc, printSep, majorSectionSeparator)
			if !sec.displayFunc(h, twc, ps) {
				printSep = false
			}
		}
	}
}
