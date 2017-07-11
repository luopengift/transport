package transport

import (
	"sync"
)

// 数据输入接口, 实现了标准io库中的ReadCloser接口
type Outputer interface {
	Init(map[string]string) error
	Write(p []byte) (n int, err error)
	Close() error
	Start() error
}

type Output struct {
	Outputer
	*sync.Mutex
}

func NewOutput(out Outputer) *Output {
	o := new(Output)
	o.Outputer = out
	o.Mutex = new(sync.Mutex)
	return o
}

func (o *Output) Set(out Outputer) error {
	o.Outputer = out
	return nil
}

func (o *Output) Write(p []byte) (int, error) {
	return o.Outputer.Write(p)
}

func (o *Output) Start() error {
	return o.Outputer.Start()
}

func (o *Output) Close() error {
	return o.Outputer.Close()
}

// 输出组件列表
var OutputPlugins = map[string]Outputer{}

func RegistOutputer(key string, out Outputer) {
	OutputPlugins[key] = out
}
