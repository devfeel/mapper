package main

import (
	"fmt"

	"github.com/devfeel/mapper"
)

// 定义测试结构体
type (
	Source struct {
		ID       int64  `mapper:"id" json:"id" form:"id"`
		Name     string `mapper:"name" json:"name" form:"name"`
		Age      int    `mapper:"age" json:"age" form:"age"`
		Email    string `mapper:"email" json:"email" form:"email"`
		Password string `mapper:"-" json:"-" form:"-"`
		Private  string `mapper:"private" json:"private" form:"private"`
	}

	Target struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		Age      int    `json:"age"`
		Email    string `json:"email"`
		Password string `json:"-"`
		Private  string `json:"private"`
	}
)

func main() {
	fmt.Println("=== Mapper 配置项完整示例 ===")
	fmt.Println()

	// 创建 Mapper 实例
	m := mapper.NewMapper()

	// ========== 配置项说明 ==========
	fmt.Println("--- Mapper 配置项 ---")
	fmt.Println()

	// 1. 类型检查
	fmt.Println("1. SetEnabledTypeChecking - 类型检查")
	fmt.Println("   开启后，映射时会检查字段类型是否一致")
	fmt.Println("   默认: false")
	m.SetEnabledTypeChecking(true)
	fmt.Printf("   当前: %v\n", m.IsEnabledTypeChecking())
	fmt.Println()

	// 2. Mapper Tag
	fmt.Println("2. SetEnabledMapperTag - mapper 标签")
	fmt.Println("   开启后，支持使用 mapper:\"field\" 标签映射")
	fmt.Println("   默认: true")
	m.SetEnabledMapperTag(false)
	fmt.Printf("   当前: %v\n", m.IsEnabledMapperTag())
	m.SetEnabledMapperTag(true) // 恢复默认
	fmt.Println()

	// 3. Json Tag
	fmt.Println("3. SetEnabledJsonTag - json 标签")
	fmt.Println("   开启后，支持使用 json:\"field\" 标签映射")
	fmt.Println("   默认: true")
	m.SetEnabledJsonTag(false)
	fmt.Printf("   当前: %v\n", m.IsEnabledJsonTag())
	m.SetEnabledJsonTag(true) // 恢复默认
	fmt.Println()

	// 4. 自动类型转换
	fmt.Println("4. SetEnabledAutoTypeConvert - 自动类型转换")
	fmt.Println("   开启后，自动支持 Time 与 int64 (Unix时间戳) 互转")
	fmt.Println("   默认: true")
	m.SetEnabledAutoTypeConvert(false)
	fmt.Printf("   当前: %v\n", m.IsEnabledAutoTypeConvert())
	m.SetEnabledAutoTypeConvert(true) // 恢复默认
	fmt.Println()

	// 5. 结构体字段映射
	fmt.Println("5. SetEnabledMapperStructField - 结构体字段映射")
	fmt.Println("   开启后，支持嵌套结构体的自动映射")
	fmt.Println("   默认: true")
	m.SetEnabledMapperStructField(false)
	fmt.Printf("   当前: %v\n", m.IsEnabledMapperStructField())
	m.SetEnabledMapperStructField(true) // 恢复默认
	fmt.Println()

	// 6. 自定义标签
	fmt.Println("6. SetEnabledCustomTag - 自定义标签")
	fmt.Println("   开启后，可以使用自定义标签名进行映射")
	fmt.Println("   默认: false")
	m.SetEnabledCustomTag(true)
	m.SetCustomTagName("form")
	fmt.Printf("   当前: enabled=%v, tag=%s\n", m.IsEnabledCustomTag(), m.GetCustomTagName())
	m.SetEnabledCustomTag(false) // 恢复默认
	fmt.Println()

	// 7. 字段忽略标签
	fmt.Println("7. SetEnableFieldIgnoreTag - 字段忽略标签")
	fmt.Println("   开启后，支持使用 mapper:\"-\" 忽略字段")
	fmt.Println("   默认: false")
	m.SetEnableFieldIgnoreTag(true)
	fmt.Printf("   当前: %v\n", m.IsEnableFieldIgnoreTag())
	m.SetEnableFieldIgnoreTag(false) // 恢复默认
	fmt.Println()

	// ========== 配置演示 ==========
	fmt.Println("--- 配置演示 ---")
	fmt.Println()

	// 重置所有配置
	m.SetEnabledMapperTag(true)
	m.SetEnabledJsonTag(true)
	m.SetEnabledAutoTypeConvert(true)
	m.SetEnabledMapperStructField(true)
	m.SetEnabledCustomTag(false)

	// 示例: 使用默认配置映射
	source1 := Source{
		ID:       1,
		Name:     "张三",
		Age:      28,
		Email:    "zhangsan@example.com",
		Password: "secret123",
		Private:  "hidden",
	}

	var target1 Target
	err := m.Mapper(source1, &target1)
	if err != nil {
		fmt.Printf("   错误: %v\n", err)
	} else {
		fmt.Printf("   源: %+v\n", source1)
		fmt.Printf("   目标: %+v\n", target1)
	}
	fmt.Println()

	// 示例: 使用自定义标签
	m.SetEnabledCustomTag(true)
	m.SetCustomTagName("form")

	type FormTarget struct {
		ID   int64  `form:"id"`
		Name string `form:"name"`
		Age  int    `form:"age"`
	}

	source2 := Source{
		ID:   2,
		Name: "李四",
		Age:  32,
	}

	var target2 FormTarget
	err = m.Mapper(source2, &target2)
	if err != nil {
		fmt.Printf("   错误: %v\n", err)
	} else {
		fmt.Printf("   使用 form 标签映射:\n")
		fmt.Printf("   源: %+v\n", source2)
		fmt.Printf("   目标: %+v\n", target2)
	}
	fmt.Println()

	// ========== 全局配置 ==========
	fmt.Println("--- 全局配置 (standardMapper) ---")
	fmt.Println()

	// 设置全局配置
	mapper.SetEnabledTypeChecking(true)
	mapper.SetEnabledMapperTag(true)
	mapper.SetEnabledJsonTag(true)

	source3 := Source{
		ID:   3,
		Name: "王五",
		Age:  35,
	}

	var target3 Target
	err = mapper.Mapper(source3, &target3)
	if err != nil {
		fmt.Printf("   错误: %v\n", err)
	} else {
		fmt.Printf("   使用全局配置映射:\n")
		fmt.Printf("   源: %+v\n", source3)
		fmt.Printf("   目标: %+v\n", target3)
	}
	fmt.Println()

	// ========== 获取类型名称 ==========
	fmt.Println("--- 其他功能 ---")
	fmt.Println()

	typeName := m.GetTypeName(source1)
	fmt.Printf("   GetTypeName: %s\n", typeName)

	fmt.Println()
	fmt.Println("=== 示例完成 ===")
}
