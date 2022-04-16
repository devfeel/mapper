package mapper

import (
	"encoding/json"
	"errors"
	"reflect"
	"sync"
	"time"
)

type mapperObject struct {
	ZeroValue                reflect.Value
	DefaultTimeWrapper       *TimeWrapper
	typeWrappers             []TypeWrapper
	timeType                 reflect.Type
	jsonTimeType             reflect.Type
	fieldNameMap             sync.Map
	registerMap              sync.Map
	enabledTypeChecking      bool
	enabledMapperStructField bool
	enabledAutoTypeConvert   bool
	enabledMapperTag         bool
	enabledJsonTag           bool
}

func NewMapper() IMapper {
	dm := mapperObject{
		ZeroValue:                reflect.Value{},
		DefaultTimeWrapper:       NewTimeWrapper(),
		typeWrappers:             []TypeWrapper{},
		timeType:                 reflect.TypeOf(time.Now()),
		jsonTimeType:             reflect.TypeOf(JSONTime(time.Now())),
		enabledTypeChecking:      false,
		enabledMapperStructField: true,
		enabledAutoTypeConvert:   true,
		enabledMapperTag:         true,
		enabledJsonTag:           true,
	}
	dm.useWrapper(dm.DefaultTimeWrapper)
	return &dm
}

// Mapper map and set value from struct fromObj to toObj
// not support auto register struct
func (dm *mapperObject) Mapper(fromObj, toObj interface{}) error {
	fromElem := reflect.ValueOf(fromObj).Elem()
	toElem := reflect.ValueOf(toObj).Elem()
	if fromElem == dm.ZeroValue {
		return errors.New("from obj is not legal value")
	}
	if toElem == dm.ZeroValue {
		return errors.New("to obj is not legal value")
	}
	return dm.elemMapper(fromElem, toElem)
}

// AutoMapper mapper and set value from struct fromObj to toObj
// support auto register struct
func (dm *mapperObject) AutoMapper(fromObj, toObj interface{}) error {
	return dm.Mapper(fromObj, toObj)
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
func (dm *mapperObject) MapperMap(fromMap map[string]interface{}, toObj interface{}) error {
	toElemType := reflect.ValueOf(toObj)
	toElem := toElemType
	if toElemType.Kind() == reflect.Ptr {
		toElem = toElemType.Elem()
	}

	if toElem == dm.ZeroValue {
		return errors.New("to obj is not legal value")
	}
	// check register flag
	// if not register, register it
	if !dm.checkIsRegister(toElem) {
		if err := dm.Register(toObj); err != nil {
			return err
		}
	}
	for k, v := range fromMap {
		fieldName := k
		// check field is exists
		realFieldName, exists := dm.CheckExistsField(toElem, fieldName)
		if !exists {
			continue
		}
		fieldInfo, exists := toElem.Type().FieldByName(realFieldName)
		if !exists {
			continue
		}

		fieldKind := fieldInfo.Type.Kind()
		fieldValue := toElem.FieldByName(realFieldName)

		if err := dm.setFieldValue(fieldValue, fieldKind, v); err != nil {
			return err
		}
	}
	return nil
}

// SetEnabledTypeChecking set enabled flag for TypeChecking
// if set true, the field type will be checked for consistency during mapping
// default is false
func (dm *mapperObject) SetEnabledTypeChecking(isEnabled bool) {
	dm.enabledTypeChecking = isEnabled
}

func (dm *mapperObject) IsEnabledTypeChecking() bool {
	return dm.enabledTypeChecking
}

// SetEnabledMapperTag set enabled flag for 'Mapper' tag check
// if set true, 'Mapper' tag will be check during mapping's GetFieldName
// default is true
func (dm *mapperObject) SetEnabledMapperTag(isEnabled bool) {
	dm.enabledMapperTag = isEnabled
}

func (dm *mapperObject) IsEnabledMapperTag() bool {
	return dm.enabledMapperTag
}

// SetEnabledJsonTag set enabled flag for 'Json' tag check
// if set true, 'Json' tag will be check during mapping's GetFieldName
// default is true
func (dm *mapperObject) SetEnabledJsonTag(isEnabled bool) {
	dm.enabledJsonTag = isEnabled
}

func (dm *mapperObject) IsEnabledJsonTag() bool {
	return dm.enabledJsonTag
}

// SetEnabledAutoTypeConvert set enabled flag for auto type convert
// if set true, field will auto convert in Time and Unix
// default is true
func (dm *mapperObject) SetEnabledAutoTypeConvert(isEnabled bool) {
	dm.enabledAutoTypeConvert = isEnabled
}

func (dm *mapperObject) IsEnabledAutoTypeConvert() bool {
	return dm.enabledAutoTypeConvert
}

// SetEnabledMapperStructField set enabled flag for MapperStructField
// if set true, the reflect.Struct field will auto mapper
// must follow premises:
// 1. fromField and toField type must be reflect.Struct and not time.Time
// 2. fromField and toField must be not same type
// default is enabled
func (dm *mapperObject) SetEnabledMapperStructField(isEnabled bool) {
	dm.enabledMapperStructField = isEnabled
}

func (dm *mapperObject) IsEnabledMapperStructField() bool {
	return dm.enabledMapperStructField
}

// GetTypeName get type name
func (dm *mapperObject) GetTypeName(obj interface{}) string {
	object := reflect.ValueOf(obj)
	return object.String()
}

// GetFieldName get fieldName with ElemValue and index
// if config tag string, return tag value
func (dm *mapperObject) GetFieldName(objElem reflect.Value, index int) string {
	fieldName := ""
	field := objElem.Type().Field(index)
	tag := dm.getStructTag(field)
	if tag != "" {
		fieldName = tag
	} else {
		fieldName = field.Name
	}
	return fieldName
}

func (dm *mapperObject) GetDefaultTimeWrapper() *TimeWrapper {
	return dm.DefaultTimeWrapper
}

// Register register struct to init Map
func (dm *mapperObject) Register(obj interface{}) error {
	objValue := reflect.ValueOf(obj)
	if objValue == dm.ZeroValue {
		return errors.New("obj value does not exist")
	}
	return dm.registerValue(objValue)
}

// MapToSlice mapper from map[string]interface{} to a slice of any type's ptr
// toSlice must be a slice of any type.
func (dm *mapperObject) MapToSlice(fromMap map[string]interface{}, toSlice interface{}) error {
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
			err = dm.MapperMap(v.(map[string]interface{}), elem.Interface())
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
func (dm *mapperObject) MapperMapSlice(fromMaps map[string]map[string]interface{}, toSlice interface{}) error {
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
	// 3 elem parse: 1.[]*type 2.*type 3.type
	if realType == reflect.Ptr {
		toElemType = toElemType.Elem()
	}
	for _, v := range fromMaps {
		elem := reflect.New(toElemType)
		err = dm.MapperMap(v, elem.Interface())
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
func (dm *mapperObject) MapperSlice(fromSlice, toSlice interface{}) error {
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
	// 3 elem parse: 1.[]*type 2.*type 3.type
	if realType == reflect.Ptr {
		elemType = elemType.Elem()
	}

	fromElems := dm.convertToSlice(fromSlice)
	for _, v := range fromElems {
		elem := reflect.New(elemType).Elem()
		if realType == reflect.Ptr {
			elem = reflect.New(elemType)
		}
		if realType == reflect.Ptr {
			err = dm.elemMapper(reflect.ValueOf(v).Elem(), elem.Elem())
		} else {
			err = dm.elemMapper(reflect.ValueOf(v), elem)
		}
		if err == nil {
			direct.Set(reflect.Append(direct, elem))
		}
	}
	return err
}

// MapToJson mapper from map[string]interface{} to json []byte
func (dm *mapperObject) MapToJson(fromMap map[string]interface{}) ([]byte, error) {
	jsonStr, err := json.Marshal(fromMap)
	if err != nil {
		return nil, err
	}
	return jsonStr, nil
}

// JsonToMap mapper from json []byte to map[string]interface{}
func (dm *mapperObject) JsonToMap(body []byte, toMap *map[string]interface{}) error {
	err := json.Unmarshal(body, toMap)
	return err
}

// CheckExistsField check field is exists by name
func (dm *mapperObject) CheckExistsField(elem reflect.Value, fieldName string) (realFieldName string, exists bool) {
	typeName := elem.Type().String()
	fileKey := typeName + nameConnector + fieldName
	realName, isOk := dm.fieldNameMap.Load(fileKey)

	if !isOk {
		return "", isOk
	} else {
		return realName.(string), isOk
	}
}

// CheckIsTypeWrapper check value is in type wrappers
func (dm *mapperObject) CheckIsTypeWrapper(value reflect.Value) bool {
	for _, w := range dm.typeWrappers {
		if w.IsType(value) {
			return true
		}
	}
	return false
}
