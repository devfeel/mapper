package main

import (
	"fmt"
	"github.com/devfeel/mapper"
)

type (
	User struct {
		Name string
		Age  int
		Id   string `mapper:"_id"`
		AA   string `json:"Score"`
	}

	Student struct {
		Name  string
		Age   int
		Id    string `mapper:"_id"`
		Score string
	}

	Teacher struct {
		Name  string
		Age   int
		Id    string `mapper:"_id"`
		Level string
	}
)

func init(){
	mapper.Register(&User{})
	mapper.Register(&Student{})
}

func main() {
	user := &User{}
	teacher:= &Teacher{}
	student := &Student{Name: "test", Age: 10, Id: "testId", Score:"100"}

	mapper.Mapper(student, user)
	mapper.AutoMapper(student, teacher)

	fmt.Println(student)
	fmt.Println(user)
	fmt.Println(teacher)
}
