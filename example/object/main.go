package main

import (
	"fmt"
	"github.com/devfeel/mapper"
)

type (
	User struct {
		Name  string `json:"name" mapper:"name"`
		Class int    `mapper:"class"`
		Age   int    `json:"age" mapper:"-"`
	}

	Student struct {
		Name  string `json:"name" mapper:"name"`
		Class int    `mapper:"class"`
		Age   []int  `json:"age" mapper:"-"`
	}
)

func main() {
	user := &User{Name: "shyandsy", Class: 1, Age: 10}
	student := &Student{}

	// create mapper object
	m := mapper.NewMapper()

	// in the version < v0.7.8, we will use field name as key when mapping structs
	// we keep it as default behavior in this version
	m.SetEnableIgnoreFieldTag(true)

	student.Age = []int{1}

	// disable the json tag
	m.SetEnabledJsonTag(false)

	// student::age should be 1
	m.Mapper(user, student)

	fmt.Println("user:")
	fmt.Println(user)
	fmt.Println("student:")
	fmt.Println(student)
}
