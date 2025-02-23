package phelp

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/twrap.mod/twrap"
)

const (
	zshCompActionNone = "none"
	zshCompActionRepl = "replace"
	zshCompActionNew  = "new"
	zshCompActionShow = "show"
)

const (
	completionFilePerms = 0o555 // r-x r-x r-x
)

// zshCompHasAction returns true if the StdHelp zsh Completion Action is not
// None, false otherwise.
func zshCompHasAction(h StdHelp) bool {
	return h.zshCompAction != zshCompActionNone
}

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
	var specs []string

	explanation := "[" + zshSafeStr(p.Description()) + "]"

	names := p.AltNames()
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

// zshWriteCompFunc writes a zsh completion function for the current executable
func zshWriteCompFunc(ps *param.PSet, w io.Writer) {
	fmt.Fprintf(w, "#compdef %s\n\n", ps.ProgBaseName())
	fmt.Fprintf(w, "function _%s {\n", ps.ProgBaseName())
	fmt.Fprintln(w, "\t_arguments -S : \\")

	var args []string

	groups := ps.GetGroups()
	for _, g := range groups {
		for _, p := range g.Params() {
			args = append(args, zshOptSpec(p)...)
		}
	}

	fmt.Fprintf(w, "\t\t%s", strings.Join(args, " \\\n\t\t"))
	fmt.Fprintln(w, "}")
}

// zshCompletionHandler performs the appropriate action according to the
// setting of the StdHelp zshCompletionAction member. It returns a suggested
// exit status.
func zshCompletionHandler(h StdHelp, twc *twrap.TWConf, ps *param.PSet) int {
	switch h.zshCompAction {
	case zshCompActionNone:
		return 0
	case zshCompActionShow:
		zshWriteCompFunc(ps, twc.W)
		return 0
	case zshCompActionNew:
		filename := zshCompFileName(h, ps)

		err := zshMakeNewCompFile(filename, ps)
		if err == nil {
			zshCompFileNotify(h, twc, filename)
		}

		return zshHandleErr(err, ps)
	case zshCompActionRepl:
		filename := zshCompFileName(h, ps)

		err := zshReplaceCompFile(filename, ps)
		if err == nil {
			zshCompFileNotify(h, twc, filename)
		}

		return zshHandleErr(err, ps)
	}

	return zshHandleErr(
		fmt.Errorf("unknown zsh completion action: %q", h.zshCompAction),
		ps)
}

// zshHandleErr will test the error, if it is non-nil it will add the error
// to the param.PSet and return a suggested exit status of 1. Otherwise it
// returns 0
func zshHandleErr(err error, ps *param.PSet) int {
	if err == nil {
		return 0
	}

	ps.AddErr("zsh completions", err)

	return 1
}

// zshCompFileName returns the name of the completions file
func zshCompFileName(h StdHelp, ps *param.PSet) string {
	return filepath.Join(h.zshCompDir, "_"+ps.ProgBaseName())
}

// zshMakeNewCompFile will construct the named file (which must not already
// exist). It will return any errors found.
func zshMakeNewCompFile(filename string, ps *param.PSet) error {
	err := filecheck.IsNew().StatusCheck(filename)
	if err != nil {
		return err
	}

	w, err := os.OpenFile(
		filename,
		os.O_WRONLY|os.O_CREATE,
		completionFilePerms)
	if err != nil {
		return err
	}

	defer w.Close()
	zshWriteCompFunc(ps, w)

	return nil
}

// zshReplaceCompFile will construct the named file (which may already
// exist). It will return any errors found.
func zshReplaceCompFile(filename string, ps *param.PSet) error {
	const userWritePerm = 0o200
	_ = os.Chmod(filename, completionFilePerms|userWritePerm)

	w, err := os.OpenFile(
		filename,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		completionFilePerms)
	if err != nil {
		return err
	}

	_ = os.Chmod(filename, completionFilePerms)

	defer w.Close()
	zshWriteCompFunc(ps, w)

	return nil
}

// zshCompFileNotify writes a notification message informing the user that
// the completion file has been successfully created.
func zshCompFileNotify(h StdHelp, twc *twrap.TWConf, filename string) {
	if h.completionsQuiet {
		return
	}

	twc.Wrap(
		"the zsh completion function has been written to "+filename+"."+
			" You will need to run compinit and possibly restart your"+
			" zsh shell for this to take effect."+
			" Please see the zsh manual for more details.",
		0)
}
