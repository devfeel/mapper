## devfeel/mapper

#### Version 0.7.5
* Feature: Support for *[] to *[] with MapperSlice
* Ops: Definitive error messages
* Tips: Merge pull request #9 from MrWormHole/master, Thanks to @MrWormHole
* 2021-01-26 12:00 in ShangHai

#### Version 0.7.4
* Feature: AutoMapper&Mapper support mapper struct to map[string]interface{}
* Refactor: set MapperMapSlice to Deprecated, will remove on v1.0
* About AutoMapper::
  ```
  func Test_AutoMapper_StructToMap(t *testing.T) {
  	from := &FromStruct{Name: "From", Sex: true, AA: "AA"}
  	to := make(map[string]interface{})
  	err := AutoMapper(from, &to)
  	if err != nil {
  		t.Error("RunResult error: mapper error", err)
  	} else {
  		if to["UserName"] == "From"{
  			t.Log("RunResult success:", to)
  		}else{
  			t.Error("RunResult failed: map[UserName]", to["UserName"])
  		}
  	}
  }
  ```
* 2020-06-07 16:00 in ShangHai

#### Version 0.7.3
* Feature: add MapToSlice to mapper from map[string]interface{} to a slice of any type's ptr
* Refactor: set MapperMapSlice to Deprecated, will remove on v1.0
* About MapToSlice::
  ```
  func Test_MapToSlice(t *testing.T) {
  	var toSlice []*testStruct
  	/*fromMaps := make(map[string]interface{})
  	for i := 0; i < 10; i++ {
  		from := new(testStruct)
  		from.Name = "s" + strconv.Itoa(i)
  		from.Sex = true
  		from.Age = i
  		fromMaps[strconv.Itoa(i)] = from
  	}*/
  	fromMaps := make(map[string]interface{})
  	for i := 0; i < 10; i++ {
  		fromMap := make(map[string]interface{})
  		fromMap["Name"] = "s" + strconv.Itoa(i)
  		fromMap["Sex"] = true
  		fromMap["Age"] = i
  		fromMaps[strconv.Itoa(i)] = fromMap
  	}
  	err := MapToSlice(fromMaps, &toSlice)
  	if err != nil {
  		t.Error(err)
  	} else {
  		t.Log(toSlice, len(toSlice))
  	}
  }
  ```
* 2020-05-02 16:00 in ChangZhou

#### Version 0.7.2
* New Feature: MapperSlice support ptr and struct
* New Feature: MapperMapSlice support ptr and struct
* Detail:
  - now support two slice's element type is ptr or struct in MapperSlice
  - now support slice's element type is ptr or struct in MapperMapSlice
  - About MapperMapSlice:
  ```golang
    //view test code in mapper_test.go:Test_MapperSlice\Test_MapperStructSlice
    //type ptr
    var toSlice []*testStruct
    //type struct
    var toSlice []testStruct
  ```
  - About MapperSlice:
  ```golang
    //view test code in mapper_test.go:Test_MapperMapSlice\Test_MapperStructMapSlice
    //type ptr
    var fromSlice []*FromStruct
 	var toSlice []*ToStruct
    //type struct
    var fromSlice []FromStruct
    var toSlice []ToStruct
  ```
* 2019-11-03 16:00 in ShangHai

#### Version 0.7.1
* New Feature: Add TypeWrapper used to register custom Type Checker
* New Feature: Add UseWrapper used to add your TypeWrapper
* Update: remove isTimeField, add TimeWrapper
* Detail:
  - now only support IsType used to check type
* 2019-02-03 16:00

#### Version 0.7
* New Feature: Add TimeToUnix\UnixToTime\TimeToUnixLocation\UnixToTimeLocation used to transform Time and Unix
* New Feature: Add SetEnabledAutoTypeConvert used to set whether or not auto do type convert when field is Time and Unix
* Detail:
  - if set enabled, field will auto convert in Time and Unix
  - it will effective when fromField is time.Time and toField is int64
  - it will effective when fromField is int64 and toField is time.Time
  - it will effective when you use Mapper or MapperMap
  - default is enabled
* Example:
```golang
    type ProductBasic struct {
        ProductId    int64
        CreateTime   time.Time
    }
    type ProductGetResponse struct {
        ProductId    int64
        CreateTime   int64
    }

    from := &ProductBasic{
        ProductId:    10001,
        CreateTime:   time.Now(),
    }
    to := &ProductGetResponse{}
    mapper.Mapper(from, to)
    fmt.println(to)
```
* Update: Added type support when use MapperMap
* 2019-01-09 12:00

#### Version 0.6.5
* New Feature: Add MapToJson to mapper from map[string]interface{} to json []byte
* New Feature: Add JsonToMap mapper from json []byte to map[string]interface{}
* Example:
```golang
    // MapToJson
    fromMap := make(map[string]interface{})
    fromMap["Name"] = "test"
    fromMap["Sex"] = true
    fromMap["Age"] = 10
    data, err := MapToJson(fromMap)
    fmt.println(data, err)

    // JsonToMap
    var retMap map[string]interface{}
    err := JsonToMap(data, &retMap)
    fmt.println(retMap, err)
```
* 2019-01-04 16:00


#### Version 0.6.4
* New Feature: Add auto mapper reflect.Struct field, fixed for #3
* New Feature: Add mapper.SetEnabledMapperStructField used to set enabled flag for MapperStructField
* Detail:
  - if set true, the reflect.Struct field will auto mapper
  - fromField and toField type must be reflect.Struct and not time.Time
  - fromField and toField must be not same type
  - default is enabled
* Example:
``` golang
type ItemStruct1 struct {
	ProductId int64
}
type ItemStruct2 struct {
	ProductId int64
}
type ProductBasic struct {
	ProductId    int64
	ProductTitle string
	Item         ItemStruct1
	CreateTime   time.Time
}
type ProductGetResponse struct {
	ProductId    int64
	ProductTitle string
	Item         ItemStruct2
	CreateTime   time.Time
}

func main() {
	from := &ProductBasic{
		ProductId:    10001,
		ProductTitle: "Test Product",
		Item:         ItemStruct1{ProductId: 20},
		CreateTime:   time.Now(),
	}
	to := &ProductGetResponse{}
	mapper.AutoMapper(from, to)
	fmt.Println(to)
}
```
* Add Demo: example/structfield
* 2018-11-29 08:00

#### Version 0.6.3
* 新增当json标签含有omitempty时，忽略多余信息，自动获取第一位tag信息，感谢 #2 from @zhangmingfeng
* 完善示例 example/main
* 2018-07-16 08:00

#### Version 0.6.2
* 新增jsontime文件，用于处理需要定制Time字段json序列化格式场景
* 2018-04-09 16:00

#### Version 0.6.1
* 新增reflectx包，增加部分便捷函数
* 2018-03-20 12:00

#### Version 0.6
* 新增MapperSlice\MapperMapSlice函数，用于处理切片类转换
* MapperSlice: 将*StructA类型的Slice转换为*StructB类型的Slice,具体使用代码可参考Test_MapperSlice
* MapperMapSlice: 将map[string]map[string]interface{}转换为*Struct类型的Slice,具体使用代码可参考Test_MapperMapSlice
* 新增PackageVersion函数，用于输出当前包版本信息
* 调整：Mapper调整为自动Register类型，无需单独Register类型代码
* 更新mapper\mapper_test.go
* 2018-03-07 13:00

#### Version 0.5
* 新增SetEnabledTypeChecking函数，用于设置是否启用字段类型一致性检查，默认为不启用
* 如果SetEnabledTypeChecking = true,则在Mapper\AutoMapper时，将对两个类型的同名字段进行类型一致性检查，如果不一致自动忽略赋值
* 更新mapper\mapper_test.go
* 更新 example/main
* 2017-11-24 11:00

#### Version 0.4
* 新增MapperMap接口，该接口支持map到struct的自动映射
* MapperMap支持自动注册struct
* 目前支持自动映射类型：
* reflect.Bool
* reflect.String
* reflect.Int8\16\32\64
* reflect.Uint8\16\32\64
* reflect.Float32\64
* time.Time：支持原生time\string\[]byte
* 更新 example/main
* 2017-11-17 09:00

#### Version 0.3
* 新增AutoMapper接口，使用该接口无需提前Register类型
* 特别的，使用该接口性能会比使用Mapper下降20%
* 更新 example/main
* 2017-11-15 10:00

#### Version 0.2
* 新增兼容Json-tag标签
* 识别顺序：私有Tag > json tag > field name
* 当tag为"-"时，将忽略tag定义，使用struct field name
* 2017-11-15 10:00

#### Version 0.1
* 初始版本
* 支持不同结构体相同名称相同类型字段自动赋值
* 支持tag标签，tag关键字为 mapper
* 2017-11-14 21:00