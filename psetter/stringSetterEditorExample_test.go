package psetter_test

import (
	"errors"
	"fmt"

	"github.com/nickwells/param.mod/v7/param"
	"github.com/nickwells/param.mod/v7/paramset"
	"github.com/nickwells/param.mod/v7/psetter"
)

type myEditor struct{}

// Edit switches on the parameter name to reset the parameter value
func (myEditor) Edit(paramName, paramVal string) (string, error) {
	switch paramName {
	case "hello":
		return "Hello, " + paramVal, nil
	case "en":
		return "Hi, " + paramVal, nil
	case "fr":
		return "Bonjour, " + paramVal, nil
	case "es":
		return "Hola, " + paramVal, nil
	case "de":
		return "Guten Tag, " + paramVal, nil
	}

	return "", errors.New("Unknown language: " + paramName)
}

// ExampleString_withEditor demonstrates the behaviour of the Editor.
func ExampleString_withEditor() {
	ps := paramset.NewNoHelpNoExitNoErrRpt() // use paramset.New()

	var s string

	var myE myEditor

	ps.Add("hello",
		psetter.String[string]{
			Value:  &s,
			Editor: myE,
		}, "help text",
		param.AltNames("en", "fr", "es", "de"),
	)

	fmt.Printf("Before parsing: s = %q\n", s)
	ps.Parse([]string{"-fr", "Nick!"})
	fmt.Printf("After  parsing: s = %q\n", s)
	// Output:
	// Before parsing: s = ""
	// After  parsing: s = "Bonjour, Nick!"
}
