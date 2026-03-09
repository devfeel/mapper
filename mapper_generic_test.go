package mapper

import (
	"testing"
)

// 测试类型定义
type (
	SourceUser struct {
		Name string
		Age  int
	}

	TargetUser struct {
		Name string
		Age  int
	}

	SourceWithTag struct {
		UserName string `mapper:"name"`
		UserAge  int    `mapper:"age"`
	}

	TargetWithTag struct {
		Name string `mapper:"name"`
		Age  int    `mapper:"age"`
	}
)

func Test_Generic_Map(t *testing.T) {
	from := &SourceUser{Name: "test", Age: 25}
	to := &TargetUser{}

	err := Map(from, to)
	if err != nil {
		t.Error("Map failed:", err)
	}

	if to.Name != from.Name || to.Age != from.Age {
		t.Error("Map result not match", to, from)
	} else {
		t.Log("Map success:", to)
	}
}

func Test_Generic_MapTo(t *testing.T) {
	from := &SourceWithTag{UserName: "test", UserAge: 25}
	to := &TargetWithTag{}

	err := MapTo(from, to)
	if err != nil {
		t.Error("MapTo failed:", err)
	}

	if to.Name != from.UserName || to.Age != from.UserAge {
		t.Error("MapTo result not match", to, from)
	} else {
		t.Log("MapTo success:", to)
	}
}

func Test_Generic_MapSlice(t *testing.T) {
	fromSlice := []SourceUser{
		{Name: "user1", Age: 10},
		{Name: "user2", Age: 20},
		{Name: "user3", Age: 30},
	}

	var toSlice []TargetUser

	err := MapSliceGeneric(fromSlice, &toSlice)
	if err != nil {
		t.Error("MapSliceGeneric failed:", err)
	}

	if len(toSlice) != 3 {
		t.Error("MapSliceGeneric length not match")
	} else {
		t.Log("MapSliceGeneric success:", toSlice)
	}
}

func Test_Generic_MapToSlice(t *testing.T) {
	fromMap := map[string]any{
		"user1": map[string]any{"Name": "user1", "Age": 10},
		"user2": map[string]any{"Name": "user2", "Age": 20},
	}

	var toSlice []TargetUser

	err := MapToSliceGeneric(fromMap, &toSlice)
	if err != nil {
		t.Error("MapToSliceGeneric failed:", err)
	}

	if len(toSlice) != 2 {
		t.Error("MapToSliceGeneric length not match")
	} else {
		t.Log("MapToSliceGeneric success:", toSlice)
	}
}

func Benchmark_Generic_Map(b *testing.B) {
	from := &SourceUser{Name: "test", Age: 25}
	to := &TargetUser{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Map(from, to)
	}
}

func Benchmark_Traditional_Map(b *testing.B) {
	from := &SourceUser{Name: "test", Age: 25}
	to := &TargetUser{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Mapper(from, to)
	}
}

func Benchmark_Generic_MapSlice(b *testing.B) {
	fromSlice := make([]SourceUser, 100)
	for i := 0; i < 100; i++ {
		fromSlice[i] = SourceUser{Name: "user", Age: i}
	}

	var toSlice []TargetUser

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapSliceGeneric(fromSlice, &toSlice)
	}
}

func Benchmark_Traditional_MapperSlice(b *testing.B) {
	fromSlice := make([]SourceUser, 100)
	for i := 0; i < 100; i++ {
		fromSlice[i] = SourceUser{Name: "user", Age: i}
	}

	toSlice := make([]TargetUser, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapperSlice(fromSlice, toSlice)
	}
}
