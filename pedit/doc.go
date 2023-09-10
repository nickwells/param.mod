/*
Package pedit defines the Editor interface and collects together some useful
editors that can be applied to parameters.

Editors are typically applied to the supplied parameter value before value
checks are applied and before setting the value. This convention is not
enforced though so check the documentation for individual Setters.

It is recommended that Editors should only perform relatively simple
transformations on the parameter value. The more work done by the Editor the
more confusing it could be. This confusion can apply to two parties. The user
trying to understand how the program works might be confused: "I typed 'X',
why is it using 'Y'?". Also the programmer trying to make changes might not
expect to have to look in the Editor code to understand the program
behaviour.
*/
package pedit
