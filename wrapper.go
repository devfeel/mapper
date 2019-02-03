package mapper

import (
	"reflect"
	"time"
)

type TypeWrapper interface {
	IsType(value reflect.Value) bool
	SetNext(m TypeWrapper)
}

type BaseTypeWrapper struct {
	next TypeWrapper
}

func (bm *BaseTypeWrapper) SetNext(m TypeWrapper) {
	bm.next = m
}

type TimeWrapper struct {
	BaseTypeWrapper
}

func (w *TimeWrapper) IsType(value reflect.Value) bool {
	if _, ok := value.Interface().(time.Time); ok {
		return true
	}
	if _, ok := value.Interface().(JSONTime); ok {
		return true
	}
	return false
}

func NewTimeWrapper() *TimeWrapper {
	return &TimeWrapper{}
}
