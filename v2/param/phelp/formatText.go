package phelp

import (
	"fmt"
	"io"
	"math"
	"strings"
)

// MinCharsToPrint is the least number of characters that formatText will
// print regardless of the depth of the indent
const MinCharsToPrint = 30.0

// TargetLineLength is teh number of characters that will be put on each line
// after any indent is taken into account
const TargetLineLength = 80

// printAsString prints a slice of runes as a string and clears the slice
func printAsString(w io.Writer, word []rune) []rune {
	if len(word) > 0 {
		fmt.Fprint(w, string(word))
		word = word[:0]
	}
	return word
}

// formatText prints the text onto the writer but wraps words and indents
// with a different indent on the first and subsequent lines. It will always
// print at least MinCharsToPrint chars but will try to fit the text into
// TargetLineLength chars.
func formatText(w io.Writer, text string, firstLineIndent, indent int) {
	firstLineMaxWidth := int(math.Max(MinCharsToPrint,
		float64(TargetLineLength-firstLineIndent)))
	nextMaxWidth := int(math.Max(MinCharsToPrint,
		float64(TargetLineLength-indent)))

	paras := strings.Split(text, "\n")

	for _, para := range paras {
		fmt.Fprint(w, strings.Repeat(" ", firstLineIndent))
		maxWidth := firstLineMaxWidth
		sep := strings.Repeat(" ", indent)

		lineLen := 0
		word := make([]rune, 0, len(para))
		spaces := make([]rune, 0, maxWidth)
		for _, r := range para {
			if r == ' ' {
				if len(word) > 0 {
					if lineLen == 0 {
						lineLen = len(word)
						word = printAsString(w, word)
						spaces = spaces[:0]
					} else if lineLen+len(word)+len(spaces) < maxWidth {
						lineLen += len(word) + len(spaces)
						spaces = printAsString(w, spaces)
						word = printAsString(w, word)
					} else {
						fmt.Fprintln(w)
						maxWidth = nextMaxWidth
						fmt.Fprint(w, sep)
						lineLen = len(word)
						word = printAsString(w, word)
						spaces = spaces[:0]
					}
				}

				spaces = append(spaces, r)
			} else {
				word = append(word, r)
			}
		}

		if len(word) > 0 {
			if lineLen == 0 {
				printAsString(w, word)
			} else if lineLen+len(word)+len(spaces) < maxWidth {
				printAsString(w, spaces)
				printAsString(w, word)
			} else {
				fmt.Fprintln(w)
				fmt.Fprint(w, sep)
				printAsString(w, word)
			}
		}

		fmt.Fprintln(w)
	}
}
