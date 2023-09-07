package psetter

import "reflect"

// bitsInType returns the number of bits needed to store this type
func bitsInType(v any) int {
	vt := reflect.TypeOf(v)
	if vt == nil {
		return 0
	}
	return int(vt.Size() * 8)
}
