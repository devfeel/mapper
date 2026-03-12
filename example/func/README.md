# Mapper 函数式泛型示例

本示例展示如何使用 `mapper_func.go` 中的函数式泛型 API 进行对象映射。

## 功能特点

- **MapDirect**: 直接返回映射结果，无需传入目标对象
- **MapDirectPtr**: 指针版本，处理指针输入
- **MapDirectSlice**: 批量映射 slice
- **MapDirectPtrSlice**: 指针切片批量映射
- **SafeMapDirect**: 带错误处理的映射
- **SafeMapDirectSlice**: 带错误处理的批量映射
- **字段映射缓存**: 重复映射时自动使用缓存提升性能

## 使用方法

```bash
# 运行示例
cd example/func
go run main.go
```

## 示例代码

```go
package main

import (
    "fmt"
    "github.com/devfeel/mapper"
)

type User struct {
    ID        int64   `mapper:"id"`
    Name      string  `mapper:"name"`
    Email     string  `mapper:"email"`
}

type UserDTO struct {
    ID        int64   `json:"id"`
    Name      string  `json:"name"`
    Email     string  `json:"email"`
}

func main() {
    // 单次映射
    user := User{ID: 1, Name: "张三", Email: "test@example.com"}
    dto := mapper.MapDirect[User, UserDTO](user)
    fmt.Printf("%+v\n", dto)

    // 批量映射
    users := []User{
        {ID: 1, Name: "用户1", Email: "user1@example.com"},
        {ID: 2, Name: "用户2", Email: "user2@example.com"},
    }
    dtos := mapper.MapDirectSlice[User, UserDTO](users)
    fmt.Printf("%+v\n", dtos)
}
```

## 性能说明

MapDirect 系列函数使用字段映射缓存：
- 首次映射时构建字段映射关系并缓存
- 后续映射自动使用缓存，减少反射调用
- 批量映射时性能提升明显（约 7x）

## 与传统方式对比

| 方式 | 代码风格 | 性能 |
|------|----------|------|
| Map (传统) | 需要传入目标对象指针 | 反射开销较大 |
| MapDirect (函数式) | 直接返回结果 | 有缓存优化，性能更好 |

## 注意事项

1. 源类型和目标类型字段名相同且类型相同才会映射
2. 支持 `mapper` tag 进行字段名映射
3. 批量操作建议使用 `MapDirectSlice` 以获得最佳性能
