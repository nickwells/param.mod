package phelp

import "strings"

// makeTextMarkdownSafe replaces some characters in the text with HTML
// excapes or prefixed with a backslash. This will make the text safe for
// presenting to a GitHub-flavoured Markdown parser
func makeTextMarkdownSafe(s string) string {
	r := strings.NewReplacer(
		`&`, "&amp;",
		`'`, "&apos;",
		`<`, "&lt;",
		`>`, "&gt;",
		`"`, "&quot;",
		`*`, `\*`,
		`_`, `\_`,
		`#`, `\#`,
		`+`, `\+`,
		`-`, `\-`,
		`.`, `\.`,
		`!`, `\!`,
		`(`, `\(`,
		`)`, `\)`,
		`{`, `\{`,
		`}`, `\}`,
		`[`, `\[`,
		`]`, `\]`,
		`\`, `\\`,
	)
	return r.Replace(s)
}
