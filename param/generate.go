//go:build generate

package param

//go:generate stringer -type=ValueReq
//go:generate mkfunccontrolparamtype -d "- this determines whether or not a ByName or ByPos param is expected to be set" -for-testing -t ShouldBeSet -v ShouldNotBeSet -v ShouldBeSet
