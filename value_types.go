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

//Support computing
func int64Func(str string) reflect.Value {
	//exp, err := parser.ParseExpr(str)
	//if err == nil {
	//	fset := token.NewFileSet()
	//	_ = ast.Print(fset, exp)
	//}
	if strings.Contains(str, "*") {
		ss := strings.Split(str, "*")
		a, _ := strconv.ParseInt(strings.TrimSpace(ss[0]), 10, 64)
		b, _ := strconv.ParseInt(strings.TrimSpace(ss[1]), 10, 64)
		return reflect.ValueOf(a * b)
	} else {
		r, _ := strconv.ParseInt(str, 10, 64)
		return reflect.ValueOf(r)
	}
}
func uint64Func(str string) reflect.Value {
	r, _ := strconv.ParseUint(str, 10, 64)
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
