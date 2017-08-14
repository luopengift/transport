package pipeline

import (
	"sync"
)

// 数据输入接口
type Inputer interface {
	Init(Configer) error
	Start() error
	Read(p []byte) (n int, err error)
	Close() error
}

type Input struct {
	Name string
	Cnt  int64 //count numbers of input message
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

func (i *Input) Read(p []byte) (int, error) {
	i.Mutex.Lock()
	i.Cnt += 1
	i.Mutex.Unlock()
	return i.Inputer.Read(p)
}

func (i *Input) Start() error {
	return i.Inputer.Start()
}

func (i *Input) Close() error {
	return i.Inputer.Close()
}
