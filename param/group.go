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

// GroupNameCheck checks that the group name is valid and returns an error if
// not
func GroupNameCheck(name string) error {
	if !groupNameCheckRE.MatchString(name) {
		return fmt.Errorf(
			"the group name '%s' is invalid. It must match: '%s'",
			name, groupNameCheckRE.String())
	}
	return nil
}

// Group holds details about a group of parameters
type Group struct {
	Name        string
	Desc        string
	Params      []*ByName
	HiddenCount int
	ConfigFiles []ConfigFileDetails
	setFrom     string
}

// AllParamsHidden returns true if all the parameters are marked as not to be
// shown in the standard usage message, false otherwise
func (g Group) AllParamsHidden() bool {
	return len(g.Params) == g.HiddenCount
}

// SetHiddenCount counts how many params have the DontShowInStdUsage
// attribute set and records this in the HiddenCount field. It also returns
// the value
func (g *Group) SetHiddenCount() int {
	g.HiddenCount = 0
	for _, p := range g.Params {
		if p.AttrIsSet(DontShowInStdUsage) {
			g.HiddenCount++
		}
	}
	return g.HiddenCount
}

// AddGroup will add a new param group to the PSet and set the
// descriptive text. It will panic if the description has already been set -
// this is to ensure that the group name is distinct. This description is
// shown when the usage message is printed. If the short-form description is
// chosen then the group name is shown instead so it's worth making it a
// useful value.
//
// A suggested standard for group names is to have parameters specific to a
// command in a group called 'cmd'. This is the default group name if none is
// set explicitly.
//
// Or for parameters specific to a package to use the package name prefixed with
//
//	'pkg.'
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
func (ps *PSet) AddGroup(name, desc string) {
	name = strings.TrimSpace(name)
	if err := GroupNameCheck(name); err != nil {
		panic("Invalid group name: " + err.Error())
	}

	g, exists := ps.groups[name]
	if exists && g.setFrom != "" {
		msg := fmt.Sprintf(
			"The description for param group %s is already set to:\n%s\nat: %s",
			name, g.Desc, g.setFrom)
		panic(msg)
	}

	stk := make([]byte, 10000)
	stkSize := runtime.Stack(stk, false)

	g = &Group{
		Name:    name,
		Desc:    desc,
		setFrom: string(stk[:stkSize]),
	}

	ps.groups[name] = g
}

// GetGroupDesc returns the description for the named group or the empty
// string if the group does not exist.
func (ps *PSet) GetGroupDesc(grpName string) string {
	g, ok := ps.groups[grpName]
	if !ok {
		return ""
	}
	return g.Desc
}

// HasGroupName returns true if the PSet has a group with the given name,
// false otherwise
func (ps *PSet) HasGroupName(grpName string) bool {
	_, ok := ps.groups[grpName]
	return ok
}
