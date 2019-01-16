package psetter

import (
	"fmt"
	"sort"
)

// AValMap - this maps allowed values for an enumerated parameter to
// explanatory text. It forms part of the usage documentation of the
// program and will appear when the -help parameter is given by the
// user.
type AValMap map[string]string

func allowedValues(av AValMap) string {
	var avals string
	var valNames = make([]string, 0, len(av))
	var maxNameLen int

	for k := range av {
		valNames = append(valNames, k)
		if len(k) > maxNameLen {
			maxNameLen = len(k)
		}
	}
	sort.Strings(valNames)
	sep := ""
	for _, v := range valNames {
		avals += sep + fmt.Sprintf("%-*s: ", maxNameLen, v) + av[v]
		sep = "\n"
	}
	return avals
}
