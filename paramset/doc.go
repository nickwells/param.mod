/*
Package paramset offers helper functions for creating a new param.PSet. It is
expected that you will use the

	paramset.NewOrPanic(...)

function which will create a new param.PSet with the standard helper set
up. Any errors detected while constructing the PSet will cause the program to
panic, this is almost certainly what you want as they constitute coding
errors. You would typically pass it a list of param.PSetOptFunc's which can
be used, for instance, to add the params to the PSet.

The advantage of using the NewOrPanic method over using the NewOrDie
method is that you can more easily write a test of your parameter setting
code. You put all the code to generate the PSet in a single method and
then call that in a panic-safe wrapper and test that it doesn't
panic. This means that any param errors will be found during testing while
still ensuring that the program will not run if they are somehow missed.

If you want to handle errors explicitly then the paramset.New(...) will
return the constructed PSet and any errors.
*/
package paramset
