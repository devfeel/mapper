package mapper

import (
	"reflect"
)

var (
	standardMapper IMapper
)

type IMapper interface {
	Mapper(fromObj, toObj interface{}) error
	AutoMapper(fromObj, toObj interface{}) error
	MapperMap(fromMap map[string]interface{}, toObj interface{}) error
	MapToSlice(fromMap map[string]interface{}, toSlice interface{}) error
	MapperMapSlice(fromMaps map[string]map[string]interface{}, toSlice interface{}) error
	MapperSlice(fromSlice, toSlice interface{}) error
	MapToJson(fromMap map[string]interface{}) ([]byte, error)
	JsonToMap(body []byte, toMap *map[string]interface{}) error

	Register(obj interface{}) error
	UseWrapper(w TypeWrapper)

	GetTypeName(obj interface{}) string
	GetFieldName(objElem reflect.Value, index int) string
	GetDefaultTimeWrapper() *TimeWrapper

	CheckExistsField(elem reflect.Value, fieldName string) (realFieldName string, exists bool)
	CheckIsTypeWrapper(value reflect.Value) bool

	SetEnabledTypeChecking(isEnabled bool)
	IsEnabledTypeChecking() bool

	SetEnabledMapperTag(isEnabled bool)
	IsEnabledMapperTag() bool

	SetEnabledJsonTag(isEnabled bool)
	IsEnabledJsonTag() bool

	SetEnabledAutoTypeConvert(isEnabled bool)
	IsEnabledAutoTypeConvert() bool

	SetEnabledMapperStructField(isEnabled bool)
	IsEnabledMapperStructField() bool
}

func init() {
	standardMapper = NewMapper()
}

func PackageVersion() string {
	return packageVersion
}

// UseWrapper register a type wrapper
func UseWrapper(w TypeWrapper) {
	standardMapper.UseWrapper(w)
}

// CheckIsTypeWrapper check value is in type wrappers
func CheckIsTypeWrapper(value reflect.Value) bool {
	return standardMapper.CheckIsTypeWrapper(value)
}

// SetEnabledTypeChecking set enabled flag for TypeChecking
// if set true, the field type will be checked for consistency during mapping
// default is false
func SetEnabledTypeChecking(isEnabled bool) {
	standardMapper.SetEnabledTypeChecking(isEnabled)
}

// SetEnabledMapperTag set enabled flag for 'Mapper' tag check
// if set true, 'Mapper' tag will be check during mapping's GetFieldName
// default is true
func SetEnabledMapperTag(isEnabled bool) {
	standardMapper.SetEnabledMapperTag(isEnabled)
}

// SetEnabledJsonTag set enabled flag for 'Json' tag check
// if set true, 'Json' tag will be check during mapping's GetFieldName
// default is true
func SetEnabledJsonTag(isEnabled bool) {
	standardMapper.SetEnabledJsonTag(isEnabled)
}

// SetEnabledAutoTypeConvert set enabled flag for auto type convert
// if set true, field will auto convert in Time and Unix
// default is true
func SetEnabledAutoTypeConvert(isEnabled bool) {
	standardMapper.SetEnabledAutoTypeConvert(isEnabled)
}

// SetEnabledMapperStructField set enabled flag for MapperStructField
// if set true, the reflect.Struct field will auto mapper
// must follow premises:
// 1. fromField and toField type must be reflect.Struct and not time.Time
// 2. fromField and toField must be not same type
// default is enabled
func SetEnabledMapperStructField(isEnabled bool) {
	standardMapper.SetEnabledMapperStructField(isEnabled)
}

// Register register struct to init Map
func Register(obj interface{}) error {
	return standardMapper.Register(obj)
}

// GetTypeName get type name
func GetTypeName(obj interface{}) string {
	object := reflect.ValueOf(obj)
	return object.String()
}

// CheckExistsField check field is exists by name
func CheckExistsField(elem reflect.Value, fieldName string) (realFieldName string, exists bool) {
	return standardMapper.CheckExistsField(elem, fieldName)
}

// GetFieldName get fieldName with ElemValue and index
// if config tag string, return tag value
func GetFieldName(objElem reflect.Value, index int) string {
	return standardMapper.GetFieldName(objElem, index)
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
	return standardMapper.MapperMap(fromMap, toObj)
}

// MapToSlice mapper from map[string]interface{} to a slice of any type's ptr
// toSlice must be a slice of any type.
func MapToSlice(fromMap map[string]interface{}, toSlice interface{}) error {
	return standardMapper.MapToSlice(fromMap, toSlice)
}

// MapperMapSlice mapper from map[string]map[string]interface{} to a slice of any type's ptr
// toSlice must be a slice of any type.
// Deprecated: will remove on v1.0, please use MapToSlice instead
func MapperMapSlice(fromMaps map[string]map[string]interface{}, toSlice interface{}) error {
	return standardMapper.MapperMapSlice(fromMaps, toSlice)
}

// MapperSlice mapper from slice of struct to a slice of any type
// fromSlice and toSlice must be a slice of any type.
func MapperSlice(fromSlice, toSlice interface{}) error {
	return standardMapper.MapperSlice(fromSlice, toSlice)
}

// MapToJson mapper from map[string]interface{} to json []byte
func MapToJson(fromMap map[string]interface{}) ([]byte, error) {
	return standardMapper.MapToJson(fromMap)
}

// JsonToMap mapper from json []byte to map[string]interface{}
func JsonToMap(body []byte, toMap *map[string]interface{}) error {
	return standardMapper.JsonToMap(body, toMap)
}

// Mapper mapper and set value from struct fromObj to toObj
// not support auto register struct
func Mapper(fromObj, toObj interface{}) error {
	return standardMapper.Mapper(fromObj, toObj)
}

// AutoMapper mapper and set value from struct fromObj to toObj
// support auto register struct
func AutoMapper(fromObj, toObj interface{}) error {
	return standardMapper.AutoMapper(fromObj, toObj)
}
