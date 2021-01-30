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

var (
	_defaultTag    = "default"
	_notDefaultTag = "-"
)
var defaultTimeLayout = "2006-01-02 15:04:05"

func SetDefaultTimeLayout(layout string) {
	mux.Lock()
	defer mux.Unlock()
	defaultTimeLayout = layout
}
func SetDefaultTag(tag string) {
	mux.Lock()
	defer mux.Unlock()
	_defaultTag = tag
}

type DefaultValueFunc func(str string) reflect.Value

var regTypeFuncMap map[reflect.Type]DefaultValueFunc
var regKindFuncMap map[reflect.Kind]DefaultValueFunc
var mux sync.RWMutex

func RegisterTypeForDefaultValue(t reflect.Type, f DefaultValueFunc) {
	mux.Lock()
	defer mux.Unlock()
	regTypeFuncMap[t] = f
}

func RegisterKindForDefaultValue(t reflect.Kind, f DefaultValueFunc) {
	mux.Lock()
	defer mux.Unlock()
	regKindFuncMap[t] = f
}

func init() {
	regTypeFuncMap = make(map[reflect.Type]DefaultValueFunc)
	regKindFuncMap = make(map[reflect.Kind]DefaultValueFunc)
	RegisterTypeForDefaultValue(reflect.TypeOf(time.Time{}), dateTimeFunc)

	RegisterKindForDefaultValue(reflect.Int64, int64Func)
	RegisterKindForDefaultValue(reflect.Uint64, uint64Func)
}
func getFuncByKind(k reflect.Kind) DefaultValueFunc {
	mux.RLock()
	defer mux.RUnlock()
	return regKindFuncMap[k]
}
func getFuncByType(k reflect.Type) DefaultValueFunc {
	mux.RLock()
	defer mux.RUnlock()
	return regTypeFuncMap[k]
}
func bindValue(v reflect.Value, tag string) {
	if v == ZeroValue || v.IsZero() {
		return
	}
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		t = v.Elem().Type()
	}

	if implDefaultInitValue(t) {
		dv, _ := v.Interface().(DefaultInitValue)
		//TODO 此处无法判断 指针类型带来的接口实现 会造成 default 执行报错
		tv := dv.Default()
		if v.Type().Kind() == reflect.Ptr {
			if tv.Type() == t {
				v.Elem().Set(tv)
				return
			} else {
				if filed, ok := v.Elem().Type().FieldByName(tv.Type().Name()); ok && canDefaultBind(filed, tag) {
					v.Elem().FieldByName(tv.Type().Name()).Set(tv)
				}
			}
		} else {
			if tv.Type() == t {
				v.Set(tv)
				return
			} else {
				log.Println(tv.Type())
			}
		}

	}
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			if v.Type().Kind() == reflect.Ptr {
				if v.Elem().Field(i).Kind() != reflect.Ptr && canDefaultBind(v.Elem().Type().Field(i), tag) && !checkSampleValue(v.Elem().Field(i), v.Elem().Type().Field(i), tag) {
					bindValue(v.Elem().Field(i), tag)
				}
			} else {
				if v.Field(i).Kind() != reflect.Ptr && canDefaultBind(v.Type().Field(i), tag) && !checkSampleValue(v.Field(i), v.Type().Field(i), tag) {
					bindValue(v.Field(i), tag)
				}
			}
		}
	}
}
func canDefaultBind(field reflect.StructField, tag string) bool {
	if t, ok := field.Tag.Lookup(tag); ok {
		return t != _notDefaultTag
	}
	return true
}
func checkSampleValue(v reflect.Value, field reflect.StructField, tag string) (ok bool) {
	if tagValue, ok := field.Tag.Lookup(tag); ok {
		if tagValue == _notDefaultTag {
			return true
		} else {
			if f := getFuncByType(v.Type()); f != nil {
				if tv := f(tagValue); tv.Type() == v.Type() {
					v.Set(tv)
					ok = true
				}
			} else {
				if f := getFuncByKind(v.Kind()); f != nil {
					t := f(tagValue)
					if v.Kind() == t.Kind() {
						v.Set(t)
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
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						v.SetUint(uint64Func(tagValue).Uint())
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
		return _defaultTag
	}
}

//default tagKey d
func BindDefaultValue(target interface{}, tags ...string) {
	if reflect.TypeOf(target).Kind() == reflect.Ptr {
		bindValue(reflect.ValueOf(target), defaultTag(tags...))
	}
}
