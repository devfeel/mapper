## devfeel/mapper

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