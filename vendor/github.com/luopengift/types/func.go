package types

import (
	"fmt"
	"reflect"
	"time"
)

/*
  Example:
    type Test struct {
        i   int
        j   int
    }

    func (t *Test) Add(i int) (int,error) {
        return t.i+t.j+i, nil
    }

    var T = &Test{1,2}

    CallMethodName(T,"Add",1) // [4 <nil>] <nil>
*/

func CallMethodName(class interface{}, fun string, args ...interface{}) (List, error) {
	value := reflect.ValueOf(class)
	method := value.MethodByName(fun)
	if bool(method.Kind() != reflect.Func) {
		return nil, fmt.Errorf("%s is %v can not callable", fun, method.Kind())
	}
	numIn := method.Type().NumIn()
	argsIn := make([]reflect.Value, numIn)
	for i := 0; i < numIn; i++ {
		argsIn[i] = reflect.ValueOf(args[i])
	}
	numOut := method.Type().NumOut()
	argsOut := method.Call(argsIn)
	rets := make(List, numOut)
	for i := 0; i < numOut; i++ {
		rets[i] = argsOut[i].Interface()
	}
	return rets, nil
}

type CallFuncType = func(interface{}, ...interface{}) (List, error)

func CallFuncName(fun interface{}, args ...interface{}) (List, error) {
	fn := reflect.ValueOf(fun)
	if fn.Kind() != reflect.Func {
		return nil, fmt.Errorf("The first argument %v is not the function", fun)
	}
	numIn := fn.Type().NumIn()
	argsIn := make([]reflect.Value, numIn)
	for i := 0; i < numIn; i++ {
		argsIn[i] = reflect.ValueOf(args[i])
	}
	numOut := fn.Type().NumOut()
	argsOut := fn.Call(argsIn)
	rets := make(List, numOut)
	for i := 0; i < numOut; i++ {
		rets[i] = argsOut[i].Interface()
	}
	return rets, nil
}

type Result struct {
	result List
	err    error
}

func FuncWithTimeout(timeout int, fun interface{}, args ...interface{}) (List, error) {
	result := make(chan Result, 1)
	go func() {
		ret, err := CallFuncName(fun, args...)
		result <- Result{ret, err}
	}()
	select {
	case res := <-result:
		return res.result, res.err
	case <-time.After(time.Duration(timeout) * time.Second):
		return nil, fmt.Errorf("timeout")
	}
}
