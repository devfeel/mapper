package main

import (
	"fmt"
	"time"

	"github.com/devfeel/mapper"
)

// Test structs with 10 fields
type Source struct {
	Name    string
	Age     int
	Email   string
	Phone   string
	Address string
	City    string
	Country string
	ZipCode string
	Status  int
	Score   float64
}

type Dest struct {
	Name    string
	Age     int
	Email   string
	Phone   string
	Address string
	City    string
	Country string
	ZipCode string
	Status  int
	Score   float64
}

func main() {
	// Create test data
	src := &Source{
		Name:    "test",
		Age:     25,
		Email:   "test@example.com",
		Phone:   "1234567890",
		Address: "123 Test St",
		City:    "TestCity",
		Country: "TestCountry",
		ZipCode: "12345",
		Status:  1,
		Score:   95.5,
	}

	// Test 1: Single mapping
	fmt.Println("=== Test 1: Single Mapping ===")
	start := time.Now()
	for i := 0; i < 1000; i++ {
		dest := &Dest{}
		mapper.Mapper(src, dest)
	}
	elapsed := time.Since(start)
	fmt.Printf("1000 single mappings: %v\n", elapsed)
	fmt.Printf("Average per mapping: %v\n", elapsed/time.Duration(1000))

	// Test 2: Slice mapping
	fmt.Println("\n=== Test 2: Slice Mapping ===")
	srcList := make([]Source, 100)
	for i := 0; i < 100; i++ {
		srcList[i] = Source{
			Name:    "test",
			Age:     25,
			Email:   "test@example.com",
			Phone:   "1234567890",
			Address: "123 Test St",
			City:    "TestCity",
			Country: "TestCountry",
			ZipCode: "12345",
			Status:  1,
			Score:   95.5,
		}
	}

	start = time.Now()
	for i := 0; i < 100; i++ {
		var destList []Dest
		mapper.MapperSlice(srcList, &destList)
	}
	elapsed = time.Since(start)
	fmt.Printf("100 slice mappings (100 items each): %v\n", elapsed)
	fmt.Printf("Average per slice: %v\n", elapsed/time.Duration(100))

	// Test 3: Different types
	fmt.Println("\n=== Test 3: Different Type Pairs ===")
	type TypeA struct{ Name string }
	type TypeB struct{ Name string }
	type TypeC struct{ Value int }
	type TypeD struct{ Value int }

	start = time.Now()
	for i := 0; i < 1000; i++ {
		a := &TypeA{Name: "test"}
		b := &TypeB{}
		mapper.Mapper(a, b)

		c := &TypeC{Value: 100}
		d := &TypeD{}
		mapper.Mapper(c, d)
	}
	elapsed = time.Since(start)
	fmt.Printf("1000 different type mappings: %v\n", elapsed)

	fmt.Println("\n=== All tests completed ===")
}
