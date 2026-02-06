package phelp

import (
	"maps"
	"slices"

	"github.com/nickwells/param.mod/v6/param"
)

// showUnusedParams will print all the parameters that have been detected in
// configuration files or set through environment variables which were not
// recognised and were ignored. Since these sources may contain parameters
// intended for other programs it is not reported as an error but you may
// want to see what has been ignored in order to detect mistakes.
func showUnusedParams(h StdHelp, ps *param.PSet) bool {
	up := ps.UnusedParams()

	h.twc.Println("Unused Parameters: ", len(up))

	paramsByName := slices.Sorted(maps.Keys(up))

	for _, pn := range paramsByName {
		h.twc.Wrap(pn, paramIndent)

		for _, loc := range up[pn] {
			h.twc.WrapPrefixed("at: ", loc, descriptionIndent)
		}
	}

	return true
}
