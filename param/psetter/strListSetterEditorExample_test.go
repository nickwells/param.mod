package psetter_test

import (
	"errors"
	"fmt"

	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

type myEditorStrList struct{}

// Edit switches on the parameter name to reset the parameter value
func (myEditorStrList) Edit(paramName, paramVal string) (string, error) {
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

// ExampleStrList_withEditor demonstrates the behaviour of the Editor.
func ExampleStrList_withEditor() {
	ps := newPSetForTesting() // use paramset.NewOrDie()

	var ss []string
	var myE myEditorStrList

	ps.Add("hello",
		psetter.StrList{
			Value:  &ss,
			Editor: myE,
		}, "help text",
		param.AltName("en"),
		param.AltName("fr"),
		param.AltName("es"),
		param.AltName("de"),
	)

	fmt.Println("Before parsing")
	for i, v := range ss {
		fmt.Printf("\tss[%d] = %q\n", i, v)
	}
	ps.Parse([]string{"-fr", "Nick,Pascal,Amelie"})
	fmt.Println("After  parsing")
	for i, v := range ss {
		fmt.Printf("\tss[%d] = %q\n", i, v)
	}
	// Output:
	// Before parsing
	// After  parsing
	//	ss[0] = "Bonjour, Nick"
	//	ss[1] = "Bonjour, Pascal"
	//	ss[2] = "Bonjour, Amelie"
}
