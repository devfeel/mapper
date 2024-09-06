package mapper

import (
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

	TagStruct struct {
		Name string `mapper:"UserName"`
		Sex  bool   `json:"UserSex"`
	}
)

var testValue reflect.Value

func init() {
	testValue = reflect.ValueOf(&testStruct{})
}

func TestPackageVersion(t *testing.T) {
	v := PackageVersion()
	if v != packageVersion {
		t.Error("PackageVersion error: not equal with packageVersion[" + packageVersion + "]")
	} else {
		t.Log("PackageVersion success")
	}
}

func Test_CheckIsTypeWrapper(t *testing.T) {
	v := TagStruct{}
	if standardMapper.CheckIsTypeWrapper(reflect.ValueOf(v)) == true {
		t.Error("CheckIsTypeWrapper error: set true but query is not true")
	} else {
		t.Log("CheckIsTypeWrapper success")
	}
}

func Test_SetEnabledMapperStructField(t *testing.T) {
	SetEnabledMapperStructField(true)
	if standardMapper.IsEnabledMapperStructField() != true {
		t.Error("SetEnabledMapperStructField error: set true but query is not true")
	} else {
		t.Log("SetEnabledMapperStructField success")
	}
}

func Test_SetEnabledAutoTypeConvert(t *testing.T) {
	SetEnabledAutoTypeConvert(true)
	if standardMapper.IsEnabledAutoTypeConvert() != true {
		t.Error("SetEnabledAutoTypeConvert error: set true but query is not true")
	} else {
		t.Log("SetEnabledAutoTypeConvert success")
	}
}

func Test_SetEnabledTypeChecking(t *testing.T) {
	SetEnabledTypeChecking(true)
	if standardMapper.IsEnabledTypeChecking() != true {
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

func Test_GetFieldNameFromMapperTag(t *testing.T) {
	v := TagStruct{}
	fieldName := GetFieldName(reflect.ValueOf(v), 0)
	if fieldName == "UserName" {
		t.Log("RunResult success:", fieldName)
	} else {
		t.Error("RunResult error: fieldName not match", fieldName)
	}
}

func Test_GetFieldNameFromJsonTag(t *testing.T) {
	v := TagStruct{}
	fieldName := GetFieldName(reflect.ValueOf(v), 1)
	if fieldName == "UserSex" {
		t.Log("RunResult success:", fieldName)
	} else {
		t.Error("RunResult error: fieldName not match", fieldName)
	}
}

func Test_SetEnableMapperTag(t *testing.T) {
	v := TagStruct{}
	SetEnabledMapperTag(false)
	fieldName := GetFieldName(reflect.ValueOf(v), 0)
	if fieldName == "Name" {
		t.Log("RunResult success:", fieldName)
	} else {
		t.Error("RunResult error: fieldName not match", fieldName)
	}
	SetEnabledMapperTag(true)
	fieldName = GetFieldName(reflect.ValueOf(v), 0)
	if fieldName == "UserName" {
		t.Log("RunResult success:", fieldName)
	} else {
		t.Error("RunResult error: fieldName not match", fieldName)
	}
}

func Test_SetEnableJsonTag(t *testing.T) {
	v := TagStruct{}
	SetEnabledJsonTag(false)
	fieldName := GetFieldName(reflect.ValueOf(v), 1)
	if fieldName == "Sex" {
		t.Log("RunResult success:", fieldName)
	} else {
		t.Error("RunResult error: fieldName not match", fieldName)
	}
	SetEnabledJsonTag(true)
	fieldName = GetFieldName(reflect.ValueOf(v), 1)
	if fieldName == "UserSex" {
		t.Log("RunResult success:", fieldName)
	} else {
		t.Error("RunResult error: fieldName not match", fieldName)
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
		for i := 0; i < len(fromSlice); i++ {
			if !reflect.DeepEqual(fromSlice[i].Name, toSlice[i].Name) ||
				!reflect.DeepEqual(fromSlice[i].Sex, toSlice[i].Sex) ||
				!reflect.DeepEqual(fromSlice[i].AA, toSlice[i].BB) {
				t.Fail()
			}
		}
	}
}

func Test_MapperSlice2(t *testing.T) {
	SetEnabledTypeChecking(true)
	var fromSlice []*FromStruct
	var toSlice []*ToStruct
	for i := 0; i < 10; i++ {
		fromSlice = append(fromSlice, &FromStruct{Name: "From" + strconv.Itoa(i), Sex: true, AA: "AA" + strconv.Itoa(i)})
	}
	err := MapperSlice(&fromSlice, &toSlice)
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

func Test_MapperStructSlice(t *testing.T) {
	SetEnabledTypeChecking(true)
	var fromSlice []FromStruct
	var toSlice []ToStruct
	for i := 0; i < 10; i++ {
		fromSlice = append(fromSlice, FromStruct{Name: "From" + strconv.Itoa(i), Sex: true, AA: "AA" + strconv.Itoa(i)})
	}
	err := MapperSlice(fromSlice, &toSlice)
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

func Test_MapperStructSlice2(t *testing.T) {
	SetEnabledTypeChecking(true)
	var fromSlice []FromStruct
	var toSlice []ToStruct
	for i := 0; i < 10; i++ {
		fromSlice = append(fromSlice, FromStruct{Name: "From" + strconv.Itoa(i), Sex: true, AA: "AA" + strconv.Itoa(i)})
	}
	err := MapperSlice(&fromSlice, &toSlice)
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

func Test_AutoMapper_StructToMap(t *testing.T) {
	from := &FromStruct{Name: "From", Sex: true, AA: "AA"}
	to := make(map[string]interface{})
	err := AutoMapper(from, &to)
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

func Test_MapToSlice(t *testing.T) {
	var toSlice []*testStruct
	/*fromMaps := make(map[string]interface{})
	for i := 0; i < 10; i++ {
		from := new(testStruct)
		from.Name = "s" + strconv.Itoa(i)
		from.Sex = true
		from.Age = i
		fromMaps[strconv.Itoa(i)] = from
	}*/
	fromMaps := make(map[string]interface{})
	for i := 0; i < 10; i++ {
		fromMap := make(map[string]interface{})
		fromMap["Name"] = "s" + strconv.Itoa(i)
		fromMap["Sex"] = true
		fromMap["Age"] = i
		fromMaps[strconv.Itoa(i)] = fromMap
	}
	err := MapToSlice(fromMaps, &toSlice)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(toSlice, len(toSlice))
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

func Test_MapperStructMapSlice(t *testing.T) {
	var toSlice []testStruct
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
	if standardMapper.GetDefaultTimeWrapper().IsType(reflect.ValueOf(t1)) {
		t.Log("check time.Now ok")
	} else {
		t.Error("check time.Now error")
	}

	var t2 JSONTime
	t2 = JSONTime(time.Now())
	if standardMapper.GetDefaultTimeWrapper().IsType(reflect.ValueOf(t2)) {
		t.Log("check mapper.Time ok")
	} else {
		t.Error("check mapper.Time error")
	}
}

func Test_MapToJson_JsonToMap(t *testing.T) {
	fromMap := createMap()
	data, err := MapToJson(fromMap)
	if err != nil {
		t.Error("MapToJson error", err)
	} else {
		var retMap map[string]interface{}
		err := JsonToMap(data, &retMap)
		if err != nil {
			t.Error("MapToJson.JsonToMap error", err)
		}
		if len(retMap) != len(fromMap) {
			t.Error("MapToJson failed, not match length")
		}
		t.Log("MapToJson success", fromMap, retMap)
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

func BenchmarkAutoMapper_Map(b *testing.B) {
	from := &FromStruct{Name: "From", Sex: true, AA: "AA"}
	to := make(map[string]interface{})

	for i := 0; i < b.N; i++ {
		Mapper(from, &to)
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

func createMap() map[string]interface{} {
	validateTime, _ := time.Parse("2006-01-02 15:04:05", "2017-01-01 10:00:00")
	fromMap := make(map[string]interface{})
	fromMap["Name"] = "test"
	fromMap["Sex"] = true
	fromMap["Age"] = 10
	fromMap["Time"] = validateTime
	return fromMap
}
