/*
Package pedit collects together useful editors that can be applied to
parameters. The definition of the Editor interface is given in package
psetter.

Editors are applied to the supplied parameter value before setting the
value. To use the editors you should check to see if the setter you are using
has an Editor and note how it is used.
*/
package pedit
