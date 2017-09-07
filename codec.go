package transport

import (
	"github.com/luopengift/golibs/channel"
	"sync"
)

type Adapter interface {
	Init(config Configer) error
	Handle(in, out []byte) (n int, err error)
}

type Codec struct {
	Name string
	Cnt  int64
	*sync.Mutex
	channel *channel.Channel
	Adapter
}

func NewCodec(name string, a Adapter, size int) *Codec {
	c := new(Codec)
	c.Cnt = 0
	c.Name = name
	c.Adapter = a
	c.Mutex = new(sync.Mutex)
	c.channel = channel.NewChannel(size)
	return c
}

func (c *Codec) Init(config Configer) error {
	return nil
}

func (c *Codec) Handle(in, out []byte) (int, error) {
	n, err := c.Adapter.Handle(in, out)
	c.Mutex.Lock()
	c.Cnt++
	c.Mutex.Unlock()
	return n, err
}
