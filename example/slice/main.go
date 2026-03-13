package main

import (
	"fmt"

	"github.com/devfeel/mapper"
)

// 定义测试结构体
type (
	// 源结构体
	SourceUser struct {
		ID        int64   `mapper:"id"`
		Name      string  `mapper:"name"`
		Email     string  `mapper:"email"`
		Age       int     `mapper:"age"`
		CreatedAt int64   `mapper:"created_at"`
		Score     float64
		Status    int
	}

	// 目标结构体 - DTO
	UserDTO struct {
		ID        int64   `json:"id"`
		Name      string  `json:"name"`
		Email     string  `json:"email"`
		Age       int     `json:"age"`
		CreatedAt int64   `json:"created_at"`
		Score     float64 `json:"score"`
		Status    int     `json:"status"`
	}

	// 目标结构体 - VO (不同字段)
	UserVO struct {
		UserID    int64   `json:"user_id"`
		UserName  string  `json:"user_name"`
		Email     string  `json:"email"`
		Age       int     `json:"age"`
		CreatedAt int64   `json:"created_at"`
		Score     float64 `json:"score"`
		IsActive  bool    `json:"is_active"`
	}
)

func main() {
	fmt.Println("=== Slice 批量映射进阶示例 ===")
	fmt.Println()

	// 准备测试数据
	users := []SourceUser{
		{ID: 1, Name: "张三", Email: "zhangsan@example.com", Age: 28, CreatedAt: 1704067200, Score: 95.5, Status: 1},
		{ID: 2, Name: "李四", Email: "lisi@example.com", Age: 32, CreatedAt: 1704153600, Score: 88.0, Status: 1},
		{ID: 3, Name: "王五", Email: "wangwu@example.com", Age: 35, CreatedAt: 1704240000, Score: 92.0, Status: 0},
		{ID: 4, Name: "赵六", Email: "zhaoliu@example.com", Age: 40, CreatedAt: 1704326400, Score: 85.5, Status: 1},
		{ID: 5, Name: "孙七", Email: "sunqi@example.com", Age: 45, CreatedAt: 1704412800, Score: 90.0, Status: 0},
	}

	// 示例1: 基本 Slice 映射
	fmt.Println("1. 基本 Slice 映射 (同构类型)")
	userCopies := mapper.MapDirectSlice[SourceUser, SourceUser](users)
	fmt.Printf("   源数量: %d, 目标数量: %d\n", len(users), len(userCopies))
	fmt.Println()

	// 示例2: 异构类型映射
	fmt.Println("2. 异构类型映射 (SourceUser -> UserDTO)")
	userDTOs := mapper.MapDirectSlice[SourceUser, UserDTO](users)
	fmt.Printf("   源数量: %d, 目标数量: %d\n", len(users), len(userDTOs))
	for i, u := range userDTOs {
		fmt.Printf("   [%d] ID=%d Name=%s Score=%.1f\n", i, u.ID, u.Name, u.Score)
	}
	fmt.Println()

	// 示例3: 指针 Slice 映射
	fmt.Println("3. 指针 Slice 映射")
	userPtrs := []*SourceUser{
		{ID: 10, Name: "周八", Email: "zhouba@example.com", Age: 50, CreatedAt: 1704499200, Score: 98.0, Status: 1},
		{ID: 11, Name: "吴九", Email: "wujiu@example.com", Age: 55, CreatedAt: 1704585600, Score: 87.5, Status: 1},
	}
	userDTOPtrs := mapper.MapDirectPtrSlice[SourceUser, UserDTO](userPtrs)
	fmt.Printf("   源数量: %d, 目标数量: %d\n", len(userPtrs), len(userDTOPtrs))
	for i, u := range userDTOPtrs {
		if u != nil {
			fmt.Printf("   [%d] ID=%d Name=%s\n", i, u.ID, u.Name)
		}
	}
	fmt.Println()

	// 示例4: 切片截取映射
	fmt.Println("4. 切片截取映射")
	largeUsers := make([]SourceUser, 100)
	for i := 0; i < 100; i++ {
		largeUsers[i] = SourceUser{
			ID:        int64(i),
			Name:      fmt.Sprintf("用户%d", i),
			Email:     fmt.Sprintf("user%d@example.com", i),
			Age:       20 + i%50,
			CreatedAt: 1704067200 + int64(i*86400),
			Score:     float64(60 + i%40),
			Status:    i % 2,
		}
	}

	// 取前10个
	subset := largeUsers[:10]
	subsetDTOs := mapper.MapDirectSlice[SourceUser, UserDTO](subset)
	fmt.Printf("   源切片: 前10个\n")
	fmt.Printf("   映射结果: %d 条\n", len(subsetDTOs))
	fmt.Println()

	// 示例5: 过滤后映射 (使用传统方式)
	fmt.Println("5. 过滤后映射 (状态为1的用户)")
	var activeUsers []SourceUser
	for _, u := range users {
		if u.Status == 1 {
			activeUsers = append(activeUsers, u)
		}
	}
	activeDTOs := mapper.MapDirectSlice[SourceUser, UserDTO](activeUsers)
	fmt.Printf("   过滤前: %d 条, 过滤后: %d 条\n", len(users), len(activeDTOs))
	for _, u := range activeDTOs {
		fmt.Printf("   - %s (Status=%d)\n", u.Name, u.Status)
	}
	fmt.Println()

	// 示例6: 大批量映射 (1000条)
	fmt.Println("6. 大批量映射 (1000条)")
	veryLargeUsers := make([]SourceUser, 1000)
	for i := 0; i < 1000; i++ {
		veryLargeUsers[i] = SourceUser{
			ID:        int64(i),
			Name:      fmt.Sprintf("用户%d", i),
			Email:     fmt.Sprintf("user%d@example.com", i),
			Age:       20 + i%50,
			CreatedAt: 1704067200,
			Score:     float64(60 + i%40),
			Status:    i % 2,
		}
	}

	// 使用并行映射 (阈值>=1000)
	result := mapper.MapDirectSlice[SourceUser, UserDTO](veryLargeUsers)
	fmt.Printf("   映射 1000 条数据，结果: %d 条\n", len(result))
	fmt.Println()

	// 示例7: 错误处理
	fmt.Println("7. 安全映射 (带错误处理)")
	usersWithError := []SourceUser{
		{ID: 1, Name: "正常用户", Email: "test@example.com", Age: 25, CreatedAt: 1704067200, Score: 90.0, Status: 1},
	}
	// 注意: SafeMapDirectSlice 会返回错误
	resultWithError, err := mapper.SafeMapDirectSlice[SourceUser, UserDTO](usersWithError)
	if err != nil {
		fmt.Printf("   错误: %v\n", err)
	} else {
		fmt.Printf("   成功映射: %d 条\n", len(resultWithError))
	}
	fmt.Println()

	fmt.Println("=== 示例完成 ===")
}
