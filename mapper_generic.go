package mapper

import (
	"errors"
	"reflect"
)

// ============ 泛型全局函数 (推荐使用) ============

// Map 同构类型映射 (泛型)
// 使用示例: mapper.Map(&User{}, &User{})
func Map[From any, To any](from *From, to *To) error {
	if from == nil || to == nil {
		return errors.New("from or to is nil")
	}
	return standardMapper.Mapper(from, to)
}

// MapTo 异构类型映射 (泛型)
// 使用示例: mapper.MapTo[Target](from, &target)
func MapTo[To any](from any, to *To) error {
	if from == nil || to == nil {
		return errors.New("from or to is nil")
	}
	return standardMapper.Mapper(from, to)
}

// MapSliceGeneric 泛型 Slice 映射
// 使用示例: mapper.MapSliceGeneric(users, &targets)
func MapSliceGeneric[From any, To any](fromSlice []From, toSlice *[]To) error {
	if fromSlice == nil || toSlice == nil {
		return errors.New("fromSlice or toSlice is nil")
	}

	result := make([]To, len(fromSlice))
	for i, v := range fromSlice {
		// 创建 From 指针
		fromPtr := reflect.New(reflect.TypeOf(v)).Interface()
		reflect.ValueOf(fromPtr).Elem().Set(reflect.ValueOf(v))
		
		var target To
		err := standardMapper.Mapper(fromPtr, &target)
		if err != nil {
			return err
		}
		result[i] = target
	}
	*toSlice = result
	return nil
}

// MapToSliceGeneric 泛型 Map 转 Slice
// 使用示例: mapper.MapToSliceGeneric[Target](mapData, &targets)
func MapToSliceGeneric[T any](fromMap map[string]any, toSlice *[]T) error {
	if fromMap == nil || toSlice == nil {
		return errors.New("fromMap or toSlice is nil")
	}

	// 创建 T 类型的空切片
	result := make([]T, 0)
	
	for _, v := range fromMap {
		if data, ok := v.(map[string]any); ok {
			var target T
			err := standardMapper.MapperMap(data, &target)
			if err != nil {
				return err
			}
			result = append(result, target)
		}
	}
	*toSlice = result
	return nil
}

// ============ 泛型 Mapper 实例 (可选) ============

// MapperGeneric 泛型 Mapper 简化实例
// 内部组合标准 mapper，复用反射逻辑
type MapperGeneric struct {
	mapper IMapper
}

// NewMapperGeneric 创建泛型 Mapper 实例
func NewMapperGeneric() *MapperGeneric {
	return &MapperGeneric{
		mapper: standardMapper,
	}
}

// Map 同构类型映射
func (m *MapperGeneric) Map(from, to any) error {
	if from == nil || to == nil {
		return errors.New("from or to is nil")
	}
	return m.mapper.Mapper(from, to)
}

// MapTo 异构类型映射 - 泛型方法
// 注意: Go 接口方法不支持泛型，因此这里使用 any 再内部转换
func (m *MapperGeneric) MapTo(to any, from any) error {
	if from == nil || to == nil {
		return errors.New("from or to is nil")
	}
	return m.mapper.Mapper(from, to)
}

// MapSlice 泛型 Slice 映射
func (m *MapperGeneric) MapSlice(fromSlice, toSlice any) error {
	if fromSlice == nil || toSlice == nil {
		return errors.New("fromSlice or toSlice is nil")
	}
	return m.mapper.MapperSlice(fromSlice, toSlice)
}

// MapperGenericInstance 全局泛型 Mapper 实例
var MapperGenericInstance = NewMapperGeneric()
