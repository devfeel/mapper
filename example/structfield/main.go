package main

import (
	"fmt"
	"github.com/devfeel/mapper"
	"time"
)

type ItemStruct1 struct {
	ProductId int64
	Name      string
	Item      string
}
type ItemStruct2 struct {
	ProductId int64
	Name      string
	Item2     string
}
type ProductBasic struct {
	ProductId    int64
	CategoryType int
	ProductTitle string
	Item         ItemStruct1
	CreateTime   time.Time
}
type ProductGetResponse struct {
	ProductId    int64
	CategoryType int
	ProductTitle string
	Item         ItemStruct2
	CreateTime   time.Time
}

func main() {
	from := &ProductBasic{
		ProductId:    10001,
		CategoryType: 1,
		ProductTitle: "Test Product",
		Item:         ItemStruct1{ProductId: 20, Name: "pro", Item: "1"},
		CreateTime:   time.Now(),
	}
	to := &ProductGetResponse{}
	mapper.AutoMapper(from, to)
	fmt.Println(to)
}
