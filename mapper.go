package mapper

import (
	"errors"
	"reflect"
	"sync"
)

var (
	ZeroValue reflect.Value
	fieldNameMap sync.Map
	registerMap sync.Map
)

const (
	PackageVersion = "0.2"
	mapperTagKey = "mapper"
	jsonTagKey = "json"
	ignoreTagValue = "-"
	nameConnector = "_"
)


func init() {
	ZeroValue = reflect.Value{}
}

//Register register struct to init Map
func Register(obj interface{}) error{
	objValue := reflect.ValueOf(obj)
	if objValue == ZeroValue{
		return errors.New("no exists this value")
	}

	typeName:= objValue.Elem().Type().String()
	objElem := objValue.Elem()
	for i := 0; i < objElem.NumField(); i++ {
		mapFieldName := typeName + nameConnector + GetFieldName(objElem, i)
		realFieldName := objElem.Type().Field(i).Name
		fieldNameMap.Store(mapFieldName, realFieldName)
	}

	//store register flag
	registerMap.Store(typeName, nil)

	return nil
}

//GetTypeName get type name
func GetTypeName(obj interface{}) string {
	object := reflect.ValueOf(obj)
	return object.String()
}

//CheckExistsField check field is exists by name
func CheckExistsField(elem reflect.Value, fieldName string) (realFieldName string, exists bool) {
	typeName := elem.Type().String()
	fileKey := typeName + nameConnector + fieldName
	realName, isOk := fieldNameMap.Load(fileKey)

	if !isOk{
		return "", isOk
	}else{
		return realName.(string), isOk
	}

}

//GetFieldName get fieldName by index
//if config tag string, return tag value
func GetFieldName(elem reflect.Value, index int) string {
	fieldName := ""
	field := elem.Type().Field(index)
	tag := getStructTag(field)
	if tag != ""{
		fieldName = tag
	}else{
		fieldName = field.Name
	}
	return fieldName
}

//Mapper mapper and set value from struct fromObj to toObj
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

//Mapper mapper and set value from struct fromObj to toObj
func AutoMapper(fromObj, toObj interface{}) error {
	fromElem := reflect.ValueOf(fromObj).Elem()
	toElem := reflect.ValueOf(toObj).Elem()
	if fromElem == ZeroValue {
		return errors.New("from obj is not legal value")
	}
	if toElem == ZeroValue {
		return errors.New("to obj is not legal value")
	}
	//check register flag
	//if not register, register it
	if !checkIsRegister(fromElem) {
		Register(fromObj)
	}
	if !checkIsRegister(toElem) {
		Register(toObj)
	}

	return elemMapper(fromElem, toElem)
}

func elemMapper(fromElem, toElem reflect.Value) error{
	for i := 0; i < fromElem.NumField(); i++ {
		fieldInfo:=fromElem.Field(i)
		fieldName := GetFieldName(fromElem, i)
		//check field is exists
		realFieldName, exists := CheckExistsField(toElem, fieldName)
		if !exists {
			continue
		}
		//TODO:check field is same type
		toElem.FieldByName(realFieldName).Set(fieldInfo)
	}
	return nil
}

func getStructTag(field reflect.StructField) string {
	tagValue := ""
	//1.check mapperTagKey
	tagValue =field.Tag.Get(mapperTagKey)
	if checkTagValidity(tagValue){
		return tagValue
	}

	//2.check jsonTagKey and ignore "-" value
	tagValue =field.Tag.Get(jsonTagKey)
	if checkTagValidity(tagValue){
		return tagValue
	}

	return ""
}

func checkTagValidity(tagValue string) bool{
	if tagValue != "" && tagValue != ignoreTagValue{
		return true
	}
	return false
}

func checkIsRegister(objElem reflect.Value) bool{
	typeName:= objElem.Type().String()
	_, isOk := registerMap.Load(typeName)
	return isOk
}