package main

import (
	"errors"
	"fmt"

	"github.com/devfeel/mapper"
)

// 定义测试结构体
type (
	User struct {
		ID   int64  `mapper:"id"`
		Name string `mapper:"name"`
		Age  int    `mapper:"age"`
	}

	UserDTO struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// 不可导出的结构体 (首字母小写)
	secretUser struct {
		ID   int64
		Name string
	}
)

func main() {
	fmt.Println("=== 错误处理场景示例 ===")
	fmt.Println()

	// 示例1: nil 值处理
	fmt.Println("1. nil 值处理")
	var nilUser *User = nil
	result := mapper.MapDirectPtr[User, UserDTO](nilUser)
	fmt.Printf("   输入: nil 指针\n")
	fmt.Printf("   输出: %v\n", result)
	fmt.Println()

	// 示例2: 空切片处理
	fmt.Println("2. 空切片处理")
	emptyUsers := []User{}
	resultSlice := mapper.MapDirectSlice[User, UserDTO](emptyUsers)
	fmt.Printf("   输入: 空切片 (len=0)\n")
	fmt.Printf("   输出: %v (len=%d)\n", resultSlice, len(resultSlice))
	fmt.Println()

	// 示例3: nil 切片处理
	fmt.Println("3. nil 切片处理")
	var nilUsers []User = nil
	resultNilSlice := mapper.MapDirectSlice[User, UserDTO](nilUsers)
	fmt.Printf("   输入: nil 切片\n")
	fmt.Printf("   输出: %v\n", resultNilSlice)
	fmt.Println()

	// 示例4: SafeMapDirect 错误处理
	fmt.Println("4. SafeMapDirect 错误处理")
	user4 := User{ID: 1, Name: "张三", Age: 28}

	// 正常情况
	dto4, err := mapper.SafeMapDirect[User, UserDTO](user4)
	if err != nil {
		fmt.Printf("   错误: %v\n", err)
	} else {
		fmt.Printf("   成功: %+v\n", dto4)
	}
	fmt.Println()

	// 示例5: SafeMapDirectSlice 错误处理
	fmt.Println("5. SafeMapDirectSlice 错误处理")
	users5 := []User{
		{ID: 1, Name: "用户1", Age: 20},
		{ID: 2, Name: "用户2", Age: 25},
	}

	dtos5, err := mapper.SafeMapDirectSlice[User, UserDTO](users5)
	if err != nil {
		fmt.Printf("   错误: %v\n", err)
	} else {
		fmt.Printf("   成功: %d 条\n", len(dtos5))
	}
	fmt.Println()

	// 示例6: 类型不匹配情况
	fmt.Println("6. 类型不匹配情况")
	type UserWrong struct {
		ID   string `mapper:"id"` // string 而不是 int64
		Name string `mapper:"name"`
		Age  int    `mapper:"age"`
	}

	type UserDTOWrong struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	userWrong := UserWrong{ID: "1", Name: "测试", Age: 20}
	// 由于类型不匹配，映射结果可能不正确
	resultWrong := mapper.MapDirect[UserWrong, UserDTOWrong](userWrong)
	fmt.Printf("   注意: 类型不匹配时，结果可能不符合预期\n")
	fmt.Printf("   源 ID 类型: %T, 值: %v\n", userWrong.ID, userWrong.ID)
	fmt.Printf("   目标 ID 类型: %T, 值: %v\n", resultWrong.ID, resultWrong.ID)
	fmt.Println()

	// 示例7: 使用传统 Mapper 的错误处理
	fmt.Println("7. 传统 Mapper 的错误处理")
	m := mapper.NewMapper()
	m.SetEnabledTypeChecking(true)

	user7 := &User{ID: 1, Name: "王五", Age: 30}
	userDTO7 := &UserDTO{}

	err = m.Mapper(user7, userDTO7)
	if err != nil {
		fmt.Printf("   错误: %v\n", err)
	} else {
		fmt.Printf("   成功: %+v\n", userDTO7)
	}
	fmt.Println()

	// 示例8: MapToSlice 错误处理
	fmt.Println("8. MapToSlice 错误处理")
	fromMaps := map[string]interface{}{
		"user1": map[string]interface{}{"id": int64(1), "name": "赵六", "age": 35},
		"user2": map[string]interface{}{"id": int64(2), "name": "孙七", "age": 40},
	}

	var toSlice []UserDTO
	err = m.MapToSlice(fromMaps, &toSlice)
	if err != nil {
		fmt.Printf("   错误: %v\n", err)
	} else {
		fmt.Printf("   成功: %d 条\n", len(toSlice))
		for _, u := range toSlice {
			fmt.Printf("   - %+v\n", u)
		}
	}
	fmt.Println()

	// 示例9: 自定义错误处理逻辑
	fmt.Println("9. 自定义错误处理")
	customErr := errors.New("自定义错误: 名称不能为空")
	fmt.Printf("   可以返回自定义错误: %v\n", customErr)
	fmt.Println()

	// 示例10: 清除缓存后重新映射
	fmt.Println("10. 清除缓存后重新映射")
	mapper.ClearFieldMappingCache()
	fmt.Println("   已清除字段映射缓存")
	fmt.Println()

	fmt.Println("=== 错误处理总结 ===")
	fmt.Println("1. MapDirect/MapDirectPtr - 不返回错误，适用于简单场景")
	fmt.Println("2. SafeMapDirect/SafeMapDirectSlice - 返回错误，适用于需要错误处理的场景")
	fmt.Println("3. 传统 Mapper 方法 - 支持更复杂的配置和错误处理")
	fmt.Println("4. 使用 nil 值时，建议使用 SafeMapDirect 系列函数")
}
