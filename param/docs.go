/*
Package param is used for setting the starting parameters for an application.
It allows developers to define parameters in their packages and then in main()
they can call Parse and command line parameters will be compared against
the defined parameters and the corresponding values will be set.

You can add parameters to the set of params to be checked with the New
function and you can add alternative names with AltName which returns an
option function that will add the alternative name to the set of ways that a
parameter can be referenced. Similarly the GroupName function allows you to
group related parameters together so that they will be reported together in
the usage message.

The way to use this package is to create a ParamSet and then to add
parameters to it and when you've set all the parameters you want, you call
Parse on the ParamSet. You can create a ParamSet with the NewSet function but
it is more convenient to use the convenience function from the paramset
package: paramset.New as this will automatically set the mandatory helper to
the Standard helper. This will provide a common set of parameters that give a
consistent interface to your command line tools.

*/
package param
