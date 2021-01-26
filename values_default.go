package mapper

import (
	"log"
	"reflect"
	"sync"
	"time"
)

type DefaultInitValue interface {
	Default() reflect.Value
}

type DefaultValueFunc func(str string) reflect.Value

var regMap map[reflect.Type]DefaultValueFunc
var mux sync.RWMutex

func RegisterDefaultValue(t reflect.Type, f DefaultValueFunc) {
	mux.Lock()
	defer mux.Unlock()
	regMap[t] = f
}

func init() {
	regMap = make(map[reflect.Type]DefaultValueFunc)
	RegisterDefaultValue(reflect.TypeOf(time.Time{}), dateTimeFunc)
}
func getFunc(k reflect.Type) DefaultValueFunc {
	mux.RLock()
	defer mux.RUnlock()
	return regMap[k]
}
func bindValue(v reflect.Value, tag string) {
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		log.Println(t)
		t = v.Elem().Type()
	}

	if implDefaultInitValue(t) {
		tv := v.Interface().(DefaultInitValue).Default()
		if tv.Kind() == reflect.Struct && tv.Type() == t {
			if v.Type().Kind() == reflect.Ptr {
				v.Elem().Set(tv)
			} else {
				v.Set(tv)
			}
			return
		}
	}
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			if v.Type().Kind() == reflect.Ptr {
				if !checkSampleValue(v.Elem().Field(i), v.Elem().Type().Field(i), tag) {
					bindValue(v.Elem().Field(i), tag)
				}
			} else {
				if !checkSampleValue(v.Field(i), v.Type().Field(i), tag) {
					bindValue(v.Field(i), tag)
				}
			}
		}
	}
}

func checkSampleValue(v reflect.Value, field reflect.StructField, tag string) (ok bool) {
	if tagValue, ok := field.Tag.Lookup(tag); ok {
		if f := getFunc(v.Type()); f != nil {
			if tv := f(tagValue); tv.Type() == v.Type() {
				v.Set(tv)
				ok = true
			}
		} else {
			switch v.Kind() {
			case reflect.String:
				v.Set(strFunc(tagValue))
				ok = true
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				v.SetInt(int64Func(tagValue).Int())
				ok = true
			case reflect.Bool:
				v.SetBool(boolFunc(tagValue).Bool())
				ok = true
			case reflect.Float32, reflect.Float64:
				v.SetFloat(float64Func(tagValue).Float())
				ok = true
			}
		}
	}
	return
}

func implDefaultInitValue(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		return t.Elem().Implements(reflect.TypeOf((*DefaultInitValue)(nil)).Elem())
	} else {
		return t.Implements(reflect.TypeOf((*DefaultInitValue)(nil)).Elem())
	}
}
func defaultTag(tags ...string) string {
	if len(tags) > 0 {
		return tags[0]
	} else {
		return "d"
	}
}

//default tagKey d
func BindDefaultValue(target interface{}, tags ...string) {
	bindValue(reflect.ValueOf(target), defaultTag(tags...))
}
