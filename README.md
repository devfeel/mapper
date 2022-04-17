# devfeel/mapper

A simple and easy go tools for auto mapper struct to map, struct to struct, slice to slice, map to slice, map to json.

## 1. Install

```
go get -u github.com/devfeel/mapper
```

## 2. Getting Started

Traditional Usage
```go
package main

import (
	"fmt"
	"github.com/devfeel/mapper"
)

type (
	User struct {
		Name string
		Age  int
		Id   string `mapper:"_id"`
		AA   string `json:"Score"`
		Time time.Time
	}

	Student struct {
		Name  string
		Age   int
		Id    string `mapper:"_id"`
		Score string
	}

	Teacher struct {
		Name  string
		Age   int
		Id    string `mapper:"_id"`
		Level string
	}
)

func init() {
	mapper.Register(&User{})
	mapper.Register(&Student{})
}

func main() {
	user := &User{}
	userMap := &User{}
	teacher := &Teacher{}
	student := &Student{Name: "test", Age: 10, Id: "testId", Score: "100"}
	valMap := make(map[string]interface{})
	valMap["Name"] = "map"
	valMap["Age"] = 10
	valMap["_id"] = "x1asd"
	valMap["Score"] = 100
	valMap["Time"] = time.Now()

	mapper.Mapper(student, user)
	mapper.AutoMapper(student, teacher)
	mapper.MapperMap(valMap, userMap)

	fmt.Println("student:", student)
	fmt.Println("user:", user)
	fmt.Println("teacher", teacher)
	fmt.Println("userMap:", userMap)
}
```

执行main，输出：
```
student: &{test 10 testId 100}
user: &{test 10 testId 100 0001-01-01 00:00:00 +0000 UTC}
teacher &{test 10 testId }
userMap: &{map 10 x1asd 100 2017-11-20 13:45:56.3972504 +0800 CST m=+0.006004001}
```

Object Usage

```go
package main

import (
  "fmt"
  "github.com/devfeel/mapper"
)

type (
  User struct {
    Name  string `json:"name" mapper:"name"`
    Class int    `mapper:"class"`
    Age   int    `json:"age" mapper:"-"`
  }

  Student struct {
    Name  string `json:"name" mapper:"name"`
    Class int    `mapper:"class"`
    Age   []int  `json:"age" mapper:"-"`
  }
)

func main() {
  user := &User{Name: "shyandsy", Class: 1, Age: 10}
  student := &Student{}

  // create mapper object
  m := mapper.NewMapper()

  // in the version < v0.7.8, we will use field name as key when mapping structs
  // we keep it as default behavior in this version
  m.SetEnableIgnoreFieldTag(true)

  student.Age = []int{1}

  // disable the json tag
  m.SetEnabledJsonTag(false)

  // student::age should be 1
  m.Mapper(user, student)

  fmt.Println("user:")
  fmt.Println(user)
  fmt.Println("student:")
  fmt.Println(student)
}
```

执行main，输出：
```
user:
&{shyandsy 1 10}
student:
&{shyandsy 1 [1]}
```



## Features

* 支持不同结构体相同名称相同类型字段自动赋值，使用Mapper
* 支持不同结构体Slice的自动赋值，使用MapperSlice
* 支持字段为结构体时的自动赋值
* 支持struct到map的自动映射赋值，使用Mapper
* 支持map到struct的自动映射赋值，使用MapperMap
* 支持map到struct slice的自动赋值，使用MapToSlice
* 支持map与json的互相转换
* 支持Time与Unix自动转换
* 支持tag标签，tag关键字为 mapper
* 兼容json-tag标签
* 当tag为"-"时，将忽略tag定义，使用struct field name
* 无需手动Register struct，内部自动识别
* 支持开启关闭
  * SetEnabledTypeChecking(bool)   // 类型检查
  * IsEnabledTypeChecking 
  * SetEnabledMapperTag            // mapper tag
  * IsEnabledMapperTag
  * SetEnabledJsonTag              // json tag
  * IsEnabledJsonTag
  * SetEnabledAutoTypeConvert      // auto type convert
  * IsEnabledAutoTypeConvert
  * SetEnabledMapperStructField    // mapper struct field
  * IsEnabledMapperStructField