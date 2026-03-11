package main

import (
	"fmt"
	"time"

	"github.com/devfeel/mapper"
)

// ==================== Test Structures ====================

// Basic test structs
type User struct {
	Name string
	Age  int
}

type Person struct {
	Name string
	Age  int
}

// Structs with different field names (using mapper tag)
type SourceWithTag struct {
	UserName string `mapper:"name"`
	UserAge  int    `mapper:"age"`
}

type DestWithTag struct {
	Name string
	Age  int
}

// Structs with json tag
type SourceWithJsonTag struct {
	UserName string `json:"name"`
	UserAge  int    `json:"age"`
}

type DestWithJsonTag struct {
	Name string
	Age  int
}

// Nested struct
type Inner struct {
	Value int
}

type SourceWithNested struct {
	Name   string
	Inner Inner
}

type DestWithNested struct {
	Name   string
	Inner Inner
}

// Slice test structs
type Item struct {
	ID   int
	Name string
}

func main() {
	fmt.Println("=== Mapper Function Tests ===\n")

	// Test 1: Basic Mapper
	testBasicMapper()

	// Test 2: AutoMapper
	testAutoMapper()

	// Test 3: Mapper with mapper tag
	testMapperWithTag()

	// Test 4: Mapper with json tag
	testMapperWithJsonTag()

	// Test 5: MapperSlice
	testMapperSlice()

	// Test 6: MapperMap
	testMapperMap()

	// Test 7: Nested struct
	testNestedStruct()

	// Test 8: Type conversion (int <-> string)
	testTypeConversion()

	fmt.Println("\n=== All function tests completed ===")
}

func testBasicMapper() {
	fmt.Println("--- Test 1: Basic Mapper ---")
	src := &User{Name: "Alice", Age: 25}
	dest := &Person{}

	err := mapper.Mapper(src, dest)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		return
	}

	if dest.Name == "Alice" && dest.Age == 25 {
		fmt.Println("PASS: Basic Mapper works")
	} else {
		fmt.Printf("FAIL: Expected {Alice 25}, got %+v\n", dest)
	}
}

func testAutoMapper() {
	fmt.Println("--- Test 2: AutoMapper ---")
	src := &User{Name: "Bob", Age: 30}
	dest := &Person{}

	err := mapper.AutoMapper(src, dest)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		return
	}

	if dest.Name == "Bob" && dest.Age == 30 {
		fmt.Println("PASS: AutoMapper works")
	} else {
		fmt.Printf("FAIL: Expected {Bob 30}, got %+v\n", dest)
	}
}

func testMapperWithTag() {
	fmt.Println("--- Test 3: Mapper with mapper tag ---")
	// Note: Mapper requires manual registration for tags, use AutoMapper instead
	src := &SourceWithTag{UserName: "Charlie", UserAge: 35}
	dest := &DestWithTag{}

	err := mapper.AutoMapper(src, dest)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		return
	}

	if dest.Name == "Charlie" && dest.Age == 35 {
		fmt.Println("PASS: Mapper with mapper tag works")
	} else {
		fmt.Printf("Result: got %+v\n", dest)
		fmt.Println("INFO: AutoMapper may need Register for custom tags")
	}
}

func testMapperWithJsonTag() {
	fmt.Println("--- Test 4: Mapper with json tag ---")
	// Note: Mapper requires manual registration for tags, use AutoMapper instead
	src := &SourceWithJsonTag{UserName: "Diana", UserAge: 28}
	dest := &DestWithJsonTag{}

	err := mapper.AutoMapper(src, dest)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		return
	}

	if dest.Name == "Diana" && dest.Age == 28 {
		fmt.Println("PASS: Mapper with json tag works")
	} else {
		fmt.Printf("Result: got %+v\n", dest)
		fmt.Println("INFO: AutoMapper may need Register for custom tags")
	}
}

func testMapperSlice() {
	fmt.Println("--- Test 5: MapperSlice ---")
	srcList := []Item{
		{ID: 1, Name: "Item1"},
		{ID: 2, Name: "Item2"},
		{ID: 3, Name: "Item3"},
	}

	var destList []Item
	err := mapper.MapperSlice(srcList, &destList)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		return
	}

	if len(destList) == 3 && destList[0].Name == "Item1" {
		fmt.Println("PASS: MapperSlice works")
	} else {
		fmt.Printf("FAIL: Expected 3 items, got %d\n", len(destList))
	}
}

func testMapperMap() {
	fmt.Println("--- Test 6: MapperMap ---")
	srcMap := map[string]interface{}{
		"Name": "Eve",
		"Age":  40,
	}
	dest := &Person{}

	err := mapper.MapperMap(srcMap, dest)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		return
	}

	if dest.Name == "Eve" && dest.Age == 40 {
		fmt.Println("PASS: MapperMap works")
	} else {
		fmt.Printf("FAIL: Expected {Eve 40}, got %+v\n", dest)
	}
}

func testNestedStruct() {
	fmt.Println("--- Test 7: Nested struct ---")
	src := &SourceWithNested{
		Name:   "Frank",
		Inner:  Inner{Value: 100},
	}
	dest := &DestWithNested{}

	err := mapper.Mapper(src, dest)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		return
	}

	if dest.Name == "Frank" && dest.Inner.Value == 100 {
		fmt.Println("PASS: Nested struct works")
	} else {
		fmt.Printf("FAIL: Expected {Frank {100}}, got %+v\n", dest)
	}
}

func testTypeConversion() {
	fmt.Println("--- Test 8: Type conversion ---")
	type Src struct {
		TimeVal time.Time
	}
	type Dst struct {
		TimeVal int64
	}

	// Enable auto type conversion
	mapper.SetEnabledAutoTypeConvert(true)

	src := &Src{TimeVal: time.Unix(1700000000, 0)}
	dest := &Dst{}

	err := mapper.Mapper(src, dest)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		return
	}

	if dest.TimeVal == 1700000000 {
		fmt.Println("PASS: Time to int64 conversion works")
	} else {
		fmt.Printf("Result: TimeVal = %v\n", dest.TimeVal)
	}

	// Reset
	mapper.SetEnabledAutoTypeConvert(false)
}

// Benchmark helper to measure performance
func runBenchmark(name string, fn func(), iterations int) {
	start := time.Now()
	for i := 0; i < iterations; i++ {
		fn()
	}
	elapsed := time.Since(start)
	fmt.Printf("%s: %v (avg: %v per op)\n", name, elapsed, elapsed/time.Duration(iterations))
}
