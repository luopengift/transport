package pipeline

import (
	"sync"
)

// 数据输入接口, 实现了标准io库中的ReadCloser接口
type Outputer interface {
	Init(Configer) error
	Write(p []byte) (n int, err error)
	Close() error
	Start() error
}

type Output struct {
	Name string
	Outputer
	*sync.Mutex
}

func NewOutput(name string, out Outputer) *Output {
	o := new(Output)
	o.Name = name
	o.Outputer = out
	o.Mutex = new(sync.Mutex)
	return o
}

func (o *Output) Set(out Outputer) error {
	if err := o.Outputer.Close(); err != nil {
		return err
	}
	o.Outputer = out
	return nil

}
/*
func (o *Output) Write(p []byte) (int, error) {
	return o.Outputer.Write(p)
}

func (o *Output) Start() error {
	return o.Outputer.Start()
}

func (o *Output) Close() error {
	return o.Outputer.Close()
}
*/
