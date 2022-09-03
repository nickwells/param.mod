/*
Package paramset offers helper functions for creating a new param.PSet. It is
expected that you will use the

	paramset.NewOrDie(...)

function which will create a new param.PSet with the standard helper set
up. Any errors detected while constructing the PSet will cause the program to
exit, this is almost certainly what you want as they constitute coding
errors. You would typically pass it a list of param.PSetOptFunc's which can
be used, for instance, to add the params to the PSet.

If you want to handle errors explicitly then the paramset.New(...) will
return the constructed PSet and any errors.
*/
package paramset
