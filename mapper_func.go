package mapper

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sync"
)

// parallelThreshold 并行处理的阈值
// 当 slice 长度 >= 此值时启用并行映射
var parallelThreshold = 1000

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

	length := len(from)
	if length == 0 {
		return []To{}
	}

	// 预获取映射关系以优化批量操作
	fromType := reflect.TypeOf((*From)(nil)).Elem()
	toType := reflect.TypeOf((*To)(nil)).Elem()

	_, hasCache := getFieldMappings(fromType, toType)
	if !hasCache {
		buildFieldMappings(fromType, toType)
	}

	// 大 slice 使用并行处理
	if length >= parallelThreshold {
		return mapDirectSliceParallel[From, To](from)
	}

	result := make([]To, length)
	for i, v := range from {
		result[i] = MapDirect[From, To](v)
	}
	return result
}

// mapDirectSliceParallel 并行批量映射
func mapDirectSliceParallel[From, To any](from []From) []To {
	length := len(from)
	result := make([]To, length)

	// 计算合适的 worker 数量
	numCPU := runtime.NumCPU()
	numWorkers := length / 100 // 每 100 个元素一个 worker
	if numWorkers < 1 {
		numWorkers = 1
	}
	if numWorkers > numCPU {
		numWorkers = numCPU
	}

	chunkSize := length / numWorkers
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for w := 0; w < numWorkers; w++ {
		start := w * chunkSize
		end := start + chunkSize
		if w == numWorkers-1 {
			end = length // 最后一个 worker 处理剩余部分
		}

		go func(s, e int) {
			defer wg.Done()
			for i := s; i < e; i++ {
				result[i] = MapDirect[From, To](from[i])
			}
		}(start, end)
	}

	wg.Wait()
	return result
}

// MapDirectPtrSlice 指针切片映射
// 使用示例: dtos := mapper.MapDirectPtrSlice[User, UserDTO](&users)
func MapDirectPtrSlice[From, To any](from []*From) []*To {
	if from == nil {
		return nil
	}

	length := len(from)
	if length == 0 {
		return []*To{}
	}

	// 预获取映射关系以优化批量操作
	fromType := reflect.TypeOf((*From)(nil)).Elem()
	toType := reflect.TypeOf((*To)(nil)).Elem()

	_, hasCache := getFieldMappings(fromType, toType)
	if !hasCache {
		buildFieldMappings(fromType, toType)
	}

	// 大 slice 使用并行处理
	if length >= parallelThreshold {
		return mapDirectPtrSliceParallel[From, To](from)
	}

	result := make([]*To, length)
	for i, v := range from {
		if v != nil {
			t := MapDirect[From, To](*v)
			result[i] = &t
		}
	}
	return result
}

// mapDirectPtrSliceParallel 并行指针切片映射
func mapDirectPtrSliceParallel[From, To any](from []*From) []*To {
	length := len(from)
	result := make([]*To, length)

	numCPU := runtime.NumCPU()
	numWorkers := length / 100
	if numWorkers < 1 {
		numWorkers = 1
	}
	if numWorkers > numCPU {
		numWorkers = numCPU
	}

	chunkSize := length / numWorkers
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for w := 0; w < numWorkers; w++ {
		start := w * chunkSize
		end := start + chunkSize
		if w == numWorkers-1 {
			end = length
		}

		go func(s, e int) {
			defer wg.Done()
			for i := s; i < e; i++ {
				if from[i] != nil {
					t := MapDirect[From, To](*from[i])
					result[i] = &t
				}
			}
		}(start, end)
	}

	wg.Wait()
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

	length := len(from)
	if length == 0 {
		return []To{}, nil
	}

	// 大 slice 使用并行处理
	if length >= parallelThreshold {
		return safeMapDirectSliceParallel[From, To](from)
	}

	result := make([]To, length)
	for i, v := range from {
		r, err := SafeMapDirect[From, To](v)
		if err != nil {
			return nil, fmt.Errorf("map slice failed at index %d", i)
		}
		result[i] = r
	}
	return result, nil
}

// safeMapDirectSliceParallel 并行安全批量映射
func safeMapDirectSliceParallel[From, To any](from []From) ([]To, error) {
	length := len(from)
	result := make([]To, length)

	numCPU := runtime.NumCPU()
	numWorkers := length / 100
	if numWorkers < 1 {
		numWorkers = 1
	}
	if numWorkers > numCPU {
		numWorkers = numCPU
	}

	chunkSize := length / numWorkers
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	errChan := make(chan error, 1) // 用于接收第一个错误

	for w := 0; w < numWorkers; w++ {
		start := w * chunkSize
		end := start + chunkSize
		if w == numWorkers-1 {
			end = length
		}

		go func(s, e int) {
			defer wg.Done()
			for i := s; i < e; i++ {
				r, err := SafeMapDirect[From, To](from[i])
				if err != nil {
					select {
					case errChan <- fmt.Errorf("map slice failed at index %d", i):
					default:
					}
					return
				}
				result[i] = r
			}
		}(start, end)
	}

	wg.Wait()

	// 检查是否有错误
	select {
	case err := <-errChan:
		return nil, err
	default:
	}

	return result, nil
}

// ClearFieldMappingCache 清除字段映射缓存
// 用于在需要重新构建映射关系时调用
func ClearFieldMappingCache() {
	fieldMappingCache = &sync.Map{}
}
