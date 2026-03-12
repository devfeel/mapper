package mapper

import (
	"errors"
	"reflect"
	"sync"
)

// ============ 字段映射缓存 ============

// fieldMappingCache 字段映射缓存
// key: fromType -> toType, value: field mappings
var fieldMappingCache = &sync.Map{}

// fieldMapping 字段映射信息
type fieldMapping struct {
	fromIndex []int  // 源字段索引路径
	toIndex   int    // 目标字段索引
}

// getFieldMappings 获取字段映射缓存
func getFieldMappings(fromType, toType reflect.Type) ([]fieldMapping, bool) {
	key := cacheKey(fromType, toType)
	if cached, ok := fieldMappingCache.Load(key); ok {
		return cached.([]fieldMapping), true
	}
	return nil, false
}

// cacheKey 生成缓存键
func cacheKey(fromType, toType reflect.Type) string {
	return fromType.String() + "->" + toType.String()
}

// buildFieldMappings 构建字段映射关系
func buildFieldMappings(fromType, toType reflect.Type) []fieldMapping {
	mappings := []fieldMapping{}
	
	for i := 0; i < fromType.NumField(); i++ {
		fromField := fromType.Field(i)
		
		// 尝试通过 mapper tag 或字段名找到对应字段
		fieldName := fromField.Name
		toField, found := toType.FieldByName(fieldName)
		if !found {
			// 尝试通过 mapper tag 查找
			if tag := fromField.Tag.Get("mapper"); tag != "" {
				if f, ok := toType.FieldByName(tag); ok {
					toField = f
					found = true
				}
			}
		}
		
		if found && fromField.Type == toField.Type {
			mappings = append(mappings, fieldMapping{
				fromIndex: []int{i},
				toIndex:   toField.Index[0],
			})
		}
	}
	
	// 存入缓存
	key := cacheKey(fromType, toType)
	fieldMappingCache.Store(key, mappings)
	
	return mappings
}

// ============ 函数式泛型 Mapper (优化版) ============

// MapDirect 同构/异构映射，直接返回结果
// 使用示例: dto := mapper.MapDirect[User, UserDTO](user)
func MapDirect[From, To any](from From) To {
	var result To
	
	fromVal := reflect.ValueOf(from)
	if !fromVal.IsValid() {
		return result
	}

	// 处理指针类型
	if fromVal.Kind() == reflect.Ptr {
		if fromVal.IsNil() {
			return result
		}
		fromVal = fromVal.Elem()
	}

	// 获取类型信息
	fromType := fromVal.Type()
	toType := reflect.TypeOf(result)
	
	if toType.Kind() != reflect.Struct {
		return result
	}

	// 尝试从缓存获取映射关系
	mappings, ok := getFieldMappings(fromType, toType)
	if !ok {
		mappings = buildFieldMappings(fromType, toType)
	}

	// 创建目标值
	toVal := reflect.New(toType).Elem()

	// 执行映射
	for _, m := range mappings {
		if len(m.fromIndex) > 0 {
			fromFieldVal := fromVal.FieldByIndex(m.fromIndex)
			toFieldVal := toVal.Field(m.toIndex)
			if toFieldVal.CanSet() {
				toFieldVal.Set(fromFieldVal)
			}
		}
	}
	
	return toVal.Interface().(To)
}

// MapDirectPtr 指针版本
// 使用示例: dto := mapper.MapDirectPtr[User, UserDTO](&user)
func MapDirectPtr[From, To any](from *From) *To {
	if from == nil {
		return nil
	}
	
	result := MapDirect[From, To](*from)
	return &result
}

// MapDirectSlice 批量映射
// 使用示例: dtos := mapper.MapDirectSlice[User, UserDTO](users)
func MapDirectSlice[From, To any](from []From) []To {
	if from == nil {
		return nil
	}
	
	// 尝试预获取映射关系以优化批量操作
	fromType := reflect.TypeOf((*From)(nil)).Elem()
	toType := reflect.TypeOf((*To)(nil)).Elem()
	
	_, hasCache := getFieldMappings(fromType, toType)
	if !hasCache {
		buildFieldMappings(fromType, toType)
	}
	
	result := make([]To, len(from))
	for i, v := range from {
		result[i] = MapDirect[From, To](v)
	}
	return result
}

// MapDirectPtrSlice 指针切片映射
// 使用示例: dtos := mapper.MapDirectPtrSlice[User, UserDTO](&users)
func MapDirectPtrSlice[From, To any](from []*From) []*To {
	if from == nil {
		return nil
	}
	
	result := make([]*To, len(from))
	for i, v := range from {
		if v != nil {
			t := MapDirect[From, To](*v)
			result[i] = &t
		}
	}
	return result
}

// ============ 错误处理函数 (优化版) ============

// SafeMapDirect 安全映射，忽略错误
// 使用示例: dto := mapper.SafeMapDirect[User, UserDTO](user)
func SafeMapDirect[From, To any](from From) (To, error) {
	var result To
	
	fromVal := reflect.ValueOf(from)
	if !fromVal.IsValid() {
		return result, errors.New("invalid from value")
	}

	// 处理指针类型
	if fromVal.Kind() == reflect.Ptr {
		if fromVal.IsNil() {
			return result, errors.New("from is nil pointer")
		}
		fromVal = fromVal.Elem()
	}

	// 获取类型信息
	fromType := fromVal.Type()
	toType := reflect.TypeOf(result)
	
	if toType.Kind() != reflect.Struct {
		return result, nil
	}

	// 尝试从缓存获取映射关系
	mappings, ok := getFieldMappings(fromType, toType)
	if !ok {
		mappings = buildFieldMappings(fromType, toType)
	}

	// 创建目标值
	toVal := reflect.New(toType).Elem()

	// 执行映射
	for _, m := range mappings {
		if len(m.fromIndex) > 0 {
			fromFieldVal := fromVal.FieldByIndex(m.fromIndex)
			toFieldVal := toVal.Field(m.toIndex)
			if toFieldVal.CanSet() {
				toFieldVal.Set(fromFieldVal)
			}
		}
	}
	
	return toVal.Interface().(To), nil
}

// SafeMapDirectSlice 安全批量映射
// 使用示例: dtos, err := mapper.SafeMapDirectSlice[User, UserDTO](users)
func SafeMapDirectSlice[From, To any](from []From) ([]To, error) {
	if from == nil {
		return nil, nil
	}
	
	result := make([]To, len(from))
	for i, v := range from {
		r, err := SafeMapDirect[From, To](v)
		if err != nil {
			return nil, errors.New("map slice failed at index " + string(rune(i+'0')))
		}
		result[i] = r
	}
	return result, nil
}

// ClearFieldMappingCache 清除字段映射缓存
// 用于在需要重新构建映射关系时调用
func ClearFieldMappingCache() {
	fieldMappingCache = &sync.Map{}
}
