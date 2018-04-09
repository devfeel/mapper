package mapper

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"
)

type (
	testStruct struct {
		Name  string
		Sex   bool
		Age   int
		Time  time.Time
		Time2 JSONTime
		Time3 JSONTime
	}

	FromStruct struct {
		Name string `mapper:"UserName"`
		Sex  bool
		AA   string `mapper:"BB"`
	}

	ToStruct struct {
		Name string `mapper:"UserName"`
		Sex  bool
		BB   string
	}
)

var testValue reflect.Value

func init() {
	testValue = reflect.ValueOf(&testStruct{})
}

func Test_SetEnabledTypeChecking(t *testing.T) {
	SetEnabledTypeChecking(true)
	if enabledTypeChecking != true {
		t.Error("SetEnabledTypeChecking error: set true but query is not true")
	} else {
		t.Log("SetEnabledTypeChecking success")
	}
}

func Test_GetTypeName(t *testing.T) {
	name := GetTypeName(&testStruct{})
	if name == "" {
		t.Error("RunResult error: name is empty")
	} else {
		t.Log("RunResult success:", name)
	}
}

func BenchmarkGetTypeName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetTypeName(&testStruct{})
	}
}

func Test_GetFieldNameWithElem(t *testing.T) {
	fieldName := GetFieldName(testValue.Elem(), 0)
	if fieldName == "Name" {
		t.Log("RunResult success:", fieldName)
	} else {
		t.Error("RunResult error: fieldName not match", fieldName)
	}
}

func BenchmarkGetFieldNameWithElem(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetFieldName(testValue.Elem(), 0)
	}
}

func Test_CheckExistsField(t *testing.T) {
	Register(&testStruct{})
	fieldName := "Name"
	_, isOk := CheckExistsField(testValue.Elem(), fieldName)
	if isOk {
		t.Log("RunResult success:", fieldName)
	} else {
		t.Error("RunResult error: fieldName not exists", fieldName)
	}
}

func BenchmarkCheckExistsField(b *testing.B) {
	Register(&testStruct{})
	elem := testValue.Elem()
	fieldName := "Name"
	for i := 0; i < b.N; i++ {
		CheckExistsField(elem, fieldName)
	}
}

func Test_Mapper(t *testing.T) {
	SetEnabledTypeChecking(true)
	from := &FromStruct{Name: "From", Sex: true, AA: "AA"}
	to := &ToStruct{}
	err := Mapper(from, to)
	if err != nil {
		t.Error("RunResult error: mapper error", err)
	} else {
		t.Log("RunResult success:", to)
	}
}

func Test_MapperSlice(t *testing.T) {
	SetEnabledTypeChecking(true)
	var fromSlice []*FromStruct
	var toSlice []*ToStruct
	for i := 0; i < 10; i++ {
		fromSlice = append(fromSlice, &FromStruct{Name: "From" + strconv.Itoa(i), Sex: true, AA: "AA" + strconv.Itoa(i)})
	}
	err := MapperSlice(fromSlice, &toSlice)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(toSlice, len(toSlice))
		for _, v := range toSlice {
			fmt.Println(v)
		}
	}
}

func BenchmarkMapperSlice(b *testing.B) {
	var fromSlice []*FromStruct
	var toSlice []*ToStruct
	for i := 0; i < 10; i++ {
		fromSlice = append(fromSlice, &FromStruct{Name: "From" + strconv.Itoa(i), Sex: true, AA: "AA" + strconv.Itoa(i)})
	}
	for i := 0; i < b.N; i++ {
		MapperSlice(fromSlice, &toSlice)
	}
}

func Test_AutoMapper(t *testing.T) {
	from := &FromStruct{Name: "From", Sex: true, AA: "AA"}
	to := &ToStruct{}
	err := AutoMapper(from, to)
	if err != nil {
		t.Error("RunResult error: mapper error", err)
	} else {
		t.Log("RunResult success:", to)
	}
}

func Test_MapperMap(t *testing.T) {
	validateTime, _ := time.Parse("2006-01-02 15:04:05", "2017-01-01 10:00:00")
	fromMap := make(map[string]interface{})
	fromMap["Name"] = "test"
	fromMap["Sex"] = true
	fromMap["Age"] = 10
	fromMap["Time"] = validateTime
	fromMap["Time2"] = validateTime
	toObj := &testStruct{}
	err := MapperMap(fromMap, toObj)
	if err != nil && toObj.Time != validateTime {
		t.Error("RunResult error: mapper error", err)
	} else {
		t.Log("RunResult success:", toObj)
	}
}

func Test_MapperMapSlice(t *testing.T) {
	var toSlice []*testStruct
	fromMaps := make(map[string]map[string]interface{})
	for i := 0; i < 10; i++ {
		fromMap := make(map[string]interface{})
		fromMap["Name"] = "s" + strconv.Itoa(i)
		fromMap["Sex"] = true
		fromMap["Age"] = i
		fromMaps[strconv.Itoa(i)] = fromMap
	}
	err := MapperMapSlice(fromMaps, &toSlice)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(toSlice, len(toSlice))
	}

}

func Test_IsTimeField(t *testing.T) {
	t1 := time.Now()
	if isTimeField(reflect.ValueOf(t1)) {
		t.Log("check time.Now ok")
	} else {
		t.Error("check time.Now error")
	}

	var t2 JSONTime
	t2 = JSONTime(time.Now())
	if isTimeField(reflect.ValueOf(t2)) {
		t.Log("check mapper.Time ok")
	} else {
		t.Error("check mapper.Time error")
	}
}

func BenchmarkMapperMapSlice(b *testing.B) {
	var s []*testStruct
	fromMaps := make(map[string]map[string]interface{})
	for i := 0; i < 10; i++ {
		fromMap := make(map[string]interface{})
		fromMap["Name"] = "s" + strconv.Itoa(i)
		fromMap["Sex"] = true
		fromMap["Age"] = i
		fromMaps[strconv.Itoa(i)] = fromMap
	}
	for i := 0; i < b.N; i++ {
		MapperMapSlice(fromMaps, s)
	}
}

func BenchmarkMapper(b *testing.B) {
	Register(&FromStruct{})
	Register(&ToStruct{})
	from := &FromStruct{Name: "From", Sex: true, AA: "AA"}
	to := &ToStruct{}

	for i := 0; i < b.N; i++ {
		Mapper(from, to)
	}
}

func BenchmarkAutoMapper(b *testing.B) {
	Register(&FromStruct{})
	Register(&ToStruct{})
	from := &FromStruct{Name: "From", Sex: true, AA: "AA"}
	to := &ToStruct{}

	for i := 0; i < b.N; i++ {
		Mapper(from, to)
	}
}

func BenchmarkMapperMap(b *testing.B) {
	Register(&testStruct{})
	fromMap := make(map[string]interface{})
	fromMap["Name"] = "test"
	fromMap["Sex"] = true
	fromMap["Age"] = 10
	fromMap["time"] = time.Now()
	toObj := &testStruct{}

	for i := 0; i < b.N; i++ {
		MapperMap(fromMap, toObj)
	}
}

func BenchmarkSyncMap(b *testing.B) {
	var sMap sync.Map
	for i := 0; i < b.N; i++ {
		sMap.Load("1")
	}
}
