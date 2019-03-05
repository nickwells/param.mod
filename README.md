# param
This provides parameter setting and value checking.

It differs from other similar packages in that it offers lots of control over
what is allowed. You can check the value of individual parameters and also
check for invalid combinations of parameters. The intention is that once
parameter parsing is complete you can be sure that the parameters have
legitimate values.

Parameters can be set not only through the command line but also, optionally,
through parameter files and environment variables.


## How to use the param package

Define a `param.PSet` and populate it, then parse the command line arguments.

```go
var param1 int64
var param2 bool

func main() {
	ps, err := paramset.New(addParams,
		param.SetProgramDescription("this program will do cool stuff"))
	if err != nil {
		log.Fatal("Couldn't construct the param set: ", err)
	}
	ps.Parse()
```

The work is done mostly in the addParam function which should take a pointer to a
`param.PSet` and return an error.

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

This illustrates the simplest use of the param package but you can specify
the behaviour much more precisely.

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

## Standard parameters
The default behaviour of the package is to add some standard
parameters. These allow the user to see a help message which is automatically
generated from the parameters added above. This can be in varying levels of
detail. If you just pass the `-help` param you will get the standard help
message but the `-help-full` parameter shows some hidden parameters and the
`-help-short` parameter gives a summary of the parameters.

Additionally the standard parameters offer the chance to examine where
parameters have been set and to control the parsing behaviour.

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
be reported together. The standard parameters offer a way of showing the help
message just for specified parameter groups. You can add a description to be
shown for the parameter group and you can have configuration files and
environment variable prefixes which are specific to just the parameters in
the group (use corresponding `SetGroup...` and `AddGroup...` functions on the
`PSet`).
