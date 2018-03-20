package reflectx

import (
	"reflect"
	"testing"
)

type Test struct {
	ID   int
	Name string
}

func TestBaseType(t *testing.T) {
	var a []int
	var test *Test
	var tests []*Test
	test = &Test{ID: 1, Name: "11"}

	destType, err := BaseType(reflect.ValueOf(a).Type(), reflect.Slice)
	if err != nil {
		t.Error(a, err)
	} else {
		t.Log(a, destType)
	}

	destType, err = BaseType(reflect.ValueOf(test).Type(), reflect.Struct)
	if err != nil {
		t.Error(*test, err)
	} else {
		t.Log(*test, destType)
	}

	destType, err = BaseType(reflect.ValueOf(tests).Type(), reflect.Slice)
	if err != nil {
		t.Error(tests, err)
	} else {
		t.Log(tests, destType)
	}

}

func TestGetSliceType(t *testing.T) {
	var tests []*Test
	destType := GetSliceType(reflect.ValueOf(tests))
	t.Log(tests, destType)
}
