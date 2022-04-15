package main

import (
	"fmt"
	"github.com/devfeel/mapper"
)

type (
	User struct {
		Name     string `json:"name" mapper:"name"`
		Age      int    `json:"age" mapper:"age"`
	}

	Student struct {
		Name  string `json:"name" mapper:"name"`
		Age   int    `json:"age" mapper:"-"`
	}
)

func main() {
	user := &User{Name: "test", Age: 10}
	student := &Student{}

	// create mapper object
	m := mapper.NewMapper()

	// enable the type checking
	m.SetEnabledTypeChecking(true)

	student.Age = 1

	// disable the json tag
	m.SetEnabledJsonTag(false)

	// student::age should be 1
	m.Mapper(user, student)

	fmt.Println(student)
}
