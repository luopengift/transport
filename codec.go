package transport

import (
	"github.com/luopengift/golibs/channel"
	"sync"
)

type Handler interface {
	Init(config Configer) error
	Handle(in, out []byte) (n int, err error)
}

type Codec struct {
	Name string
	Cnt  int64
	Handler
	*sync.Mutex
	*channel.Channel
}

func NewCodec(name string, h Handler, max int) *Codec {
	m := new(Codec)
	m.Cnt = 0
	m.Name = name
	m.Handler = h
	m.Mutex = new(sync.Mutex)
	m.Channel = channel.NewChannel(max)
	return m
}

func (m *Codec) Init(config Configer) error {
	return nil
}

func (m *Codec) Handle(in, out []byte) (int, error) {
	n, err := m.Handler.Handle(in, out)
	m.Mutex.Lock()
	m.Cnt ++
	m.Mutex.Unlock()
    return n, err
}
