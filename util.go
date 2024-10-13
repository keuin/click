package click

import (
	"reflect"
	"unsafe"
)

// convInto is a shortcut of `(*V)(unsafe.Pointer(u))`, performing raw pointer cast, bypassing static type check.
// Caller should write guard expression like `_ = U(V{})` to perform explicit type-check when compiling.
func convInto[U, V any](u *U) (v *V) {
	return (*V)(unsafe.Pointer(u))
}

func tryGetLength(v any) (int, bool) {
	switch v := v.(type) {
	case interface{ Length() int }:
		return v.Length(), true
	case interface{ Len() int }:
		return v.Len(), true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice:
		return rv.Len(), true
	default:
		return 0, false
	}
}
