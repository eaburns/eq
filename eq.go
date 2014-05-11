// Package eq defines a deep equality function that ignores unexported struct fields.
package eq

import (
	"reflect"
)

// Returns true if u are equals (==) v, ignoring unexported struct fields.
//
// Deep panics if given a channel, function, map, or unsafe pointer types.
func Deep(u, v interface{}) bool {
	if u == nil || v == nil {
		return u == v
	}
	return eq(reflect.ValueOf(u), reflect.ValueOf(v))
}

func eq(u, v reflect.Value) bool {
	if !u.IsValid() || !v.IsValid() {
		return u.IsValid() == v.IsValid()
	}
	if u.Type() != v.Type() {
		return false
	}

	switch v.Kind() {
	case reflect.Bool:
		return u.Bool() == v.Bool()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return u.Int() == v.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return u.Uint() == v.Uint()

	case reflect.Float32, reflect.Float64:
		return u.Float() == v.Float()

	case reflect.Complex64, reflect.Complex128:
		return u.Complex() == v.Complex()

	case reflect.Array, reflect.Slice:
		if u.Len() != v.Len() {
			return false
		}
		for i := 0; i < v.Len(); i++ {
			if !eq(u.Index(i), v.Index(i)) {
				return false
			}
		}
		return true

	case reflect.Interface, reflect.Ptr:
		if u.IsNil() != v.IsNil() {
			return false
		}
		return eq(u.Elem(), v.Elem())

	case reflect.String:
		return u.String() == v.String()

	case reflect.Struct:
		t := u.Type()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if exported(&f) && !eq(u.Field(i), v.Field(i)) {
				return false
			}
		}
		return true

	case reflect.Chan, reflect.Func, reflect.Map, reflect.UnsafePointer, reflect.Invalid:
		fallthrough
	default:
		panic("unsupported Kind: " + v.Kind().String())
	}
}

func exported(f *reflect.StructField) bool {
	return len(f.PkgPath) == 0
}
