package phelp

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/param.mod/v4/param"
	"github.com/nickwells/param.mod/v4/param/psetter"
	"github.com/nickwells/twrap.mod/twrap"
)

const (
	zshCompGenNone = "none"
	zshCompGenRepl = "replace"
	zshCompGenNew  = "new"
	zshCompGenShow = "show"
)

// zshSafeStr returns an edited version of the string with any characters which
// might cause problems in a zsh option spec replaced with a safe alternative
func zshSafeStr(s string) string {
	s = strings.ToValidUTF8(s, "?")
	s = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) ||
			r == '"' ||
			r == '[' ||
			r == ']' ||
			r == '(' ||
			r == ')' {
			return ' '
		}
		return r
	}, s)
	return s
}

// zshMakeAltNames constructs a list of alternative names to be ignored if
// the named parameter is matched
func zshMakeAltNames(name string, names []string) string {
	altNames := ""
	prefix := "-"
	for _, altName := range names {
		if altName != name {
			altNames += prefix + altName
			prefix = " -"
		}
	}
	if altNames == "" {
		return altNames
	}
	return "(" + altNames + ")"
}

// zshNameSuffix creates the suffix to be applied to the option name to
// indicate whether it must or should be followed by an equals sign when
// setting the parameter value
func zshNameSuffix(p *param.ByName) string {
	switch p.Setter().ValueReq() {
	case param.Mandatory:
		return "="
	case param.Optional:
		return "=-"
	}
	return ""
}

// zshMsgAction creates the message and action parts of the option spec
// depending on whether or not the parameter must or can take a following
// argument. The action part is specialised according to the type of the
// setter. The setter name is used as the message part.
func zshMsgAction(p *param.ByName) string {
	valueReq := p.Setter().ValueReq()
	if valueReq == param.None {
		return ""
	}

	msgAction := ":"
	if valueReq == param.Optional {
		msgAction += ":"
	}

	sType := fmt.Sprintf("%T", p.Setter())
	msgAction += sType + ":"
	switch sType {
	case "psetter.Bool":
		msgAction += "(true false)"
	case "psetter.Pathname":
		msgAction += "_files"
	default:
		msgAction += zshMsgActionGetAllowedVals(p)
	}
	return msgAction
}

// zshMsgActionGetAllowedVals constructs a msgAction string from any allowed
// values
func zshMsgActionGetAllowedVals(p *param.ByName) string {
	var avals []string

	if getter, ok := p.Setter().(psetter.AllowedValuesMapper); ok {
		m := getter.AllowedValuesMap()
		if m != nil {
			keys, _ := m.Keys()
			avals = append(avals, keys...)
		}
	}

	if getter, ok := p.Setter().(psetter.AllowedValuesAliasMapper); ok {
		m := getter.AllowedValuesAliasMap()
		if m != nil {
			keys, _ := m.Keys()
			avals = append(avals, keys...)
		}
	}

	if len(avals) > 0 {
		sort.Strings(avals)
		return "(" + strings.Join(avals, " ") + ")"
	}

	return ""
}

// zshOptSpec returns a string suitable to appear as an option spec for a zsh
// option completion function
func zshOptSpec(p *param.ByName) []string {
	names := p.AltNames()
	specCount := len(names)
	if p.Setter().ValueReq() == param.Optional {
		specCount *= 2
	}

	specs := make([]string, 0, specCount)

	explanation := "[" + zshSafeStr(p.Description()) + "]"
	for _, name := range names {
		altNames := zshMakeAltNames(name, names)
		specs = append(specs,
			`"`+
				altNames+
				"-"+name+
				zshNameSuffix(p)+
				explanation+
				zshMsgAction(p)+
				`"`)
	}

	return specs
}

// zshCompletions writes a zsh completion function for the current executable
func zshCompletions(ps *param.PSet, w io.Writer) {
	fmt.Fprintf(w, "#compdef %s\n\n", ps.ProgBaseName())

	fmt.Fprintf(w, "function _%s {\n", ps.ProgBaseName())
	fmt.Fprintln(w, "\t_arguments -S : \\")
	groups := ps.GetGroups()
	totArgs := len(groups)
	for _, g := range groups {
		totArgs += len(g.Params)
	}
	args := make([]string, 0, totArgs)
	for _, g := range groups {
		for _, p := range g.Params {
			args = append(args, zshOptSpec(p)...)
		}
	}
	fmt.Fprintf(w, "\t\t%s", strings.Join(args, " \\\n\t\t"))
	fmt.Fprintln(w, "}")
}

// checkZshComplFile checks that the named file satisfies the constraints
// appropriate to the way it is being generated. If it is being replaced it
// may or may not exist, if it should be new then the file must not exist.
func checkZshComplFile(h StdHelp, fileName string) error {
	fileChecks := filecheck.Provisos{
		Existence: filecheck.MustNotExist,
	}
	if h.zshMakeCompletions == zshCompGenRepl {
		fileChecks.Existence = filecheck.Optional
	}

	return fileChecks.StatusCheck(fileName)
}

// zshMakeCompFile will construct the appropriately named file containing the
// completion function(s). It will return the exit status to use
func zshMakeCompFile(h StdHelp, twc *twrap.TWConf, ps *param.PSet) int {
	if h.zshMakeCompletions == zshCompGenShow {
		zshCompletions(ps, os.Stdout)
		return 0
	}

	fileName := h.zshCompletionsDir + "/_" + ps.ProgBaseName()

	err := checkZshComplFile(h, fileName)
	if err != nil {
		ps.AddErr("zsh completions file", err)
		return 1
	}

	w, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0555)
	if err != nil {
		ps.AddErr("zsh completions file", err)
		return 1
	}

	defer w.Close()
	zshCompletions(ps, w)
	twc.Wrap(
		"the zsh completion function for "+ps.ProgBaseName()+
			" has been written to "+fileName+"."+
			" You will need to run compinit and possibly restart your"+
			" zsh shell for this to take effect."+
			" Please see the zsh manual for more details.",
		0)

	return 0
}
