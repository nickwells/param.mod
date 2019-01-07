package param

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
)

// DfltGroupName is the parameter group name that a parameter will have if no
// explicit group name is given. It is expected that parameters specific to
// the command will be in this group.
const DfltGroupName = "cmd"

var groupNameCheckRE *regexp.Regexp

// =============================================

func init() {
	groupNameCheckRE = regexp.MustCompile("^[a-zA-Z][-._a-zA-Z0-9]*$")
}

// groupNameCheck checks that the group name is valid and returns an error if
// not
func groupNameCheck(name string) error {
	if !groupNameCheckRE.MatchString(name) {
		return fmt.Errorf(
			"the group name '%s' is invalid. It must match: '%s'",
			name, groupNameCheckRE.String())
	}
	return nil
}

// GroupDesc records details about a group of parameters such as the
// description and where the description was set from (a stack trace)
type GroupDesc struct {
	Desc    string
	setFrom string
}

// SetGroupDescription will set the descriptive text for the parameter
// group. It will panic if the description has already been set - this is to
// ensure that the group name is distinct. This description is shown when the
// usage message is printed. If the short-form description is chosen then the
// group name is shown instead so it's worth making it a useful value.
//
// A suggested standard for group names is to have parameters specific to a
// command in a group called 'cmd'. This is the default group name if none is
// set explicitly.
//
// Or for parameters specific to a package to use the package name prefixed with
//
//     'pkg.'
//
// for instance: 'pkg.param'
//
// If you have several groups for the same package or command then you can
// add an additional suffix after a separating dash.
//
// for instance: 'cmd-formatting' or 'pkg.param-help'
//
// The group name will have any leading and trailing spaces deleted before
// use.
func (ps *ParamSet) SetGroupDescription(groupName, desc string) {
	groupName = strings.TrimSpace(groupName)
	if err := groupNameCheck(groupName); err != nil {
		panic("Invalid group name: " + err.Error())
	}

	pgd, exists := ps.paramGroups[groupName]
	if exists {
		msg := fmt.Sprintf(
			"description for param group %s is already set to:\n%s\nat: %s",
			groupName,
			pgd.Desc,
			pgd.setFrom)
		panic(msg)
	}

	stk := make([]byte, 10000)
	stkSize := runtime.Stack(stk, false)

	pgd.Desc = desc
	pgd.setFrom = string(stk[:stkSize])

	ps.paramGroups[groupName] = pgd
}

func (ps *ParamSet) GetGroupDesc(grpName string) string {
	g, ok := ps.paramGroups[grpName]
	if !ok {
		return ""
	}
	return g.Desc
}

func (ps *ParamSet) HasGroupName(grpName string) bool {
	_, ok := ps.paramGroups[grpName]
	return ok
}

// Groups will return a copy of the map of group names to group description
func (ps *ParamSet) Groups() map[string]GroupDesc {
	gd := make(map[string]GroupDesc)
	for k, v := range ps.paramGroups {
		gd[k] = v
	}
	return gd
}
