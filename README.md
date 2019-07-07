[![GoDoc](https://godoc.org/github.com/nickwells/param.mod?status.png)](https://godoc.org/github.com/nickwells/param.mod)

# param
This provides parameter setting and value checking.

It differs from other similar packages in that it offers lots of control over
what is allowed. You can check the value of individual parameters and also
check for invalid combinations of parameters. The intention is that once
parameter parsing is complete you can be sure that the parameters have
legitimate values.

Parameters can be set not only through the command line but also, optionally,
through parameter files and environment variables.

You should use the latest version of this module. Old versions are not
maintained and bugs are not fixed in them.


## How to use the param package

Define a `param.PSet` and populate it, then parse the command line arguments.

Here is the simplest possible way of using it:

```go
func main() {
	var who = "World!"
	ps := paramset.NewOrDie()
	ps.Add("who", psetter.String{Value: &who}, "who to greet")
	ps.Parse()
	fmt.Println("Hello,", who)
}
```
and here's a slightly bigger example using some of the common features

```go
var param1 int64
var param2 bool

func main() {
	ps := paramset.NewOrDie(addParams,
		param.SetProgramDescription("this program will do cool stuff"))
	ps.Parse()
```

The work is done mostly in the addParam function which should take a pointer to a
`param.PSet` and return an error if anything goes wrong.

```go
func addParams(ps *param.PSet) error {
	ps.Add("param-1",
		psetter.Int64{
			Value:  &param1,
			Checks: []check.Int64{check.Int64LT(42)},
		},
		"this sets the value of param1",
		param.AltName("p1"))
		
	ps.Add("param-2", psetter.Bool{Value: &param2},
		"this sets the value of param2")
		
	return nil
}
```

This illustrates a simple use of the param package with a simple boolean flag
and an integer which is checked to ensure it's less than 42. You can specify
the behaviour much more precisely if you want.

Additionally you can have positional parameters as well as named
parameters. These must come at the front of the supplied parameters and can
only be set through the command line.

You can specify a terminal parameter (by default `--`) and the remaining
parameters will be available for further processing without being parsed. If
you intend for your users to supply additional parameters you must set a
handler function using the `PSet.SetRemHandler` method. The handler can
either process the additional arguments or else, if `param.NullRemHandler`
has been given, you can process the remainder after `PSet.Parse` has
completed. If you haven't done either of these things the default behaviour
is to report the additional arguments as an error.

## Setters
Each parameter is associated with a `Setter` which sets the parameter value
either from the associated value or directly if the parameter expects no
values. There are numerous setters predefined and many of these allow you to
specify additional checks on the parameter value. See the `psetter` package
for a full list. Additionally you can write your own setter which you can use
to populate some bespoke structure. The setters and associated checks provide
much of the power and flexibility of the `param` package.

## Actions
Each parameter can have a list of associated action functions which will be
called after the value has been set. These action functions can perform
additional checks on the parameter or can be used for parameter-specific
actions. Various standard actions are provided in the `paction` package.

## Attributes
You can set additional attributes on parameters to control:
- whether a parameter can only be used on the command line
- whether the parameter must be set
- whether it can be changed once set. This could be useful to set a value in
  a system-wide file and prevent accidental changes from the system default
- hiding the parameter from the standard help message. This is used by many
  of the standard parameters to avoid showing repeatedly the same, well-known
  parameters. You can use it to hide some more obscure options from the
  standard help message.

## Standard parameters
The default behaviour of the package will add some standard
parameters. These allow the user to see a help message which is automatically
generated from the parameters added above. This can be in varying levels of
detail. If you just pass the `-help` param you will get the standard help
message but the `-help-full` parameter shows any hidden parameters and the
`-help-short` parameter gives a help message in summary form.

Additionally the standard parameters enable you to
- examine where parameters have been set which can be useful if you have
  several parameter files and environment prefixes
- list parameters from parameter files which have not been recognised
  (parameters from the command line must be recognised). Again, this can be
  useful to identify misspelled entries in parameter files
- choose to exit after parsing which can be useful when debugging the
  parameters you have added
- choose to not exit on errors which can be useful in emergencies if your
  command is being too strict.
- choose to not print the errors which might be useful in a circumstance
  where the exit status is all that is needed
- control the level of detail in the help message with the help-full and
  help-summary parameters.
- precisely control which groups of parameters you want to see or hide in
  help messages.
- provide a file holding paramters to be parsed

## The help message
The standard help message generated if the user passes the -help parameter
will show the program description and the non-hidden parameters. For each
parameter it will show:
* the parameter description
* any alternative names
* the initial value of the parameter
* the allowed values and whether there are any additional constraints

Additionally if there are any configuration files that have been specified
(use the `SetConfigFile` and `AddConfigFile` functions on the `PSet`) or
any environment variable prefixes have been given (use the `SetEnvPrefix` and
`AddEnvPrefix` functions on the `PSet`) these will be reported at the end
of the help message.

## Parameter Groups
Parameters can be grouped together so that they are reported together rather
than in alphabetical order. This is to allow logically related parameters to
be reported together and shown or hidden together. The standard parameters
offer a way of showing the help message just for specified parameter
groups. You can add a description to be shown for the parameter group and you
can have configuration files and environment variable prefixes which are
specific to just the parameters in the group (use corresponding `SetGroup...`
and `AddGroup...` functions on the `PSet`).

## Parameter files
Parameter files are intended to allow common parameters to be set once in the
file and avoid the need to repeat them each time a command is run. The files
can include other files through a `#include` directive and comments can be
added with a leading `//`. The values of the include directives and comment
introducers can be changed programmatically or they can be disallowed. White
space and blank lines are ignored and the parameters must not have their
leading `-` or `--` characters.

Files containing parameters can contain parameters which are not intended for
the current program, this is to allow a single parameter file to be used by a
collection of binaries; unrecognised parameters are silently
ignored. Alternatively, if you prefix a parameter with a program name and a `/`
then the parameter will only be used if the program name matches that of the
running program. In this latter case the parameter must be recognised and it
is an error if it is not.

Note that having parameter files, especially with the ability to include
other files can cause problems. For instance, it can be confusing to see
where a parameter has been set.  In order to help use this feature the
standard parameter `params-show-where-set` can be given which will show all
the places where a parameter has been set.

Similarly with unrecognised parameters being silently ignored so as to allow
other programs parameters to be present in the file a typo may easily go
unnoticed. In order to help use this feature the standard parameter
`params-show-unused` can be given which will show all the unused parameters
allowing you to identify any misspelled entries.

## Environment variables
You can specify that a program can initialise the parameters from environment
variables as well. You must specify the prefix to be used and then the
environment will be checked for varaibles having that prefix and with the
remainder of the variable matching the parameter name. The parameter name is
modified to change dashes to underscores when checking environment variables.
