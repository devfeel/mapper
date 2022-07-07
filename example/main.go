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

	Leader struct {
		Name      string
		LeaderAge int `mapper:"_Age" form:"Age"`
	}

	JsonUser struct {
		Name string
		Age  int
		Time mapper.JSONTime
	}
)

func init() {
}

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

	mapper.SetEnabledTypeChecking(true)

	mapper.Mapper(student, user)
	mapper.AutoMapper(student, teacher)
	mapper.MapperMap(valMap, userMap)

	fmt.Println("student:", student)
	fmt.Println("user:", user)
	fmt.Println("teacher", teacher)
	fmt.Println("userMap:", userMap)

	jsonUser := &JsonUser{
		Name: "json",
		Age:  1,
		Time: mapper.JSONTime(time.Now()),
	}

	user = &User{Name: "User2", Age: 35}
	leader := &Leader{}
	mapper.Mapper(user, leader)
	fmt.Println("leader first:", leader)
	mapper.SetCustomTagName("form")
	mapper.SetEnabledCustomTag(true)
	leader = &Leader{}
	mapper.Mapper(user, leader)
	fmt.Println("leader second:", leader)

	mapper.SetEnabledMapperTag(false)
	leader = &Leader{}
	mapper.Mapper(user, leader)
	fmt.Println("leader third:", leader)

	//mapper.SetEnabledMapperTag(false)
	mapper.SetEnabledCustomTag(false)
	leader = &Leader{}
	mapper.Mapper(user, leader)
	fmt.Println("leader fourth 2:", leader)

	fmt.Println(jsonUser)
}
