package main

import (
	"fmt"
	"github.com/devfeel/mapper"
)

// 定义测试用的结构体
type (
	// 源结构体 - 用户信息
	User struct {
		ID        int64  `mapper:"id"`
		Name      string `mapper:"name"`
		Email     string `mapper:"email"`
		Age       int    `mapper:"age"`
		CreatedAt int64  `mapper:"created_at"`
	}

	// 目标结构体 - 用户DTO
	UserDTO struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Age       int    `json:"age"`
		CreatedAt int64  `json:"created_at"`
	}

	// 目标结构体 - 用户VO (不同字段)
	UserVO struct {
		UserID   int64  `json:"user_id"`
		UserName string `json:"user_name"`
		Email    string `json:"email"`
		Age      int    `json:"age"`
	}
)

func main() {
	fmt.Println("=== Mapper 泛型示例 ===")
	fmt.Println()

	// 示例1: 使用泛型 Map 进行同构类型映射
	fmt.Println("1. 同构类型映射 (Map[From, To])")
	user1 := &User{
		ID:        1,
		Name:      "张三",
		Email:     "zhangsan@example.com",
		Age:       28,
		CreatedAt: 1704067200,
	}
	userCopy := &User{}
	err := mapper.Map(user1, userCopy)
	if err != nil {
		fmt.Printf("   映射失败: %v\n", err)
	} else {
		fmt.Printf("   原对象: %+v\n", user1)
		fmt.Printf("   复制后: %+v\n", userCopy)
	}
	fmt.Println()

	// 示例2: 使用泛型 MapTo 进行异构类型映射
	fmt.Println("2. 异构类型映射 (MapTo[To])")
	user2 := &User{
		ID:        2,
		Name:      "李四",
		Email:     "lisi@example.com",
		Age:       32,
		CreatedAt: 1704153600,
	}
	userDTO := &UserDTO{}
	err = mapper.MapTo(user2, userDTO)
	if err != nil {
		fmt.Printf("   映射失败: %v\n", err)
	} else {
		fmt.Printf("   User:   %+v\n", user2)
		fmt.Printf("   UserDTO: %+v\n", userDTO)
	}
	fmt.Println()

	// 示例3: 使用泛型 MapSliceGeneric 进行 Slice 映射
	fmt.Println("3. Slice 批量映射 (MapSliceGeneric)")
	users := []User{
		{ID: 1, Name: "用户1", Email: "user1@example.com", Age: 20, CreatedAt: 1704067200},
		{ID: 2, Name: "用户2", Email: "user2@example.com", Age: 25, CreatedAt: 1704153600},
		{ID: 3, Name: "用户3", Email: "user3@example.com", Age: 30, CreatedAt: 1704240000},
	}
	var userDTOs []UserDTO
	err = mapper.MapSliceGeneric(users, &userDTOs)
	if err != nil {
		fmt.Printf("   映射失败: %v\n", err)
	} else {
		fmt.Printf("   源数量: %d, 目标数量: %d\n", len(users), len(userDTOs))
		for i, u := range userDTOs {
			fmt.Printf("   [%d] %+v\n", i, u)
		}
	}
	fmt.Println()

	// 示例4: 使用泛型 MapToSliceGeneric 进行 Map 转 Slice
	fmt.Println("4. Map 转 Slice (MapToSliceGeneric)")
	userMap := map[string]any{
		"user1": map[string]any{"id": int64(1), "name": "王五", "email": "wangwu@example.com", "age": int(35), "created_at": int64(1704067200)},
		"user2": map[string]any{"id": int64(2), "name": "赵六", "email": "zhaoliu@example.com", "age": int(40), "created_at": int64(1704153600)},
	}
	var userDTOsFromMap []UserDTO
	err = mapper.MapToSliceGeneric(userMap, &userDTOsFromMap)
	if err != nil {
		fmt.Printf("   映射失败: %v\n", err)
	} else {
		fmt.Printf("   源数量: %d, 目标数量: %d\n", len(userMap), len(userDTOsFromMap))
		for i, u := range userDTOsFromMap {
			fmt.Printf("   [%d] %+v\n", i, u)
		}
	}
	fmt.Println()

	// 示例5: 使用 MapperGeneric 实例
	fmt.Println("5. 使用 MapperGeneric 实例")
	mp := mapper.NewMapperGeneric()
	user5 := &User{ID: 5, Name: "孙七", Email: "sunqi@example.com", Age: 45, CreatedAt: 1704326400}
	userVO := &UserVO{}
	err = mp.MapTo(userVO, user5)
	if err != nil {
		fmt.Printf("   映射失败: %v\n", err)
	} else {
		fmt.Printf("   User:  %+v\n", user5)
		fmt.Printf("   UserVO: %+v\n", userVO)
	}
	fmt.Println()

	fmt.Println("=== 示例完成 ===")
}
