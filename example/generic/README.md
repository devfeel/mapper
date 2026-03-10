# Mapper 泛型示例

本目录展示如何使用 mapper 的泛型功能。

## 泛型 API

### 1. Map[From, To] - 同构类型映射

```go
// 同类型之间复制
user1 := &User{ID: 1, Name: "张三"}
userCopy := &User{}
mapper.Map(user1, userCopy)
```

### 2. MapTo[To] - 异构类型映射

```go
// 不同类型之间映射
user := &User{ID: 1, Name: "张三"}
userDTO := &UserDTO{}
mapper.MapTo(user, userDTO)
```

### 3. MapSliceGeneric - Slice 批量映射

```go
// 批量映射
users := []User{
    {ID: 1, Name: "用户1"},
    {ID: 2, Name: "用户2"},
}
var userDTOs []UserDTO
mapper.MapSliceGeneric(users, &userDTOs)
```

### 4. MapToSliceGeneric - Map 转 Slice

```go
// Map 数据转换为 Slice
userMap := map[string]any{
    "user1": map[string]any{"id": 1, "name": "张三"},
    "user2": map[string]any{"id": 2, "name": "李四"},
}
var userDTOs []UserDTO
mapper.MapToSliceGeneric(userMap, &userDTOs)
```

### 5. MapperGeneric 实例

```go
// 创建泛型 Mapper 实例
mp := mapper.NewMapperGeneric()
mp.MapTo(target, source)
```

## 运行示例

```bash
cd example/generic
go run main.go
```

## 输出示例

```
=== Mapper 泛型示例 ===

1. 同构类型映射 (Map[From, To])
   原对象: &{1 张三 zhangsan@example.com 28 1704067200}
   复制后: &{1 张三 zhangsan@example.com 28 1704067200}

2. 异构类型映射 (MapTo[To])
   User:   &{2 李四 lisi@example.com 32 1704153600}
   UserDTO: &{2 李四 lisi@example.com 32 1704153600}

3. Slice 批量映射 (MapSliceGeneric)
   源数量: 3, 目标数量: 3
   [0] {1 用户1 user1@example.com 20 1704067200}
   [1] {2 用户2 user2@example.com 25 1704153600}
   [2] {3 用户3 user3@example.com 30 1704240000}

4. Map 转 Slice (MapToSliceGeneric)
   源数量: 2, 目标数量: 2
   [0] {1 王五 wangwu@example.com 35 1704067200}
   [1] {2 赵六 zhaoliu@example.com 40 1704153600}

5. 使用 MapperGeneric 实例
   User:  &{5 孙七 sunqi@example.com 45 1704326400}
   UserVO: &{5 孙七 sunqi@example.com 45}

=== 示例完成 ===
```

## 优势

- **类型安全**: 编译时类型检查
- **代码简洁**: 无需类型断言
- **性能更好**: 减少运行时反射开销
- **完全兼容**: 与现有 API 共存
