package itask

import (
	"errors"
	"reflect"
	"runtime"
)

var (
	ErrTypeNotFunction = errors.New("argument type not function")
	ErrInArgsMissMatch = errors.New("input arguments count not match")
	ErrOutCntMissMatch = errors.New("output parameter count not match")
	ErrExecuteTimeout  = errors.New("parallel execute timeout")
)

func GetFuncName(h *Handler) string {
	if h == nil || h.f == nil {
		return ""
	}
	return runtime.FuncForPC(reflect.ValueOf(h.f).Pointer()).Name()
}
