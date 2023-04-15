package main

import (
	"fmt"
	"github.com/devfeel/mapper"
)

// Base model
type BaseModel struct {
	Id    int `json:"id"`
	NameX string
}

// Country model
type Country struct {
	BaseModel `json:"composite-field"`
	Name      string `json:"name"`
}

type CountryRes struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	var items *[]Country = new([]Country)
	c := &Country{BaseModel{1, "1X"}, "111"}
	*items = append(*items, *c)
	var mitems *[]CountryRes = new([]CountryRes)
	mapper.MapperSlice(items, mitems)
	fmt.Println(items)
	fmt.Println(mitems)

	to := &CountryRes{}
	mapper.Mapper(c, to)
	fmt.Println(to)
}
