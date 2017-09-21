package transport

import (
	"sync"
	"sync/atomic"
)

// 数据输入接口, 实现了标准io库中的ReadCloser接口
type Outputer interface {
	Init(Configer) error
	Write(p []byte) (n int, err error)
	Start() error
	Close() error

	Version() string
}

type Output struct {
	Name string
	cnt  uint64
	*sync.Mutex
	Outputer
}

func NewOutput(name string, out Outputer) *Output {
	o := new(Output)
	o.Name = name
	o.cnt = 0
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

func (o *Output) Count() uint64 {
	return o.cnt
}

func (o *Output) Write(p []byte) (int, error) {
	n, err := o.Outputer.Write(p)
	atomic.AddUint64(&o.cnt, 1)
	return n, err
}

func (o *Output) Start() error {
	return o.Outputer.Start()
}

func (o *Output) Close() error {
	return o.Outputer.Close()
}
