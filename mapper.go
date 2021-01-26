package mapper

import (
	"encoding/json"
	"errors"
	"reflect"
	"sync"
	"time"
)

var (
	ZeroValue                reflect.Value
	DefaultTimeWrapper       *TimeWrapper
	fieldNameMap             sync.Map
	registerMap              sync.Map
	enabledTypeChecking      bool
	enabledMapperStructField bool
	enabledAutoTypeConvert   bool
	timeType                 = reflect.TypeOf(time.Now())
	jsonTimeType             = reflect.TypeOf(JSONTime(time.Now()))
	typeWrappers             []TypeWrapper
)

const (
	packageVersion = "0.7.5"
	mapperTagKey   = "mapper"
	jsonTagKey     = "json"
	IgnoreTagValue = "-"
	nameConnector  = "_"
	formatTime     = "15:04:05"
	formatDate     = "2006-01-02"
	formatDateTime = "2006-01-02 15:04:05"
)

func init() {
	ZeroValue = reflect.Value{}
	DefaultTimeWrapper = NewTimeWrapper()
	typeWrappers = []TypeWrapper{}
	UseWrapper(DefaultTimeWrapper)

	enabledTypeChecking = false
	enabledMapperStructField = true
	enabledAutoTypeConvert = true
}

func PackageVersion() string {
	return packageVersion
}

// UseWrapper register a type wrapper
func UseWrapper(w TypeWrapper) {
	if len(typeWrappers) > 0 {
		typeWrappers[len(typeWrappers)-1].SetNext(w)
	}
	typeWrappers = append(typeWrappers, w)
}

// CheckIsTypeWrapper check value is in type wrappers
func CheckIsTypeWrapper(value reflect.Value) bool {
	for _, w := range typeWrappers {
		if w.IsType(value) {
			return true
		}
	}
	return false
}

// SetEnabledTypeChecking set enabled flag for TypeChecking
// if set true, the field type will be checked for consistency during mapping
// default is false
func SetEnabledTypeChecking(isEnabled bool) {
	enabledTypeChecking = isEnabled
}

// SetEnabledAutoTypeConvert set enabled flag for auto type convert
// if set true, field will auto convert in Time and Unix
// default is true
func SetEnabledAutoTypeConvert(isEnabled bool) {
	enabledAutoTypeConvert = isEnabled
}

// SetEnabledMapperStructField set enabled flag for MapperStructField
// if set true, the reflect.Struct field will auto mapper
// must follow premises:
// 1. fromField and toField type must be reflect.Struct and not time.Time
// 2. fromField and toField must be not same type
// default is enabled
func SetEnabledMapperStructField(isEnabled bool) {
	enabledMapperStructField = isEnabled
}

// Register register struct to init Map
func Register(obj interface{}) error {
	objValue := reflect.ValueOf(obj)
	if objValue == ZeroValue {
		return errors.New("obj value does not exist")
	}
	return registerValue(objValue)
}

// registerValue register Value to init Map
func registerValue(objValue reflect.Value) error {
	regValue := objValue
	if objValue == ZeroValue {
		return errors.New("obj value does not exist")
	}

	if regValue.Type().Kind() == reflect.Ptr {
		regValue = regValue.Elem()
	}

	typeName := regValue.Type().String()
	if regValue.Type().Kind() == reflect.Struct {
		for i := 0; i < regValue.NumField(); i++ {
			mapFieldName := typeName + nameConnector + GetFieldName(regValue, i)
			realFieldName := regValue.Type().Field(i).Name
			fieldNameMap.Store(mapFieldName, realFieldName)
		}
	}

	//store register flag
	registerMap.Store(typeName, nil)
	return nil
}

// GetTypeName get type name
func GetTypeName(obj interface{}) string {
	object := reflect.ValueOf(obj)
	return object.String()
}

// CheckExistsField check field is exists by name
func CheckExistsField(elem reflect.Value, fieldName string) (realFieldName string, exists bool) {
	typeName := elem.Type().String()
	fileKey := typeName + nameConnector + fieldName
	realName, isOk := fieldNameMap.Load(fileKey)

	if !isOk {
		return "", isOk
	} else {
		return realName.(string), isOk
	}

}

// GetFieldName get fieldName with ElemValue and index
// if config tag string, return tag value
func GetFieldName(objElem reflect.Value, index int) string {
	fieldName := ""
	field := objElem.Type().Field(index)
	tag := getStructTag(field)
	if tag != "" {
		fieldName = tag
	} else {
		fieldName = field.Name
	}
	return fieldName
}

// MapperMap mapper and set value from map to object
// support auto register struct
// now support field type:
// 1.reflect.Bool
// 2.reflect.String
// 3.reflect.Int8\16\32\64
// 4.reflect.Uint8\16\32\64
// 5.reflect.Float32\64
// 6.time.Time
func MapperMap(fromMap map[string]interface{}, toObj interface{}) error {
	toElemType := reflect.ValueOf(toObj)
	toElem := toElemType
	if toElemType.Kind() == reflect.Ptr {
		toElem = toElemType.Elem()
	}

	if toElem == ZeroValue {
		return errors.New("to obj is not legal value")
	}
	//check register flag
	//if not register, register it
	if !checkIsRegister(toElem) {
		Register(toObj)
	}
	for k, v := range fromMap {
		fieldName := k
		//check field is exists
		realFieldName, exists := CheckExistsField(toElem, fieldName)
		if !exists {
			continue
		}
		fieldInfo, exists := toElem.Type().FieldByName(realFieldName)
		if !exists {
			continue
		}
		fieldKind := fieldInfo.Type.Kind()
		fieldValue := toElem.FieldByName(realFieldName)
		setFieldValue(fieldValue, fieldKind, v)
	}
	return nil
}

// MapToSlice mapper from map[string]interface{} to a slice of any type's ptr
// toSlice must be a slice of any type.
func MapToSlice(fromMap map[string]interface{}, toSlice interface{}) error {
	var err error
	toValue := reflect.ValueOf(toSlice)
	if toValue.Kind() != reflect.Ptr {
		return errors.New("toSlice must be a pointer to a slice")
	}
	if toValue.IsNil() {
		return errors.New("toSlice must not be a nil pointer")
	}

	toElemType := reflect.TypeOf(toSlice).Elem().Elem()
	realType := toElemType.Kind()
	direct := reflect.Indirect(toValue)
	if realType == reflect.Ptr {
		toElemType = toElemType.Elem()
	}
	for _, v := range fromMap {
		if reflect.TypeOf(v).Kind().String() == "map" {
			elem := reflect.New(toElemType)
			err = MapperMap(v.(map[string]interface{}), elem.Interface())
			if err == nil {
				if realType == reflect.Ptr {
					direct.Set(reflect.Append(direct, elem))
				} else {
					direct.Set(reflect.Append(direct, elem).Elem())
				}
			}
		} else {
			if realType == reflect.Ptr {
				direct.Set(reflect.Append(direct, reflect.ValueOf(v)))
			} else {
				direct.Set(reflect.Append(direct, reflect.ValueOf(v).Elem()))
			}
		}

	}
	return err
}

// MapperMapSlice mapper from map[string]map[string]interface{} to a slice of any type's ptr
// toSlice must be a slice of any type.
// Deprecated: will remove on v1.0, please use MapToSlice instead
func MapperMapSlice(fromMaps map[string]map[string]interface{}, toSlice interface{}) error {
	var err error
	toValue := reflect.ValueOf(toSlice)
	if toValue.Kind() != reflect.Ptr {
		return errors.New("toSlice must be a pointer to a slice")
	}
	if toValue.IsNil() {
		return errors.New("toSlice must not be a nil pointer")
	}

	toElemType := reflect.TypeOf(toSlice).Elem().Elem()
	realType := toElemType.Kind()
	direct := reflect.Indirect(toValue)
	//3 elem parse: 1.[]*type 2.*type 3.type
	if realType == reflect.Ptr {
		toElemType = toElemType.Elem()
	}
	for _, v := range fromMaps {
		elem := reflect.New(toElemType)
		err = MapperMap(v, elem.Interface())
		if err == nil {
			if realType == reflect.Ptr {
				direct.Set(reflect.Append(direct, elem))
			} else {
				direct.Set(reflect.Append(direct, elem.Elem()))
			}
		}
	}
	return err
}

// MapperSlice mapper from slice of struct to a slice of any type
// fromSlice and toSlice must be a slice of any type.
func MapperSlice(fromSlice, toSlice interface{}) error {
	var err error
	toValue := reflect.ValueOf(toSlice)
	if toValue.Kind() != reflect.Ptr {
		return errors.New("toSlice must be a pointer to a slice")
	}
	if toValue.IsNil() {
		return errors.New("toSlice must not be a nil pointer")
	}

	elemType := reflect.TypeOf(toSlice).Elem().Elem()
	realType := elemType.Kind()
	direct := reflect.Indirect(toValue)
	//3 elem parse: 1.[]*type 2.*type 3.type
	if realType == reflect.Ptr {
		elemType = elemType.Elem()
	}

	fromElems := convertToSlice(fromSlice)
	for _, v := range fromElems {
		elem := reflect.New(elemType).Elem()
		if realType == reflect.Ptr {
			elem = reflect.New(elemType)
		}
		if realType == reflect.Ptr {
			err = elemMapper(reflect.ValueOf(v).Elem(), elem.Elem())
		} else {
			err = elemMapper(reflect.ValueOf(v), elem)
		}
		if err == nil {
			direct.Set(reflect.Append(direct, elem))
		}
	}
	return err
}

// MapToJson mapper from map[string]interface{} to json []byte
func MapToJson(fromMap map[string]interface{}) ([]byte, error) {
	json, err := json.Marshal(fromMap)
	if err != nil {
		return nil, err
	}
	return json, nil
}

// JsonToMap mapper from json []byte to map[string]interface{}
func JsonToMap(body []byte, toMap *map[string]interface{}) error {
	err := json.Unmarshal(body, toMap)
	return err
}

// Mapper mapper and set value from struct fromObj to toObj
// not support auto register struct
func Mapper(fromObj, toObj interface{}) error {
	fromElem := reflect.ValueOf(fromObj).Elem()
	toElem := reflect.ValueOf(toObj).Elem()
	if fromElem == ZeroValue {
		return errors.New("from obj is not legal value")
	}
	if toElem == ZeroValue {
		return errors.New("to obj is not legal value")
	}
	return elemMapper(fromElem, toElem)
}

// Mapper mapper and set value from struct fromObj to toObj
// support auto register struct
func AutoMapper(fromObj, toObj interface{}) error {
	return Mapper(fromObj, toObj)
}
