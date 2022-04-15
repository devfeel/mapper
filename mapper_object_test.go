package mapper

import (
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"
)

func Test_Object_SetEnabledTypeChecking(t *testing.T) {
	m := NewMapper()
	m.SetEnabledTypeChecking(true)
	if m.IsEnabledTypeChecking() != true {
		t.Error("SetEnabledTypeChecking error: set true but query is not true")
	} else {
		t.Log("SetEnabledTypeChecking success")
	}
}

func Test_ObjectGetTypeName(t *testing.T) {
	m := NewMapper()
	name := m.GetTypeName(&testStruct{})
	if name == "" {
		t.Error("RunResult error: name is empty")
	} else {
		t.Log("RunResult success:", name)
	}
}

func BenchmarkObjectGetTypeName(b *testing.B) {
	m := NewMapper()
	for i := 0; i < b.N; i++ {
		m.GetTypeName(&testStruct{})
	}
}

func Test_Object_GetFieldNameFromMapperTag(t *testing.T) {
	m := NewMapper()
	v := TagStruct{}
	fieldName := m.GetFieldName(reflect.ValueOf(v), 0)
	if fieldName == "UserName" {
		t.Log("RunResult success:", fieldName)
	} else {
		t.Error("RunResult error: fieldName not match", fieldName)
	}
}

func Test_Object_GetFieldNameFromJsonTag(t *testing.T) {
	m := NewMapper()
	v := TagStruct{}
	fieldName := m.GetFieldName(reflect.ValueOf(v), 1)
	if fieldName == "UserSex" {
		t.Log("RunResult success:", fieldName)
	} else {
		t.Error("RunResult error: fieldName not match", fieldName)
	}
}

func Test_Object_SetEnableMapperTag(t *testing.T) {
	m := NewMapper()
	v := TagStruct{}
	m.SetEnabledMapperTag(false)
	fieldName := m.GetFieldName(reflect.ValueOf(v), 0)
	if fieldName == "Name" {
		t.Log("RunResult success:", fieldName)
	} else {
		t.Error("RunResult error: fieldName not match", fieldName)
	}
	m.SetEnabledMapperTag(true)
	fieldName = m.GetFieldName(reflect.ValueOf(v), 0)
	if fieldName == "UserName" {
		t.Log("RunResult success:", fieldName)
	} else {
		t.Error("RunResult error: fieldName not match", fieldName)
	}
}

func Test_Object_SetEnableJsonTag(t *testing.T) {
	m := NewMapper()
	v := TagStruct{}
	m.SetEnabledJsonTag(false)
	fieldName := m.GetFieldName(reflect.ValueOf(v), 1)
	if fieldName == "Sex" {
		t.Log("RunResult success:", fieldName)
	} else {
		t.Error("RunResult error: fieldName not match", fieldName)
	}
	m.SetEnabledJsonTag(true)
	fieldName = m.GetFieldName(reflect.ValueOf(v), 1)
	if fieldName == "UserSex" {
		t.Log("RunResult success:", fieldName)
	} else {
		t.Error("RunResult error: fieldName not match", fieldName)
	}
}

func Test_Object_GetFieldNameWithElem(t *testing.T) {
	m := NewMapper()
	fieldName := m.GetFieldName(testValue.Elem(), 0)
	if fieldName == "Name" {
		t.Log("RunResult success:", fieldName)
	} else {
		t.Error("RunResult error: fieldName not match", fieldName)
	}
}

func BenchmarkObjectGetFieldNameWithElem(b *testing.B) {
	m := NewMapper()
	for i := 0; i < b.N; i++ {
		m.GetFieldName(testValue.Elem(), 0)
	}
}

func Test_Object_CheckExistsField(t *testing.T) {
	m := NewMapper()
	m.Register(&testStruct{})
	fieldName := "Name"
	_, isOk := m.CheckExistsField(testValue.Elem(), fieldName)
	if isOk {
		t.Log("RunResult success:", fieldName)
	} else {
		t.Error("RunResult error: fieldName not exists", fieldName)
	}
}

func BenchmarkObjectCheckExistsField(b *testing.B) {
	m := NewMapper()
	m.Register(&testStruct{})
	elem := testValue.Elem()
	fieldName := "Name"
	for i := 0; i < b.N; i++ {
		m.CheckExistsField(elem, fieldName)
	}
}

func Test_Object_Mapper(t *testing.T) {
	m := NewMapper()
	m.SetEnabledTypeChecking(true)
	from := &FromStruct{Name: "From", Sex: true, AA: "AA"}
	to := &ToStruct{}
	err := m.Mapper(from, to)
	if err != nil {
		t.Error("RunResult error: mapper error", err)
	} else {
		t.Log("RunResult success:", to)
	}
}

func Test_Object_MapperSlice(t *testing.T) {
	m := NewMapper()
	m.SetEnabledTypeChecking(true)
	var fromSlice []*FromStruct
	var toSlice []*ToStruct
	for i := 0; i < 10; i++ {
		fromSlice = append(fromSlice, &FromStruct{Name: "From" + strconv.Itoa(i), Sex: true, AA: "AA" + strconv.Itoa(i)})
	}
	err := m.MapperSlice(fromSlice, &toSlice)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(toSlice, len(toSlice))
		for i := 0; i < len(fromSlice); i++ {
			if !reflect.DeepEqual(fromSlice[i].Name, toSlice[i].Name) ||
				!reflect.DeepEqual(fromSlice[i].Sex, toSlice[i].Sex) ||
				!reflect.DeepEqual(fromSlice[i].AA, toSlice[i].BB) {
				t.Fail()
			}
		}
	}
}

func Test_Object_MapperSlice2(t *testing.T) {
	m := NewMapper()
	m.SetEnabledTypeChecking(true)
	var fromSlice []*FromStruct
	var toSlice []*ToStruct
	for i := 0; i < 10; i++ {
		fromSlice = append(fromSlice, &FromStruct{Name: "From" + strconv.Itoa(i), Sex: true, AA: "AA" + strconv.Itoa(i)})
	}
	err := m.MapperSlice(&fromSlice, &toSlice)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(toSlice, len(toSlice))
		for i := 0; i < len(fromSlice); i++ {
			if !reflect.DeepEqual(fromSlice[i].Name, toSlice[i].Name) ||
				!reflect.DeepEqual(fromSlice[i].Sex, toSlice[i].Sex) ||
				!reflect.DeepEqual(fromSlice[i].AA, toSlice[i].BB) {
				t.Fail()
			}
		}
	}
}

func Test_Object_MapperStructSlice(t *testing.T) {
	m := NewMapper()
	m.SetEnabledTypeChecking(true)
	var fromSlice []FromStruct
	var toSlice []ToStruct
	for i := 0; i < 10; i++ {
		fromSlice = append(fromSlice, FromStruct{Name: "From" + strconv.Itoa(i), Sex: true, AA: "AA" + strconv.Itoa(i)})
	}
	err := m.MapperSlice(fromSlice, &toSlice)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(toSlice, len(toSlice))
		for i := 0; i < len(fromSlice); i++ {
			if !reflect.DeepEqual(fromSlice[i].Name, toSlice[i].Name) ||
				!reflect.DeepEqual(fromSlice[i].Sex, toSlice[i].Sex) ||
				!reflect.DeepEqual(fromSlice[i].AA, toSlice[i].BB) {
				t.Fail()
			}
		}
	}
}

func Test_Object_MapperStructSlice2(t *testing.T) {
	m := NewMapper()
	m.SetEnabledTypeChecking(true)
	var fromSlice []FromStruct
	var toSlice []ToStruct
	for i := 0; i < 10; i++ {
		fromSlice = append(fromSlice, FromStruct{Name: "From" + strconv.Itoa(i), Sex: true, AA: "AA" + strconv.Itoa(i)})
	}
	err := m.MapperSlice(&fromSlice, &toSlice)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(toSlice, len(toSlice))
		for i := 0; i < len(fromSlice); i++ {
			if !reflect.DeepEqual(fromSlice[i].Name, toSlice[i].Name) ||
				!reflect.DeepEqual(fromSlice[i].Sex, toSlice[i].Sex) ||
				!reflect.DeepEqual(fromSlice[i].AA, toSlice[i].BB) {
				t.Fail()
			}
		}
	}
}

func BenchmarkObjectMapperSlice(b *testing.B) {
	m := NewMapper()
	var fromSlice []*FromStruct
	var toSlice []*ToStruct
	for i := 0; i < 10; i++ {
		fromSlice = append(fromSlice, &FromStruct{Name: "From" + strconv.Itoa(i), Sex: true, AA: "AA" + strconv.Itoa(i)})
	}
	for i := 0; i < b.N; i++ {
		m.MapperSlice(fromSlice, &toSlice)
	}
}

func Test_Object_AutoMapper(t *testing.T) {
	m := NewMapper()
	from := &FromStruct{Name: "From", Sex: true, AA: "AA"}
	to := &ToStruct{}
	err := m.AutoMapper(from, to)
	if err != nil {
		t.Error("RunResult error: mapper error", err)
	} else {
		t.Log("RunResult success:", to)
	}
}

func Test_Object_AutoMapper_StructToMap(t *testing.T) {
	m := NewMapper()
	from := &FromStruct{Name: "From", Sex: true, AA: "AA"}
	to := make(map[string]interface{})
	err := m.AutoMapper(from, &to)
	if err != nil {
		t.Error("RunResult error: mapper error", err)
	} else {
		if to["UserName"] == "From" {
			t.Log("RunResult success:", to)
		} else {
			t.Error("RunResult failed: map[UserName]", to["UserName"])
		}
	}
}

func Test_Object_MapperMap(t *testing.T) {
	m := NewMapper()
	validateTime, _ := time.Parse("2006-01-02 15:04:05", "2017-01-01 10:00:00")
	fromMap := make(map[string]interface{})
	fromMap["Name"] = "test"
	fromMap["Sex"] = true
	fromMap["Age"] = 10
	fromMap["Time"] = validateTime
	fromMap["Time2"] = validateTime
	toObj := &testStruct{}
	err := m.MapperMap(fromMap, toObj)
	if err != nil && toObj.Time != validateTime {
		t.Error("RunResult error: mapper error", err)
	} else {
		t.Log("RunResult success:", toObj)
	}
}

func Test_Object_MapToSlice(t *testing.T) {
	m := NewMapper()
	var toSlice []*testStruct
	fromMaps := make(map[string]interface{})
	for i := 0; i < 10; i++ {
		fromMap := make(map[string]interface{})
		fromMap["Name"] = "s" + strconv.Itoa(i)
		fromMap["Sex"] = true
		fromMap["Age"] = i
		fromMaps[strconv.Itoa(i)] = fromMap
	}
	err := m.MapToSlice(fromMaps, &toSlice)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(toSlice, len(toSlice))
	}
}

func Test_Object_MapperMapSlice(t *testing.T) {
	m := NewMapper()
	var toSlice []*testStruct
	fromMaps := make(map[string]map[string]interface{})
	for i := 0; i < 10; i++ {
		fromMap := make(map[string]interface{})
		fromMap["Name"] = "s" + strconv.Itoa(i)
		fromMap["Sex"] = true
		fromMap["Age"] = i
		fromMaps[strconv.Itoa(i)] = fromMap
	}
	err := m.MapperMapSlice(fromMaps, &toSlice)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(toSlice, len(toSlice))
	}
}

func Test_Object_MapperStructMapSlice(t *testing.T) {
	m := NewMapper()
	var toSlice []testStruct
	fromMaps := make(map[string]map[string]interface{})
	for i := 0; i < 10; i++ {
		fromMap := make(map[string]interface{})
		fromMap["Name"] = "s" + strconv.Itoa(i)
		fromMap["Sex"] = true
		fromMap["Age"] = i
		fromMaps[strconv.Itoa(i)] = fromMap
	}
	err := m.MapperMapSlice(fromMaps, &toSlice)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(toSlice, len(toSlice))
	}

}

func Test_Object_IsTimeField(t *testing.T) {
	m := NewMapper()
	t1 := time.Now()
	if m.GetDefaultTimeWrapper().IsType(reflect.ValueOf(t1)) {
		t.Log("check time.Now ok")
	} else {
		t.Error("check time.Now error")
	}

	var t2 JSONTime
	t2 = JSONTime(time.Now())
	if m.GetDefaultTimeWrapper().IsType(reflect.ValueOf(t2)) {
		t.Log("check mapper.Time ok")
	} else {
		t.Error("check mapper.Time error")
	}
}

func Test_Object_MapToJson_JsonToMap(t *testing.T) {
	m := NewMapper()
	fromMap := createMap()
	data, err := m.MapToJson(fromMap)
	if err != nil {
		t.Error("MapToJson error", err)
	} else {
		var retMap map[string]interface{}
		err = m.JsonToMap(data, &retMap)
		if err != nil {
			t.Error("MapToJson.JsonToMap error", err)
		}
		if len(retMap) != len(fromMap) {
			t.Error("MapToJson failed, not match length")
		}
		t.Log("MapToJson success", fromMap, retMap)
	}
}

func BenchmarkObjectMapperMapSlice(b *testing.B) {
	m := NewMapper()
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
		m.MapperMapSlice(fromMaps, s)
	}
}

func BenchmarkObjectMapper(b *testing.B) {
	m := NewMapper()
	m.Register(&FromStruct{})
	m.Register(&ToStruct{})
	from := &FromStruct{Name: "From", Sex: true, AA: "AA"}
	to := &ToStruct{}

	for i := 0; i < b.N; i++ {
		m.Mapper(from, to)
	}
}

func BenchmarkObjectAutoMapper(b *testing.B) {
	m := NewMapper()
	m.Register(&FromStruct{})
	m.Register(&ToStruct{})
	from := &FromStruct{Name: "From", Sex: true, AA: "AA"}
	to := &ToStruct{}

	for i := 0; i < b.N; i++ {
		m.Mapper(from, to)
	}
}

func BenchmarkObjectAutoMapper_Map(b *testing.B) {
	m := NewMapper()
	from := &FromStruct{Name: "From", Sex: true, AA: "AA"}
	to := make(map[string]interface{})

	for i := 0; i < b.N; i++ {
		m.Mapper(from, &to)
	}
}

func BenchmarkObjectMapperMap(b *testing.B) {
	m := NewMapper()
	m.Register(&testStruct{})
	fromMap := make(map[string]interface{})
	fromMap["Name"] = "test"
	fromMap["Sex"] = true
	fromMap["Age"] = 10
	fromMap["time"] = time.Now()
	toObj := &testStruct{}

	for i := 0; i < b.N; i++ {
		m.MapperMap(fromMap, toObj)
	}
}

func BenchmarkObjectSyncMap(b *testing.B) {
	var sMap sync.Map
	for i := 0; i < b.N; i++ {
		sMap.Load("1")
	}
}
