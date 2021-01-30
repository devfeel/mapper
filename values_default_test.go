package mapper

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
	"time"
)

func init() {
	SetDefaultTag("default")
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

//func (s InfoForDefaultValueFunc) Default() reflect.Value {
//	return reflect.ValueOf(InfoForDefaultValueFunc{Name: "test name"})
//}

type InfoForDefaultValueFunc struct {
	Ext  string `default:"aaa"`
	Name string `default:"123"`
}
type Info struct {
	Name       string  `default:"123"`
	Ext        string  `default:"aaa"`
	BoolFalse  bool    `default:"false"`
	BoolTrue   bool    `default:"true"`
	Float      float32 `default:"12.21"`
	Float64    float64 `default:"12.51"`
	Int        int     `default:"12*11"`
	Int64      int64   `default:"12*11"`
	NumberUint uint    `default:"12"`
}

func (s AnyInfo) Default() reflect.Value {
	return reflect.ValueOf(AnyInfo{AnyInfoName: "test anyInfo name~~~~"})
}

type AnyInfo struct {
	AnyInfoName string `default:"dddddanyInfo"`
}
type ShowInfo struct {
	//Name string `default:"name"`
	Info
	*AnyInfo
	AnyInfo2                AnyInfo
	InfoForDefaultValueFunc InfoForDefaultValueFunc
	CreateTime              time.Time `default:"2020-01-02 03:04:01"`
}

func TestBindDefaultValue(t *testing.T) {
	t0 := time.Now()
	var s ShowInfo
	BindDefaultValue(&s)
	data, _ := json.Marshal(s)
	log.Printf("%s %v", data, time.Now().Sub(t0))
}

func BenchmarkBindDefaultValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var s ShowInfo
		BindDefaultValue(&s)
	}
}
func BenchmarkBindDefaultValue2(b *testing.B) {
	data := []byte(`{"Name":"123","Ext":"aaa","BoolFalse":false,"BoolTrue":true,"Float":12.21,"Float64":12.51,"Int":132,"Int64":132,"NumberUint":12,"AnyInfo2":{"AnyInfoName":"test anyInfo name~~~~"},"InfoForDefaultValueFunc":{"Ext":"aaa","Name":"123"},"CreateTime":"2020-01-02T03:04:01+08:00"}`)
	for i := 0; i < b.N; i++ {
		var s ShowInfo
		_ = json.Unmarshal(data, &s)
	}
}
