package psetter_test

import (
	"errors"
	"fmt"

	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
)

type myEditorMult struct{}

// Edit switches on the parameter name to reset the parameter value
func (myEditorMult) Edit(paramName, paramVal string) (string, error) {
	switch paramName {
	case "next":
		return paramVal, nil
	case "next2":
		return paramVal + ", " + paramVal, nil
	case "next3":
		return paramVal + ", " + paramVal + ", " + paramVal, nil
	}

	return "", errors.New("Unexpected parameter: " + paramName)
}

// ExampleStrListAppender_withEditor demonstrates the behaviour of the Editor.
func ExampleStrListAppender_withEditor() {
	ps := newPSetForTesting() // use paramset.NewOrPanic()

	ss := []string{"Hello"}

	var myE myEditorMult

	ps.Add("next",
		psetter.StrListAppender[string]{
			Value:  &ss,
			Editor: myE,
		}, "help text",
		param.AltNames("next2", "next3"),
	)

	fmt.Println("Before parsing")

	for i, v := range ss {
		fmt.Printf("\tss[%d] = %q\n", i, v)
	}

	ps.Parse([]string{"-next", "darkness", "-next3", "darkness"})
	fmt.Println("After  parsing")

	for i, v := range ss {
		fmt.Printf("\tss[%d] = %q\n", i, v)
	}
	// Output:
	// Before parsing
	//	ss[0] = "Hello"
	// After  parsing
	//	ss[0] = "Hello"
	//	ss[1] = "darkness"
	//	ss[2] = "darkness, darkness, darkness"
}
