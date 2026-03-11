package mapper

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
)

// registerValue register Value to init Map
func (dm *mapperObject) registerValue(objValue reflect.Value) error {
	regValue := objValue
	if objValue == dm.ZeroValue {
		return errors.New("obj value does not exist")
	}

	if regValue.Type().Kind() == reflect.Ptr {
		regValue = regValue.Elem()
	}

	typeName := regValue.Type().String()
	if regValue.Type().Kind() == reflect.Struct {
		for i := 0; i < regValue.NumField(); i++ {
			fieldName := dm.getFieldName(regValue, i)
			if fieldName == IgnoreTagValue {
				continue
			}
			mapFieldName := typeName + nameConnector + fieldName
			realFieldName := regValue.Type().Field(i).Name
			dm.fieldNameMap.Store(mapFieldName, realFieldName)
		}
	}

	// store register flag
	dm.registerMap.Store(typeName, nil)
	return nil
}

// cleanRegisterValue clean all register Value & field name map
func (dm *mapperObject) cleanRegisterValue() {
	dm.registerMap = *new(sync.Map)
	dm.fieldNameMap = *new(sync.Map)
}

// GetFieldName get fieldName with ElemValue and index
// if config tag string, return tag value
func (dm *mapperObject) getFieldName(objElem reflect.Value, index int) string {
	fieldName := ""
	field := objElem.Type().Field(index)
	tag := dm.getStructTag(field)
	// keeps the behavior in old version
	if tag == IgnoreTagValue && !dm.IsEnableFieldIgnoreTag() {
		tag = ""
	}

	if tag != "" {
		fieldName = tag
	} else {
		fieldName = field.Name
	}
	return fieldName
}

// UseWrapper register a type wrapper
func (dm *mapperObject) useWrapper(w TypeWrapper) {
	if len(dm.typeWrappers) > 0 {
		dm.typeWrappers[len(dm.typeWrappers)-1].SetNext(w)
	}
	dm.typeWrappers = append(dm.typeWrappers, w)
}

func (dm *mapperObject) elemMapper(fromElem, toElem reflect.Value) error {
	// check register flag
	// if not register, register it
	if !dm.checkIsRegister(fromElem) {
		if err := dm.registerValue(fromElem); err != nil {
			return err
		}
	}
	if !dm.checkIsRegister(toElem) {
		if err := dm.registerValue(toElem); err != nil {
			return err
		}
	}

	// Compute and cache field mapping for struct-to-struct
	if fromElem.Type().Kind() != reflect.Map && toElem.Type().Kind() != reflect.Map {
		dm.computeFieldMappingCache(fromElem.Type(), toElem.Type())
	}

	if toElem.Type().Kind() == reflect.Map {
		dm.elemToMap(fromElem, toElem)
	} else {
		dm.elemToStruct(fromElem, toElem)
	}

	return nil
}

// computeFieldMappingCache computes and caches field mapping between src and dest types
func (dm *mapperObject) computeFieldMappingCache(srcType, destType reflect.Type) {
	key := srcType.String() + "->" + destType.String()

	// Check if already cached
	if _, ok := dm.fieldCache.Load(key); ok {
		return
	}

	// Compute mapping: store as []int pairs [srcIdx, destIdx, ...]
	mapping := make([]int, 0, 16)

	for srcIdx := 0; srcIdx < srcType.NumField(); srcIdx++ {
		fieldName := dm.getFieldNameByType(srcType, srcIdx)

		if fieldName == IgnoreTagValue || fieldName == CompositeFieldTagValue {
			continue
		}

		// Look up dest field by name
		realFieldName, exists := dm.CheckExistsFieldByType(destType, fieldName)
		if !exists {
			continue
		}

		destField, ok := destType.FieldByName(realFieldName)
		if !ok {
			continue
		}

		// Store as pair
		mapping = append(mapping, srcIdx, destField.Index[0])
	}

	if len(mapping) > 0 {
		dm.fieldCache.Store(key, mapping)
	}
}

// getFieldMapping retrieves cached field mapping
func (dm *mapperObject) getFieldMapping(srcType, destType reflect.Type) ([]int, bool) {
	key := srcType.String() + "->" + destType.String()
	val, ok := dm.fieldCache.Load(key)
	if !ok {
		return nil, false
	}
	return val.([]int), true
}

// getFieldNameByType gets field name by type and index
func (dm *mapperObject) getFieldNameByType(objType reflect.Type, index int) string {
	field := objType.Field(index)
	tag := dm.getStructTag(field)
	if tag == IgnoreTagValue && !dm.IsEnableFieldIgnoreTag() {
		tag = ""
	}
	if tag != "" {
		return tag
	}
	return field.Name
}

// CheckExistsFieldByType checks if field exists by type (without value)
func (dm *mapperObject) CheckExistsFieldByType(elemType reflect.Type, fieldName string) (string, bool) {
	typeName := elemType.String()
	fileKey := typeName + nameConnector + fieldName
	realName, isOk := dm.fieldNameMap.Load(fileKey)
	if !isOk {
		return "", isOk
	}
	return realName.(string), isOk
}

// elemToStruct to convert the different structs for assignment.
// Uses cached field indices for better performance while keeping original logic
func (dm *mapperObject) elemToStruct(fromElem, toElem reflect.Value) {
	// Try to use cached mapping for non-nested fields
	mapping, hasCache := dm.getFieldMapping(fromElem.Type(), toElem.Type())

	if hasCache && len(mapping) > 0 {
		// Use cached indices for faster access
		for i := 0; i < len(mapping); i += 2 {
			srcIdx := mapping[i]
			destIdx := mapping[i+1]
			fromFieldInfo := fromElem.Field(srcIdx)
			toFieldInfo := toElem.Field(destIdx)
			dm.copyFieldValue(fromFieldInfo, toFieldInfo)
		}

		// Still handle nested struct fields with original logic
		for i := 0; i < fromElem.NumField(); i++ {
			fromFieldInfo := fromElem.Field(i)
			fieldName := dm.getFieldName(fromElem, i)
			if fieldName == CompositeFieldTagValue && fromFieldInfo.Kind() == reflect.Struct {
				fromNested := fromFieldInfo
				for j := 0; j < fromNested.NumField(); j++ {
					nestedFieldInfo := fromNested.Field(j)
					nestedFieldName := dm.getFieldName(fromNested, j)
					if nestedFieldName == IgnoreTagValue {
						continue
					}
					err := dm.convertstructfieldInternal(nestedFieldName, nestedFieldInfo, toElem)
					if err != nil {
						fmt.Println("auto mapper failed", nestedFieldInfo, "error", err.Error())
					}
				}
			}
		}
	} else {
		// Fallback to original logic
		dm.elemToStructOriginal(fromElem, toElem)
	}
}

// elemToStructOriginal keeps the original implementation for fallback
func (dm *mapperObject) elemToStructOriginal(fromElem, toElem reflect.Value) {
	for i := 0; i < fromElem.NumField(); i++ {
		fromFieldInfo := fromElem.Field(i)
		fieldName := dm.getFieldName(fromElem, i)
		if fieldName == IgnoreTagValue {
			continue
		}
		if fieldName == CompositeFieldTagValue &&
			fromFieldInfo.Kind() == reflect.Struct {
			//If composite fields are identified, further decomposition and judgment will be taken.
			fromNested := fromFieldInfo
			for i := 0; i < fromNested.NumField(); i++ {
				fromFieldInfo := fromNested.Field(i)
				fieldName := dm.getFieldName(fromNested, i)
				if fieldName == IgnoreTagValue {
					continue
				}
				err := dm.convertstructfieldInternal(fieldName, fromFieldInfo, toElem)
				if err != nil {
					fmt.Println("auto mapper failed", fromFieldInfo, "error", err.Error())
				}
			}
		} else {
			err := dm.convertstructfieldInternal(fieldName, fromFieldInfo, toElem)
			if err != nil {
				fmt.Println("auto mapper failed", fromFieldInfo, "error", err.Error())
			}
		}
	}
}

// copyFieldValue copies field value with type conversion support
func (dm *mapperObject) copyFieldValue(fromField, toField reflect.Value) {
	// Check type compatibility
	if dm.setting.EnabledTypeChecking {
		if fromField.Kind() != toField.Kind() {
			return
		}
	}

	// Handle nested struct
	if dm.setting.EnabledMapperStructField &&
		toField.Kind() == reflect.Struct && fromField.Kind() == reflect.Struct &&
		toField.Type() != fromField.Type() &&
		!dm.checkIsTypeWrapper(toField) && !dm.checkIsTypeWrapper(fromField) {
		x := reflect.New(toField.Type()).Elem()
		dm.elemMapper(fromField, x)
		toField.Set(x)
		return
	}

	// Handle auto type conversion (Time <-> int64)
	if dm.setting.EnabledAutoTypeConvert {
		if dm.DefaultTimeWrapper.IsType(fromField) && toField.Kind() == reflect.Int64 {
			fromTime := fromField.Interface().(time.Time)
			toField.Set(reflect.ValueOf(TimeToUnix(fromTime)))
			return
		}
		if dm.DefaultTimeWrapper.IsType(toField) && fromField.Kind() == reflect.Int64 {
			fromInt := fromField.Interface().(int64)
			toField.Set(reflect.ValueOf(UnixToTime(fromInt)))
			return
		}
	}

	// Direct copy
	toField.Set(fromField)
}

// convertstructfieldInternal to convert the fields of different structs for assignment.
func (dm *mapperObject) convertstructfieldInternal(fieldName string, fromFieldInfo, toElem reflect.Value) error {
	// check field is exists
	realFieldName, exists := dm.CheckExistsField(toElem, fieldName)
	if !exists {
		return nil
	}

	toFieldInfo := toElem.FieldByName(realFieldName)
	// check field is same type
	if dm.setting.EnabledTypeChecking {
		if fromFieldInfo.Kind() != toFieldInfo.Kind() {
			return nil
		}
	}

	if dm.setting.EnabledMapperStructField &&
		toFieldInfo.Kind() == reflect.Struct && fromFieldInfo.Kind() == reflect.Struct &&
		toFieldInfo.Type() != fromFieldInfo.Type() &&
		!dm.checkIsTypeWrapper(toFieldInfo) && !dm.checkIsTypeWrapper(fromFieldInfo) {
		x := reflect.New(toFieldInfo.Type()).Elem()
		err := dm.elemMapper(fromFieldInfo, x)
		if err != nil {
			return err
		} else {
			toFieldInfo.Set(x)
		}
	} else {
		isSet := false
		if dm.setting.EnabledAutoTypeConvert {
			if dm.DefaultTimeWrapper.IsType(fromFieldInfo) && toFieldInfo.Kind() == reflect.Int64 {
				fromTime := fromFieldInfo.Interface().(time.Time)
				toFieldInfo.Set(reflect.ValueOf(TimeToUnix(fromTime)))
				isSet = true
			} else if dm.DefaultTimeWrapper.IsType(toFieldInfo) && fromFieldInfo.Kind() == reflect.Int64 {
				fromTime := fromFieldInfo.Interface().(int64)
				toFieldInfo.Set(reflect.ValueOf(UnixToTime(fromTime)))
				isSet = true
			}
		}
		if !isSet {
			toFieldInfo.Set(fromFieldInfo)
		}
	}
	return nil
}

func (dm *mapperObject) elemToMap(fromElem, toElem reflect.Value) {
	for i := 0; i < fromElem.NumField(); i++ {
		fromFieldInfo := fromElem.Field(i)
		fieldName := dm.getFieldName(fromElem, i)
		if fieldName == IgnoreTagValue {
			continue
		}
		toElem.SetMapIndex(reflect.ValueOf(fieldName), fromFieldInfo)
	}
}

func (dm *mapperObject) setFieldValue(fieldValue reflect.Value, fieldKind reflect.Kind, value interface{}) error {
	switch fieldKind {
	case reflect.Bool:
		if value == nil {
			fieldValue.SetBool(false)
		} else if v, ok := value.(bool); ok {
			fieldValue.SetBool(v)
		} else {
			v, _ := Convert(ToString(value)).Bool()
			fieldValue.SetBool(v)
		}

	case reflect.String:
		if value == nil {
			fieldValue.SetString("")
		} else {
			fieldValue.SetString(ToString(value))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value == nil {
			fieldValue.SetInt(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fieldValue.SetInt(val.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fieldValue.SetInt(int64(val.Uint()))
			default:
				v, _ := Convert(ToString(value)).Int64()
				fieldValue.SetInt(v)
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value == nil {
			fieldValue.SetUint(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fieldValue.SetUint(uint64(val.Int()))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fieldValue.SetUint(val.Uint())
			default:
				v, _ := Convert(ToString(value)).Uint64()
				fieldValue.SetUint(v)
			}
		}
	case reflect.Float64, reflect.Float32:
		if value == nil {
			fieldValue.SetFloat(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Float64:
				fieldValue.SetFloat(val.Float())
			default:
				v, _ := Convert(ToString(value)).Float64()
				fieldValue.SetFloat(v)
			}
		}
	case reflect.Struct:
		if value == nil {
			fieldValue.Set(reflect.Zero(fieldValue.Type()))
		} else if dm.DefaultTimeWrapper.IsType(fieldValue) {
			var timeString string
			if fieldValue.Type() == dm.timeType {
				timeString = ""
				fieldValue.Set(reflect.ValueOf(value))
			}
			if fieldValue.Type() == dm.jsonTimeType {
				timeString = ""
				fieldValue.Set(reflect.ValueOf(JSONTime(value.(time.Time))))
			}
			switch d := value.(type) {
			case []byte:
				timeString = string(d)
			case string:
				timeString = d
			case int64:
				if dm.setting.EnabledAutoTypeConvert {
					// try to transform Unix time to local Time
					t, err := UnixToTimeLocation(value.(int64), time.UTC.String())
					if err != nil {
						return err
					}
					fieldValue.Set(reflect.ValueOf(t))
				}
			}
			if timeString != "" {
				if len(timeString) >= 19 {
					// 满足yyyy-MM-dd HH:mm:ss格式
					timeString = timeString[:19]
					t, err := time.ParseInLocation(formatDateTime, timeString, time.UTC)
					if err == nil {
						t = t.In(time.UTC)
						fieldValue.Set(reflect.ValueOf(t))
					}
				} else if len(timeString) >= 10 {
					// 满足yyyy-MM-dd格式
					timeString = timeString[:10]
					t, err := time.ParseInLocation(formatDate, timeString, time.UTC)
					if err == nil {
						fieldValue.Set(reflect.ValueOf(t))
					}
				}
			}
		}
	default:
		if reflect.ValueOf(value).Type() == fieldValue.Type() {
			fieldValue.Set(reflect.ValueOf(value))
		}
	}

	return nil
}

func (dm *mapperObject) getStructTag(field reflect.StructField) string {
	tagValue := ""
	// 1.check mapperTagKey
	if dm.setting.EnabledMapperTag {
		tagValue = field.Tag.Get(mapperTagKey)
		if tagValue != "" {
			return tagValue
		}
	}

	// 2.check jsonTagKey
	if dm.setting.EnabledJsonTag {
		tagValue = field.Tag.Get(jsonTagKey)
		if tagValue != "" {
			// support more tag property, as json tag omitempty 2018-07-13
			return strings.Split(tagValue, ",")[0]
		}
	}

	// 3.che
	// ck customTag
	if dm.setting.EnabledCustomTag {
		tagValue = field.Tag.Get(dm.setting.CustomTagName)
		if tagValue != "" {
			return tagValue
		}
	}

	return tagValue
}

func (dm *mapperObject) checkIsRegister(objElem reflect.Value) bool {
	typeName := objElem.Type().String()
	_, isOk := dm.registerMap.Load(typeName)
	return isOk
}

// convertToSlice convert slice interface{} to []interface{}
func (dm *mapperObject) convertToSlice(arr interface{}) []interface{} {
	v := reflect.ValueOf(arr)
	if v.Kind() == reflect.Ptr {
		if v.Elem().Kind() != reflect.Slice {
			panic("fromSlice arr is not a pointer to a slice")
		}
		v = v.Elem()
	} else {
		if v.Kind() != reflect.Slice {
			panic("fromSlice arr is not a slice")
		}
	}
	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret
}

// checkIsTypeWrapper check value is in type wrappers
func (dm *mapperObject) checkIsTypeWrapper(value reflect.Value) bool {
	for _, w := range dm.typeWrappers {
		if w.IsType(value) {
			return true
		}
	}
	return false
}
