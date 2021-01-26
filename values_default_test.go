package mapper

import (
	"log"
	"reflect"
	"testing"
	"time"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func (s Info) Default() reflect.Value {
	return reflect.ValueOf(Info{Name: "test name"})
}

//for test not same struct
type Info2 struct {
	Ext  string `d:"aaa"`
	Name string `d:"123"`
}
type Info struct {
	Ext  string `d:"aaa"`
	Name string `d:"123"`
}
type ShowInfo struct {
	Name string `d:"name"`
	Info
	CreateTime time.Time `d:"2020-01-02 03:04:01"`
}

func TestBindDefaultValue(t *testing.T) {
	var s ShowInfo
	BindDefaultValue(&s)
	log.Printf("%+v", s)
}

func BenchmarkBindDefaultValue(b *testing.B) {
	var s ShowInfo
	for i := 0; i < b.N; i++ {
		BindDefaultValue(&s)
	}
}
