package mapper

import (
	"errors"
	"reflect"
)

// ============ 函数式泛型 Mapper ============

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

	// 创建目标值
	toVal := reflect.ValueOf(&result)
	if toVal.Kind() == reflect.Ptr {
		toVal = toVal.Elem()
	}

	// 执行映射
	mapStructValue(fromVal, toVal)
	
	return result
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

// mapStructValue 结构体映射
func mapStructValue(fromVal, toVal reflect.Value) {
	if toVal.Kind() == reflect.Ptr {
		if toVal.IsNil() && toVal.CanSet() {
			toVal.Set(reflect.New(toVal.Type().Elem()))
		}
		if toVal.IsNil() {
			return
		}
		toVal = toVal.Elem()
	}

	if toVal.Kind() != reflect.Struct {
		return
	}

	// 映射同名同类型字段
	for i := 0; i < fromVal.NumField(); i++ {
		fromField := fromVal.Type().Field(i)
		toField, found := toVal.Type().FieldByName(fromField.Name)
		
		if !found {
			continue
		}

		// 类型检查
		if fromField.Type != toField.Type {
			continue
		}

		fromFieldVal := fromVal.Field(i)
		toFieldVal := toVal.FieldByName(fromField.Name)
		
		// 设置值
		if toFieldVal.CanSet() {
			toFieldVal.Set(fromFieldVal)
		}
	}
}

// ============ 错误处理函数 ============

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

	// 创建目标值
	toVal := reflect.ValueOf(&result)
	if toVal.Kind() == reflect.Ptr {
		toVal = toVal.Elem()
	}

	// 执行映射
	mapStructValue(fromVal, toVal)
	
	return result, nil
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
