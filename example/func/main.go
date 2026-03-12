package main

import (
	"fmt"
	"github.com/devfeel/mapper"
)

// 定义测试用的结构体
type (
	// 源结构体 - 用户信息
	User struct {
		ID        int64   `mapper:"id"`
		Name      string  `mapper:"name"`
		Email     string  `mapper:"email"`
		Age       int     `mapper:"age"`
		CreatedAt int64  `mapper:"created_at"`
		Score     float64
	}

	// 目标结构体 - 用户DTO
	UserDTO struct {
		ID        int64   `json:"id"`
		Name      string  `json:"name"`
		Email     string  `json:"email"`
		Age       int     `json:"age"`
		CreatedAt int64   `json:"created_at"`
		Score     float64 `json:"score"`
	}

	// 目标结构体 - 用户VO (不同字段)
	UserVO struct {
		UserID    int64   `json:"user_id"`
		UserName  string  `json:"user_name"`
		Email     string  `json:"email"`
		Age       int     `json:"age"`
		CreatedAt int64   `json:"created_at"`
		Score     float64 `json:"score"`
	}
)

func main() {
	fmt.Println("=== Mapper 函数式泛型示例 (MapDirect) ===")
	fmt.Println()

	// 示例1: MapDirect 同构类型映射
	fmt.Println("1. MapDirect 同构类型映射")
	user1 := User{
		ID:        1,
		Name:      "张三",
		Email:     "zhangsan@example.com",
		Age:       28,
		CreatedAt: 1704067200,
		Score:     95.5,
	}
	userCopy := mapper.MapDirect[User, User](user1)
	fmt.Printf("   源: %+v\n", user1)
	fmt.Printf("   副本: %+v\n", userCopy)
	fmt.Println()

	// 示例2: MapDirect 异构类型映射
	fmt.Println("2. MapDirect 异构类型映射 (User -> UserDTO)")
	user2 := User{
		ID:        2,
		Name:      "李四",
		Email:     "lisi@example.com",
		Age:       32,
		CreatedAt: 1704153600,
		Score:     88.0,
	}
	userDTO := mapper.MapDirect[User, UserDTO](user2)
	fmt.Printf("   User: %+v\n", user2)
	fmt.Printf("   UserDTO: %+v\n", userDTO)
	fmt.Println()

	// 示例3: MapDirectPtr 指针版本
	fmt.Println("3. MapDirectPtr 指针版本")
	user3 := &User{
		ID:        3,
		Name:      "王五",
		Email:     "wangwu@example.com",
		Age:       35,
		CreatedAt: 1704240000,
		Score:     92.0,
	}
	userDTO3 := mapper.MapDirectPtr[User, UserDTO](user3)
	fmt.Printf("   源指针: %+v\n", user3)
	fmt.Printf("   目标: %+v\n", userDTO3)
	fmt.Println()

	// 示例4: MapDirectSlice 批量映射
	fmt.Println("4. MapDirectSlice 批量映射")
	users := []User{
		{ID: 1, Name: "用户1", Email: "user1@example.com", Age: 20, CreatedAt: 1704067200, Score: 85.0},
		{ID: 2, Name: "用户2", Email: "user2@example.com", Age: 25, CreatedAt: 1704153600, Score: 90.0},
		{ID: 3, Name: "用户3", Email: "user3@example.com", Age: 30, CreatedAt: 1704240000, Score: 95.0},
	}
	userDTOs := mapper.MapDirectSlice[User, UserDTO](users)
	fmt.Printf("   源数量: %d, 目标数量: %d\n", len(users), len(userDTOs))
	for i, u := range userDTOs {
		fmt.Printf("   [%d] %+v\n", i, u)
	}
	fmt.Println()

	// 示例5: MapDirectPtrSlice 指针切片映射
	fmt.Println("5. MapDirectPtrSlice 指针切片映射")
	userPtrs := []*User{
		{ID: 10, Name: "赵六", Email: "zhaoliu@example.com", Age: 40, CreatedAt: 1704326400, Score: 80.0},
		{ID: 11, Name: "孙七", Email: "sunqi@example.com", Age: 45, CreatedAt: 1704412800, Score: 85.5},
	}
	userDTOPtrs := mapper.MapDirectPtrSlice[User, UserDTO](userPtrs)
	fmt.Printf("   源数量: %d, 目标数量: %d\n", len(userPtrs), len(userDTOPtrs))
	for i, u := range userDTOPtrs {
		if u != nil {
			fmt.Printf("   [%d] %+v\n", i, *u)
		}
	}
	fmt.Println()

	// 示例6: SafeMapDirect 带错误处理的映射
	fmt.Println("6. SafeMapDirect 带错误处理")
	user6 := User{
		ID:        6,
		Name:      "周八",
		Email:     "zhouba@example.com",
		Age:       50,
		CreatedAt: 1704499200,
		Score:     98.0,
	}
	userDTO6, err := mapper.SafeMapDirect[User, UserDTO](user6)
	if err != nil {
		fmt.Printf("   错误: %v\n", err)
	} else {
		fmt.Printf("   UserDTO: %+v\n", userDTO6)
	}
	fmt.Println()

	// 示例7: SafeMapDirectSlice 批量安全映射
	fmt.Println("7. SafeMapDirectSlice 批量安全映射")
	users7 := []User{
		{ID: 7, Name: "吴九", Email: "wujiu@example.com", Age: 22, CreatedAt: 1704585600, Score: 87.0},
		{ID: 8, Name: "郑十", Email: "zhengshi@example.com", Age: 27, CreatedAt: 1704672000, Score: 93.0},
	}
	userDTOs7, err := mapper.SafeMapDirectSlice[User, UserDTO](users7)
	if err != nil {
		fmt.Printf("   错误: %v\n", err)
	} else {
		fmt.Printf("   源数量: %d, 目标数量: %d\n", len(users7), len(userDTOs7))
		for i, u := range userDTOs7 {
			fmt.Printf("   [%d] %+v\n", i, u)
		}
	}
	fmt.Println()

	// 示例8: 性能对比演示
	fmt.Println("8. 性能对比 (缓存效果)")
	largeUsers := make([]User, 1000)
	for i := 0; i < 1000; i++ {
		largeUsers[i] = User{
			ID:        int64(i),
			Name:      fmt.Sprintf("User%d", i),
			Email:     fmt.Sprintf("user%d@example.com", i),
			Age:       20 + i%50,
			CreatedAt: 1704067200 + int64(i*86400),
			Score:     float64(60 + i%40),
		}
	}

	// 第一次调用（缓存未命中）
	result1 := mapper.MapDirectSlice[User, UserDTO](largeUsers)
	fmt.Printf("   首次映射 1000 条: %d 条\n", len(result1))

	// 第二次调用（缓存命中，应该更快）
	result2 := mapper.MapDirectSlice[User, UserDTO](largeUsers)
	fmt.Printf("   二次映射 1000 条: %d 条\n", len(result2))
	fmt.Println()

	fmt.Println("=== 示例完成 ===")
}
