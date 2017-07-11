package transport

import (
	"sync"
)

// 数据输入接口
type Inputer interface {
	Init(map[string]string) error
	Read(p []byte) (n int, err error)
	Close() error

	Start() error
}

type Input struct {
	Inputer
	*sync.Mutex
}

func NewInput(in Inputer) *Input {
	i := new(Input)
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
	defer i.Mutex.Unlock()

	return i.Inputer.Read(p)
}

func (i *Input) Start() error {
	return i.Inputer.Start()
}

func (i *Input) Close() error {
	return i.Inputer.Close()
}

var InputPlugins = map[string]Inputer{}

func RegistInputer(key string, out Inputer) {
	InputPlugins[key] = out
}
