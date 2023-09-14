package param

import (
	"fmt"
	"regexp"
	"runtime"
	"slices"
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
// not. A group name must start with a letter followed by zero or more
// letters, digits, dots, dashes or underscores.
func GroupNameCheck(name string) error {
	if !groupNameCheckRE.MatchString(name) {
		return fmt.Errorf(
			"the group name %q is invalid. It must match: %q",
			name, groupNameCheckRE.String())
	}
	return nil
}

// Group holds details about a group of parameters
type Group struct {
	name        string
	desc        string
	params      []*ByName
	configFiles []ConfigFileDetails
	setFrom     string
}

// Name returns the group name
func (g Group) Name() string {
	return g.name
}

// Desc returns the group description
func (g Group) Desc() string {
	return g.desc
}

// Params returns a shallow copy of the slice of parameters in the group
func (g Group) Params() []*ByName {
	return slices.Clone(g.params)
}

// HiddenCount returns the number of parameters in this group which have the
// DontShowInStdUsage attribute set.
func (g Group) HiddenCount() int {
	hiddenCount := 0
	for _, p := range g.params {
		if p.AttrIsSet(DontShowInStdUsage) {
			hiddenCount++
		}
	}
	return hiddenCount
}

// AllParamsHidden returns true if all the parameters are marked as not to be
// shown in the standard usage message, false otherwise
func (g Group) AllParamsHidden() bool {
	for _, p := range g.params {
		if !p.AttrIsSet(DontShowInStdUsage) {
			return false
		}
	}
	return true
}

// ConfigFiles returns a shallow copy of the slice of group-specific
// configuration files for this group
func (g Group) ConfigFiles() []ConfigFileDetails {
	return slices.Clone(g.configFiles)
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
			name, g.desc, g.setFrom)
		panic(msg)
	}

	stk := make([]byte, 10000)
	stkSize := runtime.Stack(stk, false)

	g = &Group{
		name:    name,
		desc:    desc,
		setFrom: string(stk[:stkSize]),
	}

	ps.groups[name] = g
}

// GetGroup returns a group pointer and a bool indicating whether the PSet
// has a group with the given name,
func (ps *PSet) GetGroup(gName string) (*Group, bool) {
	g, ok := ps.groups[gName]
	return g, ok
}
