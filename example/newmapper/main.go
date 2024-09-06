package main

import (
	"fmt"
	"github.com/devfeel/mapper"
	"time"
)

type (
	User struct {
		Name     string
		Age      int    `mapper:"_Age"`
		Id       string `mapper:"_id"`
		AA       string `json:"Score,omitempty"`
		Data     []byte
		Students []Student
		Time     time.Time
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

func main() {
	user := &User{}
	userMap := &User{}
	teacher := &Teacher{}
	student := &Student{Name: "test", Age: 10, Id: "testId", Score: "100"}
	valMap := make(map[string]interface{})
	valMap["Name"] = "map"
	valMap["Age"] = 10
	valMap["_id"] = "x1asd"
	valMap["Score"] = 100
	valMap["Data"] = []byte{1, 2, 3, 4}
	valMap["Students"] = []byte{1, 2, 3, 4} //[]Student{*student}
	valMap["Time"] = time.Now()

	mp := mapper.NewMapper(mapper.CTypeChecking(true), mapper.CMapperTag(true))

	mp.Mapper(student, user)
	mp.AutoMapper(student, teacher)
	mp.MapperMap(valMap, userMap)

	fmt.Println("student:", student)
	fmt.Println("user:", user)
	fmt.Println("teacher", teacher)
	fmt.Println("userMap:", userMap)
}
