package mapper

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func elemMapper(fromElem, toElem reflect.Value) error {
	//check register flag
	//if not register, register it
	if !checkIsRegister(fromElem) {
		registerValue(fromElem)
	}
	if !checkIsRegister(toElem) {
		registerValue(toElem)
	}
	if toElem.Type().Kind() == reflect.Map {
		elemToMap(fromElem, toElem)
	} else {
		elemToStruct(fromElem, toElem)
	}

	return nil
}

func elemToStruct(fromElem, toElem reflect.Value) {
	for i := 0; i < fromElem.NumField(); i++ {
		fromFieldInfo := fromElem.Field(i)
		fieldName := GetFieldName(fromElem, i)
		//check field is exists
		realFieldName, exists := CheckExistsField(toElem, fieldName)
		if !exists {
			continue
		}

		toFieldInfo := toElem.FieldByName(realFieldName)
		//check field is same type
		if enabledTypeChecking {
			if fromFieldInfo.Kind() != toFieldInfo.Kind() {
				continue
			}
		}

		if enabledMapperStructField &&
			toFieldInfo.Kind() == reflect.Struct && fromFieldInfo.Kind() == reflect.Struct &&
			toFieldInfo.Type() != fromFieldInfo.Type() &&
			!CheckIsTypeWrapper(toFieldInfo) && !CheckIsTypeWrapper(fromFieldInfo) {
			x := reflect.New(toFieldInfo.Type()).Elem()
			err := elemMapper(fromFieldInfo, x)
			if err != nil {
				fmt.Println("auto mapper field", fromFieldInfo, "=>", toFieldInfo, "error", err.Error())
			} else {
				toFieldInfo.Set(x)
			}
		} else {
			isSet := false
			if enabledAutoTypeConvert {
				if DefaultTimeWrapper.IsType(fromFieldInfo) && toFieldInfo.Kind() == reflect.Int64 {
					fromTime := fromFieldInfo.Interface().(time.Time)
					toFieldInfo.Set(reflect.ValueOf(TimeToUnix(fromTime)))
					isSet = true
				} else if DefaultTimeWrapper.IsType(toFieldInfo) && fromFieldInfo.Kind() == reflect.Int64 {
					fromTime := fromFieldInfo.Interface().(int64)
					toFieldInfo.Set(reflect.ValueOf(UnixToTime(fromTime)))
					isSet = true
				}
			}
			if !isSet {
				toFieldInfo.Set(fromFieldInfo)
			}
		}

	}
}

func elemToMap(fromElem, toElem reflect.Value) {
	for i := 0; i < fromElem.NumField(); i++ {
		fromFieldInfo := fromElem.Field(i)
		fieldName := GetFieldName(fromElem, i)
		toElem.SetMapIndex(reflect.ValueOf(fieldName), fromFieldInfo)
	}
}

func setFieldValue(fieldValue reflect.Value, fieldKind reflect.Kind, value interface{}) error {
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
		} else if DefaultTimeWrapper.IsType(fieldValue) {
			var timeString string
			if fieldValue.Type() == timeType {
				timeString = ""
				fieldValue.Set(reflect.ValueOf(value))
			}
			if fieldValue.Type() == jsonTimeType {
				timeString = ""
				fieldValue.Set(reflect.ValueOf(JSONTime(value.(time.Time))))
			}
			switch d := value.(type) {
			case []byte:
				timeString = string(d)
			case string:
				timeString = d
			case int64:
				if enabledAutoTypeConvert {
					//try to transform Unix time to local Time
					t, err := UnixToTimeLocation(value.(int64), time.UTC.String())
					if err != nil {
						return err
					}
					fieldValue.Set(reflect.ValueOf(t))
				}
			}
			if timeString != "" {
				if len(timeString) >= 19 {
					//满足yyyy-MM-dd HH:mm:ss格式
					timeString = timeString[:19]
					t, err := time.ParseInLocation(formatDateTime, timeString, time.UTC)
					if err == nil {
						t = t.In(time.UTC)
						fieldValue.Set(reflect.ValueOf(t))
					}
				} else if len(timeString) >= 10 {
					//满足yyyy-MM-dd格式
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

func getStructTag(field reflect.StructField) string {
	tagValue := ""
	//1.check mapperTagKey
	tagValue = field.Tag.Get(mapperTagKey)
	if checkTagValidity(tagValue) {
		return tagValue
	}

	//2.check jsonTagKey
	tagValue = field.Tag.Get(jsonTagKey)
	if checkTagValidity(tagValue) {
		// support more tag property, as json tag omitempty 2018-07-13
		return strings.Split(tagValue, ",")[0]
	}

	return ""
}

func checkTagValidity(tagValue string) bool {
	if tagValue != "" && tagValue != IgnoreTagValue {
		return true
	}
	return false
}

func checkIsRegister(objElem reflect.Value) bool {
	typeName := objElem.Type().String()
	_, isOk := registerMap.Load(typeName)
	return isOk
}

//convert slice interface{} to []interface{}
func convertToSlice(arr interface{}) []interface{} {
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
