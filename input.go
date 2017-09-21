package transport

import (
	"sync"
	"sync/atomic"
)

// 数据输入接口
type Inputer interface {
	Init(Configer) error
	Start() error
	Read(p []byte) (n int, err error)
	Close() error
	Version() string
}

type Input struct {
	Name string
	cnt  uint64 //count numbers of input message
	*sync.Mutex
	Inputer
}

func NewInput(name string, in Inputer) *Input {
	i := new(Input)
	i.Name = name
	i.Inputer = in
	i.Mutex = new(sync.Mutex)
	return i
}

func (i *Input) Set(in Inputer) error {
	i.Mutex.Lock()
	defer i.Mutex.Unlock()
	if err := i.Inputer.Close(); err != nil {
		return err
	}
	i.Inputer = in
	return nil
}

func (i *Input) Count() uint64 {
	return i.cnt
}

func (i *Input) Read(p []byte) (int, error) {
	n, err := i.Inputer.Read(p)
	atomic.AddUint64(&i.cnt, 1)
	return n, err
}

func (i *Input) Start() error {
	return i.Inputer.Start()
}

func (i *Input) Close() error {
	return i.Inputer.Close()
}
func (i *Input) Version() string {
	return ""
}
