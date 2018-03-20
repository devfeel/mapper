package reflectx

import (
	"fmt"
	"reflect"
)

// Deref is Indirect for reflect.Types
func Deref(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// BaseType get baseType from reflect.Type and check is same with expected
func BaseType(t reflect.Type, expected reflect.Kind) (reflect.Type, error) {
	t = Deref(t)
	if t.Kind() != expected {
		return nil, fmt.Errorf("expected %s but got %s", expected, t.Kind())
	}
	return t, nil
}

// GetSliceType get slice's elem type from slice
func GetSliceType(slice reflect.Value) reflect.Type {
	isPtr := slice.Kind() == reflect.Ptr
	elemType := slice.Type().Elem()
	if isPtr {
		elemType = elemType.Elem()
	}
	return elemType
}
