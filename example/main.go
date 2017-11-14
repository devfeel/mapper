package main

import (
	"fmt"
	"reflect"
	"github.com/devfeel/mapper"
)

type (
	User struct {
		Name string
		Age  int
		Id   string
		AA   string `mapper:"Score"`
		CC 	 string
	}

	Student struct {
		Name  string
		Age   int
		Id    string
		Score string
	}
)

func init(){
	mapper.Register(&User{})
	mapper.Register(&Student{})
}

func main() {
	user := &User{}
	student := &Student{Name: "test", Age: 10, Id: "testId", Score:"100"}

	mapper.Mapper(student, user)

	fmt.Println(student)
	fmt.Println(user)
}
