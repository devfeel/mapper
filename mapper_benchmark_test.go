package mapper

import (
	"testing"
)

// ============ Benchmark 测试结构体 ============

type (
	// 源结构体 - 用户信息
	BMUser struct {
		ID        int64  `mapper:"id"`
		Name      string `mapper:"name"`
		Email     string `mapper:"email"`
		Age       int    `mapper:"age"`
		CreatedAt int64  `mapper:"created_at"`
		Status    int    `mapper:"status"`
		Score     float64
	}

	// 目标结构体 - 用户DTO
	BMUserDTO struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Age       int    `json:"age"`
		CreatedAt int64  `json:"created_at"`
		Status    int    `json:"status"`
		Score     float64
	}

	// 目标结构体 - 用户VO (不同字段)
	BMUserVO struct {
		UserID    int64   `json:"user_id"`
		UserName  string  `json:"user_name"`
		Email     string  `json:"email"`
		Age       int     `json:"age"`
		CreatedAt int64   `json:"created_at"`
		Status    int     `json:"status"`
		Score     float64
	}
)

// ============ Benchmark 测试用例 ============

// Benchmark_MapDirect_100 基准测试: MapDirect 单次映射 100次
func Benchmark_MapDirect_100(b *testing.B) {
	user := BMUser{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		Age:       25,
		CreatedAt: 1704067200,
		Status:    1,
		Score:     95.5,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MapDirect[BMUser, BMUserDTO](user)
	}
}

// Benchmark_MapDirect_1K 基准测试: MapDirect 单次映射 1000次
func Benchmark_MapDirect_1K(b *testing.B) {
	user := BMUser{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		Age:       25,
		CreatedAt: 1704067200,
		Status:    1,
		Score:     95.5,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MapDirect[BMUser, BMUserDTO](user)
	}
}

// Benchmark_MapDirect_Slice_10 基准测试: MapDirectSlice 批量映射 10条
func Benchmark_MapDirectSlice_10(b *testing.B) {
	users := make([]BMUser, 10)
	for i := 0; i < 10; i++ {
		users[i] = BMUser{
			ID:        int64(i),
			Name:      "User",
			Email:     "test@example.com",
			Age:       25,
			CreatedAt: 1704067200,
			Status:    1,
			Score:     95.5,
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MapDirectSlice[BMUser, BMUserDTO](users)
	}
}

// Benchmark_MapDirectSlice_100 基准测试: MapDirectSlice 批量映射 100条
func Benchmark_MapDirectSlice_100(b *testing.B) {
	users := make([]BMUser, 100)
	for i := 0; i < 100; i++ {
		users[i] = BMUser{
			ID:        int64(i),
			Name:      "User",
			Email:     "test@example.com",
			Age:       25,
			CreatedAt: 1704067200,
			Status:    1,
			Score:     95.5,
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MapDirectSlice[BMUser, BMUserDTO](users)
	}
}

// Benchmark_MapDirectSlice_1K 基准测试: MapDirectSlice 批量映射 1000条
func Benchmark_MapDirectSlice_1K(b *testing.B) {
	users := make([]BMUser, 1000)
	for i := 0; i < 1000; i++ {
		users[i] = BMUser{
			ID:        int64(i),
			Name:      "User",
			Email:     "test@example.com",
			Age:       25,
			CreatedAt: 1704067200,
			Status:    1,
			Score:     95.5,
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MapDirectSlice[BMUser, BMUserDTO](users)
	}
}

// Benchmark_MapDirect_CacheHit 基准测试: MapDirect 缓存命中
func Benchmark_MapDirect_CacheHit(b *testing.B) {
	user := BMUser{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		Age:       25,
		CreatedAt: 1704067200,
		Status:    1,
		Score:     95.5,
	}
	
	// 预热缓存
	_ = MapDirect[BMUser, BMUserDTO](user)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MapDirect[BMUser, BMUserDTO](user)
	}
}

// Benchmark_MapDirect_DifferentTypes 基准测试: 不同类型映射
func Benchmark_MapDirect_DifferentTypes(b *testing.B) {
	user := BMUser{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		Age:       25,
		CreatedAt: 1704067200,
		Status:    1,
		Score:     95.5,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MapDirect[BMUser, BMUserVO](user)
	}
}

// Benchmark_MapDirect_Pointer 基准测试: 指针类型映射
func Benchmark_MapDirect_Pointer(b *testing.B) {
	user := &BMUser{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		Age:       25,
		CreatedAt: 1704067200,
		Status:    1,
		Score:     95.5,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MapDirectPtr[BMUser, BMUserDTO](user)
	}
}

// Benchmark_MapDirectSlice_Pointer 基准测试: 指针切片映射
func Benchmark_MapDirectSlice_Pointer(b *testing.B) {
	users := make([]*BMUser, 100)
	for i := 0; i < 100; i++ {
		users[i] = &BMUser{
			ID:        int64(i),
			Name:      "User",
			Email:     "test@example.com",
			Age:       25,
			CreatedAt: 1704067200,
			Status:    1,
			Score:     95.5,
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MapDirectPtrSlice[BMUser, BMUserDTO](users)
	}
}

// Benchmark_SafeMapDirect 基准测试: SafeMapDirect
func Benchmark_SafeMapDirect(b *testing.B) {
	user := BMUser{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		Age:       25,
		CreatedAt: 1704067200,
		Status:    1,
		Score:     95.5,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = SafeMapDirect[BMUser, BMUserDTO](user)
	}
}

// Benchmark_SafeMapDirectSlice 基准测试: SafeMapDirectSlice
func Benchmark_SafeMapDirectSlice(b *testing.B) {
	users := make([]BMUser, 100)
	for i := 0; i < 100; i++ {
		users[i] = BMUser{
			ID:        int64(i),
			Name:      "User",
			Email:     "test@example.com",
			Age:       25,
			CreatedAt: 1704067200,
			Status:    1,
			Score:     95.5,
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = SafeMapDirectSlice[BMUser, BMUserDTO](users)
	}
}

// ============ 对比基准测试 (原有实现) ============

// Benchmark_Map_Old 基准测试: 原有 Mapper 方法
func Benchmark_Map_Old(b *testing.B) {
	user := &BMUser{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		Age:       25,
		CreatedAt: 1704067200,
		Status:    1,
		Score:     95.5,
	}
	userDTO := &BMUserDTO{}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Mapper(user, userDTO)
	}
}

// Benchmark_MapSlice_Old 基准测试: 原有 MapperSlice 方法
func Benchmark_MapSlice_Old(b *testing.B) {
	users := make([]BMUser, 100)
	for i := 0; i < 100; i++ {
		users[i] = BMUser{
			ID:        int64(i),
			Name:      "User",
			Email:     "test@example.com",
			Age:       25,
			CreatedAt: 1704067200,
			Status:    1,
			Score:     95.5,
		}
	}
	userDTOs := &[]BMUserDTO{}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MapperSlice(users, userDTOs)
	}
}
