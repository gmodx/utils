package timer

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"time"
)

var (
	ErrParamsNotAdapted = errors.New("the number of params is not adapted")
	ErrNotAFunction     = errors.New("only functions can be schedule into the job queue")
)

func Tick(delay time.Duration, d time.Duration, jobFunc interface{}, jobErrCallback func(error), params ...interface{}) error {
	<-time.After(delay)

	typ := reflect.TypeOf(jobFunc)
	if typ.Kind() != reflect.Func {
		return ErrNotAFunction
	}

	ticker := time.NewTicker(d)
	go func() {
		for ; true; <-ticker.C { // invoke immediately
			_, err := invokeWithParams(jobFunc, params)
			if err != nil {
				jobErrCallback(fmt.Errorf("Call %v job error, err: %w", getFuncName(jobFunc), err))
				continue
			}
		}
	}()

	return nil
}

func invokeWithParams(jobFunc interface{}, params []interface{}) ([]reflect.Value, error) {
	f := reflect.ValueOf(jobFunc)
	if len(params) != f.Type().NumIn() {
		return nil, ErrParamsNotAdapted
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	return f.Call(in), nil
}

func getFuncName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}
