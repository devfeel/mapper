package main

import (
	"reflect"
	"time"
)

type (
	User struct {
		Name  string
		Age   int
		Score decimal
		Time  time.Time
	}

	Student struct {
		Name  string
		Age   int
		Score decimal
		Time  int64
	}
)

type decimal struct {
	value float32
}

type DecimalWrapper struct {
	// mapper.BaseTypeWrapper
}

func (w *DecimalWrapper) IsType(value reflect.Value) bool {
	if _, ok := value.Interface().(decimal); ok {
		return true
	}
	return false
}

func main() {
	/*
		mapper.UseWrapper(&DecimalWrapper{})
		user := &User{Name: "test", Age: 10, Score: decimal{value: 1}, Time: time.Now()}
		stu := &Student{}

		mapper.AutoMapper(user, stu)

		fmt.Println(user)
		fmt.Println(stu)
	*/
}
