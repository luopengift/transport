package transport

import (
	"github.com/luopengift/golibs/channel"
	"sync"
	"sync/atomic"
)

type Adapter interface {
	Init(config Configer) error
	Handle(in, out []byte) (n int, err error)
	Version() string
}

type Codec struct {
	Name string
	cnt  uint64
	*sync.Mutex
	channel *channel.Channel
	Adapter
}

func NewCodec(name string, a Adapter, size int) *Codec {
	c := new(Codec)
	c.Name = name
	c.cnt = 0
	c.Adapter = a
	c.Mutex = new(sync.Mutex)
	c.channel = channel.NewChannel(size)
	return c
}

func (c *Codec) Init(config Configer) error {
	return nil
}

func (c *Codec) Count() uint64 {
	return c.cnt
}
func (c *Codec) Handle(in, out []byte) (int, error) {
	n, err := c.Adapter.Handle(in, out)
	c.cnt = atomic.AddUint64(&c.cnt, 1)
	return n, err
}
