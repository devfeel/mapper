package main

import (
	"fmt"
	"testing"

	"github.com/devfeel/mapper"
)

// 测试用结构体
type (
	BenchUser struct {
		ID        int64   `mapper:"id"`
		Name      string  `mapper:"name"`
		Email     string  `mapper:"email"`
		Age       int     `mapper:"age"`
		CreatedAt int64   `mapper:"created_at"`
		Score     float64
	}

	BenchUserDTO struct {
		ID        int64   `json:"id"`
		Name      string  `json:"name"`
		Email     string  `json:"email"`
		Age       int     `json:"age"`
		CreatedAt int64   `json:"created_at"`
		Score     float64 `json:"score"`
	}
)

// Benchmark_MapDirect_1000 benchmark for MapDirect with 1000 iterations
func Benchmark_MapDirect_1000(b *testing.B) {
	user := BenchUser{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		Age:       25,
		CreatedAt: 1704067200,
		Score:     95.5,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapper.MapDirect[BenchUser, BenchUserDTO](user)
	}
}

// Benchmark_MapDirectSlice_10 benchmark for MapDirectSlice with 10 items
func Benchmark_MapDirectSlice_10(b *testing.B) {
	users := make([]BenchUser, 10)
	for i := 0; i < 10; i++ {
		users[i] = BenchUser{
			ID:        int64(i),
			Name:      "User",
			Email:     "test@example.com",
			Age:       25,
			CreatedAt: 1704067200,
			Score:     95.5,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapper.MapDirectSlice[BenchUser, BenchUserDTO](users)
	}
}

// Benchmark_MapDirectSlice_100 benchmark for MapDirectSlice with 100 items
func Benchmark_MapDirectSlice_100(b *testing.B) {
	users := make([]BenchUser, 100)
	for i := 0; i < 100; i++ {
		users[i] = BenchUser{
			ID:        int64(i),
			Name:      "User",
			Email:     "test@example.com",
			Age:       25,
			CreatedAt: 1704067200,
			Score:     95.5,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapper.MapDirectSlice[BenchUser, BenchUserDTO](users)
	}
}

// Benchmark_MapDirectSlice_1000 benchmark for MapDirectSlice with 1000 items
func Benchmark_MapDirectSlice_1000(b *testing.B) {
	users := make([]BenchUser, 1000)
	for i := 0; i < 1000; i++ {
		users[i] = BenchUser{
			ID:        int64(i),
			Name:      "User",
			Email:     "test@example.com",
			Age:       25,
			CreatedAt: 1704067200,
			Score:     95.5,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapper.MapDirectSlice[BenchUser, BenchUserDTO](users)
	}
}

// Benchmark_MapDirectPtrSlice_100 benchmark for MapDirectPtrSlice with 100 items
func Benchmark_MapDirectPtrSlice_100(b *testing.B) {
	users := make([]*BenchUser, 100)
	for i := 0; i < 100; i++ {
		users[i] = &BenchUser{
			ID:        int64(i),
			Name:      "User",
			Email:     "test@example.com",
			Age:       25,
			CreatedAt: 1704067200,
			Score:     95.5,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapper.MapDirectPtrSlice[BenchUser, BenchUserDTO](users)
	}
}

// Benchmark_SafeMapDirect benchmark for SafeMapDirect
func Benchmark_SafeMapDirect(b *testing.B) {
	user := BenchUser{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		Age:       25,
		CreatedAt: 1704067200,
		Score:     95.5,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mapper.SafeMapDirect[BenchUser, BenchUserDTO](user)
	}
}

// Benchmark_SafeMapDirectSlice_100 benchmark for SafeMapDirectSlice with 100 items
func Benchmark_SafeMapDirectSlice_100(b *testing.B) {
	users := make([]BenchUser, 100)
	for i := 0; i < 100; i++ {
		users[i] = BenchUser{
			ID:        int64(i),
			Name:      "User",
			Email:     "test@example.com",
			Age:       25,
			CreatedAt: 1704067200,
			Score:     95.5,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mapper.SafeMapDirectSlice[BenchUser, BenchUserDTO](users)
	}
}

// Benchmark_Old_Mapper benchmark for traditional Mapper
func Benchmark_Old_Mapper(b *testing.B) {
	user := &BenchUser{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		Age:       25,
		CreatedAt: 1704067200,
		Score:     95.5,
	}
	userDTO := &BenchUserDTO{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapper.Mapper(user, userDTO)
	}
}

// Benchmark_Old_MapperSlice benchmark for traditional MapperSlice
func Benchmark_Old_MapperSlice(b *testing.B) {
	users := make([]BenchUser, 100)
	for i := 0; i < 100; i++ {
		users[i] = BenchUser{
			ID:        int64(i),
			Name:      "User",
			Email:     "test@example.com",
			Age:       25,
			CreatedAt: 1704067200,
			Score:     95.5,
		}
	}
	userDTOs := &[]BenchUserDTO{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapper.MapperSlice(users, userDTOs)
	}
}

// 运行示例
func main() {
	fmt.Println("=== Mapper 性能测试示例 ===")
	fmt.Println()
	fmt.Println("运行基准测试:")
	fmt.Println("  go test -bench=Benchmark -benchmem ./benchmark/")
	fmt.Println()
	fmt.Println("运行特定测试:")
	fmt.Println("  go test -bench=MapDirectSlice_100 -benchmem ./benchmark/")
	fmt.Println()
	fmt.Println("示例输出:")
	fmt.Println("  Benchmark_MapDirectSlice_100-8   	  20000	     58000 ns/op")
	fmt.Println("  说明: 每次映射 100 条数据，耗时约 58 微秒")
}
