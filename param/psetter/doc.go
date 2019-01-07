/*
Package psetter contains a collection of useful types that can be used to
set parameter values of a program.

Each type satisfies the param.Setter interface. These can be used to supply
the second argument of a ParamSet Add or AddByPos method - the action
associated with the parameter. When the parameter is found while parsing the
params the Set method will be called.

A typical Setter is used to set the value of a parameter to the program.
For example below, a bool variable
    exitOnErrors
is set to true by the BoolSetter object's Set
method if the parameter
    exit-on-error
is found among the command line arguments:

    var exitOnErrors bool
    ps, err := paramset.New()
    p := ps.Add("exit-on-errors",
        psetter.BoolSetter{Value: &exitOnErrors},
        "Errors make the program exit if this flag is set to true",
        param.GroupName("MyTestGroup"))
*/
package psetter
