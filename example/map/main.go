package main

import (
	"encoding/json"
	"fmt"

	"github.com/devfeel/mapper"
)

// 定义测试结构体
type (
	// 用户结构体
	User struct {
		ID       int64  `mapper:"id"`
		Name     string `mapper:"name"`
		Email    string `mapper:"email"`
		Age      int    `mapper:"age"`
		Username string `json:"username"`
	}

	// 用户DTO
	UserDTO struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Age      int    `json:"age"`
		Username string `json:"username"`
	}
)

func main() {
	fmt.Println("=== Map 与 Struct 互转示例 ===")
	fmt.Println()

	// 示例1: Struct -> Map
	fmt.Println("1. Struct -> Map")
	user := User{
		ID:       1,
		Name:     "张三",
		Email:    "zhangsan@example.com",
		Age:      28,
		Username: "zhangsan",
	}

	userMap := make(map[string]interface{})
	err := mapper.AutoMapper(user, &userMap)
	if err != nil {
		fmt.Printf("   错误: %v\n", err)
	} else {
		fmt.Printf("   源: %+v\n", user)
		fmt.Printf("   Map: %v\n", userMap)
	}
	fmt.Println()

	// 示例2: Map -> Struct
	fmt.Println("2. Map -> Struct")
	valMap := map[string]interface{}{
		"id":       int64(2),
		"name":     "李四",
		"email":    "lisi@example.com",
		"age":      32,
		"username": "lisi",
	}

	var user2 User
	err = mapper.MapperMap(valMap, &user2)
	if err != nil {
		fmt.Printf("   错误: %v\n", err)
	} else {
		fmt.Printf("   Map: %v\n", valMap)
		fmt.Printf("   Struct: %+v\n", user2)
	}
	fmt.Println()

	// 示例3: Struct -> JSON
	fmt.Println("3. Struct -> JSON")
	user3 := User{
		ID:       3,
		Name:     "王五",
		Email:    "wangwu@example.com",
		Age:      35,
		Username: "wangwu",
	}

	jsonBytes, err := mapper.MapToJson(map[string]interface{}{
		"id":       user3.ID,
		"name":     user3.Name,
		"email":    user3.Email,
		"age":      user3.Age,
		"username": user3.Username,
	})
	if err != nil {
		fmt.Printf("   错误: %v\n", err)
	} else {
		fmt.Printf("   JSON: %s\n", string(jsonBytes))
	}
	fmt.Println()

	// 示例4: JSON -> Struct
	fmt.Println("4. JSON -> Struct")
	jsonStr := `{"id":4,"name":"赵六","email":"zhaoliu@example.com","age":40,"username":"zhaoliu"}`

	var jsonMap map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &jsonMap)
	if err != nil {
		fmt.Printf("   JSON解析错误: %v\n", err)
	} else {
		var user4 User
		err = mapper.MapperMap(jsonMap, &user4)
		if err != nil {
			fmt.Printf("   错误: %v\n", err)
		} else {
			fmt.Printf("   JSON: %s\n", jsonStr)
			fmt.Printf("   Struct: %+v\n", user4)
		}
	}
	fmt.Println()

	// 示例5: 使用 NewMapper 自定义配置
	fmt.Println("5. 使用 NewMapper 自定义配置")
	m := mapper.NewMapper()
	m.SetEnabledTypeChecking(true)
	m.SetEnabledJsonTag(true)

	user5 := User{
		ID:       5,
		Name:     "孙七",
		Email:    "sunqi@example.com",
		Age:      45,
		Username: "sunqi",
	}

	var userDTO5 UserDTO
	err = m.Mapper(user5, &userDTO5)
	if err != nil {
		fmt.Printf("   错误: %v\n", err)
	} else {
		fmt.Printf("   源: %+v\n", user5)
		fmt.Printf("   DTO: %+v\n", userDTO5)
	}
	fmt.Println()

	fmt.Println("=== 示例完成 ===")
}
