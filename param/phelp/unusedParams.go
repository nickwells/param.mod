package phelp

import (
	"sort"

	"github.com/nickwells/param.mod/v4/param"
	"github.com/nickwells/twrap.mod/twrap"
)

// showUnusedParams will print all the parameters that have been detected in
// configuration files or set through environment variables which were not
// recognised and were ignored. Since these sources may contain parameters
// intended for other programs it is not reported as an error but you may
// want to see what has been ignored in order to detect mistakes.
func showUnusedParams(_ StdHelp, twc *twrap.TWConf, ps *param.PSet) int {
	up := ps.UnusedParams()
	twc.Println("Unused Parameters: ", len(up)) //nolint: errcheck

	var paramsByName = make([]string, 0, len(up))
	for name := range up {
		paramsByName = append(paramsByName, name)
	}
	sort.Strings(paramsByName)

	for _, pn := range paramsByName {
		twc.Wrap(pn, paramIndent)
		for _, loc := range up[pn] {
			twc.WrapPrefixed("at: ", loc, descriptionIndent)
		}
	}
	return 0
}
