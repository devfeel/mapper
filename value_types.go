package mapper

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

func strFunc(str string) reflect.Value {
	return reflect.ValueOf(str)
}

func int64Func(str string) reflect.Value {
	r, _ := strconv.ParseInt(str, 10, 64)
	return reflect.ValueOf(r)
}

func boolFunc(str string) reflect.Value {
	b, _ := strconv.ParseBool(str)
	return reflect.ValueOf(b)
}

func float64Func(str string) reflect.Value {
	r, _ := strconv.ParseFloat(str, 64)
	return reflect.ValueOf(r)
}

var defaultTimeLayout = "2006-01-02 15:04:05"

func SetDefaultTimeLayout(layout string) {
	defaultTimeLayout = layout
}

func dateTimeFunc(str string) reflect.Value {
	s := strings.Split(str, ";")
	var t time.Time
	if len(s) == 1 {
		t, _ = time.ParseInLocation(defaultTimeLayout, str, time.Local)
	} else if len(s) == 2 {
		t, _ = time.ParseInLocation(s[1], s[0], time.Local)
	}
	return reflect.ValueOf(t)
}
